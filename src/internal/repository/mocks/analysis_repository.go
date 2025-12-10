package mocks

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/specvital/collector/internal/repository"
)

type MockAnalysisRepository struct {
	CreateAnalysisRecordFunc  func(ctx context.Context, params repository.CreateAnalysisRecordParams) (pgtype.UUID, error)
	RecordFailureFunc         func(ctx context.Context, analysisID pgtype.UUID, errMessage string) error
	SaveAnalysisInventoryFunc func(ctx context.Context, params repository.SaveAnalysisInventoryParams) error
	SaveAnalysisResultFunc    func(ctx context.Context, params repository.SaveAnalysisResultParams) error
}

func (m *MockAnalysisRepository) CreateAnalysisRecord(ctx context.Context, params repository.CreateAnalysisRecordParams) (pgtype.UUID, error) {
	if m.CreateAnalysisRecordFunc != nil {
		return m.CreateAnalysisRecordFunc(ctx, params)
	}
	return pgtype.UUID{Valid: true}, nil
}

func (m *MockAnalysisRepository) RecordFailure(ctx context.Context, analysisID pgtype.UUID, errMessage string) error {
	if m.RecordFailureFunc != nil {
		return m.RecordFailureFunc(ctx, analysisID, errMessage)
	}
	return nil
}

func (m *MockAnalysisRepository) SaveAnalysisInventory(ctx context.Context, params repository.SaveAnalysisInventoryParams) error {
	if m.SaveAnalysisInventoryFunc != nil {
		return m.SaveAnalysisInventoryFunc(ctx, params)
	}
	return nil
}

func (m *MockAnalysisRepository) SaveAnalysisResult(ctx context.Context, params repository.SaveAnalysisResultParams) error {
	if m.SaveAnalysisResultFunc != nil {
		return m.SaveAnalysisResultFunc(ctx, params)
	}
	return nil
}
