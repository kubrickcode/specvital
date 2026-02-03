package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

type GetAnalysisHistoryInput struct {
	Owner string
	Repo  string
}

type AnalysisHistoryItem struct {
	BranchName  *string
	CommitSHA   string
	CommittedAt *time.Time
	CompletedAt time.Time
	ID          string
	TotalTests  int
}

type GetAnalysisHistoryOutput struct {
	Items []AnalysisHistoryItem
}

type GetAnalysisHistoryUseCase struct {
	repository port.Repository
}

func NewGetAnalysisHistoryUseCase(repository port.Repository) *GetAnalysisHistoryUseCase {
	return &GetAnalysisHistoryUseCase{
		repository: repository,
	}
}

func (uc *GetAnalysisHistoryUseCase) Execute(ctx context.Context, input GetAnalysisHistoryInput) (*GetAnalysisHistoryOutput, error) {
	if input.Owner == "" || input.Repo == "" {
		return nil, fmt.Errorf("owner and repo are required: %w", domain.ErrInvalidInput)
	}

	items, err := uc.repository.GetAnalysisHistory(ctx, input.Owner, input.Repo)
	if err != nil {
		return nil, err
	}

	output := make([]AnalysisHistoryItem, len(items))
	for i, item := range items {
		output[i] = AnalysisHistoryItem{
			BranchName:  item.BranchName,
			CommitSHA:   item.CommitSHA,
			CommittedAt: item.CommittedAt,
			CompletedAt: item.CompletedAt,
			ID:          item.ID,
			TotalTests:  item.TotalTests,
		}
	}

	return &GetAnalysisHistoryOutput{Items: output}, nil
}
