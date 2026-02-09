-- Create "classification_caches" table
CREATE TABLE "public"."classification_caches" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "content_hash" bytea NOT NULL,
  "language" character varying(10) NOT NULL,
  "model_id" character varying(100) NOT NULL,
  "phase1_output" jsonb NOT NULL,
  "test_index_map" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_classification_caches_key" UNIQUE ("content_hash", "language", "model_id")
);
-- Create index "idx_classification_caches_created_at" to table: "classification_caches"
CREATE INDEX "idx_classification_caches_created_at" ON "public"."classification_caches" ("created_at");
