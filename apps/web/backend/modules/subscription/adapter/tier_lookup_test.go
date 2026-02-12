package adapter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/subscription/domain"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/subscription/domain/entity"
)

type mockSubRepo struct {
	sub *entity.SubscriptionWithPlan
	err error
}

func (m *mockSubRepo) GetPlanByTier(_ context.Context, _ entity.PlanTier) (*entity.Plan, error) {
	return nil, nil
}

func (m *mockSubRepo) GetAllPlans(_ context.Context) ([]entity.Plan, error) {
	return nil, nil
}

func (m *mockSubRepo) GetPricingPlans(_ context.Context) ([]entity.PricingPlan, error) {
	return nil, nil
}

func (m *mockSubRepo) CreateUserSubscription(_ context.Context, _, _ string, _, _ time.Time) (*entity.Subscription, error) {
	return nil, nil
}

func (m *mockSubRepo) GetActiveSubscriptionWithPlan(_ context.Context, _ string) (*entity.SubscriptionWithPlan, error) {
	return m.sub, m.err
}

func (m *mockSubRepo) GetUsersWithoutActiveSubscription(_ context.Context) ([]string, error) {
	return nil, nil
}

func TestTierLookupAdapter_GetUserTier(t *testing.T) {
	t.Run("returns empty tier for empty userID", func(t *testing.T) {
		adapter := NewTierLookupAdapter(&mockSubRepo{})

		tier, err := adapter.GetUserTier(context.Background(), "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if tier != "" {
			t.Errorf("expected empty tier, got %q", tier)
		}
	})

	t.Run("returns tier from subscription", func(t *testing.T) {
		adapter := NewTierLookupAdapter(&mockSubRepo{
			sub: &entity.SubscriptionWithPlan{
				Plan: entity.Plan{Tier: entity.PlanTierPro},
			},
		})

		tier, err := adapter.GetUserTier(context.Background(), "user-123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if tier != string(entity.PlanTierPro) {
			t.Errorf("expected %q, got %q", entity.PlanTierPro, tier)
		}
	})

	t.Run("returns empty tier when no active subscription", func(t *testing.T) {
		adapter := NewTierLookupAdapter(&mockSubRepo{
			err: domain.ErrNoActiveSubscription,
		})

		tier, err := adapter.GetUserTier(context.Background(), "user-123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if tier != "" {
			t.Errorf("expected empty tier, got %q", tier)
		}
	})

	t.Run("returns error for other repository errors", func(t *testing.T) {
		dbErr := errors.New("database error")
		adapter := NewTierLookupAdapter(&mockSubRepo{
			err: dbErr,
		})

		_, err := adapter.GetUserTier(context.Background(), "user-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !errors.Is(err, dbErr) {
			t.Errorf("expected %v, got %v", dbErr, err)
		}
	})
}
