import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { components } from "@/lib/api/generated-types";

export type UsageStatusResponse = components["schemas"]["UsageStatusResponse"];

export const fetchCurrentUsage = async (): Promise<UsageStatusResponse> => {
  const response = await apiFetch("/api/usage/current");
  return parseJsonResponse<UsageStatusResponse>(response);
};
