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

func (m *MockHandler) GetSpecCacheAvailability(_ context.Context, _ api.GetSpecCacheAvailabilityRequestObject) (api.GetSpecCacheAvailabilityResponseObject, error) {
	return api.GetSpecCacheAvailability200JSONResponse{
		Languages: map[string]bool{},
	}, nil
}

func (m *MockHandler) GetSpecCachePrediction(_ context.Context, _ api.GetSpecCachePredictionRequestObject) (api.GetSpecCachePredictionResponseObject, error) {
	return api.GetSpecCachePrediction200JSONResponse{
		TotalBehaviors:     0,
		CacheableBehaviors: 0,
		NewBehaviors:       0,
		EstimatedCost:      0,
	}, nil
}

func (m *MockHandler) GetSpecDocumentByRepository(_ context.Context, _ api.GetSpecDocumentByRepositoryRequestObject) (api.GetSpecDocumentByRepositoryResponseObject, error) {
	return api.GetSpecDocumentByRepository404ApplicationProblemPlusJSONResponse{
		NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("not found"),
	}, nil
}

func (m *MockHandler) GetVersionHistoryByRepository(_ context.Context, _ api.GetVersionHistoryByRepositoryRequestObject) (api.GetVersionHistoryByRepositoryResponseObject, error) {
	return api.GetVersionHistoryByRepository404ApplicationProblemPlusJSONResponse{
		NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("not found"),
	}, nil
}
