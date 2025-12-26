package usecase

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type ListRepositoryCardsInput struct {
	Limit int
}

type ListRepositoryCardsUseCase struct {
	repository port.Repository
}

func NewListRepositoryCardsUseCase(repository port.Repository) *ListRepositoryCardsUseCase {
	return &ListRepositoryCardsUseCase{
		repository: repository,
	}
}

func (uc *ListRepositoryCardsUseCase) Execute(ctx context.Context, input ListRepositoryCardsInput) ([]entity.RepositoryCard, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = defaultLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}

	repos, err := uc.repository.GetRecentRepositories(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("get recent repositories: %w", err)
	}

	cards := make([]entity.RepositoryCard, len(repos))
	for i, r := range repos {
		var analysis *entity.AnalysisSummary
		if r.AnalysisID != "" {
			change := 0
			prevAnalysis, err := uc.repository.GetPreviousAnalysis(ctx, r.CodebaseID, r.AnalysisID)
			if err == nil && prevAnalysis != nil {
				change = r.TotalTests - prevAnalysis.TotalTests
			}

			analysis = &entity.AnalysisSummary{
				AnalyzedAt: r.AnalyzedAt,
				Change:     change,
				CommitSHA:  r.CommitSHA,
				TestCount:  r.TotalTests,
			}
		}

		cards[i] = entity.RepositoryCard{
			FullName:       fmt.Sprintf("%s/%s", r.Owner, r.Name),
			ID:             r.CodebaseID,
			IsBookmarked:   false,
			LatestAnalysis: analysis,
			Name:           r.Name,
			Owner:          r.Owner,
			UpdateStatus:   entity.UpdateStatusUnknown,
		}
	}

	return cards, nil
}
