package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/subscription/domain"
	"github.com/specvital/web/src/backend/modules/subscription/domain/entity"
	"github.com/specvital/web/src/backend/modules/subscription/domain/port"
)

type PostgresRepository struct {
	queries *db.Queries
}

var _ port.SubscriptionRepository = (*PostgresRepository)(nil)

func NewPostgresRepository(queries *db.Queries) *PostgresRepository {
	return &PostgresRepository{queries: queries}
}

func (r *PostgresRepository) GetPlanByTier(ctx context.Context, tier entity.PlanTier) (*entity.Plan, error) {
	row, err := r.queries.GetPlanByTier(ctx, db.PlanTier(tier))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPlanNotFound
		}
		return nil, fmt.Errorf("query plan by tier: %w", err)
	}

	plan := &entity.Plan{
		ID:   uuidToString(row.ID),
		Tier: entity.PlanTier(row.Tier),
	}
	if row.SpecviewMonthlyLimit.Valid {
		plan.SpecviewMonthlyLimit = &row.SpecviewMonthlyLimit.Int32
	}
	if row.AnalysisMonthlyLimit.Valid {
		plan.AnalysisMonthlyLimit = &row.AnalysisMonthlyLimit.Int32
	}
	if row.RetentionDays.Valid {
		plan.RetentionDays = &row.RetentionDays.Int32
	}
	return plan, nil
}

func (r *PostgresRepository) GetAllPlans(ctx context.Context) ([]entity.Plan, error) {
	rows, err := r.queries.GetAllPlans(ctx)
	if err != nil {
		return nil, fmt.Errorf("query all plans: %w", err)
	}

	plans := make([]entity.Plan, len(rows))
	for i, row := range rows {
		plan := entity.Plan{
			ID:   uuidToString(row.ID),
			Tier: entity.PlanTier(row.Tier),
		}
		if row.SpecviewMonthlyLimit.Valid {
			plan.SpecviewMonthlyLimit = &row.SpecviewMonthlyLimit.Int32
		}
		if row.AnalysisMonthlyLimit.Valid {
			plan.AnalysisMonthlyLimit = &row.AnalysisMonthlyLimit.Int32
		}
		if row.RetentionDays.Valid {
			plan.RetentionDays = &row.RetentionDays.Int32
		}
		plans[i] = plan
	}

	return plans, nil
}

func (r *PostgresRepository) CreateUserSubscription(ctx context.Context, userID, planID string, periodStart, periodEnd time.Time) (*entity.Subscription, error) {
	userUUID, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID: %w", err)
	}
	planUUID, err := uuidFromString(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid planID: %w", err)
	}

	row, err := r.queries.CreateUserSubscription(ctx, db.CreateUserSubscriptionParams{
		UserID: pgtype.UUID{Bytes: userUUID, Valid: true},
		PlanID: pgtype.UUID{Bytes: planUUID, Valid: true},
		CurrentPeriodStart: pgtype.Timestamptz{
			Time:  periodStart,
			Valid: true,
		},
		CurrentPeriodEnd: pgtype.Timestamptz{
			Time:  periodEnd,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create user subscription: %w", err)
	}

	return &entity.Subscription{
		ID:                 uuidToString(row.ID),
		UserID:             uuidToString(row.UserID),
		PlanID:             uuidToString(row.PlanID),
		Status:             entity.SubscriptionStatus(row.Status),
		CurrentPeriodStart: row.CurrentPeriodStart.Time,
		CurrentPeriodEnd:   row.CurrentPeriodEnd.Time,
		CreatedAt:          row.CreatedAt.Time,
	}, nil
}

func (r *PostgresRepository) GetActiveSubscriptionWithPlan(ctx context.Context, userID string) (*entity.SubscriptionWithPlan, error) {
	userUUID, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID: %w", err)
	}

	row, err := r.queries.GetActiveSubscriptionWithPlan(ctx, pgtype.UUID{
		Bytes: userUUID,
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNoActiveSubscription
		}
		return nil, fmt.Errorf("query active subscription with plan: %w", err)
	}

	sub := &entity.SubscriptionWithPlan{
		Subscription: entity.Subscription{
			ID:                 uuidToString(row.ID),
			UserID:             uuidToString(row.UserID),
			PlanID:             uuidToString(row.PlanID),
			Status:             entity.SubscriptionStatus(row.Status),
			CurrentPeriodStart: row.CurrentPeriodStart.Time,
			CurrentPeriodEnd:   row.CurrentPeriodEnd.Time,
			CreatedAt:          row.CreatedAt.Time,
			UpdatedAt:          row.UpdatedAt.Time,
		},
		Plan: entity.Plan{
			ID:   uuidToString(row.PlanID),
			Tier: entity.PlanTier(row.PlanTier),
		},
	}

	if row.CanceledAt.Valid {
		sub.Subscription.CanceledAt = &row.CanceledAt.Time
	}

	if row.PlanSpecviewMonthlyLimit.Valid {
		sub.Plan.SpecviewMonthlyLimit = &row.PlanSpecviewMonthlyLimit.Int32
	}
	if row.PlanAnalysisMonthlyLimit.Valid {
		sub.Plan.AnalysisMonthlyLimit = &row.PlanAnalysisMonthlyLimit.Int32
	}
	if row.PlanRetentionDays.Valid {
		sub.Plan.RetentionDays = &row.PlanRetentionDays.Int32
	}

	return sub, nil
}

func (r *PostgresRepository) GetUsersWithoutActiveSubscription(ctx context.Context) ([]string, error) {
	uuids, err := r.queries.GetUsersWithoutActiveSubscription(ctx)
	if err != nil {
		return nil, fmt.Errorf("query users without subscription: %w", err)
	}

	ids := make([]string, len(uuids))
	for i, u := range uuids {
		ids[i] = uuidToString(u)
	}

	return ids, nil
}

func (r *PostgresRepository) GetPricingPlans(ctx context.Context) ([]entity.PricingPlan, error) {
	rows, err := r.queries.GetPricingPlans(ctx)
	if err != nil {
		return nil, fmt.Errorf("query pricing plans: %w", err)
	}

	plans := make([]entity.PricingPlan, len(rows))
	for i, row := range rows {
		plan := entity.PricingPlan{
			Tier: entity.PlanTier(row.Tier),
		}
		if row.MonthlyPrice.Valid {
			plan.MonthlyPrice = &row.MonthlyPrice.Int32
		}
		if row.SpecviewMonthlyLimit.Valid {
			plan.SpecviewMonthlyLimit = &row.SpecviewMonthlyLimit.Int32
		}
		if row.AnalysisMonthlyLimit.Valid {
			plan.AnalysisMonthlyLimit = &row.AnalysisMonthlyLimit.Int32
		}
		if row.RetentionDays.Valid {
			plan.RetentionDays = &row.RetentionDays.Int32
		}
		plans[i] = plan
	}

	return plans, nil
}

func uuidFromString(s string) ([16]byte, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(s); err != nil {
		return [16]byte{}, fmt.Errorf("parse UUID %q: %w", s, err)
	}
	return uuid.Bytes, nil
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", u.Bytes[0:4], u.Bytes[4:6], u.Bytes[6:8], u.Bytes[8:10], u.Bytes[10:16])
}
