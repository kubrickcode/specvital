package analyzer

import (
	"context"

	"github.com/go-chi/chi/v5"
)

// mockRepository is a test double for Repository.
type mockRepository struct {
	completedAnalysis *CompletedAnalysis
	analysisStatus    *AnalysisStatus
	createdAnalysisID string
	suitesWithCases   []TestSuiteWithCases
	err               error
	createErr         error
}

func (m *mockRepository) CreatePendingAnalysis(ctx context.Context, owner, repo string) (string, error) {
	if m.createErr != nil {
		return "", m.createErr
	}
	if m.createdAnalysisID == "" {
		return "test-analysis-id", nil
	}
	return m.createdAnalysisID, nil
}

func (m *mockRepository) GetLatestCompletedAnalysis(ctx context.Context, owner, repo string) (*CompletedAnalysis, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.completedAnalysis == nil {
		return nil, ErrNotFound
	}
	return m.completedAnalysis, nil
}

func (m *mockRepository) GetAnalysisStatus(ctx context.Context, owner, repo string) (*AnalysisStatus, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.analysisStatus == nil {
		return nil, ErrNotFound
	}
	return m.analysisStatus, nil
}

func (m *mockRepository) MarkAnalysisFailed(ctx context.Context, analysisID, errorMsg string) error {
	return nil
}

func (m *mockRepository) GetTestSuitesWithCases(ctx context.Context, analysisID string) ([]TestSuiteWithCases, error) {
	if m.suitesWithCases == nil {
		return []TestSuiteWithCases{}, nil
	}
	return m.suitesWithCases, nil
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

// setupTestHandler creates a new Handler with mock dependencies and chi router.
func setupTestHandler() (*Handler, *chi.Mux) {
	repo := &mockRepository{}
	queue := &mockQueueService{}
	handler := NewHandler(repo, queue)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	return handler, r
}

// setupTestHandlerWithMocks creates a Handler with provided mocks for more control in tests.
func setupTestHandlerWithMocks(repo *mockRepository, queue *mockQueueService) (*Handler, *chi.Mux) {
	handler := NewHandler(repo, queue)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	return handler, r
}
