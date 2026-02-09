-- Create enum type "analysis_status"
CREATE TYPE "public"."analysis_status" AS ENUM ('pending', 'running', 'completed', 'failed');
-- Create enum type "test_status"
CREATE TYPE "public"."test_status" AS ENUM ('active', 'skipped', 'todo');
-- Create "codebases" table
CREATE TABLE "public"."codebases" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "host" character varying(255) NOT NULL DEFAULT 'github.com',
  "owner" character varying(255) NOT NULL,
  "name" character varying(255) NOT NULL,
  "default_branch" character varying(100) NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_codebases_identity" UNIQUE ("host", "owner", "name")
);
-- Create index "idx_codebases_owner_name" to table: "codebases"
CREATE INDEX "idx_codebases_owner_name" ON "public"."codebases" ("owner", "name");
-- Create "analyses" table
CREATE TABLE "public"."analyses" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "codebase_id" uuid NOT NULL,
  "commit_sha" character varying(40) NOT NULL,
  "branch_name" character varying(255) NULL,
  "status" "public"."analysis_status" NOT NULL DEFAULT 'pending',
  "error_message" text NULL,
  "started_at" timestamptz NULL,
  "completed_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "total_suites" integer NOT NULL DEFAULT 0,
  "total_tests" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_analyses_commit" UNIQUE ("codebase_id", "commit_sha"),
  CONSTRAINT "fk_analyses_codebase" FOREIGN KEY ("codebase_id") REFERENCES "public"."codebases" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_analyses_codebase_status" to table: "analyses"
CREATE INDEX "idx_analyses_codebase_status" ON "public"."analyses" ("codebase_id", "status");
-- Create index "idx_analyses_created" to table: "analyses"
CREATE INDEX "idx_analyses_created" ON "public"."analyses" ("codebase_id", "created_at");
-- Create "test_suites" table
CREATE TABLE "public"."test_suites" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "analysis_id" uuid NOT NULL,
  "parent_id" uuid NULL,
  "name" character varying(500) NOT NULL,
  "file_path" character varying(1000) NOT NULL,
  "line_number" integer NULL,
  "framework" character varying(50) NULL,
  "depth" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_test_suites_analysis" FOREIGN KEY ("analysis_id") REFERENCES "public"."analyses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_test_suites_parent" FOREIGN KEY ("parent_id") REFERENCES "public"."test_suites" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_no_self_reference" CHECK (id <> parent_id)
);
-- Create index "idx_test_suites_analysis" to table: "test_suites"
CREATE INDEX "idx_test_suites_analysis" ON "public"."test_suites" ("analysis_id");
-- Create index "idx_test_suites_file" to table: "test_suites"
CREATE INDEX "idx_test_suites_file" ON "public"."test_suites" ("analysis_id", "file_path");
-- Create index "idx_test_suites_parent" to table: "test_suites"
CREATE INDEX "idx_test_suites_parent" ON "public"."test_suites" ("parent_id") WHERE (parent_id IS NOT NULL);
-- Create "test_cases" table
CREATE TABLE "public"."test_cases" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "suite_id" uuid NOT NULL,
  "name" character varying(500) NOT NULL,
  "line_number" integer NULL,
  "status" "public"."test_status" NOT NULL DEFAULT 'active',
  "tags" jsonb NOT NULL DEFAULT '[]',
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_test_cases_suite" FOREIGN KEY ("suite_id") REFERENCES "public"."test_suites" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_test_cases_status" to table: "test_cases"
CREATE INDEX "idx_test_cases_status" ON "public"."test_cases" ("suite_id", "status");
-- Create index "idx_test_cases_suite" to table: "test_cases"
CREATE INDEX "idx_test_cases_suite" ON "public"."test_cases" ("suite_id");
