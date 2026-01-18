package port

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

type SpecViewRepository interface {
	CheckAnalysisExists(ctx context.Context, analysisID string) (bool, error)
	CheckSpecDocumentExists(ctx context.Context, analysisID string) (bool, error)
	GetSpecDocument(ctx context.Context, analysisID string) (*entity.SpecDocument, error)
	GetGenerationStatus(ctx context.Context, analysisID string) (*entity.SpecGenerationStatus, error)
}

type TierLookup interface {
	GetUserTier(ctx context.Context, userID string) (string, error)
}
