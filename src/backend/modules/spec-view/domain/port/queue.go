package port

import (
	"context"

	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type QueueService interface {
	EnqueueSpecGeneration(ctx context.Context, analysisID string, language string, userID *string, tier subscription.PlanTier, forceRegenerate bool) error
}
