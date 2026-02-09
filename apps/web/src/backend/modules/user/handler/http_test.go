package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	authhandler "github.com/specvital/web/src/backend/modules/auth/handler"
	specviewhandler "github.com/specvital/web/src/backend/modules/spec-view/handler"
	domainentity "github.com/specvital/web/src/backend/modules/user/domain/entity"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
	activetaskuc "github.com/specvital/web/src/backend/modules/user/usecase/activetask"
	bookmarkuc "github.com/specvital/web/src/backend/modules/user/usecase/bookmark"
	historyuc "github.com/specvital/web/src/backend/modules/user/usecase/history"
)

func withTestUserContext(ctx context.Context, userID string) context.Context {
	claims := &entity.Claims{
		ExpiresAt: time.Now().Add(time.Hour),
		IssuedAt:  time.Now(),
		Issuer:    "specvital",
		Login:     "testuser",
		Subject:   userID,
	}
	return middleware.WithClaims(ctx, claims)
}

type mockBookmarkRepository struct {
	bookmarks  []*domainentity.BookmarkedRepository
	codebaseID string
	err        error
}

func (m *mockBookmarkRepository) AddBookmark(_ context.Context, _, _ string) error {
	return m.err
}

func (m *mockBookmarkRepository) GetCodebaseIDByOwnerRepo(_ context.Context, _, _ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.codebaseID, nil
}

func (m *mockBookmarkRepository) GetUserBookmarks(_ context.Context, _ string) ([]*domainentity.BookmarkedRepository, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.bookmarks, nil
}

func (m *mockBookmarkRepository) RemoveBookmark(_ context.Context, _, _ string) error {
	return m.err
}

type mockHistoryRepository struct {
	repos []*domainentity.AnalyzedRepository
	err   error
}

func (m *mockHistoryRepository) AddUserAnalyzedRepository(_ context.Context, _, _, _ string) (*port.AddHistoryResult, error) {
	return nil, m.err
}

func (m *mockHistoryRepository) CheckUserHistoryExists(_ context.Context, _, _, _ string) (bool, error) {
	return false, m.err
}

func (m *mockHistoryRepository) GetUserAnalyzedRepositories(_ context.Context, _ string, _ port.AnalyzedReposParams) ([]*domainentity.AnalyzedRepository, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.repos, nil
}

func setupTestRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	apiHandlers := api.NewAPIHandlers(
		&mockAnalyzerHandler{},
		handler,
		authhandler.NewMockHandler(),
		handler,
		&mockGitHubHandler{},
		&mockGitHubAppHandler{},
		&mockPricingHandler{},
		&mockRepositoryHandler{},
		specviewhandler.NewMockHandler(),
		&mockSubscriptionHandler{},
		&mockUsageHandler{},
		handler, // userActiveTasks
		nil,     // webhook
	)
	strictHandler := api.NewStrictHandler(apiHandlers, nil)
	api.HandlerFromMux(strictHandler, r)
	return r
}

type mockAnalyzerHandler struct{}

func (m *mockAnalyzerHandler) AnalyzeRepository(_ context.Context, _ api.AnalyzeRepositoryRequestObject) (api.AnalyzeRepositoryResponseObject, error) {
	return nil, nil
}

func (m *mockAnalyzerHandler) GetAnalysisHistory(_ context.Context, _ api.GetAnalysisHistoryRequestObject) (api.GetAnalysisHistoryResponseObject, error) {
	return nil, nil
}

func (m *mockAnalyzerHandler) GetAnalysisStatus(_ context.Context, _ api.GetAnalysisStatusRequestObject) (api.GetAnalysisStatusResponseObject, error) {
	return nil, nil
}

type mockRepositoryHandler struct{}

func (m *mockRepositoryHandler) GetRecentRepositories(_ context.Context, _ api.GetRecentRepositoriesRequestObject) (api.GetRecentRepositoriesResponseObject, error) {
	return nil, nil
}

func (m *mockRepositoryHandler) GetRepositoryStats(_ context.Context, _ api.GetRepositoryStatsRequestObject) (api.GetRepositoryStatsResponseObject, error) {
	return nil, nil
}

func (m *mockRepositoryHandler) GetUpdateStatus(_ context.Context, _ api.GetUpdateStatusRequestObject) (api.GetUpdateStatusResponseObject, error) {
	return nil, nil
}

func (m *mockRepositoryHandler) ReanalyzeRepository(_ context.Context, _ api.ReanalyzeRepositoryRequestObject) (api.ReanalyzeRepositoryResponseObject, error) {
	return nil, nil
}

type mockGitHubHandler struct{}

func (m *mockGitHubHandler) GetOrganizationRepositories(_ context.Context, _ api.GetOrganizationRepositoriesRequestObject) (api.GetOrganizationRepositoriesResponseObject, error) {
	return nil, nil
}

func (m *mockGitHubHandler) GetUserGitHubOrganizations(_ context.Context, _ api.GetUserGitHubOrganizationsRequestObject) (api.GetUserGitHubOrganizationsResponseObject, error) {
	return nil, nil
}

func (m *mockGitHubHandler) GetUserGitHubRepositories(_ context.Context, _ api.GetUserGitHubRepositoriesRequestObject) (api.GetUserGitHubRepositoriesResponseObject, error) {
	return nil, nil
}

type mockGitHubAppHandler struct{}

func (m *mockGitHubAppHandler) GetGitHubAppInstallURL(_ context.Context, _ api.GetGitHubAppInstallURLRequestObject) (api.GetGitHubAppInstallURLResponseObject, error) {
	return api.GetGitHubAppInstallURL200JSONResponse{InstallURL: ""}, nil
}

func (m *mockGitHubAppHandler) GetUserGitHubAppInstallations(_ context.Context, _ api.GetUserGitHubAppInstallationsRequestObject) (api.GetUserGitHubAppInstallationsResponseObject, error) {
	return api.GetUserGitHubAppInstallations200JSONResponse{Data: []api.GitHubAppInstallation{}}, nil
}

type mockPricingHandler struct{}

func (m *mockPricingHandler) GetPricing(_ context.Context, _ api.GetPricingRequestObject) (api.GetPricingResponseObject, error) {
	return api.GetPricing200JSONResponse{Data: []api.PricingPlan{}}, nil
}

type mockSubscriptionHandler struct{}

func (m *mockSubscriptionHandler) GetUserSubscription(_ context.Context, _ api.GetUserSubscriptionRequestObject) (api.GetUserSubscriptionResponseObject, error) {
	return api.GetUserSubscription200JSONResponse{}, nil
}

type mockUsageHandler struct{}

func (m *mockUsageHandler) CheckQuota(_ context.Context, _ api.CheckQuotaRequestObject) (api.CheckQuotaResponseObject, error) {
	return api.CheckQuota200JSONResponse{}, nil
}

func (m *mockUsageHandler) GetCurrentUsage(_ context.Context, _ api.GetCurrentUsageRequestObject) (api.GetCurrentUsageResponseObject, error) {
	return api.GetCurrentUsage200JSONResponse{}, nil
}

type mockActiveTaskRepository struct {
	tasks []domainentity.ActiveTask
	err   error
}

func (m *mockActiveTaskRepository) GetUserActiveTasks(_ context.Context, _ string) ([]domainentity.ActiveTask, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tasks, nil
}

func createTestHandler(t *testing.T, bookmarkRepo *mockBookmarkRepository, historyRepo *mockHistoryRepository) *Handler {
	t.Helper()

	activeTaskRepo := &mockActiveTaskRepository{tasks: []domainentity.ActiveTask{}}

	addAnalyzedRepoUC := historyuc.NewAddAnalyzedRepoUseCase(historyRepo)
	addBookmarkUC := bookmarkuc.NewAddBookmarkUseCase(bookmarkRepo)
	getBookmarksUC := bookmarkuc.NewGetBookmarksUseCase(bookmarkRepo)
	removeBookmarkUC := bookmarkuc.NewRemoveBookmarkUseCase(bookmarkRepo)
	getAnalyzedReposUC := historyuc.NewGetAnalyzedReposUseCase(historyRepo)
	getUserActiveTasksUC := activetaskuc.NewGetUserActiveTasksUseCase(activeTaskRepo)

	handler, err := NewHandler(&HandlerConfig{
		AddAnalyzedRepo:    addAnalyzedRepoUC,
		AddBookmark:        addBookmarkUC,
		GetAnalyzedRepos:   getAnalyzedReposUC,
		GetBookmarks:       getBookmarksUC,
		GetUserActiveTasks: getUserActiveTasksUC,
		Logger:             logger.New(),
		RemoveBookmark:     removeBookmarkUC,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}
	return handler
}

func TestGetUserAnalyzedRepositories_Unauthorized(t *testing.T) {
	bookmarkRepo := &mockBookmarkRepository{codebaseID: "test-id"}
	historyRepo := &mockHistoryRepository{repos: []*domainentity.AnalyzedRepository{}}
	handler := createTestHandler(t, bookmarkRepo, historyRepo)
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/user/analyzed-repositories", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestGetUserAnalyzedRepositories_Success(t *testing.T) {
	now := time.Now()
	bookmarkRepo := &mockBookmarkRepository{codebaseID: "test-id"}
	historyRepo := &mockHistoryRepository{
		repos: []*domainentity.AnalyzedRepository{
			{
				CodebaseID:  "codebase-1",
				CommitSHA:   "abc123",
				CompletedAt: now,
				HistoryID:   "history-1",
				Name:        "repo1",
				Owner:       "user1",
				TotalTests:  100,
				UpdatedAt:   now,
			},
		},
	}
	handler := createTestHandler(t, bookmarkRepo, historyRepo)
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/user/analyzed-repositories", nil)
	req = req.WithContext(withTestUserContext(req.Context(), "test-user-id"))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func TestGetUserAnalyzedRepositories_WithOwnershipFilter(t *testing.T) {
	bookmarkRepo := &mockBookmarkRepository{codebaseID: "test-id"}
	historyRepo := &mockHistoryRepository{repos: []*domainentity.AnalyzedRepository{}}
	handler := createTestHandler(t, bookmarkRepo, historyRepo)
	router := setupTestRouter(handler)

	tests := []struct {
		name       string
		queryParam string
	}{
		{"all ownership", "?ownership=all"},
		{"mine ownership", "?ownership=mine"},
		{"organization ownership", "?ownership=organization"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/user/analyzed-repositories"+tt.queryParam, nil)
			req = req.WithContext(withTestUserContext(req.Context(), "test-user-id"))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}
		})
	}
}
