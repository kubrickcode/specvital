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

enum "usage_event_type" {
  schema = schema.public
  values = ["specview", "analysis"]
}

enum "plan_tier" {
  schema = schema.public
  values = ["free", "pro", "pro_plus", "enterprise"]
}

enum "subscription_status" {
  schema = schema.public
  values = ["active", "canceled", "expired"]
}

// ==============================================================================
// System Config
// ==============================================================================

table "system_config" {
  schema = schema.public

  column "key" {
    type = varchar(100)
  }

  column "value" {
    type = text
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.key]
  }
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

  column "is_private" {
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

  index "idx_codebases_public" {
    columns = [column.is_private]
    where   = "is_private = false"
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

  column "committed_at" {
    type = timestamptz
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

  column "parser_version" {
    type    = varchar(100)
    default = "legacy"
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_analyses_codebase" {
    columns     = [column.codebase_id]
    ref_columns = [table.codebases.column.id]
    on_delete   = CASCADE
  }

  index "uq_analyses_completed_commit_version" {
    columns = [column.codebase_id, column.commit_sha, column.parser_version]
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

table "test_files" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "analysis_id" {
    type = uuid
  }

  column "file_path" {
    type = varchar(1000)
  }

  column "framework" {
    type = varchar(50)
    null = true
  }

  column "domain_hints" {
    type = jsonb
    null = true
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_test_files_analysis" {
    columns     = [column.analysis_id]
    ref_columns = [table.analyses.column.id]
    on_delete   = CASCADE
  }

  unique "uq_test_files_analysis_path" {
    columns = [column.analysis_id, column.file_path]
  }

  index "idx_test_files_analysis" {
    columns = [column.analysis_id]
  }
}

table "test_suites" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "file_id" {
    type = uuid
  }

  column "parent_id" {
    type = uuid
    null = true
  }

  column "name" {
    type = varchar(500)
  }

  column "line_number" {
    type = int
    null = true
  }

  column "depth" {
    type    = int
    default = 0
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_test_suites_file" {
    columns     = [column.file_id]
    ref_columns = [table.test_files.column.id]
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

  index "idx_test_suites_file" {
    columns = [column.file_id]
  }

  index "idx_test_suites_parent" {
    columns = [column.parent_id]
    where   = "parent_id IS NOT NULL"
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

// ==============================================================================
// Spec View - Document-based Specification Schema
// Hierarchy: spec_documents -> spec_domains -> spec_features -> spec_behaviors
// ==============================================================================

table "spec_documents" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "analysis_id" {
    type = uuid
  }

  column "content_hash" {
    type = bytea
  }

  column "language" {
    type    = varchar(10)
    default = "en"
  }

  column "executive_summary" {
    type = text
    null = true
  }

  column "model_id" {
    type = varchar(100)
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

  foreign_key "fk_spec_documents_analysis" {
    columns     = [column.analysis_id]
    ref_columns = [table.analyses.column.id]
    on_delete   = CASCADE
  }

  unique "uq_spec_documents_hash_lang_model" {
    columns = [column.content_hash, column.language, column.model_id]
  }

  index "idx_spec_documents_analysis" {
    columns = [column.analysis_id]
  }
}

table "spec_domains" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "document_id" {
    type = uuid
  }

  column "name" {
    type = varchar(255)
  }

  column "description" {
    type = text
    null = true
  }

  column "sort_order" {
    type    = int
    default = 0
  }

  column "classification_confidence" {
    type = decimal(3, 2)
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

  foreign_key "fk_spec_domains_document" {
    columns     = [column.document_id]
    ref_columns = [table.spec_documents.column.id]
    on_delete   = CASCADE
  }

  index "idx_spec_domains_document_sort" {
    columns = [column.document_id, column.sort_order]
  }
}

table "spec_features" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "domain_id" {
    type = uuid
  }

  column "name" {
    type = varchar(255)
  }

  column "description" {
    type = text
    null = true
  }

  column "sort_order" {
    type    = int
    default = 0
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

  foreign_key "fk_spec_features_domain" {
    columns     = [column.domain_id]
    ref_columns = [table.spec_domains.column.id]
    on_delete   = CASCADE
  }

  index "idx_spec_features_domain_sort" {
    columns = [column.domain_id, column.sort_order]
  }
}

table "spec_behaviors" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "feature_id" {
    type = uuid
  }

  column "source_test_case_id" {
    type = uuid
    null = true
  }

  column "original_name" {
    type = varchar(2000)
  }

  column "converted_description" {
    type = text
  }

  column "sort_order" {
    type    = int
    default = 0
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_spec_behaviors_feature" {
    columns     = [column.feature_id]
    ref_columns = [table.spec_features.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_spec_behaviors_test_case" {
    columns     = [column.source_test_case_id]
    ref_columns = [table.test_cases.column.id]
    on_delete   = SET_NULL
  }

  index "idx_spec_behaviors_feature_sort" {
    columns = [column.feature_id, column.sort_order]
  }

  index "idx_spec_behaviors_source" {
    columns = [column.source_test_case_id]
    where   = "source_test_case_id IS NOT NULL"
  }
}

// ==============================================================================
// User History Tables
// ==============================================================================

table "user_specview_history" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    type = uuid
  }

  column "document_id" {
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

  foreign_key "fk_user_specview_history_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_user_specview_history_document" {
    columns     = [column.document_id]
    ref_columns = [table.spec_documents.column.id]
    on_delete   = CASCADE
  }

  unique "uq_user_specview_history_user_document" {
    columns = [column.user_id, column.document_id]
  }

  index "idx_user_specview_history_cursor" {
    columns = [column.user_id, column.updated_at, column.id]
  }

  index "idx_user_specview_history_document" {
    columns = [column.document_id]
  }
}

// ==============================================================================
// Usage Tracking Tables
// ==============================================================================

table "usage_events" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    type = uuid
  }

  column "event_type" {
    type = enum.usage_event_type
  }

  column "analysis_id" {
    type = uuid
    null = true
  }

  column "document_id" {
    type = uuid
    null = true
  }

  column "quota_amount" {
    type = int
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_usage_events_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_usage_events_analysis" {
    columns     = [column.analysis_id]
    ref_columns = [table.analyses.column.id]
    on_delete   = SET_NULL
  }

  foreign_key "fk_usage_events_document" {
    columns     = [column.document_id]
    ref_columns = [table.spec_documents.column.id]
    on_delete   = SET_NULL
  }

  check "chk_usage_events_resource" {
    expr = "(analysis_id IS NOT NULL)::int + (document_id IS NOT NULL)::int = 1"
  }

  index "idx_usage_events_quota_lookup" {
    columns = [column.user_id, column.event_type, column.created_at]
  }

  index "idx_usage_events_analysis" {
    columns = [column.analysis_id]
    where   = "analysis_id IS NOT NULL"
  }

  index "idx_usage_events_document" {
    columns = [column.document_id]
    where   = "document_id IS NOT NULL"
  }
}

// ==============================================================================
// Subscription Tables
// ==============================================================================

table "subscription_plans" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "tier" {
    type = enum.plan_tier
  }

  column "specview_monthly_limit" {
    type = int
    null = true
  }

  column "analysis_monthly_limit" {
    type = int
    null = true
  }

  column "retention_days" {
    type = int
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  unique "uq_subscription_plans_tier" {
    columns = [column.tier]
  }
}

table "user_subscriptions" {
  schema = schema.public

  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    type = uuid
  }

  column "plan_id" {
    type = uuid
  }

  column "status" {
    type    = enum.subscription_status
    default = "active"
  }

  column "current_period_start" {
    type = timestamptz
  }

  column "current_period_end" {
    type = timestamptz
  }

  column "canceled_at" {
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

  foreign_key "fk_user_subscriptions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "fk_user_subscriptions_plan" {
    columns     = [column.plan_id]
    ref_columns = [table.subscription_plans.column.id]
    on_delete   = RESTRICT
  }

  check "chk_canceled_at_status" {
    expr = "(status = 'canceled') = (canceled_at IS NOT NULL)"
  }

  index "idx_user_subscriptions_active" {
    columns = [column.user_id]
    unique  = true
    where   = "status = 'active'"
  }

  index "idx_user_subscriptions_plan" {
    columns = [column.plan_id]
  }

  index "idx_user_subscriptions_period_end" {
    columns = [column.current_period_end]
    where   = "status = 'active'"
  }
}

// ==============================================================================
// GitHub Cache Tables
// ==============================================================================

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
