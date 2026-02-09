package usecase

import (
	"context"

	"github.com/specvital/web/src/backend/modules/subscription/domain/entity"
	"github.com/specvital/web/src/backend/modules/subscription/domain/port"
)

type GetPricingPlansUseCase struct {
	repo port.SubscriptionRepository
}

func NewGetPricingPlansUseCase(repo port.SubscriptionRepository) *GetPricingPlansUseCase {
	return &GetPricingPlansUseCase{repo: repo}
}

func (uc *GetPricingPlansUseCase) Execute(ctx context.Context) ([]entity.PricingPlan, error) {
	return uc.repo.GetPricingPlans(ctx)
}
