package usecase

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kubrickcode/specvital/apps/web/src/backend/internal/db"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/port"
	subscription "github.com/kubrickcode/specvital/apps/web/src/backend/modules/subscription/domain/entity"
	usageentity "github.com/kubrickcode/specvital/apps/web/src/backend/modules/usage/domain/entity"
	usageport "github.com/kubrickcode/specvital/apps/web/src/backend/modules/usage/domain/port"
	usageusecase "github.com/kubrickcode/specvital/apps/web/src/backend/modules/usage/usecase"
)

type RequestGenerationInput struct {
	AnalysisID string
	Language   string
	Mode       entity.GenerationMode
	Tier       subscription.PlanTier
	UserID     string
}

type RequestGenerationOutput struct {
	AnalysisID string
	Message    *string
	Status     entity.GenerationStatus
}

type RequestGenerationUseCase struct {
	checkQuota      *usageusecase.CheckQuotaUseCase
	dbPool          *pgxpool.Pool
	queue           port.QueueService
	repo            port.SpecViewRepository
	reservationRepo usageport.QuotaReservationRepository
}

func NewRequestGenerationUseCase(
	repo port.SpecViewRepository,
	queue port.QueueService,
	checkQuota *usageusecase.CheckQuotaUseCase,
	dbPool *pgxpool.Pool,
	reservationRepo usageport.QuotaReservationRepository,
) *RequestGenerationUseCase {
	return &RequestGenerationUseCase{
		checkQuota:      checkQuota,
		dbPool:          dbPool,
		queue:           queue,
		repo:            repo,
		reservationRepo: reservationRepo,
	}
}

func (uc *RequestGenerationUseCase) Execute(ctx context.Context, input RequestGenerationInput) (*RequestGenerationOutput, error) {
	if input.UserID == "" {
		return nil, domain.ErrUnauthorized
	}

	if input.AnalysisID == "" {
		return nil, domain.ErrInvalidAnalysisID
	}

	exists, err := uc.repo.CheckAnalysisExists(ctx, input.AnalysisID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrAnalysisNotFound
	}

	language := input.Language
	if language == "" {
		language = "English"
	}

	// Always check for in-progress generation (by language) - regardless of force regenerate
	// Only check the user's own generation status (per-user personalization)
	status, err := uc.repo.GetGenerationStatusByLanguage(ctx, input.UserID, input.AnalysisID, language)
	if err != nil {
		return nil, err
	}
	if status != nil {
		switch status.Status {
		case entity.StatusPending:
			return nil, domain.ErrGenerationPending
		case entity.StatusRunning:
			return nil, domain.ErrGenerationRunning
		}
	}

	// Regeneration: Worker will create a new version (no delete needed)
	// Initial generation: Check if any version already exists
	if !input.Mode.IsRegeneration() {
		docExists, err := uc.repo.CheckSpecDocumentExistsByLanguage(ctx, input.AnalysisID, language)
		if err != nil {
			return nil, err
		}
		if docExists {
			return nil, domain.ErrAlreadyExists
		}
	}

	// Get actual test count for quota calculation
	testCount, err := uc.repo.GetAnalysisTestCount(ctx, input.AnalysisID)
	if err != nil {
		return nil, fmt.Errorf("get analysis test count: %w", err)
	}

	// Quota check: validates user has remaining quota before enqueueing.
	if input.UserID != "" && uc.checkQuota != nil {
		quotaResult, err := uc.checkQuota.Execute(ctx, usageusecase.CheckQuotaInput{
			UserID:    input.UserID,
			EventType: usageentity.EventTypeSpecview,
			Amount:    testCount,
		})
		if err != nil {
			return nil, fmt.Errorf("check quota for user %s: %w", input.UserID, err)
		}
		if !quotaResult.IsAllowed {
			return nil, domain.ErrQuotaExceeded
		}
	}

	var userIDPtr *string
	if input.UserID != "" {
		userIDPtr = &input.UserID
	}

	// Enqueue with reservation if transaction support is available.
	// Reservation prevents race conditions by tracking pending usage.
	if uc.dbPool != nil && uc.reservationRepo != nil {
		if err := uc.enqueueWithReservation(ctx, input.AnalysisID, language, input.UserID, input.Tier, input.Mode, testCount); err != nil {
			return nil, err
		}
	} else {
		// Fallback: enqueue without reservation (for tests or configurations without quota tracking)
		if err := uc.queue.EnqueueSpecGeneration(ctx, input.AnalysisID, language, userIDPtr, input.Tier, input.Mode); err != nil {
			return nil, err
		}
	}

	return &RequestGenerationOutput{
		AnalysisID: input.AnalysisID,
		Status:     entity.StatusPending,
	}, nil
}

// enqueueWithReservation creates a quota reservation and enqueues the job atomically.
// If enqueue fails, the transaction is rolled back and the reservation is not created.
func (uc *RequestGenerationUseCase) enqueueWithReservation(
	ctx context.Context,
	analysisID string,
	language string,
	userID string,
	tier subscription.PlanTier,
	mode entity.GenerationMode,
	testCount int,
) error {
	tx, err := uc.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Enqueue job within transaction - get job ID
	jobID, err := uc.queue.EnqueueSpecGenerationTx(ctx, tx, analysisID, language, &userID, tier, mode)
	if err != nil {
		return err
	}

	// Create reservation with actual test count
	qtx := db.New(tx)
	if err := uc.reservationRepo.CreateReservationTx(ctx, qtx, userID, usageentity.EventTypeSpecview, int32(testCount), jobID); err != nil {
		return fmt.Errorf("create quota reservation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
