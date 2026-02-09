-- Modify "codebases" table
ALTER TABLE "public"."codebases" DROP CONSTRAINT "uq_codebases_identity", ADD COLUMN "external_repo_id" character varying(64) NOT NULL, ADD COLUMN "is_stale" boolean NOT NULL DEFAULT false;
-- Create index "idx_codebases_external_repo_id" to table: "codebases"
CREATE UNIQUE INDEX "idx_codebases_external_repo_id" ON "public"."codebases" ("host", "external_repo_id");
-- Create index "idx_codebases_identity" to table: "codebases"
CREATE UNIQUE INDEX "idx_codebases_identity" ON "public"."codebases" ("host", "owner", "name") WHERE (is_stale = false);
