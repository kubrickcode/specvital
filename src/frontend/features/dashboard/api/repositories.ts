import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type {
  AnalyzingResponse,
  QueuedResponse,
  RecentRepositoriesResponse,
  RepositoryStatsResponse,
  UpdateStatusResponse,
} from "@/lib/api/types";

export const fetchRecentRepositories = async (
  limit?: number
): Promise<RecentRepositoriesResponse> => {
  const params = limit ? `?limit=${limit}` : "";
  const response = await apiFetch(`/api/repositories/recent${params}`);
  return parseJsonResponse(response);
};

export const fetchRepositoryStats = async (): Promise<RepositoryStatsResponse> => {
  const response = await apiFetch("/api/repositories/stats");
  return parseJsonResponse(response);
};

export const checkUpdateStatus = async (
  owner: string,
  repo: string
): Promise<UpdateStatusResponse> => {
  const response = await apiFetch(`/api/repositories/${owner}/${repo}/update-status`);
  return parseJsonResponse(response);
};

export const triggerReanalyze = async (
  owner: string,
  repo: string
): Promise<QueuedResponse | AnalyzingResponse> => {
  const response = await apiFetch(`/api/repositories/${owner}/${repo}/reanalyze`, {
    method: "POST",
  });
  return parseJsonResponse(response);
};
