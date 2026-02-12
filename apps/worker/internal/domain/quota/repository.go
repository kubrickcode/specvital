package quota

import (
	"context"
	"log/slog"
	"time"
)

// CleanupTimeout is the timeout for reservation cleanup operations.
// Set to 5 seconds: sufficient for DB roundtrip, short enough not to block worker shutdown.
const CleanupTimeout = 5 * time.Second

// ReservationRepository defines the interface for quota reservation operations.
// This is used by workers to release reservations on job completion or failure.
type ReservationRepository interface {
	// DeleteByJobID removes a quota reservation by its associated River job ID.
	// Returns nil if no reservation exists (idempotent operation).
	DeleteByJobID(ctx context.Context, jobID int64) error
}

// ReleaseReservation deletes the quota reservation for the given job.
// This is called on both success and failure to ensure reservation cleanup.
// Uses background context to ensure cleanup even if job context is cancelled.
// workerName is used for log identification (e.g., "analyze", "specview").
func ReleaseReservation(repo ReservationRepository, jobID int64, workerName string) {
	if repo == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), CleanupTimeout)
	defer cancel()

	if err := repo.DeleteByJobID(ctx, jobID); err != nil {
		slog.Warn("failed to release quota reservation",
			"job_id", jobID,
			"worker", workerName,
			"error", err,
		)
	}
}
