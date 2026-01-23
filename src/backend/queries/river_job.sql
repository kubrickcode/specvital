-- name: FindActiveRiverJobByRepo :one
SELECT
    (args->>'commit_sha')::text as commit_sha,
    state::text as state,
    attempted_at
FROM river_job
WHERE
    kind = @kind::text
    AND state IN ('available', 'pending', 'retryable', 'running', 'scheduled')
    AND args->>'owner' = @owner::text
    AND args->>'repo' = @repo::text
ORDER BY created_at DESC
LIMIT 1;
