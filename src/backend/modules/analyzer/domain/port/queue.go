package port

import (
	"context"
	"time"

	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type QueueService interface {
	Enqueue(ctx context.Context, owner, repo, commitSHA string, userID *string, tier subscription.PlanTier) error
	FindTaskByRepo(ctx context.Context, owner, repo string) (*TaskInfo, error)
	Close() error
}

type TaskInfo struct {
	AttemptedAt *time.Time
	CommitSHA   string
	State       string
}
