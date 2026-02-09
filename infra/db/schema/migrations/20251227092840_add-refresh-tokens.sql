-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "token_version" integer NOT NULL DEFAULT 1;
-- Create "refresh_tokens" table
CREATE TABLE "public"."refresh_tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "token_hash" text NOT NULL,
  "family_id" uuid NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "revoked_at" timestamptz NULL,
  "replaces" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_refresh_tokens_hash" UNIQUE ("token_hash"),
  CONSTRAINT "fk_refresh_tokens_replaces" FOREIGN KEY ("replaces") REFERENCES "public"."refresh_tokens" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_refresh_tokens_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_refresh_tokens_expires" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_expires" ON "public"."refresh_tokens" ("expires_at") WHERE (revoked_at IS NULL);
-- Create index "idx_refresh_tokens_family_active" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_family_active" ON "public"."refresh_tokens" ("family_id", "created_at") WHERE (revoked_at IS NULL);
-- Create index "idx_refresh_tokens_user" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_user" ON "public"."refresh_tokens" ("user_id");
