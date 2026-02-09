-- name: GetUserByID :one
SELECT
    id,
    email,
    username,
    avatar_url,
    last_login_at,
    created_at,
    updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT
    id,
    email,
    username,
    avatar_url,
    last_login_at,
    created_at,
    updated_at
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, username, avatar_url)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login_at = now(), updated_at = now()
WHERE id = $1;

-- name: GetOAuthAccountByProvider :one
SELECT
    oa.id,
    oa.user_id,
    oa.provider,
    oa.provider_user_id,
    oa.provider_username,
    oa.access_token,
    oa.scope,
    oa.created_at,
    oa.updated_at
FROM oauth_accounts oa
WHERE oa.provider = $1 AND oa.provider_user_id = $2;

-- name: GetOAuthAccountByUserID :one
SELECT
    oa.id,
    oa.user_id,
    oa.provider,
    oa.provider_user_id,
    oa.provider_username,
    oa.access_token,
    oa.scope,
    oa.created_at,
    oa.updated_at
FROM oauth_accounts oa
WHERE oa.user_id = $1 AND oa.provider = $2;

-- name: UpsertOAuthAccount :one
INSERT INTO oauth_accounts (user_id, provider, provider_user_id, provider_username, access_token, scope)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (provider, provider_user_id) DO UPDATE SET
    provider_username = EXCLUDED.provider_username,
    access_token = EXCLUDED.access_token,
    scope = EXCLUDED.scope,
    updated_at = now()
RETURNING id;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, family_id, expires_at, replaces)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, family_id, created_at;

-- name: GetRefreshTokenByHash :one
SELECT
    id,
    user_id,
    token_hash,
    family_id,
    expires_at,
    created_at,
    revoked_at,
    replaces
FROM refresh_tokens
WHERE token_hash = $1;

-- name: RevokeRefreshToken :execrows
UPDATE refresh_tokens
SET revoked_at = now()
WHERE id = $1 AND revoked_at IS NULL;

-- name: RevokeRefreshTokenFamily :exec
UPDATE refresh_tokens
SET revoked_at = now()
WHERE family_id = $1 AND revoked_at IS NULL;

-- name: RevokeUserRefreshTokens :exec
UPDATE refresh_tokens
SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < now() - INTERVAL '7 days';
