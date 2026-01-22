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
	// Version is optional. If 0, returns the latest version. Requires Language to be set.
	Version int
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

	// Version requires valid language
	if input.Version > 0 {
		if input.Language == "" || !entity.IsValidLanguage(input.Language) {
			return nil, domain.ErrInvalidLanguage
		}
	}

	// Check for active generation first (for regeneration scenarios)
	// If generation is in progress, return generating status even if document exists
	// Skip generation check when requesting specific version (historical data)
	var status *entity.SpecGenerationStatus
	var err error
	if input.Version == 0 {
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
	}

	var doc *entity.SpecDocument
	if input.Version > 0 {
		doc, err = uc.repo.GetSpecDocumentByVersion(ctx, input.AnalysisID, input.Language, input.Version)
	} else {
		doc, err = uc.repo.GetSpecDocumentByLanguage(ctx, input.AnalysisID, input.Language)
	}
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
