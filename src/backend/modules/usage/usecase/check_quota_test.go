package usecase_test

import (
	"context"
	"testing"
	"time"

	subscriptionentity "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
	usageentity "github.com/specvital/web/src/backend/modules/usage/domain/entity"
	"github.com/specvital/web/src/backend/modules/usage/usecase"
)

type mockSubscriptionRepo struct {
	subscription *subscriptionentity.SubscriptionWithPlan
	err          error
}

func (m *mockSubscriptionRepo) GetPlanByTier(_ context.Context, _ subscriptionentity.PlanTier) (*subscriptionentity.Plan, error) {
	return nil, nil
}

func (m *mockSubscriptionRepo) GetAllPlans(_ context.Context) ([]subscriptionentity.Plan, error) {
	return nil, nil
}

func (m *mockSubscriptionRepo) CreateUserSubscription(_ context.Context, _, _ string, _, _ time.Time) (*subscriptionentity.Subscription, error) {
	return nil, nil
}

func (m *mockSubscriptionRepo) GetActiveSubscriptionWithPlan(_ context.Context, _ string) (*subscriptionentity.SubscriptionWithPlan, error) {
	return m.subscription, m.err
}

func (m *mockSubscriptionRepo) GetUsersWithoutActiveSubscription(_ context.Context) ([]string, error) {
	return nil, nil
}

func (m *mockSubscriptionRepo) GetPricingPlans(_ context.Context) ([]subscriptionentity.PricingPlan, error) {
	return nil, nil
}

type mockUsageRepo struct {
	monthlyUsage int64
	usageStats   *usageentity.UsageStats
	err          error
}

func (m *mockUsageRepo) GetMonthlyUsage(_ context.Context, _ string, _ usageentity.EventType, _, _ time.Time) (int64, error) {
	return m.monthlyUsage, m.err
}

func (m *mockUsageRepo) GetUsageByPeriod(_ context.Context, _ string, _, _ time.Time) (*usageentity.UsageStats, error) {
	return m.usageStats, m.err
}

type mockReservationRepo struct {
	reservedAmount int64
	err            error
}

func (m *mockReservationRepo) CreateReservation(_ context.Context, _ string, _ usageentity.EventType, _ int32, _ int64) error {
	return m.err
}

func (m *mockReservationRepo) GetTotalReservedAmount(_ context.Context, _ string, _ usageentity.EventType) (int64, error) {
	return m.reservedAmount, m.err
}

func (m *mockReservationRepo) DeleteReservationByJobID(_ context.Context, _ int64) error {
	return m.err
}

func TestCheckQuotaUseCase_Execute(t *testing.T) {
	now := time.Now()
	periodStart := now.AddDate(0, -1, 0)
	periodEnd := now.AddDate(0, 1, 0)

	specviewLimit := int32(5000)
	analysisLimit := int32(50)

	tests := []struct {
		name           string
		input          usecase.CheckQuotaInput
		subscription   *subscriptionentity.SubscriptionWithPlan
		monthlyUsage   int64
		reservedAmount int64
		wantIsAllowed  bool
		wantReserved   int64
		wantLimit      *int32
	}{
		{
			name: "free user under limit",
			input: usecase.CheckQuotaInput{
				UserID:    "user-1",
				EventType: usageentity.EventTypeSpecview,
				Amount:    1,
			},
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierFree,
					SpecviewMonthlyLimit: &specviewLimit,
					AnalysisMonthlyLimit: &analysisLimit,
				},
			},
			monthlyUsage:   100,
			reservedAmount: 0,
			wantIsAllowed:  true,
			wantReserved:   0,
			wantLimit:      &specviewLimit,
		},
		{
			name: "free user at limit",
			input: usecase.CheckQuotaInput{
				UserID:    "user-1",
				EventType: usageentity.EventTypeSpecview,
				Amount:    1,
			},
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierFree,
					SpecviewMonthlyLimit: &specviewLimit,
					AnalysisMonthlyLimit: &analysisLimit,
				},
			},
			monthlyUsage:   5000,
			reservedAmount: 0,
			wantIsAllowed:  false,
			wantReserved:   0,
			wantLimit:      &specviewLimit,
		},
		{
			name: "enterprise user unlimited",
			input: usecase.CheckQuotaInput{
				UserID:    "user-1",
				EventType: usageentity.EventTypeSpecview,
				Amount:    1,
			},
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierEnterprise,
					SpecviewMonthlyLimit: nil,
					AnalysisMonthlyLimit: nil,
				},
			},
			monthlyUsage:   1000000,
			reservedAmount: 500,
			wantIsAllowed:  true,
			wantReserved:   500,
			wantLimit:      nil,
		},
		{
			name: "analysis quota check",
			input: usecase.CheckQuotaInput{
				UserID:    "user-1",
				EventType: usageentity.EventTypeAnalysis,
				Amount:    1,
			},
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierFree,
					SpecviewMonthlyLimit: &specviewLimit,
					AnalysisMonthlyLimit: &analysisLimit,
				},
			},
			monthlyUsage:   49,
			reservedAmount: 0,
			wantIsAllowed:  true,
			wantReserved:   0,
			wantLimit:      &analysisLimit,
		},
		{
			name: "rejected when used + reserved exceeds limit",
			input: usecase.CheckQuotaInput{
				UserID:    "user-1",
				EventType: usageentity.EventTypeSpecview,
				Amount:    100,
			},
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierFree,
					SpecviewMonthlyLimit: &specviewLimit,
					AnalysisMonthlyLimit: &analysisLimit,
				},
			},
			monthlyUsage:   4500,
			reservedAmount: 450,
			wantIsAllowed:  false,
			wantReserved:   450,
			wantLimit:      &specviewLimit,
		},
		{
			name: "allowed when used + reserved + amount equals limit",
			input: usecase.CheckQuotaInput{
				UserID:    "user-1",
				EventType: usageentity.EventTypeSpecview,
				Amount:    50,
			},
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierFree,
					SpecviewMonthlyLimit: &specviewLimit,
					AnalysisMonthlyLimit: &analysisLimit,
				},
			},
			monthlyUsage:   4500,
			reservedAmount: 450,
			wantIsAllowed:  true,
			wantReserved:   450,
			wantLimit:      &specviewLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscriptionRepo := &mockSubscriptionRepo{subscription: tt.subscription}
			usageRepo := &mockUsageRepo{monthlyUsage: tt.monthlyUsage}
			reservationRepo := &mockReservationRepo{reservedAmount: tt.reservedAmount}

			uc := usecase.NewCheckQuotaUseCase(subscriptionRepo, usageRepo, reservationRepo)
			result, err := uc.Execute(context.Background(), tt.input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.IsAllowed != tt.wantIsAllowed {
				t.Errorf("IsAllowed = %v, want %v", result.IsAllowed, tt.wantIsAllowed)
			}

			if tt.wantLimit == nil {
				if result.Limit != nil {
					t.Errorf("Limit = %v, want nil", *result.Limit)
				}
			} else {
				if result.Limit == nil {
					t.Errorf("Limit = nil, want %v", *tt.wantLimit)
				} else if *result.Limit != *tt.wantLimit {
					t.Errorf("Limit = %v, want %v", *result.Limit, *tt.wantLimit)
				}
			}

			if result.Used != tt.monthlyUsage {
				t.Errorf("Used = %v, want %v", result.Used, tt.monthlyUsage)
			}

			if result.Reserved != tt.wantReserved {
				t.Errorf("Reserved = %v, want %v", result.Reserved, tt.wantReserved)
			}
		})
	}
}
