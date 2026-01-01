import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { RepositoryStatsResponse } from "@/lib/api/types";

export const fetchRepositoryStats = async (): Promise<RepositoryStatsResponse> => {
  const response = await apiFetch("/api/repositories/stats");
  return parseJsonResponse(response);
};
