-- Modify "spec_documents" table
ALTER TABLE "public"."spec_documents" DROP CONSTRAINT "uq_spec_documents_analysis_lang_version", DROP CONSTRAINT "uq_spec_documents_hash_lang_model_version", ADD COLUMN "user_id" uuid NOT NULL, ADD CONSTRAINT "uq_spec_documents_user_analysis_lang_version" UNIQUE ("user_id", "analysis_id", "language", "version"), ADD CONSTRAINT "uq_spec_documents_user_hash_lang_model_version" UNIQUE ("user_id", "content_hash", "language", "model_id", "version"), ADD CONSTRAINT "fk_spec_documents_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Create index "idx_spec_documents_user_created" to table: "spec_documents"
CREATE INDEX "idx_spec_documents_user_created" ON "public"."spec_documents" ("user_id", "created_at");
