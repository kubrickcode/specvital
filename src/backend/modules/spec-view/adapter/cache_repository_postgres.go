package adapter

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	params := make([]db.UpsertCachedConversionsParams, 0, len(entries))
	for _, entry := range entries {
		params = append(params, db.UpsertCachedConversionsParams{
			CacheKeyHash:   entry.CacheKeyHash,
			CodebaseID:     uuidToPgtype(entry.CodebaseID),
			ConvertedName:  entry.ConvertedName,
			FilePath:       entry.FilePath,
			Framework:      entry.Framework,
			Language:       entry.Language.String(),
			ModelID:        entry.ModelID,
			OriginalName:   entry.OriginalName,
			SuiteHierarchy: entry.SuiteHierarchy,
		})
	}

	copyCount, err := r.conn.CopyFrom(
		ctx,
		pgx.Identifier{"spec_view_cache"},
		[]string{
			"cache_key_hash",
			"codebase_id",
			"file_path",
			"framework",
			"suite_hierarchy",
			"original_name",
			"converted_name",
			"language",
			"model_id",
		},
		pgx.CopyFromSlice(len(params), func(i int) ([]any, error) {
			p := params[i]
			return []any{
				p.CacheKeyHash,
				p.CodebaseID,
				p.FilePath,
				p.Framework,
				p.SuiteHierarchy,
				p.OriginalName,
				p.ConvertedName,
				p.Language,
				p.ModelID,
			}, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("copy from: %w", err)
	}
	if copyCount != int64(len(params)) {
		return fmt.Errorf("expected %d rows, got %d", len(params), copyCount)
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
