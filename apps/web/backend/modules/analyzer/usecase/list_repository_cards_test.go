package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/port"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/usecase"
)

type mockRepository struct {
	bookmarkedIDs      []string
	getPaginatedCalled bool
	paginatedRepos     []port.PaginatedRepository
	paginationParams   port.PaginationParams
	previousAnalysis   *port.PreviousAnalysis
}

func (m *mockRepository) CheckAnalysisExistsByCommitSHA(_ context.Context, _, _, _ string) (bool, error) {
	return false, nil
}

func (m *mockRepository) FindActiveRiverJobByRepo(_ context.Context, _, _, _ string) (*port.RiverJobInfo, error) {
	return nil, nil
}

func (m *mockRepository) GetBookmarkedCodebaseIDs(_ context.Context, _ string) ([]string, error) {
	return m.bookmarkedIDs, nil
}

func (m *mockRepository) GetCodebaseID(_ context.Context, _, _ string) (string, error) {
	return "", nil
}

func (m *mockRepository) GetCompletedAnalysisByCommitSHA(_ context.Context, _, _, _ string) (*port.CompletedAnalysis, error) {
	return nil, nil
}

func (m *mockRepository) GetLatestCompletedAnalysis(_ context.Context, _, _ string) (*port.CompletedAnalysis, error) {
	return nil, nil
}

func (m *mockRepository) GetPaginatedRepositories(_ context.Context, params port.PaginationParams) ([]port.PaginatedRepository, error) {
	m.getPaginatedCalled = true
	m.paginationParams = params
	return m.paginatedRepos, nil
}

func (m *mockRepository) GetPreviousAnalysis(_ context.Context, _, _ string) (*port.PreviousAnalysis, error) {
	return m.previousAnalysis, nil
}

func (m *mockRepository) GetRepositoryStats(_ context.Context, _ string) (*entity.RepositoryStats, error) {
	return nil, nil
}

func (m *mockRepository) GetTestSuitesWithCases(_ context.Context, _ string) ([]port.TestSuiteWithCases, error) {
	return nil, nil
}

func (m *mockRepository) UpdateLastViewed(_ context.Context, _, _ string) error {
	return nil
}

func (m *mockRepository) GetAiSpecSummaries(_ context.Context, _ []string, _ string) (map[string]*entity.AiSpecSummary, error) {
	return make(map[string]*entity.AiSpecSummary), nil
}

func (m *mockRepository) GetAnalysisHistory(_ context.Context, _, _ string) ([]port.AnalysisHistoryItem, error) {
	return nil, nil
}

type mockGitClient struct {
	latestSHA string
	err       error
}

func (m *mockGitClient) GetLatestCommitSHA(_ context.Context, _, _ string) (string, error) {
	return m.latestSHA, m.err
}

func (m *mockGitClient) GetLatestCommitSHAWithToken(_ context.Context, _, _, _ string) (string, error) {
	return m.latestSHA, m.err
}

type mockTokenProvider struct{}

func (m *mockTokenProvider) GetUserGitHubToken(_ context.Context, _ string) (string, error) {
	return "", nil
}

func TestExecutePaginated_FirstPage(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	repo := &mockRepository{
		paginatedRepos: []port.PaginatedRepository{
			{
				AnalysisID:     "analysis-1",
				AnalyzedAt:     now,
				CodebaseID:     "codebase-1",
				CommitSHA:      "abc123",
				IsAnalyzedByMe: true,
				Name:           "repo1",
				Owner:          "owner1",
				TotalTests:     100,
			},
		},
	}

	uc := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, repo, &mockTokenProvider{})
	result, err := uc.ExecutePaginated(context.Background(), usecase.ListRepositoryCardsPaginatedInput{
		Limit:  20,
		SortBy: entity.SortByRecent,
		UserID: "user-1",
		View:   entity.ViewFilterAll,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.getPaginatedCalled {
		t.Error("expected GetPaginatedRepositories to be called")
	}
	if repo.paginationParams.Limit != 21 {
		t.Errorf("expected limit+1 (21), got %d", repo.paginationParams.Limit)
	}
	if len(result.Data) != 1 {
		t.Errorf("expected 1 card, got %d", len(result.Data))
	}
	if result.HasNext {
		t.Error("expected hasNext=false for single result")
	}
	if result.NextCursor != nil {
		t.Error("expected nil nextCursor when hasNext=false")
	}
}

func TestExecutePaginated_HasNextPage(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	repos := make([]port.PaginatedRepository, 21)
	for i := 0; i < 21; i++ {
		repos[i] = port.PaginatedRepository{
			AnalysisID:     "analysis-" + string(rune('a'+i)),
			AnalyzedAt:     now.Add(-time.Duration(i) * time.Hour),
			CodebaseID:     "codebase-" + string(rune('a'+i)),
			CommitSHA:      "sha" + string(rune('a'+i)),
			IsAnalyzedByMe: true,
			Name:           "repo" + string(rune('a'+i)),
			Owner:          "owner",
			TotalTests:     100 - i,
		}
	}

	repo := &mockRepository{paginatedRepos: repos}
	uc := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, repo, &mockTokenProvider{})

	result, err := uc.ExecutePaginated(context.Background(), usecase.ListRepositoryCardsPaginatedInput{
		Limit:  20,
		SortBy: entity.SortByRecent,
		UserID: "user-1",
		View:   entity.ViewFilterAll,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 20 {
		t.Errorf("expected 20 cards (limit), got %d", len(result.Data))
	}
	if !result.HasNext {
		t.Error("expected hasNext=true when more items exist")
	}
	if result.NextCursor == nil {
		t.Error("expected nextCursor when hasNext=true")
	}
}

func TestExecutePaginated_WithCursor(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	cursor := entity.EncodeCursor(entity.RepositoryCursor{
		AnalyzedAt: now,
		ID:         "codebase-1",
		Name:       "repo1",
		SortBy:     entity.SortByRecent,
		TestCount:  100,
	})

	repo := &mockRepository{
		paginatedRepos: []port.PaginatedRepository{
			{
				AnalysisID:     "analysis-2",
				AnalyzedAt:     now.Add(-time.Hour),
				CodebaseID:     "codebase-2",
				CommitSHA:      "def456",
				IsAnalyzedByMe: false,
				Name:           "repo2",
				Owner:          "owner2",
				TotalTests:     50,
			},
		},
	}

	uc := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, repo, &mockTokenProvider{})
	result, err := uc.ExecutePaginated(context.Background(), usecase.ListRepositoryCardsPaginatedInput{
		Cursor: cursor,
		Limit:  20,
		SortBy: entity.SortByRecent,
		UserID: "user-1",
		View:   entity.ViewFilterAll,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.paginationParams.Cursor == nil {
		t.Error("expected cursor to be passed to repository")
	}
	if len(result.Data) != 1 {
		t.Errorf("expected 1 card, got %d", len(result.Data))
	}
}

func TestExecutePaginated_SortByMismatch_RestartsFromBeginning(t *testing.T) {
	t.Parallel()

	cursor := entity.EncodeCursor(entity.RepositoryCursor{
		ID:     "test",
		SortBy: entity.SortByRecent,
	})

	repo := &mockRepository{}
	uc := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, repo, &mockTokenProvider{})

	// sortBy mismatch gracefully restarts pagination (no error, nil cursor)
	_, err := uc.ExecutePaginated(context.Background(), usecase.ListRepositoryCardsPaginatedInput{
		Cursor: cursor,
		SortBy: entity.SortByName, // mismatch
		UserID: "user-1",
	})

	if err != nil {
		t.Errorf("expected no error for sortBy mismatch (graceful restart), got %v", err)
	}

	// Verify cursor was treated as nil (fresh pagination)
	if repo.paginationParams.Cursor != nil {
		t.Error("expected nil cursor passed to repository (fresh start)")
	}
}

func TestExecutePaginated_DefaultValues(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{}
	uc := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, repo, &mockTokenProvider{})

	_, err := uc.ExecutePaginated(context.Background(), usecase.ListRepositoryCardsPaginatedInput{
		UserID: "user-1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.paginationParams.SortBy != entity.SortByRecent {
		t.Errorf("expected default sortBy=recent, got %s", repo.paginationParams.SortBy)
	}
	if repo.paginationParams.SortOrder != entity.SortOrderDesc {
		t.Errorf("expected default sortOrder=desc for recent, got %s", repo.paginationParams.SortOrder)
	}
	if repo.paginationParams.View != entity.ViewFilterAll {
		t.Errorf("expected default view=all, got %s", repo.paginationParams.View)
	}
	if repo.paginationParams.Limit != 21 {
		t.Errorf("expected default limit+1=21, got %d", repo.paginationParams.Limit)
	}
}

func TestExecutePaginated_BookmarkIntegration(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	repo := &mockRepository{
		paginatedRepos: []port.PaginatedRepository{
			{
				AnalysisID:     "analysis-1",
				AnalyzedAt:     now,
				CodebaseID:     "codebase-bookmarked",
				CommitSHA:      "abc123",
				IsAnalyzedByMe: true,
				Name:           "repo1",
				Owner:          "owner1",
				TotalTests:     100,
			},
			{
				AnalysisID:     "analysis-2",
				AnalyzedAt:     now,
				CodebaseID:     "codebase-not-bookmarked",
				CommitSHA:      "def456",
				IsAnalyzedByMe: true,
				Name:           "repo2",
				Owner:          "owner2",
				TotalTests:     50,
			},
		},
		bookmarkedIDs: []string{"codebase-bookmarked"},
	}

	uc := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, repo, &mockTokenProvider{})
	result, err := uc.ExecutePaginated(context.Background(), usecase.ListRepositoryCardsPaginatedInput{
		Limit:  20,
		SortBy: entity.SortByRecent,
		UserID: "user-1",
		View:   entity.ViewFilterAll,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(result.Data))
	}
	if !result.Data[0].IsBookmarked {
		t.Error("expected first card to be bookmarked")
	}
	if result.Data[1].IsBookmarked {
		t.Error("expected second card to not be bookmarked")
	}
}
