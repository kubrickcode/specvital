package adapter

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

type CacheRepositoryPostgres struct {
	conn    db.DBTX
	queries *db.Queries
}

var _ port.CacheRepository = (*CacheRepositoryPostgres)(nil)

func NewCacheRepositoryPostgres(conn db.DBTX, queries *db.Queries) *CacheRepositoryPostgres {
	if conn == nil {
		panic("conn is required")
	}
	if queries == nil {
		panic("queries is required")
	}
	return &CacheRepositoryPostgres{conn: conn, queries: queries}
}

func (r *CacheRepositoryPostgres) GetCachedConversions(ctx context.Context, keyHashes [][]byte, modelID string) (map[string]*entity.CacheEntry, error) {
	if len(keyHashes) == 0 {
		return make(map[string]*entity.CacheEntry), nil
	}

	rows, err := r.queries.GetCachedConversions(ctx, db.GetCachedConversionsParams{
		Column1: keyHashes,
		ModelID: modelID,
	})
	if err != nil {
		return nil, fmt.Errorf("get cached conversions: %w", err)
	}

	result := make(map[string]*entity.CacheEntry, len(rows))
	for _, row := range rows {
		entry := mapCacheEntryFromDB(&row)
		keyHex := hex.EncodeToString(entry.CacheKeyHash)
		result[keyHex] = entry
	}

	return result, nil
}

func (r *CacheRepositoryPostgres) UpsertCachedConversions(ctx context.Context, entries []*entity.CacheEntry) error {
	if len(entries) == 0 {
		return nil
	}

	for _, entry := range entries {
		if err := r.queries.UpsertCachedConversion(ctx, db.UpsertCachedConversionParams{
			CacheKeyHash:   entry.CacheKeyHash,
			CodebaseID:     uuidToPgtype(entry.CodebaseID),
			ConvertedName:  entry.ConvertedName,
			FilePath:       entry.FilePath,
			Framework:      entry.Framework,
			Language:       entry.Language.String(),
			ModelID:        entry.ModelID,
			OriginalName:   entry.OriginalName,
			SuiteHierarchy: entry.SuiteHierarchy,
		}); err != nil {
			return fmt.Errorf("upsert cached conversion: %w", err)
		}
	}

	return nil
}

func (r *CacheRepositoryPostgres) DeleteCodebaseCache(ctx context.Context, codebaseID uuid.UUID) error {
	if err := r.queries.DeleteCodebaseCache(ctx, uuidToPgtype(codebaseID)); err != nil {
		return fmt.Errorf("delete codebase cache: %w", err)
	}
	return nil
}

func mapCacheEntryFromDB(row *db.SpecViewCache) *entity.CacheEntry {
	return &entity.CacheEntry{
		CacheKeyHash:   row.CacheKeyHash,
		CodebaseID:     pgtypeToUUID(row.CodebaseID),
		ConvertedName:  row.ConvertedName,
		CreatedAt:      row.CreatedAt.Time,
		FilePath:       row.FilePath,
		Framework:      row.Framework,
		ID:             pgtypeToUUID(row.ID),
		Language:       entity.Language(row.Language),
		ModelID:        row.ModelID,
		OriginalName:   row.OriginalName,
		SuiteHierarchy: row.SuiteHierarchy,
	}
}

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func pgtypeToUUID(pg pgtype.UUID) uuid.UUID {
	if !pg.Valid {
		return uuid.Nil
	}
	return pg.Bytes
}
