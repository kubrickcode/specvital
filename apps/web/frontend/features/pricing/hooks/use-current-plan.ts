"use client";

import { useQuery } from "@tanstack/react-query";

import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { components } from "@/lib/api/generated-types";

type SubscriptionResponse = components["schemas"]["UserSubscriptionResponse"];
type PlanTier = components["schemas"]["PlanTier"];

const fetchCurrentPlan = async (): Promise<PlanTier> => {
  const response = await apiFetch("/api/user/subscription");
  const data = await parseJsonResponse<SubscriptionResponse>(response);
  return data.plan.tier;
};

export const useCurrentPlan = (enabled = true) => {
  return useQuery({
    enabled,
    queryFn: fetchCurrentPlan,
    queryKey: ["user", "currentPlan"],
  });
};
