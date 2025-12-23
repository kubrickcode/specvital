package mapper

import (
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
)

func ToRepositoryCard(card domain.RepositoryCard) api.RepositoryCard {
	var analysis *api.AnalysisSummary
	if card.LatestAnalysis != nil {
		analysis = &api.AnalysisSummary{
			AnalyzedAt: card.LatestAnalysis.AnalyzedAt,
			Change:     card.LatestAnalysis.Change,
			CommitSHA:  card.LatestAnalysis.CommitSHA,
			TestCount:  card.LatestAnalysis.TestCount,
		}
	}

	return api.RepositoryCard{
		FullName:       card.FullName,
		ID:             card.ID,
		IsBookmarked:   card.IsBookmarked,
		LatestAnalysis: analysis,
		Name:           card.Name,
		Owner:          card.Owner,
		UpdateStatus:   toAPIUpdateStatus(card.UpdateStatus),
	}
}

func ToRepositoryCards(cards []domain.RepositoryCard) []api.RepositoryCard {
	result := make([]api.RepositoryCard, len(cards))
	for i, card := range cards {
		result[i] = ToRepositoryCard(card)
	}
	return result
}

func ToRepositoryStatsResponse(stats *domain.RepositoryStats) api.RepositoryStatsResponse {
	return api.RepositoryStatsResponse{
		TotalRepositories: stats.TotalRepositories,
		TotalTests:        stats.TotalTests,
	}
}

func ToUpdateStatusResponse(result *domain.UpdateStatusResult) api.UpdateStatusResponse {
	resp := api.UpdateStatusResponse{
		Status: toAPIUpdateStatus(result.Status),
	}
	if result.AnalyzedCommitSHA != "" {
		resp.AnalyzedCommitSHA = &result.AnalyzedCommitSHA
	}
	if result.LatestCommitSHA != "" {
		resp.LatestCommitSHA = &result.LatestCommitSHA
	}
	return resp
}

func toAPIUpdateStatus(status domain.UpdateStatus) api.UpdateStatus {
	switch status {
	case domain.UpdateStatusNewCommits:
		return api.NewCommits
	case domain.UpdateStatusUpToDate:
		return api.UpToDate
	default:
		return api.Unknown
	}
}
