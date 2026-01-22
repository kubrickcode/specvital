package handler

import (
	"context"

	"github.com/specvital/web/src/backend/internal/api"
)

// MockHandler is a mock implementation for testing.
type MockHandler struct{}

var _ api.SpecViewHandlers = (*MockHandler)(nil)

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (m *MockHandler) GetSpecDocument(_ context.Context, _ api.GetSpecDocumentRequestObject) (api.GetSpecDocumentResponseObject, error) {
	return api.GetSpecDocument404ApplicationProblemPlusJSONResponse{
		NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("not found"),
	}, nil
}

func (m *MockHandler) GetSpecGenerationStatus(_ context.Context, _ api.GetSpecGenerationStatusRequestObject) (api.GetSpecGenerationStatusResponseObject, error) {
	return api.GetSpecGenerationStatus404ApplicationProblemPlusJSONResponse{
		NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("not found"),
	}, nil
}

func (m *MockHandler) RequestSpecGeneration(_ context.Context, _ api.RequestSpecGenerationRequestObject) (api.RequestSpecGenerationResponseObject, error) {
	return api.RequestSpecGeneration401ApplicationProblemPlusJSONResponse{
		UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("unauthorized"),
	}, nil
}

func (m *MockHandler) GetSpecVersions(_ context.Context, _ api.GetSpecVersionsRequestObject) (api.GetSpecVersionsResponseObject, error) {
	return api.GetSpecVersions404ApplicationProblemPlusJSONResponse{
		NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("not found"),
	}, nil
}
