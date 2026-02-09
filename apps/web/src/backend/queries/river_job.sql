-- name: FindActiveRiverJobByRepo :one
-- Find active (non-terminal) job for repository.
-- Terminal states (completed, cancelled, discarded) are excluded.
-- If job is cancelled, the usecase falls through to check completed analysis.
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

-- name: GetUserActiveJobs :many
-- Get all active jobs for a specific user.
-- Returns jobs in non-terminal states (available, pending, running, retryable, scheduled).
SELECT
    id,
    kind,
    state::text as state,
    args,
    created_at,
    attempted_at
FROM river_job
WHERE
    args->>'user_id' = @user_id::text
    AND state IN ('available', 'pending', 'running', 'retryable', 'scheduled')
ORDER BY created_at DESC;
