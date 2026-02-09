-- =============================================================================
-- User Personal Repositories
-- =============================================================================

-- name: UpsertUserGitHubRepository :exec
INSERT INTO user_github_repositories (
    user_id, github_repo_id, name, full_name, html_url,
    description, default_branch, language, visibility, is_private,
    archived, disabled, fork, stargazers_count, pushed_at,
    source_type, org_id, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
    'personal', NULL, now()
)
ON CONFLICT (user_id, github_repo_id) DO UPDATE SET
    name = EXCLUDED.name,
    full_name = EXCLUDED.full_name,
    html_url = EXCLUDED.html_url,
    description = EXCLUDED.description,
    default_branch = EXCLUDED.default_branch,
    language = EXCLUDED.language,
    visibility = EXCLUDED.visibility,
    is_private = EXCLUDED.is_private,
    archived = EXCLUDED.archived,
    disabled = EXCLUDED.disabled,
    fork = EXCLUDED.fork,
    stargazers_count = EXCLUDED.stargazers_count,
    pushed_at = EXCLUDED.pushed_at,
    source_type = 'personal',
    org_id = NULL,
    updated_at = now();

-- name: GetUserGitHubRepositories :many
SELECT
    github_repo_id, name, full_name, html_url, description,
    default_branch, language, visibility, is_private, archived,
    disabled, fork, stargazers_count, pushed_at
FROM user_github_repositories
WHERE user_id = $1 AND source_type = 'personal'
ORDER BY pushed_at DESC NULLS LAST, updated_at DESC;

-- name: DeleteUserGitHubRepositories :exec
DELETE FROM user_github_repositories
WHERE user_id = $1 AND source_type = 'personal';

-- name: HasUserGitHubRepositories :one
SELECT EXISTS(
    SELECT 1 FROM user_github_repositories
    WHERE user_id = $1 AND source_type = 'personal'
) AS has_repos;

-- =============================================================================
-- Organization Repositories
-- =============================================================================

-- name: UpsertOrgGitHubRepository :exec
INSERT INTO user_github_repositories (
    user_id, github_repo_id, name, full_name, html_url,
    description, default_branch, language, visibility, is_private,
    archived, disabled, fork, stargazers_count, pushed_at,
    source_type, org_id, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
    'organization', $16, now()
)
ON CONFLICT (user_id, github_repo_id) DO UPDATE SET
    name = EXCLUDED.name,
    full_name = EXCLUDED.full_name,
    html_url = EXCLUDED.html_url,
    description = EXCLUDED.description,
    default_branch = EXCLUDED.default_branch,
    language = EXCLUDED.language,
    visibility = EXCLUDED.visibility,
    is_private = EXCLUDED.is_private,
    archived = EXCLUDED.archived,
    disabled = EXCLUDED.disabled,
    fork = EXCLUDED.fork,
    stargazers_count = EXCLUDED.stargazers_count,
    pushed_at = EXCLUDED.pushed_at,
    source_type = 'organization',
    org_id = EXCLUDED.org_id,
    updated_at = now();

-- name: GetOrgGitHubRepositories :many
SELECT
    github_repo_id, name, full_name, html_url, description,
    default_branch, language, visibility, is_private, archived,
    disabled, fork, stargazers_count, pushed_at
FROM user_github_repositories
WHERE user_id = $1 AND org_id = $2
ORDER BY pushed_at DESC NULLS LAST, updated_at DESC;

-- name: DeleteOrgGitHubRepositories :exec
DELETE FROM user_github_repositories
WHERE user_id = $1 AND org_id = $2;

-- name: HasOrgGitHubRepositories :one
SELECT EXISTS(
    SELECT 1 FROM user_github_repositories
    WHERE user_id = $1 AND org_id = $2
) AS has_repos;

-- =============================================================================
-- GitHub Organizations (Global)
-- =============================================================================

-- name: UpsertGitHubOrganization :one
INSERT INTO github_organizations (
    github_org_id, login, avatar_url, html_url, description, updated_at
) VALUES ($1, $2, $3, $4, $5, now())
ON CONFLICT (github_org_id) DO UPDATE SET
    login = EXCLUDED.login,
    avatar_url = EXCLUDED.avatar_url,
    html_url = EXCLUDED.html_url,
    description = EXCLUDED.description,
    updated_at = now()
RETURNING id;

-- name: GetGitHubOrganizationByLogin :one
SELECT id, github_org_id FROM github_organizations WHERE login = $1;

-- =============================================================================
-- User Organization Memberships
-- =============================================================================

-- name: UpsertUserOrgMembership :exec
INSERT INTO user_github_org_memberships (user_id, org_id, role, updated_at)
VALUES ($1, $2, $3, now())
ON CONFLICT (user_id, org_id) DO UPDATE SET
    role = EXCLUDED.role,
    updated_at = now();

-- name: GetUserGitHubOrganizations :many
SELECT
    go.id AS org_id,
    go.github_org_id,
    go.login,
    go.avatar_url,
    go.html_url,
    go.description,
    m.role
FROM github_organizations go
JOIN user_github_org_memberships m ON go.id = m.org_id
WHERE m.user_id = $1
ORDER BY go.login;

-- name: DeleteUserOrgMemberships :exec
DELETE FROM user_github_org_memberships WHERE user_id = $1;

-- name: HasUserOrgMemberships :one
SELECT EXISTS(
    SELECT 1 FROM user_github_org_memberships WHERE user_id = $1
) AS has_orgs;
