-- name: CreateQuotaReservation :one
INSERT INTO quota_reservations (user_id, event_type, reserved_amount, job_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTotalReservedAmount :one
SELECT
    COALESCE(SUM(reserved_amount), 0)::bigint AS total
FROM quota_reservations
WHERE user_id = $1
    AND event_type = $2
    AND expires_at > now();

-- name: DeleteQuotaReservationByJobID :exec
DELETE FROM quota_reservations
WHERE job_id = $1;
