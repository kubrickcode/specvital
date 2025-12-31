package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type ListRepositoryCardsInput struct {
	Limit  int
	UserID string
}

type ListRepositoryCardsPaginatedInput struct {
	Cursor    string
	Limit     int
	SortBy    entity.SortBy
	SortOrder entity.SortOrder
	UserID    string
	View      entity.ViewFilter
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
	limit := normalizeLimit(input.Limit)

	repos, err := uc.repository.GetRecentRepositories(ctx, input.UserID, limit)
	if err != nil {
		return nil, fmt.Errorf("get recent repositories: %w", err)
	}

	bookmarkedIDs, err := uc.loadBookmarkedIDs(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return uc.buildCardsFromRecent(ctx, repos, bookmarkedIDs), nil
}

func (uc *ListRepositoryCardsUseCase) ExecutePaginated(ctx context.Context, input ListRepositoryCardsPaginatedInput) (entity.PaginatedRepositoryCards, error) {
	limit := normalizeLimit(input.Limit)
	sortBy := normalizeSortBy(input.SortBy)
	sortOrder := normalizeSortOrder(input.SortOrder, sortBy)
	view := normalizeView(input.View)

	cursor, err := entity.DecodeCursor(input.Cursor, sortBy)
	if err != nil {
		return entity.PaginatedRepositoryCards{}, err
	}

	repos, err := uc.repository.GetPaginatedRepositories(ctx, port.PaginationParams{
		Cursor:    cursor,
		Limit:     limit + 1,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		UserID:    input.UserID,
		View:      view,
	})
	if err != nil {
		return entity.PaginatedRepositoryCards{}, fmt.Errorf("get paginated repositories: %w", err)
	}

	hasNext := len(repos) > limit
	if hasNext {
		repos = repos[:limit]
	}

	bookmarkedIDs, err := uc.loadBookmarkedIDs(ctx, input.UserID)
	if err != nil {
		return entity.PaginatedRepositoryCards{}, err
	}

	cards := uc.buildCardsFromPaginated(ctx, repos, bookmarkedIDs)

	var nextCursor *string
	if hasNext && len(repos) > 0 {
		last := repos[len(repos)-1]
		encoded := entity.EncodeCursor(entity.RepositoryCursor{
			AnalyzedAt: last.AnalyzedAt,
			ID:         last.CodebaseID,
			Name:       last.Name,
			SortBy:     sortBy,
			TestCount:  last.TotalTests,
		})
		nextCursor = &encoded
	}

	return entity.PaginatedRepositoryCards{
		Data:       cards,
		HasNext:    hasNext,
		NextCursor: nextCursor,
	}, nil
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

func normalizeSortBy(sortBy entity.SortBy) entity.SortBy {
	if sortBy == "" {
		return entity.SortByRecent
	}
	return sortBy
}

func normalizeSortOrder(sortOrder entity.SortOrder, sortBy entity.SortBy) entity.SortOrder {
	if sortOrder == "" {
		return entity.DefaultSortOrder(sortBy)
	}
	return sortOrder
}

func normalizeView(view entity.ViewFilter) entity.ViewFilter {
	if view == "" {
		return entity.ViewFilterAll
	}
	return view
}

func (uc *ListRepositoryCardsUseCase) loadBookmarkedIDs(ctx context.Context, userID string) (map[string]bool, error) {
	bookmarkedIDs := make(map[string]bool)
	if userID == "" {
		return bookmarkedIDs, nil
	}

	ids, err := uc.repository.GetBookmarkedCodebaseIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get bookmarked codebase IDs: %w", err)
	}
	for _, id := range ids {
		bookmarkedIDs[id] = true
	}
	return bookmarkedIDs, nil
}

type repoData struct {
	AnalysisID     string
	AnalyzedAt     time.Time
	CodebaseID     string
	CommitSHA      string
	IsAnalyzedByMe bool
	Name           string
	Owner          string
	TotalTests     int
}

func (uc *ListRepositoryCardsUseCase) buildCard(ctx context.Context, r repoData, bookmarkedIDs map[string]bool) entity.RepositoryCard {
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

	return entity.RepositoryCard{
		FullName:       fmt.Sprintf("%s/%s", r.Owner, r.Name),
		ID:             r.CodebaseID,
		IsAnalyzedByMe: r.IsAnalyzedByMe,
		IsBookmarked:   bookmarkedIDs[r.CodebaseID],
		LatestAnalysis: analysis,
		Name:           r.Name,
		Owner:          r.Owner,
		UpdateStatus:   entity.UpdateStatusUnknown,
	}
}

func (uc *ListRepositoryCardsUseCase) buildCardsFromRecent(ctx context.Context, repos []port.RecentRepository, bookmarkedIDs map[string]bool) []entity.RepositoryCard {
	cards := make([]entity.RepositoryCard, len(repos))
	for i, r := range repos {
		cards[i] = uc.buildCard(ctx, repoData{
			AnalysisID:     r.AnalysisID,
			AnalyzedAt:     r.AnalyzedAt,
			CodebaseID:     r.CodebaseID,
			CommitSHA:      r.CommitSHA,
			IsAnalyzedByMe: r.IsAnalyzedByMe,
			Name:           r.Name,
			Owner:          r.Owner,
			TotalTests:     r.TotalTests,
		}, bookmarkedIDs)
	}
	return cards
}

func (uc *ListRepositoryCardsUseCase) buildCardsFromPaginated(ctx context.Context, repos []port.PaginatedRepository, bookmarkedIDs map[string]bool) []entity.RepositoryCard {
	cards := make([]entity.RepositoryCard, len(repos))
	for i, r := range repos {
		cards[i] = uc.buildCard(ctx, repoData{
			AnalysisID:     r.AnalysisID,
			AnalyzedAt:     r.AnalyzedAt,
			CodebaseID:     r.CodebaseID,
			CommitSHA:      r.CommitSHA,
			IsAnalyzedByMe: r.IsAnalyzedByMe,
			Name:           r.Name,
			Owner:          r.Owner,
			TotalTests:     r.TotalTests,
		}, bookmarkedIDs)
	}
	return cards
}
