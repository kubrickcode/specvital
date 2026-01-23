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

func (m *mockSpecViewRepository) CheckSpecDocumentExistsByLanguage(_ context.Context, _ string, _ string) (bool, error) {
	return m.documentExists, m.documentErr
}

func (m *mockSpecViewRepository) CheckSpecDocumentOwnership(_ context.Context, _ string) (*entity.DocumentOwnership, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) GetAvailableLanguages(_ context.Context, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) GetAvailableLanguagesByUser(_ context.Context, _ string, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) GetSpecDocumentByLanguage(_ context.Context, _ string, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) GetSpecDocumentByUser(_ context.Context, _ string, _ string, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) GetGenerationStatus(_ context.Context, _ string) (*entity.SpecGenerationStatus, error) {
	return m.status, m.statusErr
}

func (m *mockSpecViewRepository) GetGenerationStatusByLanguage(_ context.Context, _ string, _ string) (*entity.SpecGenerationStatus, error) {
	return m.status, m.statusErr
}

func (m *mockSpecViewRepository) GetSpecDocumentByVersion(_ context.Context, _ string, _ string, _ int) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) GetSpecDocumentByUserAndVersion(_ context.Context, _ string, _ string, _ string, _ int) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) GetVersionsByLanguage(_ context.Context, _ string, _ string) ([]entity.VersionInfo, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) GetVersionsByUser(_ context.Context, _ string, _ string, _ string) ([]entity.VersionInfo, error) {
	return nil, nil
}

type mockQueueService struct {
	enqueueErr error
}

func (m *mockQueueService) EnqueueSpecGeneration(_ context.Context, _ string, _ string, _ *string, _ subscriptionentity.PlanTier, _ bool) error {
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

func TestRequestGenerationUseCase_Execute_ForceRegenerate(t *testing.T) {
	tests := []struct {
		name            string
		input           RequestGenerationInput
		mockRepo        *mockSpecViewRepository
		mockQueue       *mockQueueService
		wantErr         bool
		wantErrContains string
	}{
		{
			name: "should enqueue when force regenerate with explicit language (version management)",
			input: RequestGenerationInput{
				AnalysisID:        "test-analysis-id",
				IsForceRegenerate: true,
				Language:          "Korean",
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
			},
			mockQueue: &mockQueueService{},
			wantErr:   false,
		},
		{
			name: "should use default language (English) when force regenerate without language",
			input: RequestGenerationInput{
				AnalysisID:        "test-analysis-id",
				IsForceRegenerate: true,
				Language:          "",
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
			},
			mockQueue: &mockQueueService{},
			wantErr:   false,
		},
		{
			name: "should check document exists when not force regenerate",
			input: RequestGenerationInput{
				AnalysisID:        "test-analysis-id",
				IsForceRegenerate: false,
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
				documentExists: false,
			},
			mockQueue: &mockQueueService{},
			wantErr:   false,
		},
		{
			name: "should return already exists when document exists and not force regenerate",
			input: RequestGenerationInput{
				AnalysisID:        "test-analysis-id",
				IsForceRegenerate: false,
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
				documentExists: true,
			},
			mockQueue:       &mockQueueService{},
			wantErr:         true,
			wantErrContains: domain.ErrAlreadyExists.Error(),
		},
		{
			name: "should reject force regenerate when generation is pending",
			input: RequestGenerationInput{
				AnalysisID:        "test-analysis-id",
				IsForceRegenerate: true,
				Language:          "English",
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
				status:         &entity.SpecGenerationStatus{Status: entity.StatusPending},
			},
			mockQueue:       &mockQueueService{},
			wantErr:         true,
			wantErrContains: domain.ErrGenerationPending.Error(),
		},
		{
			name: "should reject force regenerate when generation is running",
			input: RequestGenerationInput{
				AnalysisID:        "test-analysis-id",
				IsForceRegenerate: true,
				Language:          "Korean",
			},
			mockRepo: &mockSpecViewRepository{
				analysisExists: true,
				status:         &entity.SpecGenerationStatus{Status: entity.StatusRunning},
			},
			mockQueue:       &mockQueueService{},
			wantErr:         true,
			wantErrContains: domain.ErrGenerationRunning.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewRequestGenerationUseCase(tt.mockRepo, tt.mockQueue, nil)

			_, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.wantErrContains != "" && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("expected error containing %q, got %v", tt.wantErrContains, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
