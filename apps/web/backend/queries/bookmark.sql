-- name: AddBookmark :exec
INSERT INTO user_bookmarks (user_id, codebase_id)
VALUES ($1, $2)
ON CONFLICT (user_id, codebase_id) DO NOTHING;

-- name: RemoveBookmark :exec
DELETE FROM user_bookmarks
WHERE user_id = $1 AND codebase_id = $2;

-- name: IsBookmarked :one
SELECT EXISTS(
    SELECT 1 FROM user_bookmarks
    WHERE user_id = $1 AND codebase_id = $2
) AS is_bookmarked;

-- name: GetUserBookmarks :many
SELECT
    c.id AS codebase_id,
    c.owner,
    c.name,
    ub.created_at AS bookmarked_at,
    a.id AS analysis_id,
    a.commit_sha,
    a.completed_at AS analyzed_at,
    a.total_tests
FROM user_bookmarks ub
JOIN codebases c ON c.id = ub.codebase_id
LEFT JOIN LATERAL (
    SELECT id, commit_sha, completed_at, total_tests
    FROM analyses
    WHERE codebase_id = c.id AND status = 'completed'
    ORDER BY created_at DESC
    LIMIT 1
) a ON true
WHERE ub.user_id = $1 AND c.is_stale = false
ORDER BY ub.created_at DESC;

-- name: GetBookmarkedCodebaseIDsByUserID :many
SELECT codebase_id
FROM user_bookmarks
WHERE user_id = $1;

-- name: GetCodebaseByOwnerRepo :one
SELECT id
FROM codebases
WHERE host = $1 AND owner = $2 AND name = $3 AND is_stale = false;
