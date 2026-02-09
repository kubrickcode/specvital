-- Create enum type "subscription_status"
CREATE TYPE "public"."subscription_status" AS ENUM ('active', 'canceled', 'expired');
-- Create enum type "plan_tier"
CREATE TYPE "public"."plan_tier" AS ENUM ('free', 'pro', 'pro_plus', 'enterprise');
-- Create "subscription_plans" table
CREATE TABLE "public"."subscription_plans" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tier" "public"."plan_tier" NOT NULL,
  "specview_monthly_limit" integer NULL,
  "analysis_monthly_limit" integer NULL,
  "retention_days" integer NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_subscription_plans_tier" UNIQUE ("tier")
);
-- Create "user_subscriptions" table
CREATE TABLE "public"."user_subscriptions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "plan_id" uuid NOT NULL,
  "status" "public"."subscription_status" NOT NULL DEFAULT 'active',
  "current_period_start" timestamptz NOT NULL,
  "current_period_end" timestamptz NOT NULL,
  "canceled_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_subscriptions_plan" FOREIGN KEY ("plan_id") REFERENCES "public"."subscription_plans" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_user_subscriptions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_canceled_at_status" CHECK ((status = 'canceled'::public.subscription_status) = (canceled_at IS NOT NULL))
);
-- Create index "idx_user_subscriptions_active" to table: "user_subscriptions"
CREATE UNIQUE INDEX "idx_user_subscriptions_active" ON "public"."user_subscriptions" ("user_id") WHERE (status = 'active'::public.subscription_status);
-- Create index "idx_user_subscriptions_period_end" to table: "user_subscriptions"
CREATE INDEX "idx_user_subscriptions_period_end" ON "public"."user_subscriptions" ("current_period_end") WHERE (status = 'active'::public.subscription_status);
-- Create index "idx_user_subscriptions_plan" to table: "user_subscriptions"
CREATE INDEX "idx_user_subscriptions_plan" ON "public"."user_subscriptions" ("plan_id");
-- Seed subscription plans
INSERT INTO "public"."subscription_plans" (tier, specview_monthly_limit, analysis_monthly_limit, retention_days) VALUES
('free', 5000, 50, 30),
('pro', 100000, 1000, 180),
('pro_plus', 500000, 5000, 365),
('enterprise', NULL, NULL, NULL);
