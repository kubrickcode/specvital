-- Create "behavior_caches" table
CREATE TABLE "public"."behavior_caches" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "cache_key_hash" bytea NOT NULL,
  "converted_description" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_behavior_caches_key" UNIQUE ("cache_key_hash")
);
-- Create index "idx_behavior_caches_created_at" to table: "behavior_caches"
CREATE INDEX "idx_behavior_caches_created_at" ON "public"."behavior_caches" ("created_at");
