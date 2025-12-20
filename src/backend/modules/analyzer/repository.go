package analyzer

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
)

const (
	HostGitHub = "github.com"
)

type Repository interface {
	GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*CompletedAnalysis, error)
	GetTestSuitesWithCases(ctx context.Context, analysisID string) ([]TestSuiteWithCases, error)
	UpdateLastViewed(ctx context.Context, owner, repo string) error
	FindActiveRiverJobByRepo(ctx context.Context, kind, owner, repo string) (*RiverJobInfo, error)
}

type CompletedAnalysis struct {
	ID          string
	Owner       string
	Repo        string
	CommitSHA   string
	CompletedAt time.Time
	TotalSuites int
	TotalTests  int
}

type TestSuiteWithCases struct {
	FilePath  string
	Framework string
	ID        string
	Name      string
	Tests     []TestCaseRow
}

type TestCaseRow struct {
	Name   string
	Line   int
	Status string
}

type RiverJobInfo struct {
	CommitSHA string
	State     string
}

type repositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &repositoryImpl{queries: queries}
}

func (r *repositoryImpl) GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*CompletedAnalysis, error) {
	row, err := r.queries.GetLatestCompletedAnalysis(ctx, db.GetLatestCompletedAnalysisParams{
		Host:  HostGitHub,
		Owner: owner,
		Name:  repo,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.WrapNotFound(owner, repo)
		}
		return nil, fmt.Errorf("get latest completed analysis for %s/%s: %w", owner, repo, err)
	}

	return &CompletedAnalysis{
		ID:          uuidToString(row.ID),
		Owner:       row.Owner,
		Repo:        row.Repo,
		CommitSHA:   row.CommitSha,
		CompletedAt: row.CompletedAt.Time,
		TotalSuites: int(row.TotalSuites),
		TotalTests:  int(row.TotalTests),
	}, nil
}

func (r *repositoryImpl) GetTestSuitesWithCases(ctx context.Context, analysisID string) ([]TestSuiteWithCases, error) {
	uuid, err := stringToUUID(analysisID)
	if err != nil {
		return nil, fmt.Errorf("parse analysis ID: %w", err)
	}

	suiteRows, err := r.queries.GetTestSuitesByAnalysisID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get test suites: %w", err)
	}

	if len(suiteRows) == 0 {
		return []TestSuiteWithCases{}, nil
	}

	suiteIDs := make([]pgtype.UUID, len(suiteRows))
	for i, s := range suiteRows {
		suiteIDs[i] = s.ID
	}

	testRows, err := r.queries.GetTestCasesBySuiteIDs(ctx, suiteIDs)
	if err != nil {
		return nil, fmt.Errorf("get test cases: %w", err)
	}

	testsBySuite := make(map[string][]TestCaseRow)
	for _, t := range testRows {
		suiteID := uuidToString(t.SuiteID)
		line := 0
		if t.LineNumber.Valid {
			line = int(t.LineNumber.Int32)
		}
		testsBySuite[suiteID] = append(testsBySuite[suiteID], TestCaseRow{
			Name:   t.Name,
			Line:   line,
			Status: string(t.Status),
		})
	}

	suites := make([]TestSuiteWithCases, len(suiteRows))
	for i, s := range suiteRows {
		suiteID := uuidToString(s.ID)
		framework := ""
		if s.Framework.Valid {
			framework = s.Framework.String
		}
		suites[i] = TestSuiteWithCases{
			FilePath:  s.FilePath,
			Framework: framework,
			ID:        suiteID,
			Name:      s.Name,
			Tests:     testsBySuite[suiteID],
		}
	}

	return suites, nil
}

func (r *repositoryImpl) UpdateLastViewed(ctx context.Context, owner, repo string) error {
	return r.queries.UpdateCodebaseLastViewed(ctx, db.UpdateCodebaseLastViewedParams{
		Host:  HostGitHub,
		Owner: owner,
		Name:  repo,
	})
}

func (r *repositoryImpl) FindActiveRiverJobByRepo(ctx context.Context, kind, owner, repo string) (*RiverJobInfo, error) {
	row, err := r.queries.FindActiveRiverJobByRepo(ctx, db.FindActiveRiverJobByRepoParams{
		Kind:  kind,
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find active river job by repo: %w", err)
	}

	return &RiverJobInfo{
		CommitSHA: row.CommitSha,
		State:     row.State,
	}, nil
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	// pgtype.UUID.Bytes is [16]byte, format as UUID string
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func stringToUUID(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(s); err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}
