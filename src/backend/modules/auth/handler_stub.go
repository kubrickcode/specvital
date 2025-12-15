package auth

import (
	"context"

	"github.com/specvital/web/src/backend/internal/api"
)

// StubHandler is a placeholder for auth endpoints until full implementation.
type StubHandler struct{}

var _ api.AuthHandlers = (*StubHandler)(nil)

func NewStubHandler() *StubHandler {
	return &StubHandler{}
}

func (h *StubHandler) AuthCallback(ctx context.Context, request api.AuthCallbackRequestObject) (api.AuthCallbackResponseObject, error) {
	return api.AuthCallback500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("auth not implemented"),
	}, nil
}

func (h *StubHandler) AuthLogin(ctx context.Context, request api.AuthLoginRequestObject) (api.AuthLoginResponseObject, error) {
	return api.AuthLogin500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("auth not implemented"),
	}, nil
}

func (h *StubHandler) AuthLogout(ctx context.Context, request api.AuthLogoutRequestObject) (api.AuthLogoutResponseObject, error) {
	return api.AuthLogout500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("auth not implemented"),
	}, nil
}

func (h *StubHandler) AuthMe(ctx context.Context, request api.AuthMeRequestObject) (api.AuthMeResponseObject, error) {
	return api.AuthMe401ApplicationProblemPlusJSONResponse{
		UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("auth not implemented"),
	}, nil
}
