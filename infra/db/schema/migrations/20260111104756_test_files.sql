-- ==============================================================================
-- Migration: Add test_files table (schema normalization)
-- ==============================================================================

-- 1. Create test_files table
CREATE TABLE "public"."test_files" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "analysis_id" uuid NOT NULL,
  "file_path" character varying(1000) NOT NULL,
  "framework" character varying(50) NULL,
  "domain_hints" jsonb NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_test_files_analysis_path" UNIQUE ("analysis_id", "file_path"),
  CONSTRAINT "fk_test_files_analysis" FOREIGN KEY ("analysis_id") REFERENCES "public"."analyses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "idx_test_files_analysis" ON "public"."test_files" ("analysis_id");

-- 2. Migrate existing data from test_suites to test_files
INSERT INTO "public"."test_files" ("analysis_id", "file_path", "framework")
SELECT DISTINCT "analysis_id", "file_path", "framework"
FROM "public"."test_suites";

-- 3. Add file_id column to test_suites
ALTER TABLE "public"."test_suites" ADD COLUMN "file_id" uuid;

-- 4. Populate file_id by joining with test_files
UPDATE "public"."test_suites" ts
SET "file_id" = tf."id"
FROM "public"."test_files" tf
WHERE ts."analysis_id" = tf."analysis_id"
  AND ts."file_path" = tf."file_path";

-- 5. Add NOT NULL constraint and FK
ALTER TABLE "public"."test_suites" ALTER COLUMN "file_id" SET NOT NULL;

ALTER TABLE "public"."test_suites"
  ADD CONSTRAINT "fk_test_suites_file"
  FOREIGN KEY ("file_id") REFERENCES "public"."test_files" ("id")
  ON UPDATE NO ACTION ON DELETE CASCADE;

-- 6. Drop old columns and constraints
ALTER TABLE "public"."test_suites" DROP CONSTRAINT "fk_test_suites_analysis";

DROP INDEX "public"."idx_test_suites_analysis";
DROP INDEX "public"."idx_test_suites_file";

ALTER TABLE "public"."test_suites"
  DROP COLUMN "analysis_id",
  DROP COLUMN "file_path",
  DROP COLUMN "framework";

-- 7. Create new index
CREATE INDEX "idx_test_suites_file" ON "public"."test_suites" ("file_id");
