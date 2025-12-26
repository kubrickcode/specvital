package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

type GetUpdateStatusInput struct {
	Owner  string
	Repo   string
	UserID string
}

type GetUpdateStatusUseCase struct {
	gitClient     port.GitClient
	repository    port.Repository
	tokenProvider port.TokenProvider
}

func NewGetUpdateStatusUseCase(
	gitClient port.GitClient,
	repository port.Repository,
	tokenProvider port.TokenProvider,
) *GetUpdateStatusUseCase {
	return &GetUpdateStatusUseCase{
		gitClient:     gitClient,
		repository:    repository,
		tokenProvider: tokenProvider,
	}
}

func (uc *GetUpdateStatusUseCase) Execute(ctx context.Context, input GetUpdateStatusInput) (*entity.UpdateStatusResult, error) {
	if input.Owner == "" || input.Repo == "" {
		return nil, errors.New("owner and repo are required")
	}

	completed, err := uc.repository.GetLatestCompletedAnalysis(ctx, input.Owner, input.Repo)
	if err != nil {
		return nil, fmt.Errorf("get latest analysis: %w", err)
	}

	latestSHA, err := getLatestCommitWithAuth(ctx, uc.gitClient, uc.tokenProvider, input.Owner, input.Repo, input.UserID)
	if err != nil {
		return &entity.UpdateStatusResult{
			AnalyzedCommitSHA: completed.CommitSHA,
			Status:            entity.UpdateStatusUnknown,
		}, nil
	}

	status := entity.UpdateStatusUpToDate
	if latestSHA != completed.CommitSHA {
		status = entity.UpdateStatusNewCommits
	}

	return &entity.UpdateStatusResult{
		AnalyzedCommitSHA: completed.CommitSHA,
		LatestCommitSHA:   latestSHA,
		Status:            status,
	}, nil
}
