package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type ReanalyzeRepositoryInput struct {
	Owner  string
	Repo   string
	Tier   subscription.PlanTier
	UserID string
}

type ReanalyzeRepositoryUseCase struct {
	gitClient     port.GitClient
	queue         port.QueueService
	repository    port.Repository
	tokenProvider port.TokenProvider
}

func NewReanalyzeRepositoryUseCase(
	gitClient port.GitClient,
	queue port.QueueService,
	repository port.Repository,
	tokenProvider port.TokenProvider,
) *ReanalyzeRepositoryUseCase {
	return &ReanalyzeRepositoryUseCase{
		gitClient:     gitClient,
		queue:         queue,
		repository:    repository,
		tokenProvider: tokenProvider,
	}
}

func (uc *ReanalyzeRepositoryUseCase) Execute(ctx context.Context, input ReanalyzeRepositoryInput) (*AnalyzeResult, error) {
	if input.Owner == "" || input.Repo == "" {
		return nil, errors.New("owner and repo are required")
	}

	now := time.Now()

	latestSHA, err := getLatestCommitWithAuth(ctx, uc.gitClient, uc.tokenProvider, input.Owner, input.Repo, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get latest commit for %s/%s: %w", input.Owner, input.Repo, err)
	}

	var userIDPtr *string
	if input.UserID != "" {
		userIDPtr = &input.UserID
	}

	if err := uc.queue.Enqueue(ctx, input.Owner, input.Repo, latestSHA, userIDPtr, input.Tier); err != nil {
		return nil, fmt.Errorf("queue reanalysis for %s/%s: %w", input.Owner, input.Repo, err)
	}

	progress := &entity.AnalysisProgress{
		CommitSHA: latestSHA,
		CreatedAt: now,
		Status:    entity.AnalysisStatusPending,
	}
	return &AnalyzeResult{Progress: progress}, nil
}
