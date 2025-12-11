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
