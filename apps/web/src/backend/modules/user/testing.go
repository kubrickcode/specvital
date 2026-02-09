package user

import (
	"context"
	"time"

	"github.com/specvital/web/src/backend/internal/api"
)

type MockHandler struct{}

var _ api.AnalysisHistoryHandlers = (*MockHandler)(nil)
var _ api.BookmarkHandlers = (*MockHandler)(nil)
var _ api.UserActiveTasksHandlers = (*MockHandler)(nil)

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (h *MockHandler) AddBookmark(_ context.Context, _ api.AddBookmarkRequestObject) (api.AddBookmarkResponseObject, error) {
	return api.AddBookmark200JSONResponse{IsBookmarked: true, Success: true}, nil
}

func (h *MockHandler) GetUserBookmarks(_ context.Context, _ api.GetUserBookmarksRequestObject) (api.GetUserBookmarksResponseObject, error) {
	return api.GetUserBookmarks200JSONResponse{Data: []api.RepositoryCard{}}, nil
}

func (h *MockHandler) RemoveBookmark(_ context.Context, _ api.RemoveBookmarkRequestObject) (api.RemoveBookmarkResponseObject, error) {
	return api.RemoveBookmark200JSONResponse{IsBookmarked: false, Success: true}, nil
}

func (h *MockHandler) AddUserAnalyzedRepository(_ context.Context, _ api.AddUserAnalyzedRepositoryRequestObject) (api.AddUserAnalyzedRepositoryResponseObject, error) {
	return api.AddUserAnalyzedRepository200JSONResponse{AnalysisID: "mock-analysis-id", UpdatedAt: time.Now()}, nil
}

func (h *MockHandler) GetUserAnalyzedRepositories(_ context.Context, _ api.GetUserAnalyzedRepositoriesRequestObject) (api.GetUserAnalyzedRepositoriesResponseObject, error) {
	return api.GetUserAnalyzedRepositories200JSONResponse{Data: []api.RepositoryCard{}, HasNext: false}, nil
}

func (h *MockHandler) GetUserActiveTasks(_ context.Context, _ api.GetUserActiveTasksRequestObject) (api.GetUserActiveTasksResponseObject, error) {
	return api.GetUserActiveTasks200JSONResponse{Tasks: []api.ActiveTask{}}, nil
}
