package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/usage/domain/entity"
	"github.com/specvital/web/src/backend/modules/usage/domain/port"
)

type PostgresRepository struct {
	queries *db.Queries
}

var _ port.UsageRepository = (*PostgresRepository)(nil)

func NewPostgresRepository(queries *db.Queries) *PostgresRepository {
	return &PostgresRepository{queries: queries}
}

func (r *PostgresRepository) GetMonthlyUsage(ctx context.Context, userID string, eventType entity.EventType, periodStart, periodEnd time.Time) (int64, error) {
	userUUID, err := uuidFromString(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid userID: %w", err)
	}

	total, err := r.queries.GetMonthlyUsage(ctx, db.GetMonthlyUsageParams{
		UserID:    pgtype.UUID{Bytes: userUUID, Valid: true},
		EventType: db.UsageEventType(eventType),
		CreatedAt: pgtype.Timestamptz{
			Time:  periodStart,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamptz{
			Time:  periodEnd,
			Valid: true,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("query monthly usage: %w", err)
	}

	return total, nil
}

func (r *PostgresRepository) GetUsageByPeriod(ctx context.Context, userID string, periodStart, periodEnd time.Time) (*entity.UsageStats, error) {
	userUUID, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID: %w", err)
	}

	rows, err := r.queries.GetUsageByPeriod(ctx, db.GetUsageByPeriodParams{
		UserID: pgtype.UUID{Bytes: userUUID, Valid: true},
		CreatedAt: pgtype.Timestamptz{
			Time:  periodStart,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamptz{
			Time:  periodEnd,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query usage by period: %w", err)
	}

	stats := &entity.UsageStats{}
	for _, row := range rows {
		switch row.EventType {
		case db.UsageEventTypeSpecview:
			stats.Specview.Used = row.Total
		case db.UsageEventTypeAnalysis:
			stats.Analysis.Used = row.Total
		}
	}

	return stats, nil
}

func uuidFromString(s string) ([16]byte, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(s); err != nil {
		return [16]byte{}, fmt.Errorf("parse UUID %q: %w", s, err)
	}
	return uuid.Bytes, nil
}
