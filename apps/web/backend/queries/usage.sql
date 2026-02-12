-- name: GetMonthlyUsage :one
SELECT
    COALESCE(SUM(quota_amount), 0)::bigint AS total
FROM usage_events
WHERE user_id = $1
    AND event_type = $2
    AND created_at >= $3
    AND created_at < $4;

-- name: GetUsageByPeriod :many
SELECT
    event_type,
    COALESCE(SUM(quota_amount), 0)::bigint AS total
FROM usage_events
WHERE user_id = $1
    AND created_at >= $2
    AND created_at < $3
GROUP BY event_type;
