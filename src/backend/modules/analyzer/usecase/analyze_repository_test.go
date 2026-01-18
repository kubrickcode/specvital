package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
	"github.com/specvital/web/src/backend/modules/analyzer/usecase"
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type analyzeRepoMocks struct {
	gitClient     *mockGitClientForAnalyze
	queue         *mockQueueServiceForAnalyze
	repository    *mockRepositoryForAnalyze
	systemConfig  *mockSystemConfigForAnalyze
	tokenProvider *mockTokenProviderForAnalyze
}

func newAnalyzeRepoMocks() *analyzeRepoMocks {
	return &analyzeRepoMocks{
		gitClient:     &mockGitClientForAnalyze{},
		queue:         &mockQueueServiceForAnalyze{},
		repository:    &mockRepositoryForAnalyze{},
		systemConfig:  &mockSystemConfigForAnalyze{},
		tokenProvider: &mockTokenProviderForAnalyze{},
	}
}

func (m *analyzeRepoMocks) newUseCase() *usecase.AnalyzeRepositoryUseCase {
	return usecase.NewAnalyzeRepositoryUseCase(
		m.gitClient,
		m.queue,
		m.repository,
		m.systemConfig,
		m.tokenProvider,
	)
}

// mockRepositoryForAnalyze implements port.Repository for analyze tests.
type mockRepositoryForAnalyze struct {
	completedAnalysis *port.CompletedAnalysis
	completedErr      error
	suitesWithCases   []port.TestSuiteWithCases
}

func (m *mockRepositoryForAnalyze) FindActiveRiverJobByRepo(_ context.Context, _, _, _ string) (*port.RiverJobInfo, error) {
	return nil, nil
}
func (m *mockRepositoryForAnalyze) GetBookmarkedCodebaseIDs(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}
func (m *mockRepositoryForAnalyze) GetCodebaseID(_ context.Context, _, _ string) (string, error) {
	return "test-codebase-id", nil
}
func (m *mockRepositoryForAnalyze) GetLatestCompletedAnalysis(_ context.Context, _, _ string) (*port.CompletedAnalysis, error) {
	if m.completedErr != nil {
		return nil, m.completedErr
	}
	if m.completedAnalysis == nil {
		return nil, domain.ErrNotFound
	}
	return m.completedAnalysis, nil
}
func (m *mockRepositoryForAnalyze) GetPaginatedRepositories(_ context.Context, _ port.PaginationParams) ([]port.PaginatedRepository, error) {
	return nil, nil
}
func (m *mockRepositoryForAnalyze) GetPreviousAnalysis(_ context.Context, _, _ string) (*port.PreviousAnalysis, error) {
	return nil, nil
}
func (m *mockRepositoryForAnalyze) GetRepositoryStats(_ context.Context, _ string) (*entity.RepositoryStats, error) {
	return nil, nil
}
func (m *mockRepositoryForAnalyze) GetTestSuitesWithCases(_ context.Context, _ string) ([]port.TestSuiteWithCases, error) {
	return m.suitesWithCases, nil
}
func (m *mockRepositoryForAnalyze) UpdateLastViewed(_ context.Context, _, _ string) error {
	return nil
}

// mockQueueServiceForAnalyze implements port.QueueService.
type mockQueueServiceForAnalyze struct {
	enqueueCalled     bool
	enqueuedCommitSHA string
	enqueuedTier      subscription.PlanTier
	enqueueErr        error
	taskInfo          *port.TaskInfo
}

func (m *mockQueueServiceForAnalyze) Enqueue(_ context.Context, _, _, commitSHA string, _ *string, tier subscription.PlanTier) error {
	m.enqueueCalled = true
	m.enqueuedCommitSHA = commitSHA
	m.enqueuedTier = tier
	return m.enqueueErr
}
func (m *mockQueueServiceForAnalyze) FindTaskByRepo(_ context.Context, _, _ string) (*port.TaskInfo, error) {
	return m.taskInfo, nil
}
func (m *mockQueueServiceForAnalyze) Close() error {
	return nil
}

// mockGitClientForAnalyze implements port.GitClient.
type mockGitClientForAnalyze struct {
	latestSHA string
	err       error
}

func (m *mockGitClientForAnalyze) GetLatestCommitSHA(_ context.Context, _, _ string) (string, error) {
	return m.latestSHA, m.err
}
func (m *mockGitClientForAnalyze) GetLatestCommitSHAWithToken(_ context.Context, _, _, _ string) (string, error) {
	return m.latestSHA, m.err
}

// mockSystemConfigForAnalyze implements port.SystemConfigReader.
type mockSystemConfigForAnalyze struct {
	parserVersion string
	err           error
}

func (m *mockSystemConfigForAnalyze) GetParserVersion(_ context.Context) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.parserVersion, nil
}

// mockTokenProviderForAnalyze implements port.TokenProvider.
type mockTokenProviderForAnalyze struct{}

func (m *mockTokenProviderForAnalyze) GetUserGitHubToken(_ context.Context, _ string) (string, error) {
	return "", nil
}

func TestAnalyzeRepository_SameCommitSameVersion_ReturnsCached(t *testing.T) {
	t.Parallel()

	mocks := newAnalyzeRepoMocks()
	version := "v1.0.0"
	mocks.gitClient.latestSHA = "abc123"
	mocks.systemConfig.parserVersion = version
	mocks.repository.completedAnalysis = &port.CompletedAnalysis{
		CommitSHA:     "abc123",
		CompletedAt:   time.Now(),
		ID:            "analysis-1",
		Owner:         "owner",
		ParserVersion: &version,
		Repo:          "repo",
		TotalSuites:   5,
		TotalTests:    100,
	}

	uc := mocks.newUseCase()
	result, err := uc.Execute(context.Background(), usecase.AnalyzeRepositoryInput{
		Owner: "owner",
		Repo:  "repo",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Analysis == nil {
		t.Fatal("expected analysis result, got nil")
	}
	if mocks.queue.enqueueCalled {
		t.Error("expected no enqueue when commit and version match")
	}
}

func TestAnalyzeRepository_SameCommitDifferentVersion_EnqueuesReanalysis(t *testing.T) {
	t.Parallel()

	mocks := newAnalyzeRepoMocks()
	oldVersion := "v1.0.0"
	mocks.gitClient.latestSHA = "abc123"
	mocks.systemConfig.parserVersion = "v2.0.0"
	mocks.repository.completedAnalysis = &port.CompletedAnalysis{
		CommitSHA:     "abc123",
		CompletedAt:   time.Now(),
		ID:            "analysis-1",
		Owner:         "owner",
		ParserVersion: &oldVersion,
		Repo:          "repo",
		TotalSuites:   5,
		TotalTests:    100,
	}

	uc := mocks.newUseCase()
	result, err := uc.Execute(context.Background(), usecase.AnalyzeRepositoryInput{
		Owner: "owner",
		Repo:  "repo",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Progress == nil {
		t.Fatal("expected progress result, got nil")
	}
	if !mocks.queue.enqueueCalled {
		t.Error("expected enqueue when parser version differs")
	}
	if mocks.queue.enqueuedCommitSHA != "abc123" {
		t.Errorf("expected enqueued SHA abc123, got %s", mocks.queue.enqueuedCommitSHA)
	}
}

func TestAnalyzeRepository_NullParserVersion_EnqueuesReanalysis(t *testing.T) {
	t.Parallel()

	mocks := newAnalyzeRepoMocks()
	mocks.gitClient.latestSHA = "abc123"
	mocks.systemConfig.parserVersion = "v1.0.0"
	mocks.repository.completedAnalysis = &port.CompletedAnalysis{
		CommitSHA:     "abc123",
		CompletedAt:   time.Now(),
		ID:            "analysis-1",
		Owner:         "owner",
		ParserVersion: nil, // Legacy data without parser_version
		Repo:          "repo",
		TotalSuites:   5,
		TotalTests:    100,
	}

	uc := mocks.newUseCase()
	result, err := uc.Execute(context.Background(), usecase.AnalyzeRepositoryInput{
		Owner: "owner",
		Repo:  "repo",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Progress == nil {
		t.Fatal("expected progress result, got nil")
	}
	if !mocks.queue.enqueueCalled {
		t.Error("expected enqueue when parser_version is NULL")
	}
}

func TestAnalyzeRepository_DifferentCommit_EnqueuesAnalysis(t *testing.T) {
	t.Parallel()

	mocks := newAnalyzeRepoMocks()
	version := "v1.0.0"
	mocks.gitClient.latestSHA = "def456"
	mocks.systemConfig.parserVersion = version
	mocks.repository.completedAnalysis = &port.CompletedAnalysis{
		CommitSHA:     "abc123",
		CompletedAt:   time.Now(),
		ID:            "analysis-1",
		Owner:         "owner",
		ParserVersion: &version,
		Repo:          "repo",
		TotalSuites:   5,
		TotalTests:    100,
	}

	uc := mocks.newUseCase()
	result, err := uc.Execute(context.Background(), usecase.AnalyzeRepositoryInput{
		Owner: "owner",
		Repo:  "repo",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Progress == nil {
		t.Fatal("expected progress result, got nil")
	}
	if !mocks.queue.enqueueCalled {
		t.Error("expected enqueue when commit SHA differs")
	}
	if mocks.queue.enqueuedCommitSHA != "def456" {
		t.Errorf("expected enqueued SHA def456, got %s", mocks.queue.enqueuedCommitSHA)
	}
}

func TestAnalyzeRepository_SystemConfigError_ReturnsCached(t *testing.T) {
	t.Parallel()

	mocks := newAnalyzeRepoMocks()
	version := "v1.0.0"
	mocks.gitClient.latestSHA = "abc123"
	mocks.systemConfig.err = errors.New("system config unavailable")
	mocks.repository.completedAnalysis = &port.CompletedAnalysis{
		CommitSHA:     "abc123",
		CompletedAt:   time.Now(),
		ID:            "analysis-1",
		Owner:         "owner",
		ParserVersion: &version,
		Repo:          "repo",
		TotalSuites:   5,
		TotalTests:    100,
	}

	uc := mocks.newUseCase()
	result, err := uc.Execute(context.Background(), usecase.AnalyzeRepositoryInput{
		Owner: "owner",
		Repo:  "repo",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Analysis == nil {
		t.Fatal("expected cached analysis when system config unavailable")
	}
	if mocks.queue.enqueueCalled {
		t.Error("expected no enqueue when system config is unavailable")
	}
}

func TestAnalyzeRepository_NoExistingAnalysis_EnqueuesNew(t *testing.T) {
	t.Parallel()

	mocks := newAnalyzeRepoMocks()
	mocks.gitClient.latestSHA = "abc123"
	mocks.systemConfig.parserVersion = "v1.0.0"
	mocks.repository.completedErr = domain.ErrNotFound

	uc := mocks.newUseCase()
	result, err := uc.Execute(context.Background(), usecase.AnalyzeRepositoryInput{
		Owner: "owner",
		Repo:  "repo",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Progress == nil {
		t.Fatal("expected progress result, got nil")
	}
	if result.Progress.Status != entity.AnalysisStatusPending {
		t.Errorf("expected pending status, got %s", result.Progress.Status)
	}
	if !mocks.queue.enqueueCalled {
		t.Error("expected enqueue for new repository")
	}
}
