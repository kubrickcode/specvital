package retention

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/worker/internal/domain/retention"
)

type mockCleanupRepository struct {
	deleteExpiredUserAnalysisHistoryFn func(ctx context.Context, batchSize int) (retention.DeleteResult, error)
	deleteExpiredSpecDocumentsFn       func(ctx context.Context, batchSize int) (retention.DeleteResult, error)
	deleteOrphanedAnalysesFn           func(ctx context.Context, batchSize int) (retention.DeleteResult, error)
}

func (m *mockCleanupRepository) DeleteExpiredUserAnalysisHistory(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
	if m.deleteExpiredUserAnalysisHistoryFn != nil {
		return m.deleteExpiredUserAnalysisHistoryFn(ctx, batchSize)
	}
	return retention.DeleteResult{}, nil
}

func (m *mockCleanupRepository) DeleteExpiredSpecDocuments(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
	if m.deleteExpiredSpecDocumentsFn != nil {
		return m.deleteExpiredSpecDocumentsFn(ctx, batchSize)
	}
	return retention.DeleteResult{}, nil
}

func (m *mockCleanupRepository) DeleteOrphanedAnalyses(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
	if m.deleteOrphanedAnalysesFn != nil {
		return m.deleteOrphanedAnalysesFn(ctx, batchSize)
	}
	return retention.DeleteResult{}, nil
}

func TestNewCleanupUseCase(t *testing.T) {
	repo := &mockCleanupRepository{}

	t.Run("default options", func(t *testing.T) {
		uc := NewCleanupUseCase(repo)
		if uc.batchSize != retention.DefaultBatchSize {
			t.Errorf("batchSize = %d, want %d", uc.batchSize, retention.DefaultBatchSize)
		}
		if uc.batchSleep != DefaultBatchSleepDuration {
			t.Errorf("batchSleep = %v, want %v", uc.batchSleep, DefaultBatchSleepDuration)
		}
	})

	t.Run("custom batch size", func(t *testing.T) {
		uc := NewCleanupUseCase(repo, WithBatchSize(500))
		if uc.batchSize != 500 {
			t.Errorf("batchSize = %d, want 500", uc.batchSize)
		}
	})

	t.Run("invalid batch size ignored", func(t *testing.T) {
		uc := NewCleanupUseCase(repo, WithBatchSize(0))
		if uc.batchSize != retention.DefaultBatchSize {
			t.Errorf("batchSize = %d, want %d", uc.batchSize, retention.DefaultBatchSize)
		}
	})

	t.Run("custom batch sleep", func(t *testing.T) {
		uc := NewCleanupUseCase(repo, WithBatchSleep(50*time.Millisecond))
		if uc.batchSleep != 50*time.Millisecond {
			t.Errorf("batchSleep = %v, want 50ms", uc.batchSleep)
		}
	})

	t.Run("zero batch sleep allowed", func(t *testing.T) {
		uc := NewCleanupUseCase(repo, WithBatchSleep(0))
		if uc.batchSleep != 0 {
			t.Errorf("batchSleep = %v, want 0", uc.batchSleep)
		}
	})

	t.Run("negative batch sleep ignored", func(t *testing.T) {
		uc := NewCleanupUseCase(repo, WithBatchSleep(-1*time.Millisecond))
		if uc.batchSleep != DefaultBatchSleepDuration {
			t.Errorf("batchSleep = %v, want %v", uc.batchSleep, DefaultBatchSleepDuration)
		}
	})
}

func TestCleanupUseCase_Execute(t *testing.T) {
	t.Run("success - all phases complete", func(t *testing.T) {
		repo := &mockCleanupRepository{
			deleteExpiredUserAnalysisHistoryFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 10}, nil
			},
			deleteExpiredSpecDocumentsFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 5}, nil
			},
			deleteOrphanedAnalysesFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 3}, nil
			},
		}

		uc := NewCleanupUseCase(repo, WithBatchSleep(0))
		result, err := uc.Execute(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.UserAnalysisHistoryDeleted != 10 {
			t.Errorf("UserAnalysisHistoryDeleted = %d, want 10", result.UserAnalysisHistoryDeleted)
		}
		if result.SpecDocumentsDeleted != 5 {
			t.Errorf("SpecDocumentsDeleted = %d, want 5", result.SpecDocumentsDeleted)
		}
		if result.OrphanedAnalysesDeleted != 3 {
			t.Errorf("OrphanedAnalysesDeleted = %d, want 3", result.OrphanedAnalysesDeleted)
		}
		if result.TotalDeleted() != 18 {
			t.Errorf("TotalDeleted() = %d, want 18", result.TotalDeleted())
		}
	})

	t.Run("success - multiple batches", func(t *testing.T) {
		callCount := 0
		repo := &mockCleanupRepository{
			deleteExpiredUserAnalysisHistoryFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				callCount++
				if callCount <= 2 {
					return retention.DeleteResult{DeletedCount: int64(batchSize)}, nil
				}
				return retention.DeleteResult{DeletedCount: 50}, nil
			},
			deleteExpiredSpecDocumentsFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 0}, nil
			},
			deleteOrphanedAnalysesFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 0}, nil
			},
		}

		uc := NewCleanupUseCase(repo, WithBatchSize(100), WithBatchSleep(0))
		result, err := uc.Execute(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if callCount != 3 {
			t.Errorf("callCount = %d, want 3", callCount)
		}
		if result.UserAnalysisHistoryDeleted != 250 {
			t.Errorf("UserAnalysisHistoryDeleted = %d, want 250", result.UserAnalysisHistoryDeleted)
		}
	})

	t.Run("error - user analysis history deletion fails", func(t *testing.T) {
		repo := &mockCleanupRepository{
			deleteExpiredUserAnalysisHistoryFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{}, errors.New("database error")
			},
		}

		uc := NewCleanupUseCase(repo, WithBatchSleep(0))
		_, err := uc.Execute(context.Background())

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("error - spec documents deletion fails", func(t *testing.T) {
		repo := &mockCleanupRepository{
			deleteExpiredUserAnalysisHistoryFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 10}, nil
			},
			deleteExpiredSpecDocumentsFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{}, errors.New("database error")
			},
		}

		uc := NewCleanupUseCase(repo, WithBatchSleep(0))
		_, err := uc.Execute(context.Background())

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("error - orphaned analyses deletion fails", func(t *testing.T) {
		repo := &mockCleanupRepository{
			deleteExpiredUserAnalysisHistoryFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 10}, nil
			},
			deleteExpiredSpecDocumentsFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 5}, nil
			},
			deleteOrphanedAnalysesFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{}, errors.New("database error")
			},
		}

		uc := NewCleanupUseCase(repo, WithBatchSleep(0))
		_, err := uc.Execute(context.Background())

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("context cancellation - stops processing", func(t *testing.T) {
		callCount := 0
		repo := &mockCleanupRepository{
			deleteExpiredUserAnalysisHistoryFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				callCount++
				return retention.DeleteResult{DeletedCount: int64(batchSize)}, nil
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		uc := NewCleanupUseCase(repo, WithBatchSleep(0))
		_, err := uc.Execute(ctx)

		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("nothing to delete", func(t *testing.T) {
		repo := &mockCleanupRepository{
			deleteExpiredUserAnalysisHistoryFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 0}, nil
			},
			deleteExpiredSpecDocumentsFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 0}, nil
			},
			deleteOrphanedAnalysesFn: func(ctx context.Context, batchSize int) (retention.DeleteResult, error) {
				return retention.DeleteResult{DeletedCount: 0}, nil
			},
		}

		uc := NewCleanupUseCase(repo, WithBatchSleep(0))
		result, err := uc.Execute(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.TotalDeleted() != 0 {
			t.Errorf("TotalDeleted() = %d, want 0", result.TotalDeleted())
		}
	})
}

func TestCleanupResult(t *testing.T) {
	t.Run("TotalDeleted", func(t *testing.T) {
		r := CleanupResult{
			UserAnalysisHistoryDeleted: 100,
			SpecDocumentsDeleted:       50,
			OrphanedAnalysesDeleted:    25,
		}
		if got := r.TotalDeleted(); got != 175 {
			t.Errorf("TotalDeleted() = %d, want 175", got)
		}
	})

	t.Run("Duration", func(t *testing.T) {
		start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
		end := time.Date(2026, 1, 1, 12, 5, 30, 0, time.UTC)
		r := CleanupResult{
			StartedAt:   start,
			CompletedAt: end,
		}

		want := 5*time.Minute + 30*time.Second
		if got := r.Duration(); got != want {
			t.Errorf("Duration() = %v, want %v", got, want)
		}
	})
}
