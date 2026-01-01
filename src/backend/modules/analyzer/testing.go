package analyzer

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
	"github.com/specvital/web/src/backend/modules/analyzer/handler"
	"github.com/specvital/web/src/backend/modules/analyzer/usecase"
	authhandler "github.com/specvital/web/src/backend/modules/auth/handler"
	"github.com/specvital/web/src/backend/modules/user"
)

// mockRepository is a test double for port.Repository.
type mockRepository struct {
	completedAnalysis *port.CompletedAnalysis
	err               error
	lastViewedCalled  bool
	lastViewedOwner   string
	lastViewedRepo    string
	suitesWithCases   []port.TestSuiteWithCases
}

var _ port.Repository = (*mockRepository)(nil)

func (m *mockRepository) GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*port.CompletedAnalysis, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.completedAnalysis == nil {
		return nil, domain.ErrNotFound
	}
	return m.completedAnalysis, nil
}

func (m *mockRepository) GetTestSuitesWithCases(ctx context.Context, analysisID string) ([]port.TestSuiteWithCases, error) {
	if m.suitesWithCases == nil {
		return []port.TestSuiteWithCases{}, nil
	}
	return m.suitesWithCases, nil
}

func (m *mockRepository) UpdateLastViewed(ctx context.Context, owner, repo string) error {
	m.lastViewedCalled = true
	m.lastViewedOwner = owner
	m.lastViewedRepo = repo
	return nil
}

func (m *mockRepository) FindActiveRiverJobByRepo(ctx context.Context, kind, owner, repo string) (*port.RiverJobInfo, error) {
	return nil, nil
}

func (m *mockRepository) GetCodebaseID(ctx context.Context, owner, repo string) (string, error) {
	return "test-codebase-id", nil
}

func (m *mockRepository) GetRepositoryStats(ctx context.Context) (*entity.RepositoryStats, error) {
	return &entity.RepositoryStats{}, nil
}

func (m *mockRepository) GetPreviousAnalysis(ctx context.Context, codebaseID, currentAnalysisID string) (*port.PreviousAnalysis, error) {
	return nil, nil
}

func (m *mockRepository) GetBookmarkedCodebaseIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}

func (m *mockRepository) GetPaginatedRepositories(ctx context.Context, params port.PaginationParams) ([]port.PaginatedRepository, error) {
	return nil, nil
}

// mockQueueService is a test double for port.QueueService.
type mockQueueService struct {
	enqueueCalled     bool
	enqueuedOwner     string
	enqueuedRepo      string
	enqueuedCommitSHA string
	enqueuedUserID    *string
	err               error
	findTaskInfo      *port.TaskInfo
}

var _ port.QueueService = (*mockQueueService)(nil)

func (m *mockQueueService) Enqueue(ctx context.Context, owner, repo, commitSHA string, userID *string) error {
	m.enqueueCalled = true
	m.enqueuedOwner = owner
	m.enqueuedRepo = repo
	m.enqueuedCommitSHA = commitSHA
	m.enqueuedUserID = userID
	return m.err
}

func (m *mockQueueService) FindTaskByRepo(ctx context.Context, owner, repo string) (*port.TaskInfo, error) {
	return m.findTaskInfo, nil
}

func (m *mockQueueService) Close() error {
	return nil
}

// mockGitClient is a test double for port.GitClient.
type mockGitClient struct {
	commitSHA      string
	commitSHAToken string
	err            error
	errToken       error
}

var _ port.GitClient = (*mockGitClient)(nil)

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

// mockTokenProvider is a test double for port.TokenProvider.
type mockTokenProvider struct {
	token string
	err   error
}

var _ port.TokenProvider = (*mockTokenProvider)(nil)

func (m *mockTokenProvider) GetUserGitHubToken(ctx context.Context, userID string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.token, nil
}

// setupTestHandler creates a new Handler with mock dependencies and chi router.
func setupTestHandler() (*handler.Handler, *chi.Mux) {
	repo := &mockRepository{}
	queue := &mockQueueService{}
	gitClient := &mockGitClient{}
	tokenProvider := &mockTokenProvider{}
	return setupTestHandlerWithMocks(repo, queue, gitClient, tokenProvider)
}

// setupTestHandlerWithMocks creates a Handler with provided mocks for more control in tests.
func setupTestHandlerWithMocks(repo *mockRepository, queue *mockQueueService, gitClient *mockGitClient, tokenProvider port.TokenProvider) (*handler.Handler, *chi.Mux) {
	log := logger.New()

	analyzeRepositoryUC := usecase.NewAnalyzeRepositoryUseCase(gitClient, queue, repo, tokenProvider)
	getAnalysisUC := usecase.NewGetAnalysisUseCase(queue, repo)
	listRepositoryCardsUC := usecase.NewListRepositoryCardsUseCase(repo)
	getUpdateStatusUC := usecase.NewGetUpdateStatusUseCase(gitClient, repo, tokenProvider)
	getRepositoryStatsUC := usecase.NewGetRepositoryStatsUseCase(repo)
	reanalyzeRepositoryUC := usecase.NewReanalyzeRepositoryUseCase(gitClient, queue, repo, tokenProvider)

	h := handler.NewHandler(
		log,
		analyzeRepositoryUC,
		getAnalysisUC,
		listRepositoryCardsUC,
		getUpdateStatusUC,
		getRepositoryStatsUC,
		reanalyzeRepositoryUC,
	)

	r := chi.NewRouter()
	apiHandlers := api.NewAPIHandlers(h, user.NewMockHandler(), authhandler.NewMockHandler(), user.NewMockHandler(), NewMockGitHubHandler(), NewMockGitHubAppHandler(), h, nil)
	strictHandler := api.NewStrictHandler(apiHandlers, nil)
	api.HandlerFromMux(strictHandler, r)

	return h, r
}

type mockGitHubAppHandler struct{}

var _ api.GitHubAppHandlers = (*mockGitHubAppHandler)(nil)

func NewMockGitHubAppHandler() *mockGitHubAppHandler {
	return &mockGitHubAppHandler{}
}

func (m *mockGitHubAppHandler) GetGitHubAppInstallURL(_ context.Context, _ api.GetGitHubAppInstallURLRequestObject) (api.GetGitHubAppInstallURLResponseObject, error) {
	return api.GetGitHubAppInstallURL200JSONResponse{InstallURL: ""}, nil
}

func (m *mockGitHubAppHandler) GetUserGitHubAppInstallations(_ context.Context, _ api.GetUserGitHubAppInstallationsRequestObject) (api.GetUserGitHubAppInstallationsResponseObject, error) {
	return api.GetUserGitHubAppInstallations200JSONResponse{Data: []api.GitHubAppInstallation{}}, nil
}

type mockGitHubHandler struct{}

var _ api.GitHubHandlers = (*mockGitHubHandler)(nil)

func NewMockGitHubHandler() *mockGitHubHandler {
	return &mockGitHubHandler{}
}

func (m *mockGitHubHandler) GetOrganizationRepositories(_ context.Context, _ api.GetOrganizationRepositoriesRequestObject) (api.GetOrganizationRepositoriesResponseObject, error) {
	return api.GetOrganizationRepositories200JSONResponse{Data: []api.GitHubRepository{}}, nil
}

func (m *mockGitHubHandler) GetUserGitHubOrganizations(_ context.Context, _ api.GetUserGitHubOrganizationsRequestObject) (api.GetUserGitHubOrganizationsResponseObject, error) {
	return api.GetUserGitHubOrganizations200JSONResponse{Data: []api.GitHubOrganization{}}, nil
}

func (m *mockGitHubHandler) GetUserGitHubRepositories(_ context.Context, _ api.GetUserGitHubRepositoriesRequestObject) (api.GetUserGitHubRepositoriesResponseObject, error) {
	return api.GetUserGitHubRepositories200JSONResponse{Data: []api.GitHubRepository{}}, nil
}
