-- Create index "idx_spec_documents_content_hash_lang_model" to table: "spec_documents"
CREATE INDEX "idx_spec_documents_content_hash_lang_model" ON "public"."spec_documents" ("content_hash", "language", "model_id");
