package handler

import (
	"context"
	"testing"
	"time"

	"github.com/kubrickcode/specvital/apps/web/backend/common/logger"
	"github.com/kubrickcode/specvital/apps/web/backend/internal/api"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/port"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/usecase"
)

func newTestLogger() *logger.Logger {
	return logger.New()
}

func sortByPtr(s string) *api.SortByParam {
	v := api.SortByParam(s)
	return &v
}

func sortOrderPtr(s string) *api.SortOrderParam {
	v := api.SortOrderParam(s)
	return &v
}

func viewFilterPtr(s string) *api.ViewFilterParam {
	v := api.ViewFilterParam(s)
	return &v
}

type mockRepository struct {
	port.Repository
	paginatedRepos []port.PaginatedRepository
	bookmarkedIDs  []string
}

func (m *mockRepository) CheckAnalysisExistsByCommitSHA(_ context.Context, _, _, _ string) (bool, error) {
	return true, nil
}

func (m *mockRepository) GetPaginatedRepositories(_ context.Context, _ port.PaginationParams) ([]port.PaginatedRepository, error) {
	return m.paginatedRepos, nil
}

func (m *mockRepository) GetBookmarkedCodebaseIDs(_ context.Context, _ string) ([]string, error) {
	return m.bookmarkedIDs, nil
}

func (m *mockRepository) GetPreviousAnalysis(_ context.Context, _ string, _ string) (*port.PreviousAnalysis, error) {
	return nil, nil
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

func TestGetRecentRepositories_DefaultParams(t *testing.T) {
	now := time.Now()
	mock := &mockRepository{
		paginatedRepos: []port.PaginatedRepository{
			{
				AnalysisID:     "a1",
				AnalyzedAt:     now,
				CodebaseID:     "c1",
				CommitSHA:      "sha1",
				IsAnalyzedByMe: true,
				Name:           "repo1",
				Owner:          "owner1",
				TotalTests:     100,
			},
		},
	}

	log := newTestLogger()
	listUC := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, mock, &mockTokenProvider{})
	getHistoryUC := usecase.NewGetAnalysisHistoryUseCase(mock)
	h := NewHandler(log, nil, nil, getHistoryUC, listUC, nil, nil, nil, nil, nil, nil)

	req := api.GetRecentRepositoriesRequestObject{
		Params: api.GetRecentRepositoriesParams{},
	}

	resp, err := h.GetRecentRepositories(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jsonResp, ok := resp.(api.GetRecentRepositories200JSONResponse)
	if !ok {
		t.Fatalf("expected 200 response, got %T", resp)
	}

	if len(jsonResp.Data) != 1 {
		t.Errorf("expected 1 item, got %d", len(jsonResp.Data))
	}
	if jsonResp.HasNext {
		t.Error("expected hasNext=false")
	}
}

func TestGetRecentRepositories_WithPaginationParams(t *testing.T) {
	now := time.Now()
	repos := make([]port.PaginatedRepository, 21)
	for i := range repos {
		repos[i] = port.PaginatedRepository{
			AnalysisID:     "a",
			AnalyzedAt:     now.Add(-time.Duration(i) * time.Hour),
			CodebaseID:     "c",
			Name:           "repo",
			Owner:          "owner",
			IsAnalyzedByMe: true,
		}
	}

	mock := &mockRepository{paginatedRepos: repos}

	log := newTestLogger()
	listUC := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, mock, &mockTokenProvider{})
	getHistoryUC := usecase.NewGetAnalysisHistoryUseCase(mock)
	h := NewHandler(log, nil, nil, getHistoryUC, listUC, nil, nil, nil, nil, nil, nil)

	limit := 20

	req := api.GetRecentRepositoriesRequestObject{
		Params: api.GetRecentRepositoriesParams{
			Limit:     &limit,
			SortBy:    sortByPtr("recent"),
			SortOrder: sortOrderPtr("desc"),
			View:      viewFilterPtr("all"),
		},
	}

	resp, err := h.GetRecentRepositories(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jsonResp, ok := resp.(api.GetRecentRepositories200JSONResponse)
	if !ok {
		t.Fatalf("expected 200 response, got %T", resp)
	}

	if len(jsonResp.Data) != 20 {
		t.Errorf("expected 20 items, got %d", len(jsonResp.Data))
	}
	if !jsonResp.HasNext {
		t.Error("expected hasNext=true")
	}
	if jsonResp.NextCursor == nil {
		t.Error("expected nextCursor to be set")
	}
}

func TestGetRecentRepositories_InvalidCursor(t *testing.T) {
	mock := &mockRepository{}

	log := newTestLogger()
	listUC := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, mock, &mockTokenProvider{})
	getHistoryUC := usecase.NewGetAnalysisHistoryUseCase(mock)
	h := NewHandler(log, nil, nil, getHistoryUC, listUC, nil, nil, nil, nil, nil, nil)

	invalidCursor := "invalid-cursor-data"
	req := api.GetRecentRepositoriesRequestObject{
		Params: api.GetRecentRepositoriesParams{
			Cursor: &invalidCursor,
		},
	}

	resp, err := h.GetRecentRepositories(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.GetRecentRepositories400ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Fatalf("expected 400 response, got %T", resp)
	}
}

func TestGetRecentRepositories_SortByMismatch_RestartsFromBeginning(t *testing.T) {
	mock := &mockRepository{}

	log := newTestLogger()
	listUC := usecase.NewListRepositoryCardsUseCase(&mockGitClient{}, mock, &mockTokenProvider{})
	getHistoryUC := usecase.NewGetAnalysisHistoryUseCase(mock)
	h := NewHandler(log, nil, nil, getHistoryUC, listUC, nil, nil, nil, nil, nil, nil)

	cursor := entity.EncodeCursor(entity.RepositoryCursor{
		ID:         "c1",
		AnalyzedAt: time.Now(),
		SortBy:     entity.SortByName,
	})

	req := api.GetRecentRepositoriesRequestObject{
		Params: api.GetRecentRepositoriesParams{
			Cursor: &cursor,
			SortBy: sortByPtr("recent"),
		},
	}

	resp, err := h.GetRecentRepositories(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// sortBy mismatch gracefully restarts pagination from beginning (returns 200)
	_, ok := resp.(api.GetRecentRepositories200JSONResponse)
	if !ok {
		t.Fatalf("expected 200 response for sortBy mismatch (graceful restart), got %T", resp)
	}
}
