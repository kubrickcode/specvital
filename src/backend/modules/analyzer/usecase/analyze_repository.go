package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type AnalyzeRepositoryInput struct {
	Owner  string
	Repo   string
	Tier   subscription.PlanTier
	UserID string
}

type AnalyzeRepositoryUseCase struct {
	gitClient     port.GitClient
	queue         port.QueueService
	repository    port.Repository
	systemConfig  port.SystemConfigReader
	tokenProvider port.TokenProvider
}

func NewAnalyzeRepositoryUseCase(
	gitClient port.GitClient,
	queue port.QueueService,
	repository port.Repository,
	systemConfig port.SystemConfigReader,
	tokenProvider port.TokenProvider,
) *AnalyzeRepositoryUseCase {
	return &AnalyzeRepositoryUseCase{
		gitClient:     gitClient,
		queue:         queue,
		repository:    repository,
		systemConfig:  systemConfig,
		tokenProvider: tokenProvider,
	}
}

func (uc *AnalyzeRepositoryUseCase) Execute(ctx context.Context, input AnalyzeRepositoryInput) (*AnalyzeResult, error) {
	if input.Owner == "" || input.Repo == "" {
		return nil, errors.New("owner and repo are required")
	}

	now := time.Now()

	latestSHA, err := getLatestCommitWithAuth(ctx, uc.gitClient, uc.tokenProvider, input.Owner, input.Repo, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get latest commit for %s/%s: %w", input.Owner, input.Repo, err)
	}

	completed, err := uc.repository.GetLatestCompletedAnalysis(ctx, input.Owner, input.Repo)
	if err == nil {
		if uc.shouldReturnCachedAnalysis(completed) {
			analysis, buildErr := buildAnalysisFromCompleted(ctx, uc.repository, completed)
			if buildErr != nil {
				return nil, fmt.Errorf("build analysis for %s/%s: %w", input.Owner, input.Repo, buildErr)
			}
			// Non-critical: UpdateLastViewed failure doesn't affect main flow
			_ = uc.repository.UpdateLastViewed(ctx, input.Owner, input.Repo)
			return &AnalyzeResult{Analysis: analysis}, nil
		}
	}

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("get analysis for %s/%s: %w", input.Owner, input.Repo, err)
	}

	taskInfo, err := uc.queue.FindTaskByRepo(ctx, input.Owner, input.Repo)
	// Non-critical: queue search failure doesn't block new task creation
	_ = err
	if taskInfo != nil && taskInfo.CommitSHA == latestSHA {
		progress := &entity.AnalysisProgress{
			CommitSHA: taskInfo.CommitSHA,
			CreatedAt: now,
			StartedAt: taskInfo.AttemptedAt,
			Status:    mapQueueStateToAnalysisStatus(taskInfo.State),
		}
		return &AnalyzeResult{Progress: progress}, nil
	}

	var userIDPtr *string
	if input.UserID != "" {
		userIDPtr = &input.UserID
	}

	if err := uc.queue.Enqueue(ctx, input.Owner, input.Repo, latestSHA, userIDPtr, input.Tier); err != nil {
		return nil, fmt.Errorf("queue analysis for %s/%s: %w", input.Owner, input.Repo, err)
	}

	progress := &entity.AnalysisProgress{
		CommitSHA: latestSHA,
		CreatedAt: now,
		Status:    entity.AnalysisStatusPending,
	}
	return &AnalyzeResult{Progress: progress}, nil
}

// shouldReturnCachedAnalysis determines if the cached analysis can be returned.
// Cache-first policy: Returns cached analysis even with new commits or parser updates.
// Returns false (needs re-analysis) only when:
// - NULL parser_version (legacy data before version tracking)
//
// Note: Neither commit SHA difference nor parser version mismatch triggers re-analysis.
// Users can manually trigger reanalysis via the update banner.
func (uc *AnalyzeRepositoryUseCase) shouldReturnCachedAnalysis(completed *port.CompletedAnalysis) bool {
	// Legacy data without parser_version → needs re-analysis
	return completed.ParserVersion != nil
}
