package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/subscription/domain"
	"github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

func TestGetUserSubscriptionUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		setupRepo func() *mockSubscriptionRepository
		wantErr   bool
		wantTier  entity.PlanTier
	}{
		{
			name:   "should return active subscription",
			userID: "user-123",
			setupRepo: func() *mockSubscriptionRepository {
				return &mockSubscriptionRepository{
					getActiveSubscriptionWithPlanFunc: func(_ context.Context, userID string) (*entity.SubscriptionWithPlan, error) {
						if userID == "user-123" {
							return &entity.SubscriptionWithPlan{
								Subscription: entity.Subscription{
									ID:                 "sub-id",
									UserID:             userID,
									PlanID:             "plan-free",
									Status:             entity.SubscriptionStatusActive,
									CurrentPeriodStart: time.Now(),
									CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0),
								},
								Plan: entity.Plan{
									ID:   "plan-free",
									Tier: entity.PlanTierFree,
								},
							}, nil
						}
						return nil, domain.ErrNoActiveSubscription
					},
				}
			},
			wantErr:  false,
			wantTier: entity.PlanTierFree,
		},
		{
			name:   "should return error when no active subscription",
			userID: "user-no-sub",
			setupRepo: func() *mockSubscriptionRepository {
				return &mockSubscriptionRepository{
					getActiveSubscriptionWithPlanFunc: func(_ context.Context, _ string) (*entity.SubscriptionWithPlan, error) {
						return nil, domain.ErrNoActiveSubscription
					},
				}
			},
			wantErr: true,
		},
		{
			name:   "should return error on repository failure",
			userID: "user-error",
			setupRepo: func() *mockSubscriptionRepository {
				return &mockSubscriptionRepository{
					getActiveSubscriptionWithPlanFunc: func(_ context.Context, _ string) (*entity.SubscriptionWithPlan, error) {
						return nil, errors.New("db connection error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo()
			uc := NewGetUserSubscriptionUseCase(repo)

			result, err := uc.Execute(context.Background(), tt.userID)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantErr && result.Plan.Tier != tt.wantTier {
				t.Errorf("expected tier %s, got %s", tt.wantTier, result.Plan.Tier)
			}
		})
	}
}
