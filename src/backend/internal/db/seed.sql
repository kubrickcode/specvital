-- Seed data for local development
-- This file is executed by `just migrate-local`

INSERT INTO "public"."subscription_plans" (tier, specview_monthly_limit, analysis_monthly_limit, retention_days) VALUES
('free', 5000, 50, 30),
('pro', 100000, 1000, 180),
('pro_plus', 500000, 5000, 365),
('enterprise', NULL, NULL, NULL);
