import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type {
  BookmarkedRepositoriesResponse,
  PaginatedRepositoriesResponse,
  UserAnalyzedRepositoriesResponse,
} from "@/lib/api/types";

const SEARCH_DATA_LIMIT = 50;

export const fetchUserAnalyzedRepositories =
  async (): Promise<UserAnalyzedRepositoriesResponse> => {
    const response = await apiFetch(`/api/user/analyzed-repositories?limit=${SEARCH_DATA_LIMIT}`);
    return parseJsonResponse(response);
  };

export const fetchUserBookmarks = async (): Promise<BookmarkedRepositoriesResponse> => {
  const response = await apiFetch("/api/user/bookmarks");
  return parseJsonResponse(response);
};

export const fetchCommunityRepositories = async (): Promise<PaginatedRepositoriesResponse> => {
  const response = await apiFetch(
    `/api/repositories/recent?limit=${SEARCH_DATA_LIMIT}&view=community`
  );
  return parseJsonResponse(response);
};
