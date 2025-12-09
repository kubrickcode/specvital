-- name: GetLatestCompletedAnalysis :one
SELECT
    a.id,
    a.commit_sha,
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
VALUES ($1, '', 'pending')
RETURNING id;

-- name: MarkAnalysisFailed :exec
UPDATE analyses
SET status = 'failed', error_message = $2
WHERE id = $1;

-- name: GetTestSuitesByAnalysisID :many
SELECT
    ts.id,
    ts.file_path,
    ts.framework
FROM test_suites ts
WHERE ts.analysis_id = $1
ORDER BY ts.file_path;

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
