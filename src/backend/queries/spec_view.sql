-- name: GetCachedConversions :many
SELECT
    id,
    cache_key_hash,
    codebase_id,
    file_path,
    framework,
    suite_hierarchy,
    original_name,
    converted_name,
    language,
    model_id,
    created_at
FROM spec_view_cache
WHERE cache_key_hash = ANY($1::bytea[])
  AND model_id = $2;

-- name: UpsertCachedConversions :copyfrom
INSERT INTO spec_view_cache (
    cache_key_hash,
    codebase_id,
    file_path,
    framework,
    suite_hierarchy,
    original_name,
    converted_name,
    language,
    model_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: DeleteCodebaseCache :exec
DELETE FROM spec_view_cache
WHERE codebase_id = $1;
