package handler

import (
	"context"

	"github.com/specvital/web/src/backend/internal/api"
)

type MockHandler struct{}

var _ api.UsageHandlers = (*MockHandler)(nil)

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (m *MockHandler) CheckQuota(_ context.Context, _ api.CheckQuotaRequestObject) (api.CheckQuotaResponseObject, error) {
	return api.CheckQuota200JSONResponse{}, nil
}

func (m *MockHandler) GetCurrentUsage(_ context.Context, _ api.GetCurrentUsageRequestObject) (api.GetCurrentUsageResponseObject, error) {
	return api.GetCurrentUsage200JSONResponse{}, nil
}
