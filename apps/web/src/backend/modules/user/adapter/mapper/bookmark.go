package mapper

import (
	"fmt"

	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/user/domain/entity"
)

func ToBookmarkedRepositoriesResponse(repos []*entity.BookmarkedRepository) api.BookmarkedRepositoriesResponse {
	cards := make([]api.RepositoryCard, 0, len(repos))

	for _, repo := range repos {
		cards = append(cards, ToRepositoryCard(repo))
	}

	return api.BookmarkedRepositoriesResponse{
		Data: cards,
	}
}

func ToRepositoryCard(repo *entity.BookmarkedRepository) api.RepositoryCard {
	card := api.RepositoryCard{
		FullName:     fmt.Sprintf("%s/%s", repo.Owner, repo.Name),
		ID:           repo.CodebaseID,
		IsBookmarked: true,
		Name:         repo.Name,
		Owner:        repo.Owner,
		UpdateStatus: api.Unknown,
	}

	if repo.AnalyzedAt != nil {
		card.LatestAnalysis = &api.AnalysisSummary{
			AnalyzedAt: *repo.AnalyzedAt,
			Change:     repo.Change,
			CommitSHA:  repo.CommitSHA,
			TestCount:  repo.TotalTests,
		}
	}

	return card
}

func ToBookmarkResponse(isBookmarked bool) api.BookmarkResponse {
	return api.BookmarkResponse{
		IsBookmarked: isBookmarked,
		Success:      true,
	}
}
