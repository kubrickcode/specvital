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

  primary_key {
    columns = [column.id]
  }

  unique "uq_codebases_identity" {
    columns = [column.host, column.owner, column.name]
  }

  index "idx_codebases_owner_name" {
    columns = [column.owner, column.name]
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
