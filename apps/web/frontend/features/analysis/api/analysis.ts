import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { AnalysisResponse } from "@/lib/api/types";

type FetchAnalysisOptions = {
  commit?: string | null;
};

export async function fetchAnalysis(
  owner: string,
  repo: string,
  options?: FetchAnalysisOptions
): Promise<AnalysisResponse> {
  const url = new URL(`/api/analyze/${owner}/${repo}`, "http://dummy");
  if (options?.commit) {
    url.searchParams.set("commit", options.commit);
  }
  const response = await apiFetch(`${url.pathname}${url.search}`);
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
