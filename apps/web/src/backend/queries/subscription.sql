-- name: GetPlanByTier :one
SELECT
    id,
    tier,
    specview_monthly_limit,
    analysis_monthly_limit,
    retention_days,
    created_at
FROM subscription_plans
WHERE tier = $1;

-- name: GetAllPlans :many
SELECT
    id,
    tier,
    specview_monthly_limit,
    analysis_monthly_limit,
    retention_days,
    created_at
FROM subscription_plans
ORDER BY
    CASE tier
        WHEN 'free' THEN 1
        WHEN 'pro' THEN 2
        WHEN 'pro_plus' THEN 3
        WHEN 'enterprise' THEN 4
    END;

-- name: GetPricingPlans :many
SELECT
    tier,
    monthly_price,
    specview_monthly_limit,
    analysis_monthly_limit,
    retention_days
FROM subscription_plans
ORDER BY
    CASE tier
        WHEN 'free' THEN 1
        WHEN 'pro' THEN 2
        WHEN 'pro_plus' THEN 3
        WHEN 'enterprise' THEN 4
    END;

-- name: CreateUserSubscription :one
INSERT INTO user_subscriptions (user_id, plan_id, current_period_start, current_period_end)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, plan_id, status, current_period_start, current_period_end, created_at;

-- name: GetActiveSubscriptionWithPlan :one
SELECT
    us.id,
    us.user_id,
    us.plan_id,
    us.status,
    us.current_period_start,
    us.current_period_end,
    us.canceled_at,
    us.created_at,
    us.updated_at,
    sp.tier AS plan_tier,
    sp.specview_monthly_limit AS plan_specview_monthly_limit,
    sp.analysis_monthly_limit AS plan_analysis_monthly_limit,
    sp.retention_days AS plan_retention_days
FROM user_subscriptions us
JOIN subscription_plans sp ON us.plan_id = sp.id
WHERE us.user_id = $1 AND us.status = 'active';

-- name: GetUsersWithoutActiveSubscription :many
SELECT id
FROM users
WHERE id NOT IN (
    SELECT user_id FROM user_subscriptions WHERE status = 'active'
);
