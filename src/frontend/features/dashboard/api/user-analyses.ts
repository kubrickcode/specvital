import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { OwnershipFilter, UserAnalyzedRepositoriesResponse } from "@/lib/api/types";

type FetchUserAnalyzedRepositoriesParams = {
  cursor?: string;
  limit?: number;
  ownership?: OwnershipFilter;
};

export const fetchUserAnalyzedRepositories = async (
  params?: FetchUserAnalyzedRepositoriesParams
): Promise<UserAnalyzedRepositoriesResponse> => {
  const searchParams = new URLSearchParams();

  if (params?.cursor) {
    searchParams.set("cursor", params.cursor);
  }
  if (params?.limit) {
    searchParams.set("limit", String(params.limit));
  }
  if (params?.ownership && params.ownership !== "all") {
    searchParams.set("ownership", params.ownership);
  }

  const queryString = searchParams.toString();
  const url = `/api/user/analyzed-repositories${queryString ? `?${queryString}` : ""}`;

  const response = await apiFetch(url);
  return parseJsonResponse(response);
};
