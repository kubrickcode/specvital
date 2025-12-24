-- Drop index "idx_user_analysis_history_user" from table: "user_analysis_history"
DROP INDEX "public"."idx_user_analysis_history_user";
-- Create index "idx_user_analysis_history_cursor" to table: "user_analysis_history"
CREATE INDEX "idx_user_analysis_history_cursor" ON "public"."user_analysis_history" ("user_id", "updated_at", "id");
