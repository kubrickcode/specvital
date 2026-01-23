package usecase

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type ListRepositoryCardsPaginatedInput struct {
	Cursor    string
	Limit     int
	Ownership entity.OwnershipFilter
	SortBy    entity.SortBy
	SortOrder entity.SortOrder
	UserID    string
	View      entity.ViewFilter
}

type ListRepositoryCardsUseCase struct {
	gitClient     port.GitClient
	repository    port.Repository
	tokenProvider port.TokenProvider
}

func NewListRepositoryCardsUseCase(
	gitClient port.GitClient,
	repository port.Repository,
	tokenProvider port.TokenProvider,
) *ListRepositoryCardsUseCase {
	return &ListRepositoryCardsUseCase{
		gitClient:     gitClient,
		repository:    repository,
		tokenProvider: tokenProvider,
	}
}

func (uc *ListRepositoryCardsUseCase) ExecutePaginated(ctx context.Context, input ListRepositoryCardsPaginatedInput) (entity.PaginatedRepositoryCards, error) {
	limit := normalizeLimit(input.Limit)
	sortBy := normalizeSortBy(input.SortBy)
	sortOrder := normalizeSortOrder(input.SortOrder, sortBy)
	view := normalizeView(input.View)
	ownership := normalizeOwnership(input.Ownership, input.UserID)

	cursor, err := entity.DecodeCursor(input.Cursor, sortBy)
	if err != nil {
		return entity.PaginatedRepositoryCards{}, err
	}

	repos, err := uc.repository.GetPaginatedRepositories(ctx, port.PaginationParams{
		Cursor:    cursor,
		Limit:     limit + 1,
		Ownership: ownership,
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

	cards := uc.buildCardsFromPaginated(ctx, repos, bookmarkedIDs, input.UserID)

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

func normalizeOwnership(ownership entity.OwnershipFilter, userID string) entity.OwnershipFilter {
	if ownership == "" {
		return entity.OwnershipFilterAll
	}
	if userID == "" && ownership != entity.OwnershipFilterAll {
		return entity.OwnershipFilterAll
	}
	return ownership
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
	ActiveCount    int
	AiSpecSummary  *entity.AiSpecSummary
	AnalysisID     string
	AnalyzedAt     time.Time
	CodebaseID     string
	CommitSHA      string
	FocusedCount   int
	IsAnalyzedByMe bool
	Name           string
	Owner          string
	SkippedCount   int
	TodoCount      int
	TotalTests     int
	XfailCount     int
}

func (uc *ListRepositoryCardsUseCase) buildCard(ctx context.Context, r repoData, bookmarkedIDs map[string]bool, status entity.UpdateStatus) entity.RepositoryCard {
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
			TestSummary: &entity.TestStatusSummary{
				Active:  r.ActiveCount,
				Focused: r.FocusedCount,
				Skipped: r.SkippedCount,
				Todo:    r.TodoCount,
				Xfail:   r.XfailCount,
			},
		}
	}

	return entity.RepositoryCard{
		AiSpecSummary:  r.AiSpecSummary,
		FullName:       fmt.Sprintf("%s/%s", r.Owner, r.Name),
		ID:             r.CodebaseID,
		IsAnalyzedByMe: r.IsAnalyzedByMe,
		IsBookmarked:   bookmarkedIDs[r.CodebaseID],
		LatestAnalysis: analysis,
		Name:           r.Name,
		Owner:          r.Owner,
		UpdateStatus:   status,
	}
}

func (uc *ListRepositoryCardsUseCase) buildCardsFromPaginated(ctx context.Context, repos []port.PaginatedRepository, bookmarkedIDs map[string]bool, userID string) []entity.RepositoryCard {
	cards := make([]entity.RepositoryCard, len(repos))
	statuses := make([]entity.UpdateStatus, len(repos))
	var aiSpecSummaries map[string]*entity.AiSpecSummary

	g, gCtx := errgroup.WithContext(ctx)

	for i, r := range repos {
		g.Go(func() error {
			statuses[i] = uc.getUpdateStatus(gCtx, r.Owner, r.Name, r.CommitSHA, userID)
			return nil
		})
	}

	g.Go(func() error {
		analysisIDs := make([]string, 0, len(repos))
		for _, r := range repos {
			if r.AnalysisID != "" {
				analysisIDs = append(analysisIDs, r.AnalysisID)
			}
		}
		summaries, err := uc.repository.GetAiSpecSummaries(gCtx, analysisIDs, userID)
		if err != nil {
			// Intentional graceful degradation: AI Spec badge is non-critical
			// Dashboard should load even if AI Spec query fails
			return nil
		}
		aiSpecSummaries = summaries
		return nil
	})

	_ = g.Wait()

	for i, r := range repos {
		var aiSpec *entity.AiSpecSummary
		if aiSpecSummaries != nil {
			aiSpec = aiSpecSummaries[r.AnalysisID]
		}
		cards[i] = uc.buildCard(ctx, repoData{
			ActiveCount:    r.ActiveCount,
			AiSpecSummary:  aiSpec,
			AnalysisID:     r.AnalysisID,
			AnalyzedAt:     r.AnalyzedAt,
			CodebaseID:     r.CodebaseID,
			CommitSHA:      r.CommitSHA,
			FocusedCount:   r.FocusedCount,
			IsAnalyzedByMe: r.IsAnalyzedByMe,
			Name:           r.Name,
			Owner:          r.Owner,
			SkippedCount:   r.SkippedCount,
			TodoCount:      r.TodoCount,
			TotalTests:     r.TotalTests,
			XfailCount:     r.XfailCount,
		}, bookmarkedIDs, statuses[i])
	}
	return cards
}

func (uc *ListRepositoryCardsUseCase) getUpdateStatus(ctx context.Context, owner, repo, analyzedSHA, userID string) entity.UpdateStatus {
	if analyzedSHA == "" {
		return entity.UpdateStatusUnknown
	}

	latestSHA, err := getLatestCommitWithAuth(ctx, uc.gitClient, uc.tokenProvider, owner, repo, userID)
	if err != nil {
		return entity.UpdateStatusUnknown
	}

	if latestSHA == analyzedSHA {
		return entity.UpdateStatusUpToDate
	}
	return entity.UpdateStatusNewCommits
}
