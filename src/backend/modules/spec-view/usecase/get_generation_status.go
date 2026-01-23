package usecase

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

type GetGenerationStatusInput struct {
	AnalysisID string
	Language   string // Optional: if specified, returns status for that specific language
	UserID     string // Required: user must be authenticated
}

type GetGenerationStatusOutput struct {
	Status *entity.SpecGenerationStatus
}

type GetGenerationStatusUseCase struct {
	repo port.SpecViewRepository
}

func NewGetGenerationStatusUseCase(repo port.SpecViewRepository) *GetGenerationStatusUseCase {
	return &GetGenerationStatusUseCase{repo: repo}
}

func (uc *GetGenerationStatusUseCase) Execute(ctx context.Context, input GetGenerationStatusInput) (*GetGenerationStatusOutput, error) {
	if input.UserID == "" {
		return nil, domain.ErrUnauthorized
	}

	if input.AnalysisID == "" {
		return nil, domain.ErrInvalidAnalysisID
	}

	if input.Language != "" && !entity.IsValidLanguage(input.Language) {
		return nil, domain.ErrInvalidLanguage
	}

	// Check ownership: if document exists and belongs to different user, deny access.
	// ownership == nil means no document exists yet, which is allowed to support
	// viewing generation-in-progress status for own requests.
	ownership, err := uc.repo.CheckSpecDocumentOwnership(ctx, input.AnalysisID)
	if err != nil {
		return nil, err
	}
	if ownership != nil && ownership.UserID != input.UserID {
		return nil, domain.ErrForbidden
	}

	var status *entity.SpecGenerationStatus

	if input.Language != "" {
		status, err = uc.repo.GetGenerationStatusByLanguage(ctx, input.AnalysisID, input.Language)
	} else {
		status, err = uc.repo.GetGenerationStatus(ctx, input.AnalysisID)
	}
	if err != nil {
		return nil, err
	}

	if status == nil {
		return &GetGenerationStatusOutput{
			Status: &entity.SpecGenerationStatus{
				AnalysisID: input.AnalysisID,
				Status:     entity.StatusNotFound,
			},
		}, nil
	}

	return &GetGenerationStatusOutput{Status: status}, nil
}
