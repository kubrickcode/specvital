package usecase_test

import (
	"context"
	"testing"
	"time"

	subscriptionentity "github.com/kubrickcode/specvital/apps/web/backend/modules/subscription/domain/entity"
	usageentity "github.com/kubrickcode/specvital/apps/web/backend/modules/usage/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/usage/usecase"
)

func TestGetCurrentUsageUseCase_Execute(t *testing.T) {
	now := time.Now()
	periodStart := now.AddDate(0, -1, 0)
	periodEnd := now.AddDate(0, 1, 0)

	specviewLimit := int32(5000)
	analysisLimit := int32(50)
	retentionDays := int32(30)

	tests := []struct {
		name                 string
		subscription         *subscriptionentity.SubscriptionWithPlan
		usageStats           *usageentity.UsageStats
		reservedAmount       int64
		wantSpecviewUsed     int64
		wantSpecviewReserved int64
		wantAnalysisUsed     int64
		wantAnalysisReserved int64
		wantTier             subscriptionentity.PlanTier
		wantNilLimits        bool
		wantNilPercentage    bool
	}{
		{
			name: "free user with usage",
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierFree,
					SpecviewMonthlyLimit: &specviewLimit,
					AnalysisMonthlyLimit: &analysisLimit,
					RetentionDays:        &retentionDays,
				},
			},
			usageStats: &usageentity.UsageStats{
				Specview: usageentity.UsageMetric{Used: 2500},
				Analysis: usageentity.UsageMetric{Used: 25},
			},
			reservedAmount:       0,
			wantSpecviewUsed:     2500,
			wantSpecviewReserved: 0,
			wantAnalysisUsed:     25,
			wantAnalysisReserved: 0,
			wantTier:             subscriptionentity.PlanTierFree,
			wantNilLimits:        false,
			wantNilPercentage:    false,
		},
		{
			name: "free user with reserved quota",
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierFree,
					SpecviewMonthlyLimit: &specviewLimit,
					AnalysisMonthlyLimit: &analysisLimit,
					RetentionDays:        &retentionDays,
				},
			},
			usageStats: &usageentity.UsageStats{
				Specview: usageentity.UsageMetric{Used: 1000},
				Analysis: usageentity.UsageMetric{Used: 10},
			},
			reservedAmount:       500,
			wantSpecviewUsed:     1000,
			wantSpecviewReserved: 500,
			wantAnalysisUsed:     10,
			wantAnalysisReserved: 500,
			wantTier:             subscriptionentity.PlanTierFree,
			wantNilLimits:        false,
			wantNilPercentage:    false,
		},
		{
			name: "enterprise user unlimited",
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierEnterprise,
					SpecviewMonthlyLimit: nil,
					AnalysisMonthlyLimit: nil,
					RetentionDays:        nil,
				},
			},
			usageStats: &usageentity.UsageStats{
				Specview: usageentity.UsageMetric{Used: 100000},
				Analysis: usageentity.UsageMetric{Used: 5000},
			},
			reservedAmount:       0,
			wantSpecviewUsed:     100000,
			wantSpecviewReserved: 0,
			wantAnalysisUsed:     5000,
			wantAnalysisReserved: 0,
			wantTier:             subscriptionentity.PlanTierEnterprise,
			wantNilLimits:        true,
			wantNilPercentage:    true,
		},
		{
			name: "enterprise user with reserved quota",
			subscription: &subscriptionentity.SubscriptionWithPlan{
				Subscription: subscriptionentity.Subscription{
					CurrentPeriodStart: periodStart,
					CurrentPeriodEnd:   periodEnd,
				},
				Plan: subscriptionentity.Plan{
					Tier:                 subscriptionentity.PlanTierEnterprise,
					SpecviewMonthlyLimit: nil,
					AnalysisMonthlyLimit: nil,
					RetentionDays:        nil,
				},
			},
			usageStats: &usageentity.UsageStats{
				Specview: usageentity.UsageMetric{Used: 50000},
				Analysis: usageentity.UsageMetric{Used: 2500},
			},
			reservedAmount:       1000,
			wantSpecviewUsed:     50000,
			wantSpecviewReserved: 1000,
			wantAnalysisUsed:     2500,
			wantAnalysisReserved: 1000,
			wantTier:             subscriptionentity.PlanTierEnterprise,
			wantNilLimits:        true,
			wantNilPercentage:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscriptionRepo := &mockSubscriptionRepo{subscription: tt.subscription}
			usageRepo := &mockUsageRepo{usageStats: tt.usageStats}
			reservationRepo := &mockReservationRepo{reservedAmount: tt.reservedAmount}

			uc := usecase.NewGetCurrentUsageUseCase(subscriptionRepo, usageRepo, reservationRepo)
			result, err := uc.Execute(context.Background(), usecase.GetCurrentUsageInput{
				UserID: "user-1",
			})

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Specview.Used != tt.wantSpecviewUsed {
				t.Errorf("Specview.Used = %v, want %v", result.Specview.Used, tt.wantSpecviewUsed)
			}

			if result.Specview.Reserved != tt.wantSpecviewReserved {
				t.Errorf("Specview.Reserved = %v, want %v", result.Specview.Reserved, tt.wantSpecviewReserved)
			}

			if result.Analysis.Used != tt.wantAnalysisUsed {
				t.Errorf("Analysis.Used = %v, want %v", result.Analysis.Used, tt.wantAnalysisUsed)
			}

			if result.Analysis.Reserved != tt.wantAnalysisReserved {
				t.Errorf("Analysis.Reserved = %v, want %v", result.Analysis.Reserved, tt.wantAnalysisReserved)
			}

			if result.Plan.Tier != tt.wantTier {
				t.Errorf("Plan.Tier = %v, want %v", result.Plan.Tier, tt.wantTier)
			}

			if tt.wantNilLimits {
				if result.Specview.Limit != nil {
					t.Errorf("Specview.Limit = %v, want nil", *result.Specview.Limit)
				}
				if result.Analysis.Limit != nil {
					t.Errorf("Analysis.Limit = %v, want nil", *result.Analysis.Limit)
				}
			} else {
				if result.Specview.Limit == nil {
					t.Error("Specview.Limit = nil, want non-nil")
				}
				if result.Analysis.Limit == nil {
					t.Error("Analysis.Limit = nil, want non-nil")
				}
			}

			if tt.wantNilPercentage {
				if result.Specview.Percentage != nil {
					t.Errorf("Specview.Percentage = %v, want nil", *result.Specview.Percentage)
				}
				if result.Analysis.Percentage != nil {
					t.Errorf("Analysis.Percentage = %v, want nil", *result.Analysis.Percentage)
				}
			} else {
				if result.Specview.Percentage == nil {
					t.Error("Specview.Percentage = nil, want non-nil")
				}
				if result.Analysis.Percentage == nil {
					t.Error("Analysis.Percentage = nil, want non-nil")
				}
			}
		})
	}
}

func TestBuildUsageMetric_Percentage(t *testing.T) {
	now := time.Now()
	periodStart := now.AddDate(0, -1, 0)
	periodEnd := now.AddDate(0, 1, 0)

	limit := int32(100)

	subscription := &subscriptionentity.SubscriptionWithPlan{
		Subscription: subscriptionentity.Subscription{
			CurrentPeriodStart: periodStart,
			CurrentPeriodEnd:   periodEnd,
		},
		Plan: subscriptionentity.Plan{
			Tier:                 subscriptionentity.PlanTierFree,
			SpecviewMonthlyLimit: &limit,
			AnalysisMonthlyLimit: &limit,
		},
	}

	usageStats := &usageentity.UsageStats{
		Specview: usageentity.UsageMetric{Used: 50},
		Analysis: usageentity.UsageMetric{Used: 75},
	}

	subscriptionRepo := &mockSubscriptionRepo{subscription: subscription}
	usageRepo := &mockUsageRepo{usageStats: usageStats}
	reservationRepo := &mockReservationRepo{reservedAmount: 0}

	uc := usecase.NewGetCurrentUsageUseCase(subscriptionRepo, usageRepo, reservationRepo)
	result, err := uc.Execute(context.Background(), usecase.GetCurrentUsageInput{
		UserID: "user-1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Specview.Percentage == nil {
		t.Fatal("Specview.Percentage = nil")
	}
	if *result.Specview.Percentage != 50 {
		t.Errorf("Specview.Percentage = %v, want 50", *result.Specview.Percentage)
	}

	if result.Analysis.Percentage == nil {
		t.Fatal("Analysis.Percentage = nil")
	}
	if *result.Analysis.Percentage != 75 {
		t.Errorf("Analysis.Percentage = %v, want 75", *result.Analysis.Percentage)
	}
}
