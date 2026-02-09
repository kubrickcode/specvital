/**
 * Mock data for subscription and usage APIs
 * Endpoints:
 * - /api/user/subscription
 * - /api/usage/current
 */

import type {
  UserSubscriptionResponse,
  UsageStatusResponse,
  PlanInfo,
} from "./types";

const now = new Date();
const periodStart = new Date(now.getFullYear(), now.getMonth(), 1);
const periodEnd = new Date(now.getFullYear(), now.getMonth() + 1, 0, 23, 59, 59);

// Plan configurations
const freePlan: PlanInfo = {
  tier: "free",
  specviewMonthlyLimit: 10,
  analysisMonthlyLimit: 50,
  retentionDays: 30,
};

const proPlan: PlanInfo = {
  tier: "pro",
  specviewMonthlyLimit: 100,
  analysisMonthlyLimit: 500,
  retentionDays: 90,
};

const proPlusPlan: PlanInfo = {
  tier: "pro_plus",
  specviewMonthlyLimit: 500,
  analysisMonthlyLimit: 2000,
  retentionDays: 365,
};

const enterprisePlan: PlanInfo = {
  tier: "enterprise",
  specviewMonthlyLimit: null,
  analysisMonthlyLimit: null,
  retentionDays: null,
};

// Subscription responses
export const mockSubscriptionFree: UserSubscriptionResponse = {
  plan: freePlan,
  currentPeriodStart: periodStart.toISOString(),
  currentPeriodEnd: periodEnd.toISOString(),
};

export const mockSubscriptionPro: UserSubscriptionResponse = {
  plan: proPlan,
  currentPeriodStart: periodStart.toISOString(),
  currentPeriodEnd: periodEnd.toISOString(),
};

export const mockSubscriptionProPlus: UserSubscriptionResponse = {
  plan: proPlusPlan,
  currentPeriodStart: periodStart.toISOString(),
  currentPeriodEnd: periodEnd.toISOString(),
};

export const mockSubscriptionEnterprise: UserSubscriptionResponse = {
  plan: enterprisePlan,
  currentPeriodStart: periodStart.toISOString(),
  currentPeriodEnd: periodEnd.toISOString(),
};

// Usage responses - Normal state (30% usage)
export const mockUsageNormal: UsageStatusResponse = {
  specview: {
    used: 30,
    limit: 100,
    percentage: 30,
  },
  analysis: {
    used: 150,
    limit: 500,
    percentage: 30,
  },
  resetAt: periodEnd.toISOString(),
  plan: proPlan,
};

// Usage responses - Warning state (70% usage)
export const mockUsageWarning: UsageStatusResponse = {
  specview: {
    used: 70,
    limit: 100,
    percentage: 70,
  },
  analysis: {
    used: 350,
    limit: 500,
    percentage: 70,
  },
  resetAt: periodEnd.toISOString(),
  plan: proPlan,
};

// Usage responses - Critical state (90% usage)
export const mockUsageCritical: UsageStatusResponse = {
  specview: {
    used: 92,
    limit: 100,
    percentage: 92,
  },
  analysis: {
    used: 460,
    limit: 500,
    percentage: 92,
  },
  resetAt: periodEnd.toISOString(),
  plan: proPlan,
};

// Usage responses - Exceeded state (100%+ usage)
export const mockUsageExceeded: UsageStatusResponse = {
  specview: {
    used: 105,
    limit: 100,
    percentage: 105,
  },
  analysis: {
    used: 520,
    limit: 500,
    percentage: 104,
  },
  resetAt: periodEnd.toISOString(),
  plan: proPlan,
};

// Usage responses - Free tier
export const mockUsageFree: UsageStatusResponse = {
  specview: {
    used: 5,
    limit: 10,
    percentage: 50,
  },
  analysis: {
    used: 25,
    limit: 50,
    percentage: 50,
  },
  resetAt: periodEnd.toISOString(),
  plan: freePlan,
};

// Usage responses - Enterprise (unlimited)
export const mockUsageEnterprise: UsageStatusResponse = {
  specview: {
    used: 500,
    limit: null,
    percentage: null,
  },
  analysis: {
    used: 2500,
    limit: null,
    percentage: null,
  },
  resetAt: periodEnd.toISOString(),
  plan: enterprisePlan,
};
