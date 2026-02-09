-- Modify "codebases" table
ALTER TABLE "public"."codebases" ADD COLUMN "last_viewed_at" timestamptz NULL;
-- Create index "idx_codebases_last_viewed" to table: "codebases"
CREATE INDEX "idx_codebases_last_viewed" ON "public"."codebases" ("last_viewed_at") WHERE (last_viewed_at IS NOT NULL);
