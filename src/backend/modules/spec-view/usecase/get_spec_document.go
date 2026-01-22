package usecase

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

type GetSpecDocumentInput struct {
	AnalysisID string
	// Language is optional. If empty, returns the most recent document.
	Language string
}

type GetSpecDocumentOutput struct {
	Document         *entity.SpecDocument
	GenerationStatus *entity.SpecGenerationStatus
}

type GetSpecDocumentUseCase struct {
	repo port.SpecViewRepository
}

func NewGetSpecDocumentUseCase(repo port.SpecViewRepository) *GetSpecDocumentUseCase {
	return &GetSpecDocumentUseCase{repo: repo}
}

func (uc *GetSpecDocumentUseCase) Execute(ctx context.Context, input GetSpecDocumentInput) (*GetSpecDocumentOutput, error) {
	if input.AnalysisID == "" {
		return nil, domain.ErrInvalidAnalysisID
	}

	// Check for active generation first (for regeneration scenarios)
	// If generation is in progress, return generating status even if document exists
	var status *entity.SpecGenerationStatus
	var err error
	if input.Language != "" {
		status, err = uc.repo.GetGenerationStatusByLanguage(ctx, input.AnalysisID, input.Language)
	} else {
		status, err = uc.repo.GetGenerationStatus(ctx, input.AnalysisID)
	}
	if err != nil {
		return nil, err
	}

	if status != nil && (status.Status == entity.StatusPending || status.Status == entity.StatusRunning) {
		return &GetSpecDocumentOutput{GenerationStatus: status}, nil
	}

	doc, err := uc.repo.GetSpecDocumentByLanguage(ctx, input.AnalysisID, input.Language)
	if err != nil {
		return nil, err
	}

	if doc != nil {
		availableLanguages, err := uc.repo.GetAvailableLanguages(ctx, input.AnalysisID)
		if err != nil {
			return nil, err
		}
		doc.AvailableLanguages = availableLanguages
		return &GetSpecDocumentOutput{Document: doc}, nil
	}

	if status != nil && status.Status != entity.StatusNotFound {
		return &GetSpecDocumentOutput{GenerationStatus: status}, nil
	}

	return nil, domain.ErrDocumentNotFound
}
