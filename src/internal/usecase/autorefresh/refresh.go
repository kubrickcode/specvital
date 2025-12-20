package autorefresh

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/specvital/collector/internal/domain/analysis"
)

// Tuned for 1h cron interval:
// - 1 failure: network blip
// - 2 failures: transient issue
// - 3 failures: persistent problem, stop to prevent cascade
const maxConsecutiveEnqueueFailures = 3

var ErrCircuitBreakerOpen = errors.New("circuit breaker: too many consecutive enqueue failures")

type AutoRefreshUseCase struct {
	repository analysis.AutoRefreshRepository
	taskQueue  analysis.TaskQueue
	vcs        analysis.VCS
}

func NewAutoRefreshUseCase(
	repository analysis.AutoRefreshRepository,
	taskQueue analysis.TaskQueue,
	vcs analysis.VCS,
) *AutoRefreshUseCase {
	return &AutoRefreshUseCase{
		repository: repository,
		taskQueue:  taskQueue,
		vcs:        vcs,
	}
}

// Returns ErrCircuitBreakerOpen if too many consecutive enqueue failures occur.
func (uc *AutoRefreshUseCase) Execute(ctx context.Context) error {
	codebases, err := uc.repository.GetCodebasesForAutoRefresh(ctx)
	if err != nil {
		return err
	}

	if len(codebases) == 0 {
		slog.InfoContext(ctx, "no codebases eligible for auto-refresh")
		return nil
	}

	now := time.Now()
	var enqueued int
	var consecutiveFailures int

	for _, codebase := range codebases {
		if consecutiveFailures >= maxConsecutiveEnqueueFailures {
			slog.ErrorContext(ctx, "circuit breaker open, aborting auto-refresh",
				"consecutive_failures", consecutiveFailures,
				"enqueued_before_abort", enqueued,
			)
			return fmt.Errorf("%w: %d failures", ErrCircuitBreakerOpen, consecutiveFailures)
		}

		if !analysis.ShouldRefreshAt(
			codebase.LastViewedAt,
			codebase.LastCompletedAt,
			codebase.ConsecutiveFailures,
			now,
		) {
			continue
		}

		repoURL := fmt.Sprintf("https://%s/%s/%s", codebase.Host, codebase.Owner, codebase.Name)
		commitSHA, err := uc.vcs.GetHeadCommit(ctx, repoURL, nil)
		if err != nil {
			consecutiveFailures++
			slog.ErrorContext(ctx, "failed to get head commit for auto-refresh",
				"owner", codebase.Owner,
				"repo", codebase.Name,
				"consecutive_failures", consecutiveFailures,
				"error", err,
			)
			continue
		}

		if codebase.LastCommitSHA == commitSHA {
			slog.DebugContext(ctx, "skipping auto-refresh: no new commits",
				"owner", codebase.Owner,
				"repo", codebase.Name,
				"commit", commitSHA,
			)
			continue
		}

		if err := uc.taskQueue.EnqueueAnalysis(ctx, codebase.Owner, codebase.Name, commitSHA); err != nil {
			consecutiveFailures++
			slog.ErrorContext(ctx, "failed to enqueue auto-refresh task",
				"owner", codebase.Owner,
				"repo", codebase.Name,
				"consecutive_failures", consecutiveFailures,
				"error", err,
			)
			continue
		}

		consecutiveFailures = 0
		enqueued++
		slog.DebugContext(ctx, "enqueued auto-refresh task",
			"owner", codebase.Owner,
			"repo", codebase.Name,
		)
	}

	slog.InfoContext(ctx, "auto-refresh execution completed",
		"total_candidates", len(codebases),
		"enqueued", enqueued,
	)

	return nil
}
