"use client";

import { useQuery } from "@tanstack/react-query";

import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { components } from "@/lib/api/generated-types";

type PricingResponse = components["schemas"]["PricingResponse"];

export const usePricing = () => {
  return useQuery({
    queryFn: async () => {
      const response = await apiFetch("/api/pricing");
      const data = await parseJsonResponse<PricingResponse>(response);
      return data.data;
    },
    gcTime: 30 * 60 * 1000,
    queryKey: ["pricing"],
    staleTime: 30 * 60 * 1000,
  });
};
