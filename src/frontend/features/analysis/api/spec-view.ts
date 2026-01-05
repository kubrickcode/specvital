import { apiFetch, parseJsonResponse } from "@/lib/api/client";

import type { ConversionLanguage, ConvertSpecViewResponse } from "../types";

const SPEC_VIEW_TIMEOUT_MS = 5 * 60 * 1000; // 5 minutes for AI conversion

type ConvertSpecViewParams = {
  commitSha: string;
  isForceRefresh?: boolean;
  language: ConversionLanguage;
  owner: string;
  repo: string;
};

export const convertSpecView = async ({
  commitSha,
  isForceRefresh = false,
  language,
  owner,
  repo,
}: ConvertSpecViewParams): Promise<ConvertSpecViewResponse> => {
  const response = await apiFetch(`/api/spec-view/convert/${owner}/${repo}/${commitSha}`, {
    body: JSON.stringify({ isForceRefresh, language }),
    method: "POST",
    timeoutMs: SPEC_VIEW_TIMEOUT_MS,
  });
  return parseJsonResponse<ConvertSpecViewResponse>(response);
};
