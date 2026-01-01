package mapper

import (
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
)

func ToRepositoryCard(card entity.RepositoryCard) api.RepositoryCard {
	var analysis *api.AnalysisSummary
	if card.LatestAnalysis != nil {
		var testSummary *api.TestStatusSummary
		if card.LatestAnalysis.TestSummary != nil {
			testSummary = &api.TestStatusSummary{
				Active:  card.LatestAnalysis.TestSummary.Active,
				Focused: card.LatestAnalysis.TestSummary.Focused,
				Skipped: card.LatestAnalysis.TestSummary.Skipped,
				Todo:    card.LatestAnalysis.TestSummary.Todo,
				Xfail:   card.LatestAnalysis.TestSummary.Xfail,
			}
		}
		analysis = &api.AnalysisSummary{
			AnalyzedAt:  card.LatestAnalysis.AnalyzedAt,
			Change:      card.LatestAnalysis.Change,
			CommitSHA:   card.LatestAnalysis.CommitSHA,
			TestCount:   card.LatestAnalysis.TestCount,
			TestSummary: testSummary,
		}
	}

	return api.RepositoryCard{
		FullName:       card.FullName,
		ID:             card.ID,
		IsAnalyzedByMe: card.IsAnalyzedByMe,
		IsBookmarked:   card.IsBookmarked,
		LatestAnalysis: analysis,
		Name:           card.Name,
		Owner:          card.Owner,
		UpdateStatus:   toAPIUpdateStatus(card.UpdateStatus),
	}
}

func ToRepositoryCards(cards []entity.RepositoryCard) []api.RepositoryCard {
	result := make([]api.RepositoryCard, len(cards))
	for i, card := range cards {
		result[i] = ToRepositoryCard(card)
	}
	return result
}

func ToRepositoryStatsResponse(stats *entity.RepositoryStats) api.RepositoryStatsResponse {
	if stats == nil {
		return api.RepositoryStatsResponse{}
	}
	return api.RepositoryStatsResponse{
		TotalRepositories: stats.TotalRepositories,
		TotalTests:        stats.TotalTests,
	}
}

func ToUpdateStatusResponse(result *entity.UpdateStatusResult) api.UpdateStatusResponse {
	if result == nil {
		return api.UpdateStatusResponse{Status: api.Unknown}
	}
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

func ToPaginatedRepositoriesResponse(result entity.PaginatedRepositoryCards) api.PaginatedRepositoriesResponse {
	return api.PaginatedRepositoriesResponse{
		Data:       ToRepositoryCards(result.Data),
		HasNext:    result.HasNext,
		NextCursor: result.NextCursor,
	}
}

func toAPIUpdateStatus(status entity.UpdateStatus) api.UpdateStatus {
	switch status {
	case entity.UpdateStatusNewCommits:
		return api.NewCommits
	case entity.UpdateStatusUpToDate:
		return api.UpToDate
	default:
		return api.Unknown
	}
}
