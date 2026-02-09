import type { PlanTier } from "../types";

export type PlanConfig = {
  cta: "getStarted" | "startFree" | "contactUs";
  highlighted?: boolean;
  tier: PlanTier;
};

export const PLAN_CONFIGS: Record<PlanTier, PlanConfig> = {
  enterprise: { cta: "contactUs", tier: "enterprise" },
  free: { cta: "getStarted", tier: "free" },
  pro: { cta: "startFree", highlighted: true, tier: "pro" },
  pro_plus: { cta: "startFree", tier: "pro_plus" },
};

export const FAQ_ITEMS = ["specview", "analysis", "paymentLive", "retention"] as const;

export type FaqKey = (typeof FAQ_ITEMS)[number];
