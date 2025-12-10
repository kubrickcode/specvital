import type { AnalysisResponse } from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000";
const DEFAULT_TIMEOUT_MS = 30000;

const getApiUrl = (path: string): string => {
  if (typeof window === "undefined") {
    return `${API_URL}${path}`;
  }
  return path;
};

export const fetchAnalysis = async (
  owner: string,
  repo: string,
  timeoutMs = DEFAULT_TIMEOUT_MS
): Promise<AnalysisResponse> => {
  const url = getApiUrl(`/api/analyze/${owner}/${repo}`);

  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs);

  let response: Response;

  try {
    response = await fetch(url, {
      cache: "no-store",
      headers: { Accept: "application/json" },
      signal: controller.signal,
    });
  } catch (error) {
    if (error instanceof Error && error.name === "AbortError") {
      throw new Error("Request timed out");
    }
    throw new Error(
      `Failed to fetch analysis: ${error instanceof Error ? error.message : "Network error"}`
    );
  } finally {
    clearTimeout(timeoutId);
  }

  // 202 Accepted is valid for queued/analyzing status
  if (!response.ok && response.status !== 202) {
    let errorMessage = `API request failed: ${response.statusText}`;
    try {
      const errorBody = await response.json();
      if (errorBody.detail) {
        errorMessage = errorBody.detail;
      }
    } catch {
      // Ignore JSON parse error
    }
    throw new Error(errorMessage);
  }

  let data: unknown;
  try {
    data = await response.json();
  } catch (error) {
    throw new Error(
      `Failed to parse response as JSON: ${error instanceof Error ? error.message : "Parse error"}`
    );
  }

  return data as AnalysisResponse;
};
