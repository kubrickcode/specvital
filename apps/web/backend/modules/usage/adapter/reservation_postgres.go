package adapter

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kubrickcode/specvital/apps/web/backend/internal/db"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/usage/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/usage/domain/port"
)

type QuotaReservationPostgresRepository struct {
	queries *db.Queries
}

var _ port.QuotaReservationRepository = (*QuotaReservationPostgresRepository)(nil)

func NewQuotaReservationPostgresRepository(queries *db.Queries) *QuotaReservationPostgresRepository {
	return &QuotaReservationPostgresRepository{queries: queries}
}

func (r *QuotaReservationPostgresRepository) CreateReservation(ctx context.Context, userID string, eventType entity.EventType, amount int32, jobID int64) error {
	return r.CreateReservationTx(ctx, r.queries, userID, eventType, amount, jobID)
}

func (r *QuotaReservationPostgresRepository) CreateReservationTx(ctx context.Context, qtx *db.Queries, userID string, eventType entity.EventType, amount int32, jobID int64) error {
	userUUID, err := uuidFromString(userID)
	if err != nil {
		return fmt.Errorf("invalid userID: %w", err)
	}

	_, err = qtx.CreateQuotaReservation(ctx, db.CreateQuotaReservationParams{
		UserID:         pgtype.UUID{Bytes: userUUID, Valid: true},
		EventType:      db.UsageEventType(eventType),
		ReservedAmount: amount,
		JobID:          jobID,
	})
	if err != nil {
		return fmt.Errorf("create quota reservation: %w", err)
	}

	return nil
}

func (r *QuotaReservationPostgresRepository) GetTotalReservedAmount(ctx context.Context, userID string, eventType entity.EventType) (int64, error) {
	userUUID, err := uuidFromString(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid userID: %w", err)
	}

	total, err := r.queries.GetTotalReservedAmount(ctx, db.GetTotalReservedAmountParams{
		UserID:    pgtype.UUID{Bytes: userUUID, Valid: true},
		EventType: db.UsageEventType(eventType),
	})
	if err != nil {
		return 0, fmt.Errorf("query reserved amount: %w", err)
	}

	return total, nil
}

func (r *QuotaReservationPostgresRepository) DeleteReservationByJobID(ctx context.Context, jobID int64) error {
	if err := r.queries.DeleteQuotaReservationByJobID(ctx, jobID); err != nil {
		return fmt.Errorf("delete quota reservation: %w", err)
	}
	return nil
}
