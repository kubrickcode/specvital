-- Modify "analyses" table
ALTER TABLE "public"."analyses" ADD COLUMN "committed_at" timestamptz NULL;
