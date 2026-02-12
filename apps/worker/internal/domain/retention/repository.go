package retention

import "context"

// CleanupRepository defines operations for retention-based data cleanup.
type CleanupRepository interface {
	// DeleteExpiredUserAnalysisHistory removes user analysis history records
	// that have exceeded their retention period.
	// Returns the number of deleted records.
	DeleteExpiredUserAnalysisHistory(ctx context.Context, batchSize int) (DeleteResult, error)

	// DeleteExpiredSpecDocuments removes spec documents
	// that have exceeded their retention period.
	// Returns the number of deleted records.
	DeleteExpiredSpecDocuments(ctx context.Context, batchSize int) (DeleteResult, error)

	// DeleteOrphanedAnalyses removes analyses that have no references
	// in user_analysis_history.
	// Returns the number of deleted records.
	DeleteOrphanedAnalyses(ctx context.Context, batchSize int) (DeleteResult, error)
}

// DeleteResult holds the outcome of a deletion operation.
type DeleteResult struct {
	DeletedCount int64
}

// HasMore returns true if more records may exist for deletion.
func (r DeleteResult) HasMore(batchSize int) bool {
	return r.DeletedCount >= int64(batchSize)
}
