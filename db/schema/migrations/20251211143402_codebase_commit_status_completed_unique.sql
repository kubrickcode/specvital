-- Modify "analyses" table
ALTER TABLE "public"."analyses" DROP CONSTRAINT "uq_analyses_commit";
-- Create index "uq_analyses_completed_commit" to table: "analyses"
CREATE UNIQUE INDEX "uq_analyses_completed_commit" ON "public"."analyses" ("codebase_id", "commit_sha") WHERE (status = 'completed'::public.analysis_status);
