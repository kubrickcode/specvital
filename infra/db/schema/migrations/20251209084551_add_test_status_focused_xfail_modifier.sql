-- Add value to enum type: "test_status"
ALTER TYPE "public"."test_status" ADD VALUE 'focused';
-- Add value to enum type: "test_status"
ALTER TYPE "public"."test_status" ADD VALUE 'xfail';
-- Modify "test_cases" table
ALTER TABLE "public"."test_cases" ADD COLUMN "modifier" character varying(50) NULL;
