package analyzer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/specvital/web/src/backend/internal/db"
)

const (
	HostGitHub = "github.com"
)

var (
	ErrNotFound = errors.New("analysis not found")
)

type Repository interface {
	CreatePendingAnalysis(ctx context.Context, owner, repo string) (string, error)
	GetAnalysisStatus(ctx context.Context, owner, repo string) (*AnalysisStatus, error)
	GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*CompletedAnalysis, error)
	GetTestSuitesWithCases(ctx context.Context, analysisID string) ([]TestSuiteWithCases, error)
	MarkAnalysisFailed(ctx context.Context, analysisID, errorMsg string) error
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
	ID        string
	FilePath  string
	Framework string
	Tests     []TestCaseRow
}

type TestCaseRow struct {
	Name   string
	Line   int
	Status string
}

type AnalysisStatus struct {
	ID           string
	Status       string
	ErrorMessage *string
	CreatedAt    time.Time
	CompletedAt  *time.Time
}

type repositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &repositoryImpl{queries: queries}
}

func (r *repositoryImpl) CreatePendingAnalysis(ctx context.Context, owner, repo string) (string, error) {
	codebaseID, err := r.queries.UpsertCodebase(ctx, db.UpsertCodebaseParams{
		Host:  HostGitHub,
		Owner: owner,
		Name:  repo,
	})
	if err != nil {
		return "", fmt.Errorf("upsert codebase for %s/%s: %w", owner, repo, err)
	}

	analysisID, err := r.queries.CreatePendingAnalysis(ctx, codebaseID)
	if err != nil {
		return "", fmt.Errorf("create pending analysis for %s/%s: %w", owner, repo, err)
	}

	return uuidToString(analysisID), nil
}

func (r *repositoryImpl) GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*CompletedAnalysis, error) {
	row, err := r.queries.GetLatestCompletedAnalysis(ctx, db.GetLatestCompletedAnalysisParams{
		Host:  HostGitHub,
		Owner: owner,
		Name:  repo,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
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

func (r *repositoryImpl) GetAnalysisStatus(ctx context.Context, owner, repo string) (*AnalysisStatus, error) {
	row, err := r.queries.GetAnalysisStatus(ctx, db.GetAnalysisStatusParams{
		Host:  HostGitHub,
		Owner: owner,
		Name:  repo,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get analysis status for %s/%s: %w", owner, repo, err)
	}

	status := &AnalysisStatus{
		ID:        uuidToString(row.ID),
		Status:    string(row.Status),
		CreatedAt: row.CreatedAt.Time,
	}

	if row.ErrorMessage.Valid {
		status.ErrorMessage = &row.ErrorMessage.String
	}
	if row.CompletedAt.Valid {
		status.CompletedAt = &row.CompletedAt.Time
	}

	return status, nil
}

func (r *repositoryImpl) MarkAnalysisFailed(ctx context.Context, analysisID, errorMsg string) error {
	uuid, err := stringToUUID(analysisID)
	if err != nil {
		return fmt.Errorf("parse analysis ID: %w", err)
	}

	err = r.queries.MarkAnalysisFailed(ctx, db.MarkAnalysisFailedParams{
		ID:           uuid,
		ErrorMessage: pgtype.Text{String: errorMsg, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("mark analysis failed: %w", err)
	}
	return nil
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
			ID:        suiteID,
			FilePath:  s.FilePath,
			Framework: framework,
			Tests:     testsBySuite[suiteID],
		}
	}

	return suites, nil
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
