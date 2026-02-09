import type { components } from "@/lib/api/generated-types";

export type PlanTier = components["schemas"]["PlanTier"];

export type PricingPlan = {
  cta: "getStarted" | "startFree" | "contactUs";
  features: PlanFeature[];
  highlighted?: boolean;
  monthlyPrice: number | null;
  tier: PlanTier;
};

export type PlanFeature = {
  label: string;
  value: string | number | null;
};
