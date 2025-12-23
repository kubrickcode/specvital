package analyzer

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/auth"
)

// mockRepository is a test double for Repository.
type mockRepository struct {
	completedAnalysis *CompletedAnalysis
	err               error
	lastViewedCalled  bool
	lastViewedOwner   string
	lastViewedRepo    string
	suitesWithCases   []TestSuiteWithCases
}

func (m *mockRepository) GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*CompletedAnalysis, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.completedAnalysis == nil {
		return nil, domain.ErrNotFound
	}
	return m.completedAnalysis, nil
}

func (m *mockRepository) GetTestSuitesWithCases(ctx context.Context, analysisID string) ([]TestSuiteWithCases, error) {
	if m.suitesWithCases == nil {
		return []TestSuiteWithCases{}, nil
	}
	return m.suitesWithCases, nil
}

func (m *mockRepository) UpdateLastViewed(ctx context.Context, owner, repo string) error {
	m.lastViewedCalled = true
	m.lastViewedOwner = owner
	m.lastViewedRepo = repo
	return nil
}

func (m *mockRepository) FindActiveRiverJobByRepo(ctx context.Context, kind, owner, repo string) (*RiverJobInfo, error) {
	return nil, nil
}

func (m *mockRepository) GetCodebaseID(ctx context.Context, owner, repo string) (string, error) {
	return "test-codebase-id", nil
}

func (m *mockRepository) GetRecentRepositories(ctx context.Context, limit int) ([]RecentRepository, error) {
	return nil, nil
}

func (m *mockRepository) GetRepositoryStats(ctx context.Context) (*domain.RepositoryStats, error) {
	return &domain.RepositoryStats{}, nil
}

func (m *mockRepository) GetPreviousAnalysis(ctx context.Context, codebaseID, currentAnalysisID string) (*PreviousAnalysis, error) {
	return nil, nil
}

// mockQueueService is a test double for QueueService.
type mockQueueService struct {
	enqueueCalled     bool
	enqueuedOwner     string
	enqueuedRepo      string
	enqueuedCommitSHA string
	enqueuedUserID    *string
	err               error
	findTaskInfo      *TaskInfo
}

func (m *mockQueueService) Enqueue(ctx context.Context, owner, repo, commitSHA string, userID *string) error {
	m.enqueueCalled = true
	m.enqueuedOwner = owner
	m.enqueuedRepo = repo
	m.enqueuedCommitSHA = commitSHA
	m.enqueuedUserID = userID
	return m.err
}

func (m *mockQueueService) FindTaskByRepo(ctx context.Context, owner, repo string) (*TaskInfo, error) {
	return m.findTaskInfo, nil
}

func (m *mockQueueService) Close() error {
	return nil
}

// mockGitClient is a test double for client.GitClient.
type mockGitClient struct {
	commitSHA      string
	commitSHAToken string
	err            error
	errToken       error
}

func (m *mockGitClient) GetLatestCommitSHA(ctx context.Context, owner, repo string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	if m.commitSHA == "" {
		return "test-commit-sha", nil
	}
	return m.commitSHA, nil
}

func (m *mockGitClient) GetLatestCommitSHAWithToken(ctx context.Context, owner, repo, token string) (string, error) {
	if m.errToken != nil {
		return "", m.errToken
	}
	if m.commitSHAToken == "" {
		return "test-commit-sha-with-token", nil
	}
	return m.commitSHAToken, nil
}

// Ensure mockGitClient implements client.GitClient.
var _ client.GitClient = (*mockGitClient)(nil)

// mockTokenProvider is a test double for TokenProvider.
type mockTokenProvider struct {
	token string
	err   error
}

func (m *mockTokenProvider) GetUserGitHubToken(ctx context.Context, userID string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.token, nil
}

// setupTestHandler creates a new AnalyzerHandler with mock dependencies and chi router.
func setupTestHandler() (*AnalyzerHandler, *chi.Mux) {
	repo := &mockRepository{}
	queue := &mockQueueService{}
	gitClient := &mockGitClient{}
	tokenProvider := &mockTokenProvider{}
	return setupTestHandlerWithMocks(repo, queue, gitClient, tokenProvider)
}

// setupTestHandlerWithMocks creates an AnalyzerHandler with provided mocks for more control in tests.
func setupTestHandlerWithMocks(repo *mockRepository, queue *mockQueueService, gitClient *mockGitClient, tokenProvider TokenProvider) (*AnalyzerHandler, *chi.Mux) {
	log := logger.New()
	service := NewAnalyzerService(log, repo, queue, gitClient, tokenProvider)
	handler := NewAnalyzerHandler(log, service)

	r := chi.NewRouter()
	repoHandler := NewMockRepositoryHandler()
	apiHandlers := api.NewAPIHandlers(handler, auth.NewMockHandler(), auth.NewMockBookmarkHandler(), repoHandler)
	strictHandler := api.NewStrictHandler(apiHandlers, nil)
	api.HandlerFromMux(strictHandler, r)

	return handler, r
}

type mockRepositoryHandler struct{}

var _ api.RepositoryHandlers = (*mockRepositoryHandler)(nil)

func NewMockRepositoryHandler() *mockRepositoryHandler {
	return &mockRepositoryHandler{}
}

func (m *mockRepositoryHandler) GetRecentRepositories(ctx context.Context, request api.GetRecentRepositoriesRequestObject) (api.GetRecentRepositoriesResponseObject, error) {
	return api.GetRecentRepositories200JSONResponse{Data: []api.RepositoryCard{}}, nil
}

func (m *mockRepositoryHandler) GetRepositoryStats(ctx context.Context, request api.GetRepositoryStatsRequestObject) (api.GetRepositoryStatsResponseObject, error) {
	return api.GetRepositoryStats200JSONResponse{TotalRepositories: 0, TotalTests: 0}, nil
}

func (m *mockRepositoryHandler) GetUpdateStatus(ctx context.Context, request api.GetUpdateStatusRequestObject) (api.GetUpdateStatusResponseObject, error) {
	return api.GetUpdateStatus200JSONResponse{Status: api.Unknown}, nil
}

func (m *mockRepositoryHandler) ReanalyzeRepository(ctx context.Context, request api.ReanalyzeRepositoryRequestObject) (api.ReanalyzeRepositoryResponseObject, error) {
	return api.ReanalyzeRepository500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("mock"),
	}, nil
}
