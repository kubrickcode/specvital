-- name: UpsertGitHubAppInstallation :exec
INSERT INTO github_app_installations (
    installation_id,
    account_type,
    account_id,
    account_login,
    account_avatar_url,
    installer_user_id,
    suspended_at,
    created_at,
    updated_at
) VALUES (
    @installation_id,
    @account_type,
    @account_id,
    @account_login,
    @account_avatar_url,
    @installer_user_id,
    @suspended_at,
    now(),
    now()
)
ON CONFLICT (installation_id) DO UPDATE SET
    account_type = EXCLUDED.account_type,
    account_login = EXCLUDED.account_login,
    account_avatar_url = EXCLUDED.account_avatar_url,
    suspended_at = EXCLUDED.suspended_at,
    updated_at = now();

-- name: GetGitHubAppInstallationByID :one
SELECT * FROM github_app_installations
WHERE installation_id = @installation_id;

-- name: GetGitHubAppInstallationByAccountID :one
SELECT * FROM github_app_installations
WHERE account_id = @account_id;

-- name: ListGitHubAppInstallationsByUserID :many
SELECT * FROM github_app_installations
WHERE installer_user_id = @user_id
ORDER BY created_at DESC;

-- name: DeleteGitHubAppInstallation :exec
DELETE FROM github_app_installations
WHERE installation_id = @installation_id;

-- name: UpdateGitHubAppInstallationSuspended :exec
UPDATE github_app_installations
SET suspended_at = @suspended_at, updated_at = now()
WHERE installation_id = @installation_id;

-- name: ListGitHubAppInstallationsByAccountIDs :many
SELECT * FROM github_app_installations
WHERE account_id = ANY(@account_ids::bigint[])
ORDER BY account_id;

-- name: ListGitHubAppOrganizationInstallations :many
SELECT * FROM github_app_installations
WHERE account_type = 'organization'
ORDER BY account_login;
