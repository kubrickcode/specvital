-- Modify "codebases" table
ALTER TABLE "public"."codebases" ADD COLUMN "is_private" boolean NOT NULL DEFAULT false;
-- Create index "idx_codebases_public" to table: "codebases"
CREATE INDEX "idx_codebases_public" ON "public"."codebases" ("is_private") WHERE (is_private = false);
