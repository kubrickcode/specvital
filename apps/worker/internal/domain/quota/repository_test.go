package quota

import (
	"context"
	"errors"
	"testing"
)

type mockReservationRepository struct {
	deletedJobIDs []int64
	deleteErr     error
}

func (m *mockReservationRepository) DeleteByJobID(ctx context.Context, jobID int64) error {
	m.deletedJobIDs = append(m.deletedJobIDs, jobID)
	return m.deleteErr
}

func TestReleaseReservation(t *testing.T) {
	t.Run("should call DeleteByJobID with correct job ID", func(t *testing.T) {
		repo := &mockReservationRepository{}
		jobID := int64(12345)

		ReleaseReservation(repo, jobID, "test-worker")

		if len(repo.deletedJobIDs) != 1 {
			t.Errorf("expected 1 delete call, got %d", len(repo.deletedJobIDs))
		}
		if repo.deletedJobIDs[0] != jobID {
			t.Errorf("expected job ID %d, got %d", jobID, repo.deletedJobIDs[0])
		}
	})

	t.Run("should handle nil repository gracefully", func(t *testing.T) {
		// Should not panic
		ReleaseReservation(nil, 12345, "test-worker")
	})

	t.Run("should not propagate delete errors", func(t *testing.T) {
		repo := &mockReservationRepository{
			deleteErr: errors.New("database connection failed"),
		}

		// Should not panic, error is logged but not returned
		ReleaseReservation(repo, 12345, "test-worker")

		// Verify delete was still attempted
		if len(repo.deletedJobIDs) != 1 {
			t.Errorf("expected delete to be attempted even on error")
		}
	})

	t.Run("should use background context with timeout", func(t *testing.T) {
		var capturedCtx context.Context
		repo := &contextCapturingRepo{
			captureCtx: func(ctx context.Context) {
				capturedCtx = ctx
			},
		}

		ReleaseReservation(repo, 12345, "test-worker")

		if capturedCtx == nil {
			t.Fatal("expected context to be captured")
		}

		// Verify context has a deadline (from timeout)
		if _, ok := capturedCtx.Deadline(); !ok {
			t.Error("expected context to have deadline from timeout")
		}
	})
}

type contextCapturingRepo struct {
	captureCtx func(ctx context.Context)
}

func (r *contextCapturingRepo) DeleteByJobID(ctx context.Context, jobID int64) error {
	if r.captureCtx != nil {
		r.captureCtx(ctx)
	}
	return nil
}
