package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/specvital/web/src/backend/internal/db"
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

func (m *mockSpecViewRepository) GetGenerationStatus(_ context.Context, _ string, _ string) (*entity.SpecGenerationStatus, error) {
	return m.status, m.statusErr
}

func (m *mockSpecViewRepository) GetGenerationStatusByLanguage(_ context.Context, _ string, _ string, _ string) (*entity.SpecGenerationStatus, error) {
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

func (m *mockSpecViewRepository) HasPreviousSpecByLanguage(_ context.Context, _ string, _ string, _ string) (bool, error) {
	return false, nil
}

func (m *mockSpecViewRepository) GetLanguagesWithPreviousSpec(_ context.Context, _ string, _ string) ([]string, error) {
	return nil, nil
}

func (m *mockSpecViewRepository) CheckCodebaseExists(_ context.Context, _, _ string) (bool, error) {
	return false, nil
}
func (m *mockSpecViewRepository) GetSpecDocumentByRepository(_ context.Context, _, _, _, _ string) (*entity.RepoSpecDocument, error) {
	return nil, nil
}
func (m *mockSpecViewRepository) GetSpecDocumentByRepositoryAndVersion(_ context.Context, _, _, _, _ string, _ int) (*entity.RepoSpecDocument, error) {
	return nil, nil
}
func (m *mockSpecViewRepository) GetVersionHistoryByRepository(_ context.Context, _, _, _, _ string) ([]entity.RepoVersionInfo, error) {
	return nil, nil
}
func (m *mockSpecViewRepository) GetAvailableLanguagesByRepository(_ context.Context, _, _, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}

type mockQueueService struct {
	enqueueErr error
	jobID      int64
}

func (m *mockQueueService) EnqueueSpecGeneration(_ context.Context, _ string, _ string, _ *string, _ subscriptionentity.PlanTier, _ entity.GenerationMode) error {
	return m.enqueueErr
}

func (m *mockQueueService) EnqueueSpecGenerationTx(_ context.Context, _ pgx.Tx, _ string, _ string, _ *string, _ subscriptionentity.PlanTier, _ entity.GenerationMode) (int64, error) {
	if m.enqueueErr != nil {
		return 0, m.enqueueErr
	}
	if m.jobID == 0 {
		return 1, nil
	}
	return m.jobID, nil
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

type mockReservationRepository struct {
	reservedAmount int64
	err            error
}

func (m *mockReservationRepository) CreateReservation(_ context.Context, _ string, _ usageentity.EventType, _ int32, _ int64) error {
	return m.err
}

func (m *mockReservationRepository) CreateReservationTx(_ context.Context, _ *db.Queries, _ string, _ usageentity.EventType, _ int32, _ int64) error {
	return m.err
}

func (m *mockReservationRepository) GetTotalReservedAmount(_ context.Context, _ string, _ usageentity.EventType) (int64, error) {
	return m.reservedAmount, m.err
}

func (m *mockReservationRepository) DeleteReservationByJobID(_ context.Context, _ int64) error {
	return m.err
}

func TestRequestGenerationUseCase_Execute_QuotaCheck(t *testing.T) {
	now := time.Now()
	periodEnd := now.AddDate(0, 1, 0)

	tests := []struct {
		name                string
		input               RequestGenerationInput
		mockRepo            *mockSpecViewRepository
		mockQueue           *mockQueueService
		mockSubRepo         *mockSubscriptionRepository
		mockUsageRepo       *mockUsageRepository
		mockReservationRepo *mockReservationRepository
		expectedErr         error
		wantErrContains     string
		expectedOutput      bool
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
			mockReservationRepo: &mockReservationRepository{
				reservedAmount: 0,
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
			mockReservationRepo: &mockReservationRepository{
				reservedAmount: 0,
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
			mockReservationRepo: &mockReservationRepository{
				reservedAmount: 500,
			},
			expectedErr:    nil,
			expectedOutput: true,
		},
		{
			name: "should reject when user ID is empty",
			input: RequestGenerationInput{
				AnalysisID: "test-analysis-id",
				UserID:     "",
			},
			mockRepo:            &mockSpecViewRepository{},
			mockQueue:           &mockQueueService{},
			mockSubRepo:         nil,
			mockUsageRepo:       nil,
			mockReservationRepo: nil,
			expectedErr:         domain.ErrUnauthorized,
			expectedOutput:      false,
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
			mockUsageRepo:       &mockUsageRepository{},
			mockReservationRepo: &mockReservationRepository{},
			wantErrContains:     "database connection failed",
			expectedOutput:      false,
		},
		{
			name: "should reject when used + reserved exceeds quota",
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
				monthlyUsage: 4500,
			},
			mockReservationRepo: &mockReservationRepository{
				reservedAmount: 501,
			},
			expectedErr:    domain.ErrQuotaExceeded,
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var checkQuotaUC *usageusecase.CheckQuotaUseCase
			if tt.mockSubRepo != nil && tt.mockUsageRepo != nil && tt.mockReservationRepo != nil {
				checkQuotaUC = usageusecase.NewCheckQuotaUseCase(tt.mockSubRepo, tt.mockUsageRepo, tt.mockReservationRepo)
			}

			// Pass nil for dbPool and reservationRepo to use fallback path without transaction
			uc := NewRequestGenerationUseCase(tt.mockRepo, tt.mockQueue, checkQuotaUC, nil, nil)

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
				AnalysisID: "test-analysis-id",
				Mode:       entity.GenerationModeRegenerateFresh,
				Language:   "Korean",
				UserID:     "test-user-id",
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
				AnalysisID: "test-analysis-id",
				Mode:       entity.GenerationModeRegenerateFresh,
				Language:   "",
				UserID:     "test-user-id",
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
				AnalysisID: "test-analysis-id",
				Mode:       entity.GenerationModeInitial,
				UserID:     "test-user-id",
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
				AnalysisID: "test-analysis-id",
				Mode:       entity.GenerationModeInitial,
				UserID:     "test-user-id",
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
				AnalysisID: "test-analysis-id",
				Mode:       entity.GenerationModeRegenerateFresh,
				Language:   "English",
				UserID:     "test-user-id",
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
				AnalysisID: "test-analysis-id",
				Mode:       entity.GenerationModeRegenerateFresh,
				Language:   "Korean",
				UserID:     "test-user-id",
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
			uc := NewRequestGenerationUseCase(tt.mockRepo, tt.mockQueue, nil, nil, nil)

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
