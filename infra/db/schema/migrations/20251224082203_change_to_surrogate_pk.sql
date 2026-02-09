-- Modify "user_analysis_history" table
ALTER TABLE "public"."user_analysis_history" DROP CONSTRAINT "user_analysis_history_pkey", ADD COLUMN "id" uuid NOT NULL DEFAULT gen_random_uuid(), ADD PRIMARY KEY ("id"), ADD CONSTRAINT "uq_user_analysis_history_user_analysis" UNIQUE ("user_id", "analysis_id");
-- Modify "user_bookmarks" table
ALTER TABLE "public"."user_bookmarks" DROP CONSTRAINT "user_bookmarks_pkey", ADD COLUMN "id" uuid NOT NULL DEFAULT gen_random_uuid(), ADD PRIMARY KEY ("id"), ADD CONSTRAINT "uq_user_bookmarks_user_codebase" UNIQUE ("user_id", "codebase_id");
