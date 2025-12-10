package mocks

import (
	"context"

	"github.com/specvital/collector/internal/service"
)

type MockAnalysisService struct {
	AnalyzeFunc func(ctx context.Context, req service.AnalyzeRequest) error
}

func (m *MockAnalysisService) Analyze(ctx context.Context, req service.AnalyzeRequest) error {
	if m.AnalyzeFunc != nil {
		return m.AnalyzeFunc(ctx, req)
	}
	return nil
}
