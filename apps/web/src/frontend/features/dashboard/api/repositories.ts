import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type {
  AnalyzingResponse,
  OwnershipFilterParam,
  PaginatedRepositoriesResponse,
  QueuedResponse,
  SortByParam,
  SortOrderParam,
  UpdateStatusResponse,
  ViewFilterParam,
} from "@/lib/api/types";

export type PaginatedRepositoriesParams = {
  cursor?: string;
  limit?: number;
  ownership?: OwnershipFilterParam;
  sortBy?: SortByParam;
  sortOrder?: SortOrderParam;
  view?: ViewFilterParam;
};

export const fetchPaginatedRepositories = async (
  params: PaginatedRepositoriesParams = {}
): Promise<PaginatedRepositoriesResponse> => {
  const searchParams = new URLSearchParams();

  if (params.cursor) searchParams.set("cursor", params.cursor);
  if (params.limit) searchParams.set("limit", String(params.limit));
  if (params.ownership) searchParams.set("ownership", params.ownership);
  if (params.sortBy) searchParams.set("sortBy", params.sortBy);
  if (params.sortOrder) searchParams.set("sortOrder", params.sortOrder);
  if (params.view) searchParams.set("view", params.view);

  const queryString = searchParams.toString();
  const url = queryString ? `/api/repositories/recent?${queryString}` : "/api/repositories/recent";

  const response = await apiFetch(url);
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
