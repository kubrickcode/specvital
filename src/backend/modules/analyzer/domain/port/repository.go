package port

import (
	"context"
	"time"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
)

type Repository interface {
	CheckAnalysisExistsByCommitSHA(ctx context.Context, owner, repo, commitSHA string) (bool, error)
	FindActiveRiverJobByRepo(ctx context.Context, kind, owner, repo string) (*RiverJobInfo, error)
	GetAiSpecSummaries(ctx context.Context, codebaseIDs []string, userID string) (map[string]*entity.AiSpecSummary, error)
	GetBookmarkedCodebaseIDs(ctx context.Context, userID string) ([]string, error)
	GetCodebaseID(ctx context.Context, owner, repo string) (string, error)
	GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*CompletedAnalysis, error)
	GetPaginatedRepositories(ctx context.Context, params PaginationParams) ([]PaginatedRepository, error)
	GetPreviousAnalysis(ctx context.Context, codebaseID, currentAnalysisID string) (*PreviousAnalysis, error)
	GetRepositoryStats(ctx context.Context, userID string) (*entity.RepositoryStats, error)
	GetTestSuitesWithCases(ctx context.Context, analysisID string) ([]TestSuiteWithCases, error)
	UpdateLastViewed(ctx context.Context, owner, repo string) error
}

type PaginationParams struct {
	Cursor    *entity.RepositoryCursor
	Limit     int
	Ownership entity.OwnershipFilter
	SortBy    entity.SortBy
	SortOrder entity.SortOrder
	UserID    string
	View      entity.ViewFilter
}

type PaginatedRepository struct {
	ActiveCount    int
	AnalysisID     string
	AnalyzedAt     time.Time
	CodebaseID     string
	CommitSHA      string
	FocusedCount   int
	IsAnalyzedByMe bool
	Name           string
	Owner          string
	SkippedCount   int
	TodoCount      int
	TotalTests     int
	XfailCount     int
}

type CompletedAnalysis struct {
	BranchName    *string
	CommitSHA     string
	CommittedAt   *time.Time
	CompletedAt   time.Time
	ID            string
	Owner         string
	ParserVersion *string
	Repo          string
	TotalSuites   int
	TotalTests    int
}

type TestSuiteWithCases struct {
	FilePath  string
	Framework string
	ID        string
	Name      string
	Tests     []TestCaseRow
}

type TestCaseRow struct {
	Line   int
	Name   string
	Status string
}

type RiverJobInfo struct {
	AttemptedAt *time.Time
	CommitSHA   string
	State       string
}

type PreviousAnalysis struct {
	CommitSHA  string
	ID         string
	TotalTests int
}

type HistoryChecker interface {
	CheckUserHistoryExists(ctx context.Context, userID, owner, repo string) (bool, error)
}

type TierLookup interface {
	GetUserTier(ctx context.Context, userID string) (string, error)
}
