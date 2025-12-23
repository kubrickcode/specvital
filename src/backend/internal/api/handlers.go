package api

import "context"

type AnalyzerHandlers interface {
	AnalyzeRepository(ctx context.Context, request AnalyzeRepositoryRequestObject) (AnalyzeRepositoryResponseObject, error)
	GetAnalysisStatus(ctx context.Context, request GetAnalysisStatusRequestObject) (GetAnalysisStatusResponseObject, error)
}

type AuthHandlers interface {
	AuthCallback(ctx context.Context, request AuthCallbackRequestObject) (AuthCallbackResponseObject, error)
	AuthLogin(ctx context.Context, request AuthLoginRequestObject) (AuthLoginResponseObject, error)
	AuthLogout(ctx context.Context, request AuthLogoutRequestObject) (AuthLogoutResponseObject, error)
	AuthMe(ctx context.Context, request AuthMeRequestObject) (AuthMeResponseObject, error)
}

type BookmarkHandlers interface {
	AddBookmark(ctx context.Context, request AddBookmarkRequestObject) (AddBookmarkResponseObject, error)
	GetUserBookmarks(ctx context.Context, request GetUserBookmarksRequestObject) (GetUserBookmarksResponseObject, error)
	RemoveBookmark(ctx context.Context, request RemoveBookmarkRequestObject) (RemoveBookmarkResponseObject, error)
}

type RepositoryHandlers interface {
	GetRecentRepositories(ctx context.Context, request GetRecentRepositoriesRequestObject) (GetRecentRepositoriesResponseObject, error)
	GetRepositoryStats(ctx context.Context, request GetRepositoryStatsRequestObject) (GetRepositoryStatsResponseObject, error)
	GetUpdateStatus(ctx context.Context, request GetUpdateStatusRequestObject) (GetUpdateStatusResponseObject, error)
	ReanalyzeRepository(ctx context.Context, request ReanalyzeRepositoryRequestObject) (ReanalyzeRepositoryResponseObject, error)
}

type APIHandlers struct {
	analyzer   AnalyzerHandlers
	auth       AuthHandlers
	bookmark   BookmarkHandlers
	repository RepositoryHandlers
}

var _ StrictServerInterface = (*APIHandlers)(nil)

func NewAPIHandlers(analyzer AnalyzerHandlers, auth AuthHandlers, bookmark BookmarkHandlers, repository RepositoryHandlers) *APIHandlers {
	return &APIHandlers{
		analyzer:   analyzer,
		auth:       auth,
		bookmark:   bookmark,
		repository: repository,
	}
}

func (h *APIHandlers) AnalyzeRepository(ctx context.Context, request AnalyzeRepositoryRequestObject) (AnalyzeRepositoryResponseObject, error) {
	return h.analyzer.AnalyzeRepository(ctx, request)
}

func (h *APIHandlers) GetAnalysisStatus(ctx context.Context, request GetAnalysisStatusRequestObject) (GetAnalysisStatusResponseObject, error) {
	return h.analyzer.GetAnalysisStatus(ctx, request)
}

func (h *APIHandlers) AuthCallback(ctx context.Context, request AuthCallbackRequestObject) (AuthCallbackResponseObject, error) {
	return h.auth.AuthCallback(ctx, request)
}

func (h *APIHandlers) AuthLogin(ctx context.Context, request AuthLoginRequestObject) (AuthLoginResponseObject, error) {
	return h.auth.AuthLogin(ctx, request)
}

func (h *APIHandlers) AuthLogout(ctx context.Context, request AuthLogoutRequestObject) (AuthLogoutResponseObject, error) {
	return h.auth.AuthLogout(ctx, request)
}

func (h *APIHandlers) AuthMe(ctx context.Context, request AuthMeRequestObject) (AuthMeResponseObject, error) {
	return h.auth.AuthMe(ctx, request)
}

func (h *APIHandlers) AddBookmark(ctx context.Context, request AddBookmarkRequestObject) (AddBookmarkResponseObject, error) {
	return h.bookmark.AddBookmark(ctx, request)
}

func (h *APIHandlers) GetRecentRepositories(ctx context.Context, request GetRecentRepositoriesRequestObject) (GetRecentRepositoriesResponseObject, error) {
	return h.repository.GetRecentRepositories(ctx, request)
}

func (h *APIHandlers) GetRepositoryStats(ctx context.Context, request GetRepositoryStatsRequestObject) (GetRepositoryStatsResponseObject, error) {
	return h.repository.GetRepositoryStats(ctx, request)
}

func (h *APIHandlers) GetUpdateStatus(ctx context.Context, request GetUpdateStatusRequestObject) (GetUpdateStatusResponseObject, error) {
	return h.repository.GetUpdateStatus(ctx, request)
}

func (h *APIHandlers) ReanalyzeRepository(ctx context.Context, request ReanalyzeRepositoryRequestObject) (ReanalyzeRepositoryResponseObject, error) {
	return h.repository.ReanalyzeRepository(ctx, request)
}

func (h *APIHandlers) GetUserBookmarks(ctx context.Context, request GetUserBookmarksRequestObject) (GetUserBookmarksResponseObject, error) {
	return h.bookmark.GetUserBookmarks(ctx, request)
}

func (h *APIHandlers) RemoveBookmark(ctx context.Context, request RemoveBookmarkRequestObject) (RemoveBookmarkResponseObject, error) {
	return h.bookmark.RemoveBookmark(ctx, request)
}
