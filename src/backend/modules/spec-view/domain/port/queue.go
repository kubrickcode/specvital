package port

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type QueueService interface {
	EnqueueSpecGeneration(ctx context.Context, analysisID string, language string, userID *string, tier subscription.PlanTier, mode entity.GenerationMode) error
}
