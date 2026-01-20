-- Modify "spec_documents" table
ALTER TABLE "public"."spec_documents" ADD COLUMN "version" integer NOT NULL DEFAULT 1;
-- Create index "idx_spec_documents_latest_version" to table: "spec_documents"
CREATE INDEX "idx_spec_documents_latest_version" ON "public"."spec_documents" ("analysis_id", "language", "version");
