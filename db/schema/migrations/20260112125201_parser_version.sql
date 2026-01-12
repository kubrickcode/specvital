-- Drop index "uq_analyses_completed_commit" from table: "analyses"
DROP INDEX "public"."uq_analyses_completed_commit";
-- Modify "analyses" table
ALTER TABLE "public"."analyses" ADD COLUMN "parser_version" character varying(100) NOT NULL DEFAULT 'legacy';
-- Create index "uq_analyses_completed_commit_version" to table: "analyses"
CREATE UNIQUE INDEX "uq_analyses_completed_commit_version" ON "public"."analyses" ("codebase_id", "commit_sha", "parser_version") WHERE (status = 'completed'::public.analysis_status);
-- Create "system_config" table
CREATE TABLE "public"."system_config" (
  "key" character varying(100) NOT NULL,
  "value" text NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("key")
);
