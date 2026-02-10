package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/domain/port"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/usecase"
	subscription "github.com/kubrickcode/specvital/apps/web/src/backend/modules/subscription/domain/entity"
)

type getAnalysisMocks struct {
	queue      *mockQueueServiceForGetAnalysis
	repository *mockRepositoryForGetAnalysis
}

func newGetAnalysisMocks() *getAnalysisMocks {
	return &getAnalysisMocks{
		queue:      &mockQueueServiceForGetAnalysis{},
		repository: &mockRepositoryForGetAnalysis{},
	}
}

func (m *getAnalysisMocks) newUseCase() *usecase.GetAnalysisUseCase {
	return usecase.NewGetAnalysisUseCase(m.queue, m.repository)
}

// mockRepositoryForGetAnalysis implements port.Repository for get analysis tests.
type mockRepositoryForGetAnalysis struct {
	completedAnalysis      *port.CompletedAnalysis
	completedAnalysisBySHA *port.CompletedAnalysis
	completedErr           error
	completedBySHAErr      error
	suitesWithCases        []port.TestSuiteWithCases
	lastViewedCalled       bool
}

func (m *mockRepositoryForGetAnalysis) CheckAnalysisExistsByCommitSHA(_ context.Context, _, _, _ string) (bool, error) {
	return false, nil
}
func (m *mockRepositoryForGetAnalysis) FindActiveRiverJobByRepo(_ context.Context, _, _, _ string) (*port.RiverJobInfo, error) {
	return nil, nil
}
func (m *mockRepositoryForGetAnalysis) GetAiSpecSummaries(_ context.Context, _ []string, _ string) (map[string]*entity.AiSpecSummary, error) {
	return make(map[string]*entity.AiSpecSummary), nil
}
func (m *mockRepositoryForGetAnalysis) GetAnalysisHistory(_ context.Context, _, _ string) ([]port.AnalysisHistoryItem, error) {
	return nil, nil
}
func (m *mockRepositoryForGetAnalysis) GetBookmarkedCodebaseIDs(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}
func (m *mockRepositoryForGetAnalysis) GetCodebaseID(_ context.Context, _, _ string) (string, error) {
	return "test-codebase-id", nil
}
func (m *mockRepositoryForGetAnalysis) GetCompletedAnalysisByCommitSHA(_ context.Context, _, _, commitSHA string) (*port.CompletedAnalysis, error) {
	if m.completedBySHAErr != nil {
		return nil, m.completedBySHAErr
	}
	if m.completedAnalysisBySHA != nil && m.completedAnalysisBySHA.CommitSHA == commitSHA {
		return m.completedAnalysisBySHA, nil
	}
	return nil, domain.ErrNotFound
}
func (m *mockRepositoryForGetAnalysis) GetLatestCompletedAnalysis(_ context.Context, _, _ string) (*port.CompletedAnalysis, error) {
	if m.completedErr != nil {
		return nil, m.completedErr
	}
	if m.completedAnalysis == nil {
		return nil, domain.ErrNotFound
	}
	return m.completedAnalysis, nil
}
func (m *mockRepositoryForGetAnalysis) GetPaginatedRepositories(_ context.Context, _ port.PaginationParams) ([]port.PaginatedRepository, error) {
	return nil, nil
}
func (m *mockRepositoryForGetAnalysis) GetPreviousAnalysis(_ context.Context, _, _ string) (*port.PreviousAnalysis, error) {
	return nil, nil
}
func (m *mockRepositoryForGetAnalysis) GetRepositoryStats(_ context.Context, _ string) (*entity.RepositoryStats, error) {
	return nil, nil
}
func (m *mockRepositoryForGetAnalysis) GetTestSuitesWithCases(_ context.Context, _ string) ([]port.TestSuiteWithCases, error) {
	return m.suitesWithCases, nil
}
func (m *mockRepositoryForGetAnalysis) UpdateLastViewed(_ context.Context, _, _ string) error {
	m.lastViewedCalled = true
	return nil
}

// mockQueueServiceForGetAnalysis implements port.QueueService for get analysis tests.
type mockQueueServiceForGetAnalysis struct {
	taskInfo *port.TaskInfo
}

func (m *mockQueueServiceForGetAnalysis) Enqueue(_ context.Context, _, _, _ string, _ *string, _ subscription.PlanTier) error {
	return nil
}
func (m *mockQueueServiceForGetAnalysis) EnqueueTx(_ context.Context, _ pgx.Tx, _, _, _ string, _ *string, _ subscription.PlanTier) (int64, error) {
	return 0, nil
}
func (m *mockQueueServiceForGetAnalysis) FindTaskByRepo(_ context.Context, _, _ string) (*port.TaskInfo, error) {
	return m.taskInfo, nil
}
func (m *mockQueueServiceForGetAnalysis) Close() error {
	return nil
}

func TestGetAnalysisUseCase_ExecuteByCommitSHA(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	testAnalysis := &port.CompletedAnalysis{
		ID:          "analysis-123",
		CommitSHA:   "abc1234567890",
		Owner:       "testowner",
		Repo:        "testrepo",
		CompletedAt: now,
		TotalTests:  10,
		TotalSuites: 2,
	}

	t.Run("returns analysis for valid commit SHA", func(t *testing.T) {
		t.Parallel()

		mocks := newGetAnalysisMocks()
		mocks.repository.completedAnalysisBySHA = testAnalysis
		uc := mocks.newUseCase()

		result, err := uc.Execute(context.Background(), usecase.GetAnalysisInput{
			Owner:     "testowner",
			Repo:      "testrepo",
			CommitSHA: "abc1234567890",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Analysis == nil {
			t.Fatal("expected analysis to be returned")
		}
		if result.Analysis.CommitSHA != "abc1234567890" {
			t.Errorf("expected commit SHA abc1234567890, got %s", result.Analysis.CommitSHA)
		}
		if !mocks.repository.lastViewedCalled {
			t.Error("expected UpdateLastViewed to be called")
		}
	})

	t.Run("returns ErrNotFound for unknown commit", func(t *testing.T) {
		t.Parallel()

		mocks := newGetAnalysisMocks()
		mocks.repository.completedAnalysisBySHA = testAnalysis
		uc := mocks.newUseCase()

		_, err := uc.Execute(context.Background(), usecase.GetAnalysisInput{
			Owner:     "testowner",
			Repo:      "testrepo",
			CommitSHA: "unknown123",
		})

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("propagates repository errors", func(t *testing.T) {
		t.Parallel()

		mocks := newGetAnalysisMocks()
		mocks.repository.completedBySHAErr = errors.New("database connection failed")
		uc := mocks.newUseCase()

		_, err := uc.Execute(context.Background(), usecase.GetAnalysisInput{
			Owner:     "testowner",
			Repo:      "testrepo",
			CommitSHA: "abc1234",
		})

		if err == nil {
			t.Fatal("expected error to be returned")
		}
		if errors.Is(err, domain.ErrNotFound) {
			t.Error("should not be ErrNotFound for database errors")
		}
	})

	t.Run("does not check queue when commit SHA is provided", func(t *testing.T) {
		t.Parallel()

		mocks := newGetAnalysisMocks()
		mocks.repository.completedAnalysisBySHA = testAnalysis
		mocks.queue.taskInfo = &port.TaskInfo{
			CommitSHA: "different123",
			State:     "running",
		}
		uc := mocks.newUseCase()

		result, err := uc.Execute(context.Background(), usecase.GetAnalysisInput{
			Owner:     "testowner",
			Repo:      "testrepo",
			CommitSHA: "abc1234567890",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Analysis == nil {
			t.Fatal("expected analysis to be returned")
		}
		if result.Progress != nil {
			t.Error("should not return progress when commit SHA is specified")
		}
	})
}
