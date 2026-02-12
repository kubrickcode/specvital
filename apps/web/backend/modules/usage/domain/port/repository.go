package port

import (
	"context"
	"time"

	"github.com/kubrickcode/specvital/apps/web/backend/internal/db"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/usage/domain/entity"
)

type UsageRepository interface {
	GetMonthlyUsage(ctx context.Context, userID string, eventType entity.EventType, periodStart, periodEnd time.Time) (int64, error)
	GetUsageByPeriod(ctx context.Context, userID string, periodStart, periodEnd time.Time) (*entity.UsageStats, error)
}

type QuotaReservationRepository interface {
	CreateReservation(ctx context.Context, userID string, eventType entity.EventType, amount int32, jobID int64) error
	// CreateReservationTx creates a reservation within a transaction.
	CreateReservationTx(ctx context.Context, qtx *db.Queries, userID string, eventType entity.EventType, amount int32, jobID int64) error
	GetTotalReservedAmount(ctx context.Context, userID string, eventType entity.EventType) (int64, error)
	DeleteReservationByJobID(ctx context.Context, jobID int64) error
}
