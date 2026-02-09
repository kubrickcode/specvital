-- Modify "spec_documents" table
ALTER TABLE "public"."spec_documents" ADD CONSTRAINT "chk_retention_days_positive" CHECK ((retention_days_at_creation IS NULL) OR (retention_days_at_creation > 0)), ADD COLUMN "retention_days_at_creation" integer NULL;
-- Create index "idx_spec_documents_retention_cleanup" to table: "spec_documents"
CREATE INDEX "idx_spec_documents_retention_cleanup" ON "public"."spec_documents" ("created_at") WHERE (retention_days_at_creation IS NOT NULL);
-- Modify "user_analysis_history" table
ALTER TABLE "public"."user_analysis_history" ADD CONSTRAINT "chk_retention_days_positive" CHECK ((retention_days_at_creation IS NULL) OR (retention_days_at_creation > 0)), ADD COLUMN "retention_days_at_creation" integer NULL;
-- Create index "idx_user_analysis_history_retention_cleanup" to table: "user_analysis_history"
CREATE INDEX "idx_user_analysis_history_retention_cleanup" ON "public"."user_analysis_history" ("created_at") WHERE (retention_days_at_creation IS NOT NULL);
