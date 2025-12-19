-- name: GetRiverJobByAnalysisID :one
SELECT
    args->>'analysis_id' as analysis_id,
    state::text as state
FROM river_job
WHERE
    kind = @kind::text
    AND args->>'analysis_id' = @analysis_id::text
ORDER BY created_at DESC
LIMIT 1;

-- name: FindActiveRiverJobByRepo :one
SELECT
    args->>'analysis_id' as analysis_id,
    state::text as state
FROM river_job
WHERE
    kind = @kind::text
    AND state IN ('available', 'pending', 'retryable', 'running', 'scheduled')
    AND args->>'owner' = @owner::text
    AND args->>'repo' = @repo::text
ORDER BY created_at DESC
LIMIT 1;
