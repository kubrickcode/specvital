package port

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type QueueService interface {
	Enqueue(ctx context.Context, owner, repo, commitSHA string, userID *string, tier subscription.PlanTier) error
	// EnqueueTx enqueues an analysis job within a transaction.
	// Returns the job ID for quota reservation tracking.
	EnqueueTx(ctx context.Context, tx pgx.Tx, owner, repo, commitSHA string, userID *string, tier subscription.PlanTier) (int64, error)
	FindTaskByRepo(ctx context.Context, owner, repo string) (*TaskInfo, error)
	Close() error
}

type TaskInfo struct {
	AttemptedAt *time.Time
	CommitSHA   string
	State       string
}
