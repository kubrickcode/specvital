import { apiFetch } from "@/lib/api/client";

import type {
  RequestSpecGenerationRequest,
  RequestSpecGenerationResponse,
  SpecDocumentResponse,
  SpecGenerationStatusResponse,
  SpecLanguage,
} from "../types";

export class QuotaExceededError extends Error {
  limit: number;
  used: number;

  constructor(message: string, used: number, limit: number) {
    super(message);
    this.name = "QuotaExceededError";
    this.used = used;
    this.limit = limit;
  }
}

export async function fetchSpecDocument(
  analysisId: string,
  language?: SpecLanguage
): Promise<SpecDocumentResponse> {
  const url = language
    ? `/api/spec-view/${analysisId}?language=${encodeURIComponent(language)}`
    : `/api/spec-view/${analysisId}`;
  const response = await apiFetch(url);

  if (!response.ok && response.status !== 404) {
    const errorBody = await response.json().catch(() => ({}));
    throw new Error(errorBody.detail || response.statusText);
  }

  if (response.status === 404) {
    // 문서가 없으면 큐 상태 확인
    const statusResponse = await apiFetch(`/api/spec-view/status/${analysisId}`);
    if (statusResponse.ok) {
      const statusData: SpecGenerationStatusResponse = await statusResponse.json();
      return {
        generationStatus: { status: statusData.status },
        status: "generating",
      };
    }
    // 큐에도 없으면 not_found
    return {
      generationStatus: { status: "not_found" },
      status: "generating",
    };
  }

  return response.json();
}

export async function requestSpecGeneration(
  request: RequestSpecGenerationRequest
): Promise<RequestSpecGenerationResponse> {
  const response = await apiFetch("/api/spec-view/generate", {
    body: JSON.stringify(request),
    method: "POST",
  });

  if (response.status === 429) {
    const errorBody = await response.json().catch(() => ({}));
    throw new QuotaExceededError(
      errorBody.detail || "Quota exceeded",
      errorBody.used ?? 0,
      errorBody.limit ?? 0
    );
  }

  if (!response.ok && response.status !== 202 && response.status !== 409) {
    const errorBody = await response.json().catch(() => ({}));
    throw new Error(errorBody.detail || response.statusText);
  }

  return response.json();
}

export async function fetchGenerationStatus(
  analysisId: string
): Promise<SpecGenerationStatusResponse> {
  const response = await apiFetch(`/api/spec-view/status/${analysisId}`);

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}));
    throw new Error(errorBody.detail || response.statusText);
  }

  return response.json();
}
