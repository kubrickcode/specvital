package usecase

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
	usageentity "github.com/specvital/web/src/backend/modules/usage/domain/entity"
	usageusecase "github.com/specvital/web/src/backend/modules/usage/usecase"
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
	checkQuota *usageusecase.CheckQuotaUseCase
	queue      port.QueueService
	repo       port.SpecViewRepository
}

func NewRequestGenerationUseCase(
	repo port.SpecViewRepository,
	queue port.QueueService,
	checkQuota *usageusecase.CheckQuotaUseCase,
) *RequestGenerationUseCase {
	return &RequestGenerationUseCase{
		checkQuota: checkQuota,
		queue:      queue,
		repo:       repo,
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

	// Quota check: validates user has remaining quota before enqueueing.
	// Note: This check is not atomic with queue enrollment. Concurrent requests may
	// pass validation before usage is recorded, allowing marginal over-usage.
	// This is acceptable as usage recording happens asynchronously in the worker.
	if input.UserID != "" && uc.checkQuota != nil {
		quotaResult, err := uc.checkQuota.Execute(ctx, usageusecase.CheckQuotaInput{
			UserID:    input.UserID,
			EventType: usageentity.EventTypeSpecview,
			Amount:    1,
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

	if err := uc.queue.EnqueueSpecGeneration(ctx, input.AnalysisID, language, userIDPtr, input.Tier, input.Mode); err != nil {
		return nil, err
	}

	return &RequestGenerationOutput{
		AnalysisID: input.AnalysisID,
		Status:     entity.StatusPending,
	}, nil
}
