package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/kubrickcode/specvital/apps/web/src/backend/common/queue"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/domain/port"
	subscription "github.com/kubrickcode/specvital/apps/web/src/backend/modules/subscription/domain/entity"
)

var _ port.QueueService = (*RiverQueueService)(nil)

const (
	TypeAnalyze = "analysis:analyze"

	maxRetries     = 3
	enqueueTimeout = 5 * time.Second
)

type AnalyzeArgs struct {
	CommitSHA string  `json:"commit_sha" river:"unique"`
	Owner     string  `json:"owner" river:"unique"`
	Repo      string  `json:"repo" river:"unique"`
	UserID    *string `json:"user_id,omitempty"`
}

func (AnalyzeArgs) Kind() string { return TypeAnalyze }

type RiverQueueService struct {
	client *river.Client[pgx.Tx]
	repo   port.Repository
}

func NewRiverQueueService(client *river.Client[pgx.Tx], repo port.Repository) *RiverQueueService {
	return &RiverQueueService{client: client, repo: repo}
}

func (s *RiverQueueService) Enqueue(ctx context.Context, owner, repo, commitSHA string, userID *string, tier subscription.PlanTier) error {
	ctx, cancel := context.WithTimeout(ctx, enqueueTimeout)
	defer cancel()

	args := AnalyzeArgs{
		CommitSHA: commitSHA,
		Owner:     owner,
		Repo:      repo,
		UserID:    userID,
	}

	targetQueue := queue.SelectQueueForAnalysis(tier, false)

	_, err := s.client.Insert(ctx, args, &river.InsertOpts{
		MaxAttempts: maxRetries,
		Queue:       targetQueue,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
			ByState: []rivertype.JobState{
				rivertype.JobStateAvailable,
				rivertype.JobStatePending,
				rivertype.JobStateRunning,
				rivertype.JobStateRetryable,
				rivertype.JobStateScheduled,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("enqueue task for %s/%s: %w", owner, repo, err)
	}

	return nil
}

func (s *RiverQueueService) EnqueueTx(ctx context.Context, tx pgx.Tx, owner, repo, commitSHA string, userID *string, tier subscription.PlanTier) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, enqueueTimeout)
	defer cancel()

	args := AnalyzeArgs{
		CommitSHA: commitSHA,
		Owner:     owner,
		Repo:      repo,
		UserID:    userID,
	}

	targetQueue := queue.SelectQueueForAnalysis(tier, false)

	result, err := s.client.InsertTx(ctx, tx, args, &river.InsertOpts{
		MaxAttempts: maxRetries,
		Queue:       targetQueue,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
			ByState: []rivertype.JobState{
				rivertype.JobStateAvailable,
				rivertype.JobStatePending,
				rivertype.JobStateRunning,
				rivertype.JobStateRetryable,
				rivertype.JobStateScheduled,
			},
		},
	})
	if err != nil {
		return 0, fmt.Errorf("enqueue task for %s/%s: %w", owner, repo, err)
	}

	return result.Job.ID, nil
}

func (s *RiverQueueService) FindTaskByRepo(ctx context.Context, owner, repo string) (*port.TaskInfo, error) {
	info, err := s.repo.FindActiveRiverJobByRepo(ctx, TypeAnalyze, owner, repo)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	return &port.TaskInfo{
		AttemptedAt: info.AttemptedAt,
		CommitSHA:   info.CommitSHA,
		State:       info.State,
	}, nil
}

// Close is a no-op as River client lifecycle is managed by the application container.
func (s *RiverQueueService) Close() error {
	return nil
}
