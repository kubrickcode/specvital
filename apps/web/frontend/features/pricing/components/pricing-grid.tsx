"use client";

import { Skeleton } from "@/components/ui/skeleton";
import { useAuth } from "@/features/auth";

import { usePricing } from "../api/use-pricing";
import { PLAN_CONFIGS } from "../constants/plans";
import { useCurrentPlan } from "../hooks/use-current-plan";
import type { PricingPlan } from "../types";
import { PricingCard } from "./pricing-card";

const PricingGridSkeleton = () => (
  <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
    {[1, 2, 3, 4].map((i) => (
      <div className="flex flex-col gap-4 rounded-lg border p-6" key={i}>
        <Skeleton className="mx-auto h-6 w-24" />
        <Skeleton className="mx-auto h-4 w-32" />
        <Skeleton className="mx-auto h-10 w-20" />
        <div className="space-y-3">
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-full" />
        </div>
        <Skeleton className="mt-4 h-10 w-full" />
      </div>
    ))}
  </div>
);

export const PricingGrid = () => {
  const { isAuthenticated } = useAuth();
  const { data: currentTier } = useCurrentPlan(isAuthenticated);
  const { data: pricingData, isLoading } = usePricing();

  if (isLoading || !pricingData) {
    return <PricingGridSkeleton />;
  }

  const plans: PricingPlan[] = pricingData.map((apiPlan) => {
    const config = PLAN_CONFIGS[apiPlan.tier];
    return {
      cta: config.cta,
      features: [
        { label: "specview", value: apiPlan.specviewMonthlyLimit ?? null },
        { label: "analysis", value: apiPlan.analysisMonthlyLimit ?? null },
        { label: "retention", value: apiPlan.retentionDays ?? null },
      ],
      highlighted: config.highlighted,
      monthlyPrice: apiPlan.monthlyPrice ?? null,
      tier: apiPlan.tier,
    };
  });

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
      {plans.map((plan) => (
        <PricingCard currentTier={currentTier} key={plan.tier} plan={plan} />
      ))}
    </div>
  );
};
