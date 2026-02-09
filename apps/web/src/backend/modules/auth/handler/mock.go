package handler

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

func (h *MockHandler) AuthRefresh(_ context.Context, _ api.AuthRefreshRequestObject) (api.AuthRefreshResponseObject, error) {
	return api.AuthRefresh500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("mock: not implemented"),
	}, nil
}

func (h *MockHandler) AuthDevLogin(_ context.Context, _ api.AuthDevLoginRequestObject) (api.AuthDevLoginResponseObject, error) {
	return api.AuthDevLogin403ApplicationProblemPlusJSONResponse{
		ForbiddenApplicationProblemPlusJSONResponse: api.NewForbidden("mock: dev login disabled"),
	}, nil
}
