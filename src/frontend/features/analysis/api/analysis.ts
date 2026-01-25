import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { AnalysisResponse } from "@/lib/api/types";

export async function fetchAnalysis(owner: string, repo: string): Promise<AnalysisResponse> {
  const response = await apiFetch(`/api/analyze/${owner}/${repo}`);
  if (!response.ok && response.status !== 202) {
    const errorBody = await response.json().catch(() => ({}));
    throw new Error(errorBody.detail || response.statusText);
  }
  return response.json();
}

export async function fetchAnalysisStatus(owner: string, repo: string): Promise<AnalysisResponse> {
  const response = await apiFetch(`/api/analyze/${owner}/${repo}/status`);
  return parseJsonResponse(response);
}
