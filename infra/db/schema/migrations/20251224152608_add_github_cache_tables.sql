-- Create "github_organizations" table
CREATE TABLE "public"."github_organizations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "github_org_id" bigint NOT NULL,
  "login" character varying(255) NOT NULL,
  "avatar_url" text NULL,
  "html_url" text NULL,
  "description" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_github_organizations_github_org_id" UNIQUE ("github_org_id")
);
-- Create index "idx_github_organizations_login" to table: "github_organizations"
CREATE INDEX "idx_github_organizations_login" ON "public"."github_organizations" ("login");
-- Create "user_github_org_memberships" table
CREATE TABLE "public"."user_github_org_memberships" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "org_id" uuid NOT NULL,
  "role" character varying(50) NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_user_github_org_memberships_user_org" UNIQUE ("user_id", "org_id"),
  CONSTRAINT "fk_user_github_org_memberships_org" FOREIGN KEY ("org_id") REFERENCES "public"."github_organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_user_github_org_memberships_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_user_github_org_memberships_org" to table: "user_github_org_memberships"
CREATE INDEX "idx_user_github_org_memberships_org" ON "public"."user_github_org_memberships" ("org_id");
-- Create index "idx_user_github_org_memberships_user" to table: "user_github_org_memberships"
CREATE INDEX "idx_user_github_org_memberships_user" ON "public"."user_github_org_memberships" ("user_id");
-- Create "user_github_repositories" table
CREATE TABLE "public"."user_github_repositories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "github_repo_id" bigint NOT NULL,
  "name" character varying(255) NOT NULL,
  "full_name" character varying(500) NOT NULL,
  "html_url" text NOT NULL,
  "description" text NULL,
  "default_branch" character varying(100) NULL,
  "language" character varying(50) NULL,
  "visibility" character varying(20) NOT NULL DEFAULT 'public',
  "is_private" boolean NOT NULL DEFAULT false,
  "archived" boolean NOT NULL DEFAULT false,
  "disabled" boolean NOT NULL DEFAULT false,
  "fork" boolean NOT NULL DEFAULT false,
  "stargazers_count" integer NOT NULL DEFAULT 0,
  "pushed_at" timestamptz NULL,
  "source_type" character varying(20) NOT NULL DEFAULT 'personal',
  "org_id" uuid NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_user_github_repositories_user_repo" UNIQUE ("user_id", "github_repo_id"),
  CONSTRAINT "fk_user_github_repositories_org" FOREIGN KEY ("org_id") REFERENCES "public"."github_organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_user_github_repositories_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_user_github_repositories_language" to table: "user_github_repositories"
CREATE INDEX "idx_user_github_repositories_language" ON "public"."user_github_repositories" ("user_id", "language") WHERE (language IS NOT NULL);
-- Create index "idx_user_github_repositories_org" to table: "user_github_repositories"
CREATE INDEX "idx_user_github_repositories_org" ON "public"."user_github_repositories" ("user_id", "org_id") WHERE (org_id IS NOT NULL);
-- Create index "idx_user_github_repositories_source" to table: "user_github_repositories"
CREATE INDEX "idx_user_github_repositories_source" ON "public"."user_github_repositories" ("user_id", "source_type");
-- Create index "idx_user_github_repositories_user" to table: "user_github_repositories"
CREATE INDEX "idx_user_github_repositories_user" ON "public"."user_github_repositories" ("user_id", "updated_at");
