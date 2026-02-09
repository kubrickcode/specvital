package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/specvital/web/src/backend/modules/subscription/domain/entity"
	"github.com/specvital/web/src/backend/modules/subscription/domain/port"
)

type AssignDefaultPlanUseCase struct {
	defaultTier entity.PlanTier
	repo        port.SubscriptionRepository
}

var _ port.Subscriber = (*AssignDefaultPlanUseCase)(nil)

func NewAssignDefaultPlanUseCase(defaultTier entity.PlanTier, repo port.SubscriptionRepository) *AssignDefaultPlanUseCase {
	if !defaultTier.IsValid() {
		panic(fmt.Sprintf("invalid default tier: %s", defaultTier))
	}
	if repo == nil {
		panic("repo is required")
	}
	return &AssignDefaultPlanUseCase{
		defaultTier: defaultTier,
		repo:        repo,
	}
}

func (uc *AssignDefaultPlanUseCase) AssignDefaultPlan(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("assign default plan: userID is required")
	}

	plan, err := uc.repo.GetPlanByTier(ctx, uc.defaultTier)
	if err != nil {
		return fmt.Errorf("get plan by tier %s: %w", uc.defaultTier, err)
	}

	now := time.Now().UTC()
	periodStart := now
	periodEnd := now.AddDate(0, 1, 0).Add(-time.Nanosecond)

	_, err = uc.repo.CreateUserSubscription(ctx, userID, plan.ID, periodStart, periodEnd)
	if err != nil {
		return fmt.Errorf("create subscription: %w", err)
	}

	return nil
}
