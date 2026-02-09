package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/subscription/domain"
	"github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type mockSubscriptionRepository struct {
	getPlanByTierFunc                     func(ctx context.Context, tier entity.PlanTier) (*entity.Plan, error)
	getAllPlansFunc                       func(ctx context.Context) ([]entity.Plan, error)
	createUserSubscriptionFunc            func(ctx context.Context, userID, planID string, periodStart, periodEnd time.Time) (*entity.Subscription, error)
	getActiveSubscriptionWithPlanFunc     func(ctx context.Context, userID string) (*entity.SubscriptionWithPlan, error)
	getUsersWithoutActiveSubscriptionFunc func(ctx context.Context) ([]string, error)
}

func (m *mockSubscriptionRepository) GetPlanByTier(ctx context.Context, tier entity.PlanTier) (*entity.Plan, error) {
	if m.getPlanByTierFunc != nil {
		return m.getPlanByTierFunc(ctx, tier)
	}
	return nil, domain.ErrPlanNotFound
}

func (m *mockSubscriptionRepository) GetAllPlans(ctx context.Context) ([]entity.Plan, error) {
	if m.getAllPlansFunc != nil {
		return m.getAllPlansFunc(ctx)
	}
	return nil, nil
}

func (m *mockSubscriptionRepository) CreateUserSubscription(ctx context.Context, userID, planID string, periodStart, periodEnd time.Time) (*entity.Subscription, error) {
	if m.createUserSubscriptionFunc != nil {
		return m.createUserSubscriptionFunc(ctx, userID, planID, periodStart, periodEnd)
	}
	return &entity.Subscription{
		ID:                 "sub-id",
		UserID:             userID,
		PlanID:             planID,
		Status:             entity.SubscriptionStatusActive,
		CurrentPeriodStart: periodStart,
		CurrentPeriodEnd:   periodEnd,
		CreatedAt:          time.Now(),
	}, nil
}

func (m *mockSubscriptionRepository) GetActiveSubscriptionWithPlan(ctx context.Context, userID string) (*entity.SubscriptionWithPlan, error) {
	if m.getActiveSubscriptionWithPlanFunc != nil {
		return m.getActiveSubscriptionWithPlanFunc(ctx, userID)
	}
	return nil, domain.ErrNoActiveSubscription
}

func (m *mockSubscriptionRepository) GetUsersWithoutActiveSubscription(ctx context.Context) ([]string, error) {
	if m.getUsersWithoutActiveSubscriptionFunc != nil {
		return m.getUsersWithoutActiveSubscriptionFunc(ctx)
	}
	return nil, nil
}

func (m *mockSubscriptionRepository) GetPricingPlans(_ context.Context) ([]entity.PricingPlan, error) {
	return nil, nil
}

func TestAssignDefaultPlanUseCase_AssignDefaultPlan(t *testing.T) {
	tests := []struct {
		name      string
		tier      entity.PlanTier
		userID    string
		setupRepo func() *mockSubscriptionRepository
		wantErr   bool
	}{
		{
			name:   "should assign free plan to user",
			tier:   entity.PlanTierFree,
			userID: "user-123",
			setupRepo: func() *mockSubscriptionRepository {
				return &mockSubscriptionRepository{
					getPlanByTierFunc: func(_ context.Context, tier entity.PlanTier) (*entity.Plan, error) {
						if tier == entity.PlanTierFree {
							return &entity.Plan{
								ID:   "plan-free",
								Tier: entity.PlanTierFree,
							}, nil
						}
						return nil, domain.ErrPlanNotFound
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "should assign promo plan (pro) when configured",
			tier:   entity.PlanTierPro,
			userID: "user-456",
			setupRepo: func() *mockSubscriptionRepository {
				return &mockSubscriptionRepository{
					getPlanByTierFunc: func(_ context.Context, tier entity.PlanTier) (*entity.Plan, error) {
						if tier == entity.PlanTierPro {
							return &entity.Plan{
								ID:   "plan-pro",
								Tier: entity.PlanTierPro,
							}, nil
						}
						return nil, domain.ErrPlanNotFound
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "should return error when plan not found",
			tier:   entity.PlanTierFree,
			userID: "user-789",
			setupRepo: func() *mockSubscriptionRepository {
				return &mockSubscriptionRepository{
					getPlanByTierFunc: func(_ context.Context, _ entity.PlanTier) (*entity.Plan, error) {
						return nil, domain.ErrPlanNotFound
					},
				}
			},
			wantErr: true,
		},
		{
			name:   "should return error when subscription creation fails",
			tier:   entity.PlanTierFree,
			userID: "user-error",
			setupRepo: func() *mockSubscriptionRepository {
				return &mockSubscriptionRepository{
					getPlanByTierFunc: func(_ context.Context, _ entity.PlanTier) (*entity.Plan, error) {
						return &entity.Plan{ID: "plan-free", Tier: entity.PlanTierFree}, nil
					},
					createUserSubscriptionFunc: func(_ context.Context, _, _ string, _, _ time.Time) (*entity.Subscription, error) {
						return nil, errors.New("db error")
					},
				}
			},
			wantErr: true,
		},
		{
			name:   "should return error when userID is empty",
			tier:   entity.PlanTierFree,
			userID: "",
			setupRepo: func() *mockSubscriptionRepository {
				return &mockSubscriptionRepository{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo()
			uc := NewAssignDefaultPlanUseCase(tt.tier, repo)

			err := uc.AssignDefaultPlan(context.Background(), tt.userID)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
