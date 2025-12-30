package port

import (
	"context"
	"time"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
)

type Repository interface {
	FindActiveRiverJobByRepo(ctx context.Context, kind, owner, repo string) (*RiverJobInfo, error)
	GetCodebaseID(ctx context.Context, owner, repo string) (string, error)
	GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*CompletedAnalysis, error)
	GetPreviousAnalysis(ctx context.Context, codebaseID, currentAnalysisID string) (*PreviousAnalysis, error)
	GetRecentRepositories(ctx context.Context, limit int) ([]RecentRepository, error)
	GetRepositoryStats(ctx context.Context) (*entity.RepositoryStats, error)
	GetTestSuitesWithCases(ctx context.Context, analysisID string) ([]TestSuiteWithCases, error)
	UpdateLastViewed(ctx context.Context, owner, repo string) error
}

type CompletedAnalysis struct {
	BranchName  *string
	CommitSHA   string
	CommittedAt *time.Time
	CompletedAt time.Time
	ID          string
	Owner       string
	Repo        string
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
	Line   int
	Name   string
	Status string
}

type RiverJobInfo struct {
	CommitSHA string
	State     string
}

type RecentRepository struct {
	AnalysisID string
	AnalyzedAt time.Time
	CodebaseID string
	CommitSHA  string
	Name       string
	Owner      string
	TotalTests int
}

type PreviousAnalysis struct {
	CommitSHA  string
	ID         string
	TotalTests int
}
