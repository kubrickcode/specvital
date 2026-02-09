-- Create "quota_reservations" table
CREATE TABLE "public"."quota_reservations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "event_type" "public"."usage_event_type" NOT NULL,
  "reserved_amount" integer NOT NULL,
  "job_id" bigint NOT NULL,
  "expires_at" timestamptz NOT NULL DEFAULT (now() + '01:00:00'::interval),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_quota_reservations_job_id" UNIQUE ("job_id"),
  CONSTRAINT "fk_quota_reservations_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_reserved_amount_positive" CHECK (reserved_amount > 0)
);
-- Create index "idx_quota_reservations_expires" to table: "quota_reservations"
CREATE INDEX "idx_quota_reservations_expires" ON "public"."quota_reservations" ("expires_at");
-- Create index "idx_quota_reservations_user_event" to table: "quota_reservations"
CREATE INDEX "idx_quota_reservations_user_event" ON "public"."quota_reservations" ("user_id", "event_type");
