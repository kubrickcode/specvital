-- name: GetSystemConfigValue :one
SELECT value
FROM system_config
WHERE key = $1;
