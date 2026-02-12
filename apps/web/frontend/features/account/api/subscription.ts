import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { components } from "@/lib/api/generated-types";

export type UserSubscriptionResponse = components["schemas"]["UserSubscriptionResponse"];

export const fetchSubscription = async (): Promise<UserSubscriptionResponse> => {
  const response = await apiFetch("/api/user/subscription");
  return parseJsonResponse<UserSubscriptionResponse>(response);
};
