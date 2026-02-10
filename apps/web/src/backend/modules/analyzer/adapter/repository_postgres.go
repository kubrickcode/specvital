package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kubrickcode/specvital/apps/web/src/backend/internal/db"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/domain/port"
)

var _ port.Repository = (*PostgresRepository)(nil)

const HostGitHub = "github.com"

type PostgresRepository struct {
	queries *db.Queries
}

func NewPostgresRepository(queries *db.Queries) *PostgresRepository {
	return &PostgresRepository{queries: queries}
}

func (r *PostgresRepository) CheckAnalysisExistsByCommitSHA(ctx context.Context, owner, repo, commitSHA string) (bool, error) {
	exists, err := r.queries.CheckAnalysisExistsByCommitSHA(ctx, db.CheckAnalysisExistsByCommitSHAParams{
		Host:      HostGitHub,
		Owner:     owner,
		Name:      repo,
		CommitSha: commitSHA,
	})
	if err != nil {
		return false, fmt.Errorf("check analysis exists by commit SHA: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) GetAiSpecSummaries(ctx context.Context, codebaseIDs []string, userID string) (map[string]*entity.AiSpecSummary, error) {
	if len(codebaseIDs) == 0 || userID == "" {
		return make(map[string]*entity.AiSpecSummary), nil
	}

	userUUID, err := stringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	codebaseUUIDs := make([]pgtype.UUID, len(codebaseIDs))
	for i, id := range codebaseIDs {
		uuid, err := stringToUUID(id)
		if err != nil {
			return nil, fmt.Errorf("parse codebase ID %s: %w", id, err)
		}
		codebaseUUIDs[i] = uuid
	}

	rows, err := r.queries.GetAiSpecSummariesByCodebaseIDs(ctx, db.GetAiSpecSummariesByCodebaseIDsParams{
		CodebaseIds: codebaseUUIDs,
		UserID:      userUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("get AI Spec summaries: %w", err)
	}

	result := make(map[string]*entity.AiSpecSummary, len(rows))
	for _, row := range rows {
		codebaseID := uuidToString(row.CodebaseID)
		summary := &entity.AiSpecSummary{
			HasSpec:       true,
			LanguageCount: int(row.LanguageCount),
		}
		// LatestGeneratedAt from sqlc is interface{} - handle both time.Time and pgtype.Timestamptz
		switch ts := row.LatestGeneratedAt.(type) {
		case time.Time:
			summary.LatestGeneratedAt = &ts
		case pgtype.Timestamptz:
			if ts.Valid {
				summary.LatestGeneratedAt = &ts.Time
			}
		}
		result[codebaseID] = summary
	}

	return result, nil
}

func (r *PostgresRepository) FindActiveRiverJobByRepo(ctx context.Context, kind, owner, repo string) (*port.RiverJobInfo, error) {
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

	info := &port.RiverJobInfo{
		CommitSHA: row.CommitSha,
		State:     row.State,
	}
	if row.AttemptedAt.Valid {
		info.AttemptedAt = &row.AttemptedAt.Time
	}
	return info, nil
}

func (r *PostgresRepository) GetAnalysisHistory(ctx context.Context, owner, repo string) ([]port.AnalysisHistoryItem, error) {
	rows, err := r.queries.GetCompletedAnalysesByCodebase(ctx, db.GetCompletedAnalysesByCodebaseParams{
		Host:  HostGitHub,
		Owner: owner,
		Name:  repo,
	})
	if err != nil {
		return nil, fmt.Errorf("get analysis history for %s/%s: %w", owner, repo, err)
	}

	items := make([]port.AnalysisHistoryItem, len(rows))
	for i, row := range rows {
		var branchName *string
		if row.BranchName.Valid {
			branchName = &row.BranchName.String
		}
		var committedAt *time.Time
		if row.CommittedAt.Valid {
			t := row.CommittedAt.Time
			committedAt = &t
		}

		items[i] = port.AnalysisHistoryItem{
			BranchName:  branchName,
			CommitSHA:   row.CommitSha,
			CommittedAt: committedAt,
			CompletedAt: row.CompletedAt.Time,
			ID:          uuidToString(row.ID),
			TotalTests:  int(row.TotalTests),
		}
	}

	return items, nil
}

func (r *PostgresRepository) GetBookmarkedCodebaseIDs(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return nil, nil
	}

	uuid, err := stringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	rows, err := r.queries.GetBookmarkedCodebaseIDsByUserID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get bookmarked codebase IDs: %w", err)
	}

	ids := make([]string, len(rows))
	for i, row := range rows {
		ids[i] = uuidToString(row)
	}
	return ids, nil
}

func (r *PostgresRepository) GetCodebaseID(ctx context.Context, owner, repo string) (string, error) {
	id, err := r.queries.GetCodebaseIDByOwnerRepo(ctx, db.GetCodebaseIDByOwnerRepoParams{
		Host:  HostGitHub,
		Owner: owner,
		Name:  repo,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.WrapNotFound(owner, repo)
		}
		return "", fmt.Errorf("get codebase ID for %s/%s: %w", owner, repo, err)
	}
	return uuidToString(id), nil
}

func (r *PostgresRepository) GetCompletedAnalysisByCommitSHA(ctx context.Context, owner, repo, commitSHA string) (*port.CompletedAnalysis, error) {
	row, err := r.queries.GetCompletedAnalysisByCommitSHA(ctx, db.GetCompletedAnalysisByCommitSHAParams{
		Host:            HostGitHub,
		Owner:           owner,
		Name:            repo,
		CommitShaPrefix: commitSHA,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.WrapNotFound(owner, repo)
		}
		return nil, fmt.Errorf("get completed analysis by commit SHA for %s/%s@%s: %w", owner, repo, commitSHA, err)
	}

	var branchName *string
	if row.BranchName.Valid {
		branchName = &row.BranchName.String
	}
	var committedAt *time.Time
	if row.CommittedAt.Valid {
		t := row.CommittedAt.Time
		committedAt = &t
	}
	var parserVersion *string
	if row.ParserVersion != "" {
		parserVersion = &row.ParserVersion
	}

	return &port.CompletedAnalysis{
		BranchName:    branchName,
		CommitSHA:     row.CommitSha,
		CommittedAt:   committedAt,
		CompletedAt:   row.CompletedAt.Time,
		ID:            uuidToString(row.ID),
		Owner:         row.Owner,
		ParserVersion: parserVersion,
		Repo:          row.Repo,
		TotalSuites:   int(row.TotalSuites),
		TotalTests:    int(row.TotalTests),
	}, nil
}

func (r *PostgresRepository) GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*port.CompletedAnalysis, error) {
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

	var branchName *string
	if row.BranchName.Valid {
		branchName = &row.BranchName.String
	}
	var committedAt *time.Time
	if row.CommittedAt.Valid {
		t := row.CommittedAt.Time
		committedAt = &t
	}
	var parserVersion *string
	if row.ParserVersion != "" {
		parserVersion = &row.ParserVersion
	}

	return &port.CompletedAnalysis{
		BranchName:    branchName,
		CommitSHA:     row.CommitSha,
		CommittedAt:   committedAt,
		CompletedAt:   row.CompletedAt.Time,
		ID:            uuidToString(row.ID),
		Owner:         row.Owner,
		ParserVersion: parserVersion,
		Repo:          row.Repo,
		TotalSuites:   int(row.TotalSuites),
		TotalTests:    int(row.TotalTests),
	}, nil
}

func (r *PostgresRepository) GetPreviousAnalysis(ctx context.Context, codebaseID, currentAnalysisID string) (*port.PreviousAnalysis, error) {
	codebaseUUID, err := stringToUUID(codebaseID)
	if err != nil {
		return nil, fmt.Errorf("parse codebase ID: %w", err)
	}
	analysisUUID, err := stringToUUID(currentAnalysisID)
	if err != nil {
		return nil, fmt.Errorf("parse analysis ID: %w", err)
	}

	row, err := r.queries.GetPreviousAnalysis(ctx, db.GetPreviousAnalysisParams{
		CodebaseID: codebaseUUID,
		ID:         analysisUUID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get previous analysis: %w", err)
	}

	return &port.PreviousAnalysis{
		CommitSHA:  row.CommitSha,
		ID:         uuidToString(row.ID),
		TotalTests: int(row.TotalTests),
	}, nil
}

func (r *PostgresRepository) GetRepositoryStats(ctx context.Context, userID string) (*entity.RepositoryStats, error) {
	uuid, err := stringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	row, err := r.queries.GetRepositoryStats(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get repository stats: %w", err)
	}
	return &entity.RepositoryStats{
		TotalRepositories: int(row.TotalRepositories),
		TotalTests:        int(row.TotalTests),
	}, nil
}

func (r *PostgresRepository) GetTestSuitesWithCases(ctx context.Context, analysisID string) ([]port.TestSuiteWithCases, error) {
	uuid, err := stringToUUID(analysisID)
	if err != nil {
		return nil, fmt.Errorf("parse analysis ID: %w", err)
	}

	suiteRows, err := r.queries.GetTestSuitesByAnalysisID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get test suites: %w", err)
	}

	if len(suiteRows) == 0 {
		return []port.TestSuiteWithCases{}, nil
	}

	suiteIDs := make([]pgtype.UUID, len(suiteRows))
	for i, s := range suiteRows {
		suiteIDs[i] = s.ID
	}

	testRows, err := r.queries.GetTestCasesBySuiteIDs(ctx, suiteIDs)
	if err != nil {
		return nil, fmt.Errorf("get test cases: %w", err)
	}

	testsBySuite := make(map[string][]port.TestCaseRow)
	for _, t := range testRows {
		suiteID := uuidToString(t.SuiteID)
		line := 0
		if t.LineNumber.Valid {
			line = int(t.LineNumber.Int32)
		}
		testsBySuite[suiteID] = append(testsBySuite[suiteID], port.TestCaseRow{
			Line:   line,
			Name:   t.Name,
			Status: string(t.Status),
		})
	}

	suites := make([]port.TestSuiteWithCases, len(suiteRows))
	for i, s := range suiteRows {
		suiteID := uuidToString(s.ID)
		framework := ""
		if s.Framework.Valid {
			framework = s.Framework.String
		}
		suites[i] = port.TestSuiteWithCases{
			FilePath:  s.FilePath,
			Framework: framework,
			ID:        suiteID,
			Name:      s.Name,
			Tests:     testsBySuite[suiteID],
		}
	}

	return suites, nil
}

func (r *PostgresRepository) GetPaginatedRepositories(ctx context.Context, params port.PaginationParams) ([]port.PaginatedRepository, error) {
	var userUUID pgtype.UUID
	if params.UserID != "" {
		var err error
		userUUID, err = stringToUUID(params.UserID)
		if err != nil {
			return nil, fmt.Errorf("parse user ID: %w", err)
		}
	}

	var cursorID pgtype.UUID
	var err error
	if params.Cursor != nil {
		cursorID, err = stringToUUID(params.Cursor.ID)
		if err != nil {
			return nil, fmt.Errorf("parse cursor ID: %w", err)
		}
	}

	switch params.SortBy {
	case entity.SortByRecent:
		return r.getPaginatedByRecent(ctx, userUUID, cursorID, params)
	case entity.SortByName:
		return r.getPaginatedByName(ctx, userUUID, cursorID, params)
	case entity.SortByTests:
		return r.getPaginatedByTests(ctx, userUUID, cursorID, params)
	default:
		return r.getPaginatedByRecent(ctx, userUUID, cursorID, params)
	}
}

func (r *PostgresRepository) getPaginatedByRecent(ctx context.Context, userUUID, cursorID pgtype.UUID, params port.PaginationParams) ([]port.PaginatedRepository, error) {
	var cursorAnalyzedAt pgtype.Timestamptz
	if params.Cursor != nil {
		cursorAnalyzedAt = pgtype.Timestamptz{Time: params.Cursor.AnalyzedAt, Valid: true}
	}

	rows, err := r.queries.GetPaginatedRepositoriesByRecent(ctx, db.GetPaginatedRepositoriesByRecentParams{
		UserID:           userUUID,
		ViewFilter:       params.View.String(),
		OwnershipFilter:  params.Ownership.String(),
		CursorAnalyzedAt: cursorAnalyzedAt,
		SortOrder:        params.SortOrder.String(),
		CursorID:         cursorID,
		PageLimit:        int32(params.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("get paginated repositories by recent: %w", err)
	}

	repos := make([]port.PaginatedRepository, len(rows))
	for i, row := range rows {
		repos[i] = port.PaginatedRepository{
			ActiveCount:    int(row.ActiveCount),
			AnalysisID:     uuidToString(row.AnalysisID),
			AnalyzedAt:     row.AnalyzedAt.Time,
			CodebaseID:     uuidToString(row.CodebaseID),
			CommitSHA:      row.CommitSha,
			FocusedCount:   int(row.FocusedCount),
			IsAnalyzedByMe: row.IsAnalyzedByMe,
			Name:           row.Name,
			Owner:          row.Owner,
			SkippedCount:   int(row.SkippedCount),
			TodoCount:      int(row.TodoCount),
			TotalTests:     int(row.TotalTests),
			XfailCount:     int(row.XfailCount),
		}
	}
	return repos, nil
}

func (r *PostgresRepository) getPaginatedByName(ctx context.Context, userUUID, cursorID pgtype.UUID, params port.PaginationParams) ([]port.PaginatedRepository, error) {
	var cursorName string
	if params.Cursor != nil {
		cursorName = params.Cursor.Name
	}

	rows, err := r.queries.GetPaginatedRepositoriesByName(ctx, db.GetPaginatedRepositoriesByNameParams{
		UserID:          userUUID,
		ViewFilter:      params.View.String(),
		OwnershipFilter: params.Ownership.String(),
		CursorName:      cursorName,
		SortOrder:       params.SortOrder.String(),
		CursorID:        cursorID,
		PageLimit:       int32(params.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("get paginated repositories by name: %w", err)
	}

	repos := make([]port.PaginatedRepository, len(rows))
	for i, row := range rows {
		repos[i] = port.PaginatedRepository{
			ActiveCount:    int(row.ActiveCount),
			AnalysisID:     uuidToString(row.AnalysisID),
			AnalyzedAt:     row.AnalyzedAt.Time,
			CodebaseID:     uuidToString(row.CodebaseID),
			CommitSHA:      row.CommitSha,
			FocusedCount:   int(row.FocusedCount),
			IsAnalyzedByMe: row.IsAnalyzedByMe,
			Name:           row.Name,
			Owner:          row.Owner,
			SkippedCount:   int(row.SkippedCount),
			TodoCount:      int(row.TodoCount),
			TotalTests:     int(row.TotalTests),
			XfailCount:     int(row.XfailCount),
		}
	}
	return repos, nil
}

func (r *PostgresRepository) getPaginatedByTests(ctx context.Context, userUUID, cursorID pgtype.UUID, params port.PaginationParams) ([]port.PaginatedRepository, error) {
	var cursorTestCount int32
	if params.Cursor != nil {
		cursorTestCount = int32(params.Cursor.TestCount)
	}

	rows, err := r.queries.GetPaginatedRepositoriesByTests(ctx, db.GetPaginatedRepositoriesByTestsParams{
		UserID:          userUUID,
		ViewFilter:      params.View.String(),
		OwnershipFilter: params.Ownership.String(),
		CursorTestCount: cursorTestCount,
		SortOrder:       params.SortOrder.String(),
		CursorID:        cursorID,
		PageLimit:       int32(params.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("get paginated repositories by tests: %w", err)
	}

	repos := make([]port.PaginatedRepository, len(rows))
	for i, row := range rows {
		repos[i] = port.PaginatedRepository{
			ActiveCount:    int(row.ActiveCount),
			AnalysisID:     uuidToString(row.AnalysisID),
			AnalyzedAt:     row.AnalyzedAt.Time,
			CodebaseID:     uuidToString(row.CodebaseID),
			CommitSHA:      row.CommitSha,
			FocusedCount:   int(row.FocusedCount),
			IsAnalyzedByMe: row.IsAnalyzedByMe,
			Name:           row.Name,
			Owner:          row.Owner,
			SkippedCount:   int(row.SkippedCount),
			TodoCount:      int(row.TodoCount),
			TotalTests:     int(row.TotalTests),
			XfailCount:     int(row.XfailCount),
		}
	}
	return repos, nil
}

func (r *PostgresRepository) UpdateLastViewed(ctx context.Context, owner, repo string) error {
	if err := r.queries.UpdateCodebaseLastViewed(ctx, db.UpdateCodebaseLastViewedParams{
		Host:  HostGitHub,
		Owner: owner,
		Name:  repo,
	}); err != nil {
		return fmt.Errorf("update last viewed for %s/%s: %w", owner, repo, err)
	}
	return nil
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
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
