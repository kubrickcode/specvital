package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

type CacheRepository interface {
	GetCachedConversions(ctx context.Context, keyHashes [][]byte, modelID string) (map[string]*entity.CacheEntry, error)
	UpsertCachedConversions(ctx context.Context, entries []*entity.CacheEntry) error
	DeleteCodebaseCache(ctx context.Context, codebaseID uuid.UUID) error
}
