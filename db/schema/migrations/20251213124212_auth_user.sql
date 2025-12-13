-- Create enum type "oauth_provider"
CREATE TYPE "public"."oauth_provider" AS ENUM ('github');
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "email" character varying(255) NULL,
  "username" character varying(255) NOT NULL,
  "avatar_url" text NULL,
  "last_login_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email") WHERE (email IS NOT NULL);
-- Create index "idx_users_username" to table: "users"
CREATE INDEX "idx_users_username" ON "public"."users" ("username");
-- Create "oauth_accounts" table
CREATE TABLE "public"."oauth_accounts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "provider" "public"."oauth_provider" NOT NULL,
  "provider_user_id" character varying(255) NOT NULL,
  "provider_username" character varying(255) NULL,
  "access_token" text NULL,
  "scope" character varying(500) NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_oauth_provider_user" UNIQUE ("provider", "provider_user_id"),
  CONSTRAINT "fk_oauth_accounts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_oauth_accounts_user_id" to table: "oauth_accounts"
CREATE INDEX "idx_oauth_accounts_user_id" ON "public"."oauth_accounts" ("user_id");
-- Create index "idx_oauth_accounts_user_provider" to table: "oauth_accounts"
CREATE INDEX "idx_oauth_accounts_user_provider" ON "public"."oauth_accounts" ("user_id", "provider");
