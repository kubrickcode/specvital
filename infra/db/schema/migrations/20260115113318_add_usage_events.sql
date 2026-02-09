-- Create enum type "usage_event_type"
CREATE TYPE "public"."usage_event_type" AS ENUM ('specview', 'analysis');
-- Create "usage_events" table
CREATE TABLE "public"."usage_events" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "event_type" "public"."usage_event_type" NOT NULL,
  "analysis_id" uuid NULL,
  "document_id" uuid NULL,
  "quota_amount" integer NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_usage_events_analysis" FOREIGN KEY ("analysis_id") REFERENCES "public"."analyses" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_usage_events_document" FOREIGN KEY ("document_id") REFERENCES "public"."spec_documents" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_usage_events_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_usage_events_resource" CHECK ((((analysis_id IS NOT NULL))::integer + ((document_id IS NOT NULL))::integer) = 1)
);
-- Create index "idx_usage_events_analysis" to table: "usage_events"
CREATE INDEX "idx_usage_events_analysis" ON "public"."usage_events" ("analysis_id") WHERE (analysis_id IS NOT NULL);
-- Create index "idx_usage_events_document" to table: "usage_events"
CREATE INDEX "idx_usage_events_document" ON "public"."usage_events" ("document_id") WHERE (document_id IS NOT NULL);
-- Create index "idx_usage_events_quota_lookup" to table: "usage_events"
CREATE INDEX "idx_usage_events_quota_lookup" ON "public"."usage_events" ("user_id", "event_type", "created_at");
