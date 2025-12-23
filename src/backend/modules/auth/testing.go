package auth

import (
	"context"

	"github.com/specvital/web/src/backend/internal/api"
)

// MockHandler is a test double for auth endpoints.
type MockHandler struct{}

var _ api.AuthHandlers = (*MockHandler)(nil)

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (h *MockHandler) AuthCallback(_ context.Context, _ api.AuthCallbackRequestObject) (api.AuthCallbackResponseObject, error) {
	return api.AuthCallback500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("mock: not implemented"),
	}, nil
}

func (h *MockHandler) AuthLogin(_ context.Context, _ api.AuthLoginRequestObject) (api.AuthLoginResponseObject, error) {
	return api.AuthLogin500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("mock: not implemented"),
	}, nil
}

func (h *MockHandler) AuthLogout(_ context.Context, _ api.AuthLogoutRequestObject) (api.AuthLogoutResponseObject, error) {
	return api.AuthLogout500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("mock: not implemented"),
	}, nil
}

func (h *MockHandler) AuthMe(_ context.Context, _ api.AuthMeRequestObject) (api.AuthMeResponseObject, error) {
	return api.AuthMe401ApplicationProblemPlusJSONResponse{
		UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("mock: not implemented"),
	}, nil
}

// MockBookmarkHandler is a test double for bookmark endpoints.
type MockBookmarkHandler struct{}

var _ api.BookmarkHandlers = (*MockBookmarkHandler)(nil)

func NewMockBookmarkHandler() *MockBookmarkHandler {
	return &MockBookmarkHandler{}
}

func (h *MockBookmarkHandler) AddBookmark(_ context.Context, _ api.AddBookmarkRequestObject) (api.AddBookmarkResponseObject, error) {
	return api.AddBookmark500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("mock: not implemented"),
	}, nil
}

func (h *MockBookmarkHandler) GetUserBookmarks(_ context.Context, _ api.GetUserBookmarksRequestObject) (api.GetUserBookmarksResponseObject, error) {
	return api.GetUserBookmarks500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("mock: not implemented"),
	}, nil
}

func (h *MockBookmarkHandler) RemoveBookmark(_ context.Context, _ api.RemoveBookmarkRequestObject) (api.RemoveBookmarkResponseObject, error) {
	return api.RemoveBookmark500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("mock: not implemented"),
	}, nil
}
