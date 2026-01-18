package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	subscriptionentity "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
	usageentity "github.com/specvital/web/src/backend/modules/usage/domain/entity"
	usageusecase "github.com/specvital/web/src/backend/modules/usage/usecase"
)

type mockSpecViewRepository struct {
	analysisExists bool
	documentExists bool
	status         *entity.SpecGenerationStatus
	analysisErr    error
	documentErr    error
	statusErr      error
}

func (m *mockSpecViewRepository) CheckAnalysisExists(_ context.Context, _ string) (bool, error) {
	return m.analysisExists, m.analysisErr
}

func (m *mockSpecViewRepository) CheckSpecDocumentExists(_ context.Context, _ string) (bool, error) {
	return m.documentExists, m.documentErr
}

func (m *mockSpecViewRepository) GetSpecDocument(_ context.Context, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) GetGenerationStatus(_ context.Context, _ string) (*entity.SpecGenerationStatus, error) {
	return m.status, m.statusErr
}

type mockQueueService struct {
	enqueueErr error
}

func (m *mockQueueService) EnqueueSpecGeneration(_ context.Context, _ string, _ string, _ *string, _ subscriptionentity.PlanTier) error {
	return m.enqueueErr
}

type mockSubscriptionRepository struct {
	subscription *subscriptionentity.SubscriptionWithPlan
	err          error
}

func (m *mockSubscriptionRepository) GetPlanByTier(_ context.Context, _ subscriptionentity.PlanTier) (*subscriptionentity.Plan, error) {
	return nil, nil
}

func (m *mockSubscriptionRepository) GetAllPlans(_ context.Context) ([]subscriptionentity.Plan, error) {
	return nil, nil
}

func (m *mockSubscriptionRepository) CreateUserSubscription(_ context.Context, _, _ string, _, _ time.Time) (*subscriptionentity.Subscription, error) {
	return nil, nil
}

func (m *mockSubscriptionRepository) GetActiveSubscriptionWithPlan(_ context.Context, _ string) (*subscriptionentity.SubscriptionWithPlan, error) {
	return m.subscription, m.err
}

func (m *mockSubscriptionRepository) GetUsersWithoutActiveSubscription(_ context.Context) ([]string, error) {
	return nil, nil
}

func (m *mockSubscriptionRepository) GetPricingPlans(_ context.Context) ([]subscriptionentity.PricingPlan, error) {
	return nil, nil
}

type mockUsageRepository struct {
	monthlyUsage int64
	err          error
}

func (m *mockUsageRepository) GetMonthlyUsage(_ context.Context, _ string, _ usageentity.EventType, _, _ time.Time) (int64, error) {
	return m.monthlyUsage, m.err
}

func (m *mockUsageRepository) GetUsageByPeriod(_ context.Context, _ string, _, _ time.Time) (*usageentity.UsageStats, error) {
	return nil, nil
}

func TestRequestGenerationUseCase_Execute_QuotaCheck(t *testing.T) {
	now := time.Now()
	periodEnd := now.AddDate(0, 1, 0)

	tests := []struct {
		name            string
		input           RequestGenerationInput
		mockRepo        *mockSpecViewRepository
		mockQueue       *mockQueueService
		mockSubRepo     *mockSubscriptionRepository
		mockUsageRepo   *mockUsageRepository
		expectedErr     error
		wantErrContains string
		expectedOutput  bool
	}{
		{
			name: "should allow when user is under quota",
			input: RequestGenerationInput{
				AnalysisID: "test-analysis-id",
				UserID:     "test-user-id",
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
				documentExists: false,
				status:         nil,
			},
			mockQueue: &mockQueueService{},
			mockSubRepo: &mockSubscriptionRepository{
				subscription: &subscriptionentity.SubscriptionWithPlan{
					Subscription: subscriptionentity.Subscription{
						CurrentPeriodStart: now,
						CurrentPeriodEnd:   periodEnd,
					},
					Plan: subscriptionentity.Plan{
						SpecviewMonthlyLimit: ptr(int32(5000)),
					},
				},
			},
			mockUsageRepo: &mockUsageRepository{
				monthlyUsage: 100,
			},
			expectedErr:    nil,
			expectedOutput: true,
		},
		{
			name: "should reject when user exceeds quota",
			input: RequestGenerationInput{
				AnalysisID: "test-analysis-id",
				UserID:     "test-user-id",
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
				documentExists: false,
				status:         nil,
			},
			mockQueue: &mockQueueService{},
			mockSubRepo: &mockSubscriptionRepository{
				subscription: &subscriptionentity.SubscriptionWithPlan{
					Subscription: subscriptionentity.Subscription{
						CurrentPeriodStart: now,
						CurrentPeriodEnd:   periodEnd,
					},
					Plan: subscriptionentity.Plan{
						SpecviewMonthlyLimit: ptr(int32(5000)),
					},
				},
			},
			mockUsageRepo: &mockUsageRepository{
				monthlyUsage: 5000,
			},
			expectedErr:    domain.ErrQuotaExceeded,
			expectedOutput: false,
		},
		{
			name: "should allow enterprise user (unlimited)",
			input: RequestGenerationInput{
				AnalysisID: "test-analysis-id",
				UserID:     "test-user-id",
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
				documentExists: false,
				status:         nil,
			},
			mockQueue: &mockQueueService{},
			mockSubRepo: &mockSubscriptionRepository{
				subscription: &subscriptionentity.SubscriptionWithPlan{
					Subscription: subscriptionentity.Subscription{
						CurrentPeriodStart: now,
						CurrentPeriodEnd:   periodEnd,
					},
					Plan: subscriptionentity.Plan{
						Tier:                 subscriptionentity.PlanTierEnterprise,
						SpecviewMonthlyLimit: nil, // unlimited
					},
				},
			},
			mockUsageRepo: &mockUsageRepository{
				monthlyUsage: 100000,
			},
			expectedErr:    nil,
			expectedOutput: true,
		},
		{
			name: "should skip quota check when no user ID",
			input: RequestGenerationInput{
				AnalysisID: "test-analysis-id",
				UserID:     "",
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
				documentExists: false,
				status:         nil,
			},
			mockQueue:      &mockQueueService{},
			mockSubRepo:    nil,
			mockUsageRepo:  nil,
			expectedErr:    nil,
			expectedOutput: true,
		},
		{
			name: "should propagate error when subscription lookup fails",
			input: RequestGenerationInput{
				AnalysisID: "test-analysis-id",
				UserID:     "test-user-id",
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
				documentExists: false,
				status:         nil,
			},
			mockQueue: &mockQueueService{},
			mockSubRepo: &mockSubscriptionRepository{
				err: errors.New("database connection failed"),
			},
			mockUsageRepo:   &mockUsageRepository{},
			wantErrContains: "database connection failed",
			expectedOutput:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var checkQuotaUC *usageusecase.CheckQuotaUseCase
			if tt.mockSubRepo != nil && tt.mockUsageRepo != nil {
				checkQuotaUC = usageusecase.NewCheckQuotaUseCase(tt.mockSubRepo, tt.mockUsageRepo)
			}

			uc := NewRequestGenerationUseCase(tt.mockRepo, tt.mockQueue, checkQuotaUC)

			result, err := uc.Execute(context.Background(), tt.input)

			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedErr)
					return
				}
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				if result != nil {
					t.Errorf("expected nil result, got %v", result)
				}
			} else if tt.wantErrContains != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantErrContains)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("expected error containing %q, got %v", tt.wantErrContains, err)
				}
				if result != nil {
					t.Errorf("expected nil result, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if tt.expectedOutput && result == nil {
					t.Errorf("expected result, got nil")
				}
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
