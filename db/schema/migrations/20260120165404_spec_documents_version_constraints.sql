-- Drop index "idx_spec_documents_latest_version" from table: "spec_documents"
DROP INDEX "public"."idx_spec_documents_latest_version";
-- Modify "spec_documents" table
ALTER TABLE "public"."spec_documents" DROP CONSTRAINT "uq_spec_documents_hash_lang_model", ADD CONSTRAINT "uq_spec_documents_analysis_lang_version" UNIQUE ("analysis_id", "language", "version"), ADD CONSTRAINT "uq_spec_documents_hash_lang_model_version" UNIQUE ("content_hash", "language", "model_id", "version");
