package analyzer

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
	"github.com/specvital/web/src/backend/modules/analyzer/handler"
	"github.com/specvital/web/src/backend/modules/analyzer/usecase"
	authhandler "github.com/specvital/web/src/backend/modules/auth/handler"
	specviewhandler "github.com/specvital/web/src/backend/modules/spec-view/handler"
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
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

func (m *mockRepository) CheckAnalysisExistsByCommitSHA(ctx context.Context, owner, repo, commitSHA string) (bool, error) {
	if m.completedAnalysis != nil && m.completedAnalysis.CommitSHA == commitSHA {
		return true, nil
	}
	return false, nil
}

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

func (m *mockRepository) GetCompletedAnalysisByCommitSHA(ctx context.Context, owner, repo, commitSHA string) (*port.CompletedAnalysis, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.completedAnalysis != nil && m.completedAnalysis.CommitSHA == commitSHA {
		return m.completedAnalysis, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockRepository) GetRepositoryStats(ctx context.Context, userID string) (*entity.RepositoryStats, error) {
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

func (m *mockRepository) GetAiSpecSummaries(ctx context.Context, analysisIDs []string, userID string) (map[string]*entity.AiSpecSummary, error) {
	return make(map[string]*entity.AiSpecSummary), nil
}

func (m *mockRepository) GetAnalysisHistory(ctx context.Context, owner, repo string) ([]port.AnalysisHistoryItem, error) {
	return nil, nil
}

// mockSystemConfigReader is a test double for port.SystemConfigReader.
type mockSystemConfigReader struct {
	parserVersion string
	err           error
}

var _ port.SystemConfigReader = (*mockSystemConfigReader)(nil)

func (m *mockSystemConfigReader) GetParserVersion(ctx context.Context) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.parserVersion, nil
}

// mockQueueService is a test double for port.QueueService.
type mockQueueService struct {
	enqueueCalled     bool
	enqueuedOwner     string
	enqueuedRepo      string
	enqueuedCommitSHA string
	enqueuedTier      subscription.PlanTier
	enqueuedUserID    *string
	err               error
	findTaskInfo      *port.TaskInfo
}

var _ port.QueueService = (*mockQueueService)(nil)

func (m *mockQueueService) Enqueue(ctx context.Context, owner, repo, commitSHA string, userID *string, tier subscription.PlanTier) error {
	m.enqueueCalled = true
	m.enqueuedOwner = owner
	m.enqueuedRepo = repo
	m.enqueuedCommitSHA = commitSHA
	m.enqueuedUserID = userID
	m.enqueuedTier = tier
	return m.err
}

func (m *mockQueueService) EnqueueTx(ctx context.Context, tx pgx.Tx, owner, repo, commitSHA string, userID *string, tier subscription.PlanTier) (int64, error) {
	m.enqueueCalled = true
	m.enqueuedOwner = owner
	m.enqueuedRepo = repo
	m.enqueuedCommitSHA = commitSHA
	m.enqueuedUserID = userID
	m.enqueuedTier = tier
	if m.err != nil {
		return 0, m.err
	}
	return 1, nil
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
	systemConfig := &mockSystemConfigReader{parserVersion: "v1.0.0"}

	analyzeRepositoryUC := usecase.NewAnalyzeRepositoryUseCase(gitClient, queue, repo, systemConfig, tokenProvider, nil, nil)
	getAnalysisUC := usecase.NewGetAnalysisUseCase(queue, repo)
	getAnalysisHistoryUC := usecase.NewGetAnalysisHistoryUseCase(repo)
	listRepositoryCardsUC := usecase.NewListRepositoryCardsUseCase(gitClient, repo, tokenProvider)
	getUpdateStatusUC := usecase.NewGetUpdateStatusUseCase(gitClient, repo, systemConfig, tokenProvider)
	getRepositoryStatsUC := usecase.NewGetRepositoryStatsUseCase(repo)
	reanalyzeRepositoryUC := usecase.NewReanalyzeRepositoryUseCase(gitClient, queue, repo, tokenProvider)

	h := handler.NewHandler(
		log,
		analyzeRepositoryUC,
		getAnalysisUC,
		getAnalysisHistoryUC,
		listRepositoryCardsUC,
		getUpdateStatusUC,
		getRepositoryStatsUC,
		reanalyzeRepositoryUC,
		nil,
		nil,
		nil,
	)

	r := chi.NewRouter()
	apiHandlers := api.NewAPIHandlers(h, user.NewMockHandler(), authhandler.NewMockHandler(), user.NewMockHandler(), NewMockGitHubHandler(), NewMockGitHubAppHandler(), NewMockPricingHandler(), h, specviewhandler.NewMockHandler(), NewMockSubscriptionHandler(), NewMockUsageHandler(), user.NewMockHandler(), nil)
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

type mockUsageHandler struct{}

var _ api.UsageHandlers = (*mockUsageHandler)(nil)

func NewMockUsageHandler() *mockUsageHandler {
	return &mockUsageHandler{}
}

func (m *mockUsageHandler) CheckQuota(_ context.Context, _ api.CheckQuotaRequestObject) (api.CheckQuotaResponseObject, error) {
	return api.CheckQuota200JSONResponse{}, nil
}

func (m *mockUsageHandler) GetCurrentUsage(_ context.Context, _ api.GetCurrentUsageRequestObject) (api.GetCurrentUsageResponseObject, error) {
	return api.GetCurrentUsage200JSONResponse{}, nil
}

type mockSubscriptionHandler struct{}

var _ api.SubscriptionHandlers = (*mockSubscriptionHandler)(nil)

func NewMockSubscriptionHandler() *mockSubscriptionHandler {
	return &mockSubscriptionHandler{}
}

func (m *mockSubscriptionHandler) GetUserSubscription(_ context.Context, _ api.GetUserSubscriptionRequestObject) (api.GetUserSubscriptionResponseObject, error) {
	return api.GetUserSubscription200JSONResponse{}, nil
}

type mockPricingHandler struct{}

var _ api.PricingHandlers = (*mockPricingHandler)(nil)

func NewMockPricingHandler() *mockPricingHandler {
	return &mockPricingHandler{}
}

func (m *mockPricingHandler) GetPricing(_ context.Context, _ api.GetPricingRequestObject) (api.GetPricingResponseObject, error) {
	return api.GetPricing200JSONResponse{Data: []api.PricingPlan{}}, nil
}
