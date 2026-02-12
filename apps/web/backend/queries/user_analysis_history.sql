-- name: GetUserAnalyzedRepositories :many
SELECT
    uah.id AS history_id,
    uah.updated_at,
    a.id AS analysis_id,
    a.commit_sha,
    a.completed_at,
    a.total_tests,
    c.id AS codebase_id,
    c.owner,
    c.name,
    u.username AS user_username
FROM user_analysis_history uah
JOIN analyses a ON a.id = uah.analysis_id
JOIN codebases c ON c.id = a.codebase_id
JOIN users u ON u.id = uah.user_id
LEFT JOIN user_github_org_memberships ugom ON ugom.user_id = uah.user_id
LEFT JOIN github_organizations go ON go.id = ugom.org_id AND go.login = c.owner
WHERE uah.user_id = sqlc.arg(user_id)
  AND a.status = 'completed'
  AND c.is_stale = false
  AND (
    sqlc.arg(cursor_time)::timestamptz IS NULL
    OR (uah.updated_at, uah.id) < (sqlc.arg(cursor_time)::timestamptz, sqlc.arg(cursor_id)::uuid)
  )
  AND (
    sqlc.arg(ownership)::text = 'all'
    OR (sqlc.arg(ownership)::text = 'mine' AND c.owner = u.username)
    OR (sqlc.arg(ownership)::text = 'organization' AND go.id IS NOT NULL)
    OR (sqlc.arg(ownership)::text = 'others' AND c.owner != u.username AND go.id IS NULL)
  )
ORDER BY uah.updated_at DESC, uah.id DESC
LIMIT sqlc.arg(page_limit);

-- name: AddUserAnalyzedRepository :one
INSERT INTO user_analysis_history (user_id, analysis_id)
VALUES (sqlc.arg(user_id), sqlc.arg(analysis_id))
ON CONFLICT (user_id, analysis_id) DO UPDATE SET updated_at = now()
RETURNING id, analysis_id, updated_at;

-- name: GetLatestAnalysisIDByOwnerRepo :one
SELECT a.id
FROM analyses a
JOIN codebases c ON c.id = a.codebase_id
WHERE c.host = 'github.com'
  AND c.owner = sqlc.arg(owner)
  AND c.name = sqlc.arg(repo)
  AND a.status = 'completed'
ORDER BY a.completed_at DESC
LIMIT 1;

-- name: CheckUserHistoryExists :one
SELECT EXISTS (
    SELECT 1 FROM user_analysis_history uah
    JOIN analyses a ON a.id = uah.analysis_id
    JOIN codebases c ON c.id = a.codebase_id
    WHERE uah.user_id = sqlc.arg(user_id)
      AND c.host = 'github.com'
      AND c.owner = sqlc.arg(owner)
      AND c.name = sqlc.arg(repo)
      AND a.status = 'completed'
) AS exists;
