package retention

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/specvital/worker/internal/domain/retention"
)

const (
	DefaultBatchSleepDuration = 100 * time.Millisecond
)

// CleanupUseCase orchestrates retention-based data cleanup.
type CleanupUseCase struct {
	batchSize      int
	batchSleep     time.Duration
	cleanupRepo    retention.CleanupRepository
}

// Option configures CleanupUseCase.
type Option func(*CleanupUseCase)

// WithBatchSize sets the batch size for cleanup operations.
func WithBatchSize(size int) Option {
	return func(uc *CleanupUseCase) {
		if size > 0 {
			uc.batchSize = size
		}
	}
}

// WithBatchSleep sets the sleep duration between batches.
func WithBatchSleep(d time.Duration) Option {
	return func(uc *CleanupUseCase) {
		if d >= 0 {
			uc.batchSleep = d
		}
	}
}

// NewCleanupUseCase creates a CleanupUseCase with the given repository.
func NewCleanupUseCase(repo retention.CleanupRepository, opts ...Option) *CleanupUseCase {
	uc := &CleanupUseCase{
		batchSize:   retention.DefaultBatchSize,
		batchSleep:  DefaultBatchSleepDuration,
		cleanupRepo: repo,
	}

	for _, opt := range opts {
		opt(uc)
	}

	return uc
}

// CleanupResult aggregates the outcome of a complete cleanup run.
type CleanupResult struct {
	UserAnalysisHistoryDeleted int64
	SpecDocumentsDeleted       int64
	OrphanedAnalysesDeleted    int64
	StartedAt                  time.Time
	CompletedAt                time.Time
}

// TotalDeleted returns the total number of records deleted.
func (r CleanupResult) TotalDeleted() int64 {
	return r.UserAnalysisHistoryDeleted + r.SpecDocumentsDeleted + r.OrphanedAnalysesDeleted
}

// Duration returns how long the cleanup took.
func (r CleanupResult) Duration() time.Duration {
	return r.CompletedAt.Sub(r.StartedAt)
}

// Execute performs the two-phase cleanup process.
// Phase 1: Delete expired user data (user_analysis_history, spec_documents)
// Phase 2: Delete orphaned analyses (no references in user_analysis_history)
func (uc *CleanupUseCase) Execute(ctx context.Context) (CleanupResult, error) {
	result := CleanupResult{
		StartedAt: time.Now(),
	}

	slog.InfoContext(ctx, "starting retention cleanup",
		"batch_size", uc.batchSize,
	)

	// Phase 1: Delete expired user analysis history
	historyDeleted, err := uc.deleteInBatches(ctx, "user_analysis_history", uc.cleanupRepo.DeleteExpiredUserAnalysisHistory)
	if err != nil {
		return result, fmt.Errorf("delete expired user analysis history: %w", err)
	}
	result.UserAnalysisHistoryDeleted = historyDeleted

	// Phase 1: Delete expired spec documents
	specDocsDeleted, err := uc.deleteInBatches(ctx, "spec_documents", uc.cleanupRepo.DeleteExpiredSpecDocuments)
	if err != nil {
		return result, fmt.Errorf("delete expired spec documents: %w", err)
	}
	result.SpecDocumentsDeleted = specDocsDeleted

	// Phase 2: Delete orphaned analyses
	orphansDeleted, err := uc.deleteInBatches(ctx, "orphaned_analyses", uc.cleanupRepo.DeleteOrphanedAnalyses)
	if err != nil {
		return result, fmt.Errorf("delete orphaned analyses: %w", err)
	}
	result.OrphanedAnalysesDeleted = orphansDeleted

	result.CompletedAt = time.Now()

	slog.InfoContext(ctx, "retention cleanup completed",
		"user_analysis_history_deleted", result.UserAnalysisHistoryDeleted,
		"spec_documents_deleted", result.SpecDocumentsDeleted,
		"orphaned_analyses_deleted", result.OrphanedAnalysesDeleted,
		"total_deleted", result.TotalDeleted(),
		"duration", result.Duration(),
	)

	return result, nil
}

type deleteFunc func(ctx context.Context, batchSize int) (retention.DeleteResult, error)

func (uc *CleanupUseCase) deleteInBatches(ctx context.Context, target string, deleteFn deleteFunc) (int64, error) {
	var totalDeleted int64

	for {
		select {
		case <-ctx.Done():
			return totalDeleted, ctx.Err()
		default:
		}

		result, err := deleteFn(ctx, uc.batchSize)
		if err != nil {
			return totalDeleted, err
		}

		totalDeleted += result.DeletedCount

		if result.DeletedCount > 0 {
			slog.DebugContext(ctx, "batch deleted",
				"target", target,
				"deleted", result.DeletedCount,
				"total", totalDeleted,
			)
		}

		if !result.HasMore(uc.batchSize) {
			break
		}

		if uc.batchSleep > 0 {
			select {
			case <-ctx.Done():
				return totalDeleted, ctx.Err()
			case <-time.After(uc.batchSleep):
			}
		}
	}

	return totalDeleted, nil
}
