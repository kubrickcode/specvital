-- Modify "subscription_plans" table
ALTER TABLE "public"."subscription_plans" ADD COLUMN "monthly_price" integer NULL;

-- Set pricing data
UPDATE "public"."subscription_plans" SET monthly_price = 0 WHERE tier = 'free';
UPDATE "public"."subscription_plans" SET monthly_price = 15 WHERE tier = 'pro';
UPDATE "public"."subscription_plans" SET monthly_price = 59 WHERE tier = 'pro_plus';
-- enterprise: NULL (contact sales)
