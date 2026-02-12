package usecase

import (
	"context"
	"fmt"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/port"
)

type GetRepositoryStatsUseCase struct {
	repository port.Repository
}

func NewGetRepositoryStatsUseCase(repository port.Repository) *GetRepositoryStatsUseCase {
	return &GetRepositoryStatsUseCase{
		repository: repository,
	}
}

type GetRepositoryStatsInput struct {
	UserID string
}

func (uc *GetRepositoryStatsUseCase) Execute(ctx context.Context, input GetRepositoryStatsInput) (*entity.RepositoryStats, error) {
	stats, err := uc.repository.GetRepositoryStats(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get repository stats: %w", err)
	}

	return stats, nil
}
