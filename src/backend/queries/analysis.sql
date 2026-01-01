-- name: GetLatestCompletedAnalysis :one
SELECT
    a.id,
    a.commit_sha,
    a.branch_name,
    a.committed_at,
    a.completed_at,
    a.total_suites,
    a.total_tests,
    c.owner,
    c.name as repo
FROM analyses a
JOIN codebases c ON c.id = a.codebase_id
WHERE c.host = $1 AND c.owner = $2 AND c.name = $3
  AND a.status = 'completed'
ORDER BY a.created_at DESC
LIMIT 1;

-- name: GetAnalysisStatus :one
SELECT
    a.id,
    a.status,
    a.error_message,
    a.created_at,
    a.completed_at
FROM analyses a
JOIN codebases c ON c.id = a.codebase_id
WHERE c.host = $1 AND c.owner = $2 AND c.name = $3
ORDER BY a.created_at DESC
LIMIT 1;

-- name: UpsertCodebase :one
INSERT INTO codebases (host, owner, name)
VALUES ($1, $2, $3)
ON CONFLICT (host, owner, name) DO UPDATE SET updated_at = now()
RETURNING id;

-- name: CreatePendingAnalysis :one
INSERT INTO analyses (codebase_id, commit_sha, status)
VALUES ($1, $2, 'pending')
RETURNING id;

-- name: MarkAnalysisFailed :exec
UPDATE analyses
SET status = 'failed', error_message = $2
WHERE id = $1;

-- name: GetTestSuitesByAnalysisID :many
SELECT
    ts.id,
    ts.file_path,
    ts.framework,
    ts.name
FROM test_suites ts
WHERE ts.analysis_id = $1
ORDER BY ts.file_path, ts.depth, ts.line_number;

-- name: GetTestCasesBySuiteIDs :many
SELECT
    tc.id,
    tc.suite_id,
    tc.name,
    tc.line_number,
    tc.status
FROM test_cases tc
WHERE tc.suite_id = ANY($1::uuid[])
ORDER BY tc.suite_id, tc.line_number;

-- name: UpdateCodebaseLastViewed :exec
UPDATE codebases
SET last_viewed_at = now()
WHERE host = $1 AND owner = $2 AND name = $3;

-- name: GetRepositoryStats :one
SELECT
    COUNT(DISTINCT c.id) AS total_repositories,
    COALESCE(SUM(a.total_tests), 0)::bigint AS total_tests
FROM codebases c
LEFT JOIN LATERAL (
    SELECT total_tests
    FROM analyses
    WHERE codebase_id = c.id AND status = 'completed'
    ORDER BY created_at DESC
    LIMIT 1
) a ON true
WHERE c.last_viewed_at IS NOT NULL AND c.is_stale = false;

-- name: GetPreviousAnalysis :one
SELECT
    id,
    commit_sha,
    completed_at,
    total_tests
FROM analyses
WHERE codebase_id = $1
  AND status = 'completed'
  AND id != $2
ORDER BY created_at DESC
LIMIT 1;

-- name: GetCodebaseIDByOwnerRepo :one
SELECT id
FROM codebases
WHERE host = $1 AND owner = $2 AND name = $3 AND is_stale = false;

-- name: GetPaginatedRepositoriesByRecent :many
SELECT
    c.id AS codebase_id,
    c.owner,
    c.name,
    a.id AS analysis_id,
    a.commit_sha,
    a.completed_at AS analyzed_at,
    a.total_tests,
    a.active_count,
    a.focused_count,
    a.skipped_count,
    a.todo_count,
    a.xfail_count,
    EXISTS(
        SELECT 1 FROM user_analysis_history uah
        WHERE uah.analysis_id = a.id AND uah.user_id = sqlc.arg(user_id)::uuid
    ) AS is_analyzed_by_me
FROM codebases c
JOIN LATERAL (
    SELECT
        an.id,
        an.commit_sha,
        an.completed_at,
        an.total_tests,
        COALESCE(tc_summary.active_count, 0)::int AS active_count,
        COALESCE(tc_summary.focused_count, 0)::int AS focused_count,
        COALESCE(tc_summary.skipped_count, 0)::int AS skipped_count,
        COALESCE(tc_summary.todo_count, 0)::int AS todo_count,
        COALESCE(tc_summary.xfail_count, 0)::int AS xfail_count
    FROM analyses an
    LEFT JOIN LATERAL (
        SELECT
            COUNT(*) FILTER (WHERE tc.status = 'active') AS active_count,
            COUNT(*) FILTER (WHERE tc.status = 'focused') AS focused_count,
            COUNT(*) FILTER (WHERE tc.status = 'skipped') AS skipped_count,
            COUNT(*) FILTER (WHERE tc.status = 'todo') AS todo_count,
            COUNT(*) FILTER (WHERE tc.status = 'xfail') AS xfail_count
        FROM test_cases tc
        JOIN test_suites ts ON ts.id = tc.suite_id
        WHERE ts.analysis_id = an.id
    ) tc_summary ON true
    WHERE an.codebase_id = c.id AND an.status = 'completed'
    ORDER BY an.created_at DESC
    LIMIT 1
) a ON true
WHERE c.last_viewed_at IS NOT NULL
  AND c.is_stale = false
  AND (
    sqlc.arg(view_filter)::text = 'all'
    OR (sqlc.arg(view_filter)::text = 'my' AND EXISTS(
        SELECT 1 FROM user_analysis_history uah
        WHERE uah.analysis_id = a.id AND uah.user_id = sqlc.arg(user_id)::uuid
    ))
    OR (sqlc.arg(view_filter)::text = 'community' AND NOT EXISTS(
        SELECT 1 FROM user_analysis_history uah
        WHERE uah.analysis_id = a.id AND uah.user_id = sqlc.arg(user_id)::uuid
    ))
  )
  AND (
    sqlc.arg(cursor_analyzed_at)::timestamptz IS NULL
    OR (
      (sqlc.arg(sort_order)::text = 'desc' AND (a.completed_at, c.id) < (sqlc.arg(cursor_analyzed_at)::timestamptz, sqlc.arg(cursor_id)::uuid))
      OR (sqlc.arg(sort_order)::text = 'asc' AND (a.completed_at, c.id) > (sqlc.arg(cursor_analyzed_at)::timestamptz, sqlc.arg(cursor_id)::uuid))
    )
  )
ORDER BY
  CASE WHEN sqlc.arg(sort_order)::text = 'desc' THEN a.completed_at END DESC,
  CASE WHEN sqlc.arg(sort_order)::text = 'asc' THEN a.completed_at END ASC,
  CASE WHEN sqlc.arg(sort_order)::text = 'desc' THEN c.id END DESC,
  CASE WHEN sqlc.arg(sort_order)::text = 'asc' THEN c.id END ASC
LIMIT sqlc.arg(page_limit);

-- name: GetPaginatedRepositoriesByName :many
SELECT
    c.id AS codebase_id,
    c.owner,
    c.name,
    a.id AS analysis_id,
    a.commit_sha,
    a.completed_at AS analyzed_at,
    a.total_tests,
    a.active_count,
    a.focused_count,
    a.skipped_count,
    a.todo_count,
    a.xfail_count,
    EXISTS(
        SELECT 1 FROM user_analysis_history uah
        WHERE uah.analysis_id = a.id AND uah.user_id = sqlc.arg(user_id)::uuid
    ) AS is_analyzed_by_me
FROM codebases c
JOIN LATERAL (
    SELECT
        an.id,
        an.commit_sha,
        an.completed_at,
        an.total_tests,
        COALESCE(tc_summary.active_count, 0)::int AS active_count,
        COALESCE(tc_summary.focused_count, 0)::int AS focused_count,
        COALESCE(tc_summary.skipped_count, 0)::int AS skipped_count,
        COALESCE(tc_summary.todo_count, 0)::int AS todo_count,
        COALESCE(tc_summary.xfail_count, 0)::int AS xfail_count
    FROM analyses an
    LEFT JOIN LATERAL (
        SELECT
            COUNT(*) FILTER (WHERE tc.status = 'active') AS active_count,
            COUNT(*) FILTER (WHERE tc.status = 'focused') AS focused_count,
            COUNT(*) FILTER (WHERE tc.status = 'skipped') AS skipped_count,
            COUNT(*) FILTER (WHERE tc.status = 'todo') AS todo_count,
            COUNT(*) FILTER (WHERE tc.status = 'xfail') AS xfail_count
        FROM test_cases tc
        JOIN test_suites ts ON ts.id = tc.suite_id
        WHERE ts.analysis_id = an.id
    ) tc_summary ON true
    WHERE an.codebase_id = c.id AND an.status = 'completed'
    ORDER BY an.created_at DESC
    LIMIT 1
) a ON true
WHERE c.last_viewed_at IS NOT NULL
  AND c.is_stale = false
  AND (
    sqlc.arg(view_filter)::text = 'all'
    OR (sqlc.arg(view_filter)::text = 'my' AND EXISTS(
        SELECT 1 FROM user_analysis_history uah
        WHERE uah.analysis_id = a.id AND uah.user_id = sqlc.arg(user_id)::uuid
    ))
    OR (sqlc.arg(view_filter)::text = 'community' AND NOT EXISTS(
        SELECT 1 FROM user_analysis_history uah
        WHERE uah.analysis_id = a.id AND uah.user_id = sqlc.arg(user_id)::uuid
    ))
  )
  AND (
    sqlc.arg(cursor_name)::text IS NULL
    OR (
      (sqlc.arg(sort_order)::text = 'asc' AND (c.name, c.id) > (sqlc.arg(cursor_name)::text, sqlc.arg(cursor_id)::uuid))
      OR (sqlc.arg(sort_order)::text = 'desc' AND (c.name, c.id) < (sqlc.arg(cursor_name)::text, sqlc.arg(cursor_id)::uuid))
    )
  )
ORDER BY
  CASE WHEN sqlc.arg(sort_order)::text = 'asc' THEN c.name END ASC,
  CASE WHEN sqlc.arg(sort_order)::text = 'desc' THEN c.name END DESC,
  CASE WHEN sqlc.arg(sort_order)::text = 'asc' THEN c.id END ASC,
  CASE WHEN sqlc.arg(sort_order)::text = 'desc' THEN c.id END DESC
LIMIT sqlc.arg(page_limit);

-- name: GetPaginatedRepositoriesByTests :many
SELECT
    c.id AS codebase_id,
    c.owner,
    c.name,
    a.id AS analysis_id,
    a.commit_sha,
    a.completed_at AS analyzed_at,
    a.total_tests,
    a.active_count,
    a.focused_count,
    a.skipped_count,
    a.todo_count,
    a.xfail_count,
    EXISTS(
        SELECT 1 FROM user_analysis_history uah
        WHERE uah.analysis_id = a.id AND uah.user_id = sqlc.arg(user_id)::uuid
    ) AS is_analyzed_by_me
FROM codebases c
JOIN LATERAL (
    SELECT
        an.id,
        an.commit_sha,
        an.completed_at,
        an.total_tests,
        COALESCE(tc_summary.active_count, 0)::int AS active_count,
        COALESCE(tc_summary.focused_count, 0)::int AS focused_count,
        COALESCE(tc_summary.skipped_count, 0)::int AS skipped_count,
        COALESCE(tc_summary.todo_count, 0)::int AS todo_count,
        COALESCE(tc_summary.xfail_count, 0)::int AS xfail_count
    FROM analyses an
    LEFT JOIN LATERAL (
        SELECT
            COUNT(*) FILTER (WHERE tc.status = 'active') AS active_count,
            COUNT(*) FILTER (WHERE tc.status = 'focused') AS focused_count,
            COUNT(*) FILTER (WHERE tc.status = 'skipped') AS skipped_count,
            COUNT(*) FILTER (WHERE tc.status = 'todo') AS todo_count,
            COUNT(*) FILTER (WHERE tc.status = 'xfail') AS xfail_count
        FROM test_cases tc
        JOIN test_suites ts ON ts.id = tc.suite_id
        WHERE ts.analysis_id = an.id
    ) tc_summary ON true
    WHERE an.codebase_id = c.id AND an.status = 'completed'
    ORDER BY an.created_at DESC
    LIMIT 1
) a ON true
WHERE c.last_viewed_at IS NOT NULL
  AND c.is_stale = false
  AND (
    sqlc.arg(view_filter)::text = 'all'
    OR (sqlc.arg(view_filter)::text = 'my' AND EXISTS(
        SELECT 1 FROM user_analysis_history uah
        WHERE uah.analysis_id = a.id AND uah.user_id = sqlc.arg(user_id)::uuid
    ))
    OR (sqlc.arg(view_filter)::text = 'community' AND NOT EXISTS(
        SELECT 1 FROM user_analysis_history uah
        WHERE uah.analysis_id = a.id AND uah.user_id = sqlc.arg(user_id)::uuid
    ))
  )
  AND (
    sqlc.arg(cursor_test_count)::int IS NULL
    OR (
      (sqlc.arg(sort_order)::text = 'desc' AND (a.total_tests, c.id) < (sqlc.arg(cursor_test_count)::int, sqlc.arg(cursor_id)::uuid))
      OR (sqlc.arg(sort_order)::text = 'asc' AND (a.total_tests, c.id) > (sqlc.arg(cursor_test_count)::int, sqlc.arg(cursor_id)::uuid))
    )
  )
ORDER BY
  CASE WHEN sqlc.arg(sort_order)::text = 'desc' THEN a.total_tests END DESC,
  CASE WHEN sqlc.arg(sort_order)::text = 'asc' THEN a.total_tests END ASC,
  CASE WHEN sqlc.arg(sort_order)::text = 'desc' THEN c.id END DESC,
  CASE WHEN sqlc.arg(sort_order)::text = 'asc' THEN c.id END ASC
LIMIT sqlc.arg(page_limit);
