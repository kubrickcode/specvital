package usecase

import (
	"context"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/spec-view/domain"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/spec-view/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/spec-view/domain/port"
)

type GetCachePredictionInput struct {
	AnalysisID string
	Language   string
	UserID     string
}

type GetCachePredictionOutput struct {
	// TotalBehaviors is the total number of behaviors from the previous spec
	TotalBehaviors int
	// CacheableBehaviors is the number of behaviors that can be cached
	// (matching test case name + file path exists in current analysis)
	CacheableBehaviors int
	// NewBehaviors is the number of new behaviors that need to be generated
	// (tests in current analysis without cached behavior)
	NewBehaviors int
	// EstimatedCost is the estimated quota usage (equals NewBehaviors)
	EstimatedCost int
}

type GetCachePredictionUseCase struct {
	repo port.SpecViewRepository
}

func NewGetCachePredictionUseCase(repo port.SpecViewRepository) *GetCachePredictionUseCase {
	return &GetCachePredictionUseCase{repo: repo}
}

func (uc *GetCachePredictionUseCase) Execute(ctx context.Context, input GetCachePredictionInput) (*GetCachePredictionOutput, error) {
	if input.UserID == "" {
		return nil, domain.ErrUnauthorized
	}

	if input.AnalysisID == "" {
		return nil, domain.ErrInvalidAnalysisID
	}

	if input.Language == "" || !entity.IsValidLanguage(input.Language) {
		return nil, domain.ErrInvalidLanguage
	}

	exists, err := uc.repo.CheckAnalysisExists(ctx, input.AnalysisID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrAnalysisNotFound
	}

	currentTestCount, err := uc.repo.GetCurrentAnalysisTestCount(ctx, input.AnalysisID)
	if err != nil {
		return nil, err
	}

	predictionData, err := uc.repo.GetCachePredictionData(ctx, input.UserID, input.AnalysisID, input.Language)
	if err != nil {
		return nil, err
	}

	// Calculate new behaviors:
	// If there's no previous spec, all tests need new behavior generation
	// Otherwise, new behaviors = current tests - cacheable behaviors
	var newBehaviors int
	if predictionData.TotalBehaviors == 0 {
		newBehaviors = currentTestCount
	} else {
		newBehaviors = currentTestCount - predictionData.CacheableBehaviors
		if newBehaviors < 0 {
			newBehaviors = 0
		}
	}

	return &GetCachePredictionOutput{
		TotalBehaviors:     predictionData.TotalBehaviors,
		CacheableBehaviors: predictionData.CacheableBehaviors,
		NewBehaviors:       newBehaviors,
		EstimatedCost:      newBehaviors,
	}, nil
}
