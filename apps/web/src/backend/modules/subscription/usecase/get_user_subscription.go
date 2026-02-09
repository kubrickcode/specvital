package usecase

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/modules/subscription/domain/entity"
	"github.com/specvital/web/src/backend/modules/subscription/domain/port"
)

type GetUserSubscriptionUseCase struct {
	repo port.SubscriptionRepository
}

func NewGetUserSubscriptionUseCase(repo port.SubscriptionRepository) *GetUserSubscriptionUseCase {
	if repo == nil {
		panic("repo is required")
	}
	return &GetUserSubscriptionUseCase{repo: repo}
}

func (uc *GetUserSubscriptionUseCase) Execute(ctx context.Context, userID string) (*entity.SubscriptionWithPlan, error) {
	sub, err := uc.repo.GetActiveSubscriptionWithPlan(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get active subscription: %w", err)
	}
	return sub, nil
}
