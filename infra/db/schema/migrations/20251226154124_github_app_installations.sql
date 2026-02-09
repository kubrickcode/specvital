-- Create enum type "github_account_type"
CREATE TYPE "public"."github_account_type" AS ENUM ('organization', 'user');
-- Create "github_app_installations" table
CREATE TABLE "public"."github_app_installations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "installation_id" bigint NOT NULL,
  "account_type" "public"."github_account_type" NOT NULL,
  "account_id" bigint NOT NULL,
  "account_login" character varying(255) NOT NULL,
  "account_avatar_url" text NULL,
  "installer_user_id" uuid NULL,
  "suspended_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_github_app_installations_account" UNIQUE ("account_type", "account_id"),
  CONSTRAINT "uq_github_app_installations_installation_id" UNIQUE ("installation_id"),
  CONSTRAINT "fk_github_app_installations_installer" FOREIGN KEY ("installer_user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_github_app_installations_installer" to table: "github_app_installations"
CREATE INDEX "idx_github_app_installations_installer" ON "public"."github_app_installations" ("installer_user_id") WHERE (installer_user_id IS NOT NULL);
