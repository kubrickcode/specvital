schema "public" {}

// ==============================================================================
// Enums
// ==============================================================================

enum "analysis_status" {
  schema = schema.public
  values = ["pending", "running", "completed", "failed"]
}

enum "test_status" {
  schema = schema.public
  values = ["active", "skipped", "todo", "focused", "xfail"]
}

enum "oauth_provider" {
  schema = schema.public
  values = ["github"]
}

enum "github_account_type" {
  schema = schema.public
  values = ["organization", "user"]
}

// ==============================================================================
// Tables
// ==============================================================================

table "codebases" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "host" {
    type    = varchar(255)
    default = "github.com"
  }

  column "owner" {
    type = varchar(255)
  }

  column "name" {
    type = varchar(255)
  }

  column "external_repo_id" {
    type = varchar(64)
  }

  column "is_stale" {
    type    = bool
    default = false
  }

  column "default_branch" {
    type = varchar(100)
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "last_viewed_at" {
    type = timestamptz
    null = true
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_codebases_identity" {
    columns = [column.host, column.owner, column.name]
    unique  = true
    where   = "is_stale = false"
  }

  index "idx_codebases_external_repo_id" {
    columns = [column.host, column.external_repo_id]
    unique  = true
  }

  index "idx_codebases_owner_name" {
    columns = [column.owner, column.name]
  }

  index "idx_codebases_last_viewed" {
    columns = [column.last_viewed_at]
    where   = "last_viewed_at IS NOT NULL"
  }
}

table "analyses" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "codebase_id" {
    type = uuid
  }

  column "commit_sha" {
    type = varchar(40)
  }

  column "branch_name" {
    type = varchar(255)
    null = true
  }

  column "status" {
    type    = enum.analysis_status
    default = "pending"
  }

  column "error_message" {
    type = text
    null = true
  }

  column "started_at" {
    type = timestamptz
    null = true
  }

  column "completed_at" {
    type = timestamptz
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "total_suites" {
    type    = int
    default = 0
  }

  column "total_tests" {
    type    = int
    default = 0
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_analyses_codebase" {
    columns     = [column.codebase_id]
    ref_columns = [table.codebases.column.id]
    on_delete   = CASCADE
  }

  index "uq_analyses_completed_commit" {
    columns = [column.codebase_id, column.commit_sha]
    unique  = true
    where   = "status = 'completed'"
  }

  index "idx_analyses_codebase_status" {
    columns = [column.codebase_id, column.status]
  }

  index "idx_analyses_created" {
    columns = [column.codebase_id, column.created_at]
  }
}

table "test_suites" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "analysis_id" {
    type = uuid
  }

  column "parent_id" {
    type = uuid
    null = true
  }

  column "name" {
    type = varchar(500)
  }

  column "file_path" {
    type = varchar(1000)
  }

  column "line_number" {
    type = int
    null = true
  }

  column "framework" {
    type = varchar(50)
    null = true
  }

  column "depth" {
    type    = int
    default = 0
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_test_suites_analysis" {
    columns     = [column.analysis_id]
    ref_columns = [table.analyses.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_test_suites_parent" {
    columns     = [column.parent_id]
    ref_columns = [table.test_suites.column.id]
    on_delete   = CASCADE
  }

  check "chk_no_self_reference" {
    expr = "id != parent_id"
  }

  index "idx_test_suites_analysis" {
    columns = [column.analysis_id]
  }

  index "idx_test_suites_parent" {
    columns = [column.parent_id]
    where   = "parent_id IS NOT NULL"
  }

  index "idx_test_suites_file" {
    columns = [column.analysis_id, column.file_path]
  }
}

table "test_cases" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "suite_id" {
    type = uuid
  }

  column "name" {
    type = varchar(2000)
  }

  column "line_number" {
    type = int
    null = true
  }

  column "status" {
    type    = enum.test_status
    default = "active"
  }

  column "modifier" {
    type = varchar(50)
    null = true
  }

  column "tags" {
    type    = jsonb
    default = "[]"
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_test_cases_suite" {
    columns     = [column.suite_id]
    ref_columns = [table.test_suites.column.id]
    on_delete   = CASCADE
  }

  index "idx_test_cases_suite" {
    columns = [column.suite_id]
  }

  index "idx_test_cases_status" {
    columns = [column.suite_id, column.status]
  }
}

// ==============================================================================
// Auth Tables
// ==============================================================================

table "users" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "email" {
    type = varchar(255)
    null = true
  }

  column "username" {
    type = varchar(255)
  }

  column "avatar_url" {
    type = text
    null = true
  }

  column "last_login_at" {
    type = timestamptz
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "token_version" {
    type    = int
    default = 1
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_users_email" {
    columns = [column.email]
    unique  = true
    where   = "email IS NOT NULL"
  }

  index "idx_users_username" {
    columns = [column.username]
  }
}

table "refresh_tokens" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    type = uuid
  }

  column "token_hash" {
    type = text
  }

  column "family_id" {
    type = uuid
  }

  column "expires_at" {
    type = timestamptz
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "revoked_at" {
    type = timestamptz
    null = true
  }

  column "replaces" {
    type = uuid
    null = true
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_refresh_tokens_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_refresh_tokens_replaces" {
    columns     = [column.replaces]
    ref_columns = [table.refresh_tokens.column.id]
    on_delete   = SET_NULL
  }

  unique "uq_refresh_tokens_hash" {
    columns = [column.token_hash]
  }

  index "idx_refresh_tokens_user" {
    columns = [column.user_id]
  }

  index "idx_refresh_tokens_family_active" {
    columns = [column.family_id, column.created_at]
    where   = "revoked_at IS NULL"
  }

  index "idx_refresh_tokens_expires" {
    columns = [column.expires_at]
    where   = "revoked_at IS NULL"
  }
}

table "oauth_accounts" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    type = uuid
  }

  column "provider" {
    type = enum.oauth_provider
  }

  column "provider_user_id" {
    type = varchar(255)
  }

  column "provider_username" {
    type = varchar(255)
    null = true
  }

  // Note: Encrypted at application level (AES-256-GCM)
  column "access_token" {
    type = text
    null = true
  }

  column "scope" {
    type = varchar(500)
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_oauth_accounts_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  unique "uq_oauth_provider_user" {
    columns = [column.provider, column.provider_user_id]
  }

  index "idx_oauth_accounts_user_id" {
    columns = [column.user_id]
  }

  index "idx_oauth_accounts_user_provider" {
    columns = [column.user_id, column.provider]
  }
}

table "user_bookmarks" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    type = uuid
  }

  column "codebase_id" {
    type = uuid
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_user_bookmarks_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_user_bookmarks_codebase" {
    columns     = [column.codebase_id]
    ref_columns = [table.codebases.column.id]
    on_delete   = CASCADE
  }

  unique "uq_user_bookmarks_user_codebase" {
    columns = [column.user_id, column.codebase_id]
  }

  index "idx_user_bookmarks_user" {
    columns = [column.user_id, column.created_at]
  }
}

table "user_analysis_history" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    type = uuid
  }

  column "analysis_id" {
    type = uuid
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_user_analysis_history_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_user_analysis_history_analysis" {
    columns     = [column.analysis_id]
    ref_columns = [table.analyses.column.id]
    on_delete   = CASCADE
  }

  unique "uq_user_analysis_history_user_analysis" {
    columns = [column.user_id, column.analysis_id]
  }

  index "idx_user_analysis_history_cursor" {
    columns = [column.user_id, column.updated_at, column.id]
  }

  index "idx_user_analysis_history_analysis" {
    columns = [column.analysis_id]
  }
}

// ==============================================================================
// GitHub App Tables
// ==============================================================================

table "github_app_installations" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "installation_id" {
    type = bigint
  }

  column "account_type" {
    type = enum.github_account_type
  }

  column "account_id" {
    type = bigint
  }

  column "account_login" {
    type = varchar(255)
  }

  column "account_avatar_url" {
    type = text
    null = true
  }

  column "installer_user_id" {
    type = uuid
    null = true
  }

  column "suspended_at" {
    type = timestamptz
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_github_app_installations_installer" {
    columns     = [column.installer_user_id]
    ref_columns = [table.users.column.id]
    on_delete   = SET_NULL
  }

  unique "uq_github_app_installations_installation_id" {
    columns = [column.installation_id]
  }

  unique "uq_github_app_installations_account" {
    columns = [column.account_type, column.account_id]
  }

  index "idx_github_app_installations_installer" {
    columns = [column.installer_user_id]
    where   = "installer_user_id IS NOT NULL"
  }
}

// ==============================================================================
// GitHub Cache Tables
// ==============================================================================

table "user_github_repositories" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    type = uuid
  }

  column "github_repo_id" {
    type = bigint
  }

  column "name" {
    type = varchar(255)
  }

  column "full_name" {
    type = varchar(500)
  }

  column "html_url" {
    type = text
  }

  column "description" {
    type = text
    null = true
  }

  column "default_branch" {
    type = varchar(100)
    null = true
  }

  column "language" {
    type = varchar(50)
    null = true
  }

  column "visibility" {
    type    = varchar(20)
    default = "public"
  }

  column "is_private" {
    type    = bool
    default = false
  }

  column "archived" {
    type    = bool
    default = false
  }

  column "disabled" {
    type    = bool
    default = false
  }

  column "fork" {
    type    = bool
    default = false
  }

  column "stargazers_count" {
    type    = int
    default = 0
  }

  column "pushed_at" {
    type = timestamptz
    null = true
  }

  column "source_type" {
    type    = varchar(20)
    default = "personal"
  }

  column "org_id" {
    type = uuid
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_user_github_repositories_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_user_github_repositories_org" {
    columns     = [column.org_id]
    ref_columns = [table.github_organizations.column.id]
    on_delete   = CASCADE
  }

  unique "uq_user_github_repositories_user_repo" {
    columns = [column.user_id, column.github_repo_id]
  }

  index "idx_user_github_repositories_user" {
    columns = [column.user_id, column.updated_at]
  }

  index "idx_user_github_repositories_language" {
    columns = [column.user_id, column.language]
    where   = "language IS NOT NULL"
  }

  index "idx_user_github_repositories_org" {
    columns = [column.user_id, column.org_id]
    where   = "org_id IS NOT NULL"
  }

  index "idx_user_github_repositories_source" {
    columns = [column.user_id, column.source_type]
  }
}

table "github_organizations" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "github_org_id" {
    type = bigint
  }

  column "login" {
    type = varchar(255)
  }

  column "avatar_url" {
    type = text
    null = true
  }

  column "html_url" {
    type = text
    null = true
  }

  column "description" {
    type = text
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  unique "uq_github_organizations_github_org_id" {
    columns = [column.github_org_id]
  }

  index "idx_github_organizations_login" {
    columns = [column.login]
  }
}

table "user_github_org_memberships" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    type = uuid
  }

  column "org_id" {
    type = uuid
  }

  column "role" {
    type = varchar(50)
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_user_github_org_memberships_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_user_github_org_memberships_org" {
    columns     = [column.org_id]
    ref_columns = [table.github_organizations.column.id]
    on_delete   = CASCADE
  }

  unique "uq_user_github_org_memberships_user_org" {
    columns = [column.user_id, column.org_id]
  }

  index "idx_user_github_org_memberships_user" {
    columns = [column.user_id]
  }

  index "idx_user_github_org_memberships_org" {
    columns = [column.org_id]
  }
}
