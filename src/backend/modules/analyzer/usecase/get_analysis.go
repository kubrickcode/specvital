package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

type GetAnalysisInput struct {
	Owner string
	Repo  string
}

type GetAnalysisUseCase struct {
	queue      port.QueueService
	repository port.Repository
}

func NewGetAnalysisUseCase(
	queue port.QueueService,
	repository port.Repository,
) *GetAnalysisUseCase {
	return &GetAnalysisUseCase{
		queue:      queue,
		repository: repository,
	}
}

func (uc *GetAnalysisUseCase) Execute(ctx context.Context, input GetAnalysisInput) (*AnalyzeResult, error) {
	if input.Owner == "" || input.Repo == "" {
		return nil, errors.New("owner and repo are required")
	}

	now := time.Now()

	// Check for recent jobs FIRST - if there's a job in queue, return its progress
	// This prevents returning old completed analysis while new analysis is pending
	taskInfo, err := uc.queue.FindTaskByRepo(ctx, input.Owner, input.Repo)
	if err != nil {
		slog.DebugContext(ctx, "queue search failed, continuing with cache check", "owner", input.Owner, "repo", input.Repo, "error", err)
	}
	if taskInfo != nil {
		progress := &entity.AnalysisProgress{
			CommitSHA: taskInfo.CommitSHA,
			CreatedAt: now,
			StartedAt: taskInfo.AttemptedAt,
			Status:    mapQueueStateToAnalysisStatus(taskInfo.State),
		}
		return &AnalyzeResult{Progress: progress}, nil
	}

	// No recent job - check for completed analysis
	completed, err := uc.repository.GetLatestCompletedAnalysis(ctx, input.Owner, input.Repo)
	if err == nil {
		analysis, buildErr := buildAnalysisFromCompleted(ctx, uc.repository, completed)
		if buildErr != nil {
			return nil, fmt.Errorf("build analysis for %s/%s: %w", input.Owner, input.Repo, buildErr)
		}
		// Non-critical: UpdateLastViewed failure doesn't affect main flow
		_ = uc.repository.UpdateLastViewed(ctx, input.Owner, input.Repo)
		return &AnalyzeResult{Analysis: analysis}, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("get analysis for %s/%s: %w", input.Owner, input.Repo, err)
	}

	return nil, domain.ErrNotFound
}
