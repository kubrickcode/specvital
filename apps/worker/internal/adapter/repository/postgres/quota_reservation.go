package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kubrickcode/specvital/apps/worker/internal/domain/quota"
	"github.com/kubrickcode/specvital/apps/worker/internal/infra/db"
)

var _ quota.ReservationRepository = (*QuotaReservationRepository)(nil)

// QuotaReservationRepository implements quota.ReservationRepository using PostgreSQL.
type QuotaReservationRepository struct {
	pool *pgxpool.Pool
}

// NewQuotaReservationRepository creates a new QuotaReservationRepository.
func NewQuotaReservationRepository(pool *pgxpool.Pool) *QuotaReservationRepository {
	return &QuotaReservationRepository{pool: pool}
}

// DeleteByJobID removes a quota reservation by its associated River job ID.
// Returns nil if no reservation exists (idempotent operation).
func (r *QuotaReservationRepository) DeleteByJobID(ctx context.Context, jobID int64) error {
	return db.New(r.pool).DeleteQuotaReservationByJobID(ctx, jobID)
}
