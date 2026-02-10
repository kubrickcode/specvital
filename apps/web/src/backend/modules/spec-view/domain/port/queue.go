package port

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/entity"
	subscription "github.com/kubrickcode/specvital/apps/web/src/backend/modules/subscription/domain/entity"
)

type QueueService interface {
	EnqueueSpecGeneration(ctx context.Context, analysisID string, language string, userID *string, tier subscription.PlanTier, mode entity.GenerationMode) error
	// EnqueueSpecGenerationTx enqueues a spec generation job within a transaction.
	// Returns the job ID for quota reservation tracking.
	EnqueueSpecGenerationTx(ctx context.Context, tx pgx.Tx, analysisID string, language string, userID *string, tier subscription.PlanTier, mode entity.GenerationMode) (int64, error)
}
