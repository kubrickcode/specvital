package analyzer

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
)

// mockRepository is a test double for Repository.
type mockRepository struct {
	analysisStatus    *AnalysisStatus
	completedAnalysis *CompletedAnalysis
	createErr         error
	createdAnalysisID string
	err               error
	suitesWithCases   []TestSuiteWithCases
}

func (m *mockRepository) CreatePendingAnalysis(ctx context.Context, owner, repo, commitSHA string) (string, error) {
	if m.createErr != nil {
		return "", m.createErr
	}
	if m.createdAnalysisID == "" {
		return "test-analysis-id", nil
	}
	return m.createdAnalysisID, nil
}

func (m *mockRepository) GetAnalysisStatus(ctx context.Context, owner, repo string) (*AnalysisStatus, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.analysisStatus == nil {
		return nil, domain.ErrNotFound
	}
	return m.analysisStatus, nil
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

func (m *mockRepository) MarkAnalysisFailed(ctx context.Context, analysisID, errorMsg string) error {
	return nil
}

// mockQueueService is a test double for QueueService.
type mockQueueService struct {
	enqueueCalled      bool
	enqueuedAnalysisID string
	enqueuedOwner      string
	enqueuedRepo       string
	err                error
}

func (m *mockQueueService) Enqueue(ctx context.Context, analysisID, owner, repo string) error {
	m.enqueueCalled = true
	m.enqueuedAnalysisID = analysisID
	m.enqueuedOwner = owner
	m.enqueuedRepo = repo
	return m.err
}

func (m *mockQueueService) Close() error {
	return nil
}

// mockGitClient is a test double for client.GitClient.
type mockGitClient struct {
	commitSHA string
	err       error
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

// Ensure mockGitClient implements client.GitClient.
var _ client.GitClient = (*mockGitClient)(nil)

// setupTestHandler creates a new AnalyzerHandler with mock dependencies and chi router.
func setupTestHandler() (*AnalyzerHandler, *chi.Mux) {
	repo := &mockRepository{}
	queue := &mockQueueService{}
	gitClient := &mockGitClient{}
	return setupTestHandlerWithMocks(repo, queue, gitClient)
}

// setupTestHandlerWithMocks creates an AnalyzerHandler with provided mocks for more control in tests.
func setupTestHandlerWithMocks(repo *mockRepository, queue *mockQueueService, gitClient *mockGitClient) (*AnalyzerHandler, *chi.Mux) {
	log := logger.New()
	service := NewAnalyzerService(log, repo, queue, gitClient)
	handler := NewAnalyzerHandler(log, service)

	r := chi.NewRouter()
	strictHandler := api.NewStrictHandler(handler, nil)
	api.HandlerFromMux(strictHandler, r)

	return handler, r
}
