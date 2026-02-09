import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { components } from "@/lib/api/generated-types";

export type AnalysisHistoryItem = components["schemas"]["AnalysisHistoryItem"];
export type AnalysisHistoryResponse = components["schemas"]["AnalysisHistoryResponse"];

export async function fetchAnalysisHistory(
  owner: string,
  repo: string
): Promise<AnalysisHistoryResponse> {
  const response = await apiFetch(`/api/analyze/${owner}/${repo}/history`);
  return parseJsonResponse<AnalysisHistoryResponse>(response);
}
