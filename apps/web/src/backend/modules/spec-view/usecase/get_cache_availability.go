package usecase

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

type GetCacheAvailabilityInput struct {
	AnalysisID string
	UserID     string
}

type GetCacheAvailabilityOutput struct {
	// Languages maps each supported language to its cache availability.
	// True means a previous spec exists for the same codebase in that language.
	Languages map[string]bool
}

type GetCacheAvailabilityUseCase struct {
	repo port.SpecViewRepository
}

func NewGetCacheAvailabilityUseCase(repo port.SpecViewRepository) *GetCacheAvailabilityUseCase {
	return &GetCacheAvailabilityUseCase{repo: repo}
}

func (uc *GetCacheAvailabilityUseCase) Execute(ctx context.Context, input GetCacheAvailabilityInput) (*GetCacheAvailabilityOutput, error) {
	if input.UserID == "" {
		return nil, domain.ErrUnauthorized
	}

	if input.AnalysisID == "" {
		return nil, domain.ErrInvalidAnalysisID
	}

	// Check analysis exists
	exists, err := uc.repo.CheckAnalysisExists(ctx, input.AnalysisID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrAnalysisNotFound
	}

	// Get all languages with previous spec in single batch query
	langsWithCache, err := uc.repo.GetLanguagesWithPreviousSpec(ctx, input.UserID, input.AnalysisID)
	if err != nil {
		return nil, err
	}

	// Build lookup set for O(1) access
	cacheSet := make(map[string]struct{}, len(langsWithCache))
	for _, lang := range langsWithCache {
		cacheSet[lang] = struct{}{}
	}

	// Build result map for all supported languages
	languages := make(map[string]bool, len(entity.SupportedLanguages))
	for lang := range entity.SupportedLanguages {
		_, hasPrev := cacheSet[lang]
		languages[lang] = hasPrev
	}

	return &GetCacheAvailabilityOutput{
		Languages: languages,
	}, nil
}
