package usecase

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

type GetRepositoryStatsUseCase struct {
	repository port.Repository
}

func NewGetRepositoryStatsUseCase(repository port.Repository) *GetRepositoryStatsUseCase {
	return &GetRepositoryStatsUseCase{
		repository: repository,
	}
}

func (uc *GetRepositoryStatsUseCase) Execute(ctx context.Context) (*entity.RepositoryStats, error) {
	stats, err := uc.repository.GetRepositoryStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get repository stats: %w", err)
	}

	return stats, nil
}
