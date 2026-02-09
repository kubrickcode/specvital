package mapper

import (
	"fmt"

	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/user/domain/entity"
)

func ToUserAnalyzedRepositoriesResponse(result *entity.AnalyzedReposResult) api.UserAnalyzedRepositoriesResponse {
	cards := make([]api.RepositoryCard, 0, len(result.Data))

	for _, repo := range result.Data {
		cards = append(cards, ToAnalyzedRepositoryCard(repo))
	}

	return api.UserAnalyzedRepositoriesResponse{
		Data:       cards,
		HasNext:    result.HasNext,
		NextCursor: result.NextCursor,
	}
}

func ToAnalyzedRepositoryCard(repo *entity.AnalyzedRepository) api.RepositoryCard {
	card := api.RepositoryCard{
		FullName:     fmt.Sprintf("%s/%s", repo.Owner, repo.Name),
		ID:           repo.CodebaseID,
		IsBookmarked: false,
		Name:         repo.Name,
		Owner:        repo.Owner,
		UpdateStatus: api.Unknown,
	}

	if !repo.CompletedAt.IsZero() {
		card.LatestAnalysis = &api.AnalysisSummary{
			AnalyzedAt: repo.CompletedAt,
			Change:     0,
			CommitSHA:  repo.CommitSHA,
			TestCount:  repo.TotalTests,
		}
	}

	return card
}
