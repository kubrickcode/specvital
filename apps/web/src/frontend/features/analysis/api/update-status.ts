import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { AnalyzingResponse, QueuedResponse, UpdateStatusResponse } from "@/lib/api/types";

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
