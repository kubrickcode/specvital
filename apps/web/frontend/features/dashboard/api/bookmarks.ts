import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { BookmarkResponse } from "@/lib/api/types";

export const addBookmark = async (owner: string, repo: string): Promise<BookmarkResponse> => {
  const response = await apiFetch(`/api/repositories/${owner}/${repo}/bookmark`, {
    method: "POST",
  });
  return parseJsonResponse(response);
};

export const removeBookmark = async (owner: string, repo: string): Promise<BookmarkResponse> => {
  const response = await apiFetch(`/api/repositories/${owner}/${repo}/bookmark`, {
    method: "DELETE",
  });
  return parseJsonResponse(response);
};
