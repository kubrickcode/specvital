package port

import (
	"context"
	"time"

	"github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type SubscriptionRepository interface {
	GetPlanByTier(ctx context.Context, tier entity.PlanTier) (*entity.Plan, error)
	GetAllPlans(ctx context.Context) ([]entity.Plan, error)
	GetPricingPlans(ctx context.Context) ([]entity.PricingPlan, error)
	CreateUserSubscription(ctx context.Context, userID, planID string, periodStart, periodEnd time.Time) (*entity.Subscription, error)
	GetActiveSubscriptionWithPlan(ctx context.Context, userID string) (*entity.SubscriptionWithPlan, error)
	GetUsersWithoutActiveSubscription(ctx context.Context) ([]string, error)
}
