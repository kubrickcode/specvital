import { getApiUrl, parseJsonResponse } from "@/lib/api/client";
import type { components } from "@/lib/api/generated-types";

type AddAnalyzedRepositoryRequest = components["schemas"]["AddAnalyzedRepositoryRequest"];
type AddAnalyzedRepositoryResponse = components["schemas"]["AddAnalyzedRepositoryResponse"];

export const addToHistory = async (
  owner: string,
  repo: string
): Promise<AddAnalyzedRepositoryResponse> => {
  const body: AddAnalyzedRepositoryRequest = { owner, repo };

  const response = await fetch(getApiUrl("/api/user/analyzed-repositories"), {
    body: JSON.stringify(body),
    cache: "no-store",
    credentials: "include",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    method: "POST",
  });

  return parseJsonResponse<AddAnalyzedRepositoryResponse>(response);
};
