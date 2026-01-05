-- Create "spec_view_cache" table
CREATE TABLE "public"."spec_view_cache" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "cache_key_hash" bytea NOT NULL,
  "codebase_id" uuid NOT NULL,
  "file_path" text NOT NULL,
  "framework" character varying(50) NOT NULL,
  "suite_hierarchy" text NOT NULL,
  "original_name" text NOT NULL,
  "converted_name" text NOT NULL,
  "language" character varying(10) NOT NULL DEFAULT 'en',
  "model_id" character varying(100) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_spec_view_cache_key_model" UNIQUE ("cache_key_hash", "model_id"),
  CONSTRAINT "fk_spec_view_cache_codebase" FOREIGN KEY ("codebase_id") REFERENCES "public"."codebases" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_spec_view_cache_codebase" to table: "spec_view_cache"
CREATE INDEX "idx_spec_view_cache_codebase" ON "public"."spec_view_cache" ("codebase_id");
-- Create index "idx_spec_view_cache_lookup" to table: "spec_view_cache"
CREATE INDEX "idx_spec_view_cache_lookup" ON "public"."spec_view_cache" ("cache_key_hash", "model_id");
