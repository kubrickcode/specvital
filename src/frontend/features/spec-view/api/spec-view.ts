import { apiFetch } from "@/lib/api/client";

import type {
  CacheAvailabilityResponse,
  RepoSpecDocumentResponse,
  RepoVersionHistoryResponse,
  RequestSpecGenerationRequest,
  RequestSpecGenerationResponse,
  SpecDocumentResponse,
  SpecGenerationStatusResponse,
  SpecLanguage,
  VersionHistoryResponse,
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

export class NoSubscriptionError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "NoSubscriptionError";
  }
}

export class UnauthorizedError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "UnauthorizedError";
  }
}

export class ForbiddenError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ForbiddenError";
  }
}

type FetchSpecDocumentOptions = {
  language?: SpecLanguage;
  version?: number;
};

export const fetchSpecDocument = async (
  analysisId: string,
  options?: FetchSpecDocumentOptions
): Promise<SpecDocumentResponse> => {
  const params = new URLSearchParams();
  if (options?.language) {
    params.set("language", options.language);
  }
  if (options?.version !== undefined) {
    params.set("version", String(options.version));
  }
  const queryString = params.toString();
  const url = queryString
    ? `/api/spec-view/${analysisId}?${queryString}`
    : `/api/spec-view/${analysisId}`;
  const response = await apiFetch(url);

  if (!response.ok && response.status !== 404) {
    const errorBody = await response.json().catch(() => ({}));

    if (response.status === 401) {
      throw new UnauthorizedError(errorBody.detail || "Authentication required");
    }
    if (response.status === 403) {
      throw new ForbiddenError(errorBody.detail || "Access denied");
    }

    throw new Error(errorBody.detail || response.statusText);
  }

  if (response.status === 404) {
    const statusUrl = options?.language
      ? `/api/spec-view/status/${analysisId}?language=${encodeURIComponent(options.language)}`
      : `/api/spec-view/status/${analysisId}`;
    const statusResponse = await apiFetch(statusUrl);

    if (statusResponse.ok) {
      const statusData: SpecGenerationStatusResponse = await statusResponse.json();
      return {
        generationStatus: { status: statusData.status },
        status: "generating",
      };
    }

    return {
      generationStatus: { status: "not_found" },
      status: "generating",
    };
  }

  return response.json();
};

export const requestSpecGeneration = async (
  request: RequestSpecGenerationRequest
): Promise<RequestSpecGenerationResponse> => {
  const response = await apiFetch("/api/spec-view/generate", {
    body: JSON.stringify(request),
    method: "POST",
  });

  if (response.status === 403) {
    const errorBody = await response.json().catch(() => ({}));
    throw new NoSubscriptionError(errorBody.detail || "Active subscription required");
  }

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
};

export const fetchGenerationStatus = async (
  analysisId: string,
  language?: SpecLanguage
): Promise<SpecGenerationStatusResponse> => {
  const url = language
    ? `/api/spec-view/status/${analysisId}?language=${encodeURIComponent(language)}`
    : `/api/spec-view/status/${analysisId}`;
  const response = await apiFetch(url);

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}));

    if (response.status === 401) {
      throw new UnauthorizedError(errorBody.detail || "Authentication required");
    }
    if (response.status === 403) {
      throw new ForbiddenError(errorBody.detail || "Access denied");
    }

    throw new Error(errorBody.detail || response.statusText);
  }

  return response.json();
};

export const fetchVersionHistory = async (
  analysisId: string,
  language: SpecLanguage
): Promise<VersionHistoryResponse> => {
  const url = `/api/spec-view/${analysisId}/versions?language=${encodeURIComponent(language)}`;
  const response = await apiFetch(url);

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}));

    if (response.status === 401) {
      throw new UnauthorizedError(errorBody.detail || "Authentication required");
    }
    if (response.status === 403) {
      throw new ForbiddenError(errorBody.detail || "Access denied");
    }

    throw new Error(errorBody.detail || response.statusText);
  }

  return response.json();
};

export const fetchCacheAvailability = async (
  analysisId: string
): Promise<CacheAvailabilityResponse> => {
  const url = `/api/spec-view/${analysisId}/cache-availability`;
  const response = await apiFetch(url);

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}));

    if (response.status === 401) {
      throw new UnauthorizedError(errorBody.detail || "Authentication required");
    }

    throw new Error(errorBody.detail || response.statusText);
  }

  return response.json();
};

// Repository-based API functions (cross-analysis version access)

type FetchRepoSpecDocumentOptions = {
  documentId?: string;
  language?: SpecLanguage;
  version?: number;
};

export const fetchRepoSpecDocument = async (
  owner: string,
  repo: string,
  options?: FetchRepoSpecDocumentOptions
): Promise<RepoSpecDocumentResponse> => {
  const params = new URLSearchParams();
  if (options?.documentId) {
    params.set("documentId", options.documentId);
  }
  if (options?.language) {
    params.set("language", options.language);
  }
  if (options?.version !== undefined) {
    params.set("version", String(options.version));
  }
  const queryString = params.toString();
  const url = queryString
    ? `/api/spec-view/repository/${encodeURIComponent(owner)}/${encodeURIComponent(repo)}?${queryString}`
    : `/api/spec-view/repository/${encodeURIComponent(owner)}/${encodeURIComponent(repo)}`;

  const response = await apiFetch(url);

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}));

    if (response.status === 401) {
      throw new UnauthorizedError(errorBody.detail || "Authentication required");
    }
    if (response.status === 403) {
      throw new ForbiddenError(errorBody.detail || "Access denied");
    }

    // 404 = No spec document found for this repository → return empty status
    if (response.status === 404) {
      return { status: "empty" };
    }

    throw new Error(errorBody.detail || response.statusText);
  }

  return response.json();
};

export const fetchRepoVersionHistory = async (
  owner: string,
  repo: string,
  language: SpecLanguage
): Promise<RepoVersionHistoryResponse> => {
  const url = `/api/spec-view/repository/${encodeURIComponent(owner)}/${encodeURIComponent(repo)}/versions?language=${encodeURIComponent(language)}`;
  const response = await apiFetch(url);

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}));

    if (response.status === 401) {
      throw new UnauthorizedError(errorBody.detail || "Authentication required");
    }
    if (response.status === 403) {
      throw new ForbiddenError(errorBody.detail || "Access denied");
    }

    throw new Error(errorBody.detail || response.statusText);
  }

  return response.json();
};
