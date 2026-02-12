import { isTokenExpiredResponse, refreshToken, shouldSkipRefresh } from "./token-refresh";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000";
const DEFAULT_TIMEOUT_MS = 30000;

export const getApiUrl = (path: string): string => {
  if (typeof window === "undefined") {
    return `${API_URL}${path}`;
  }
  return path;
};

type FetchOptions = {
  body?: string;
  method?: "DELETE" | "GET" | "POST";
  skipRefresh?: boolean;
  timeoutMs?: number;
};

const executeFetch = async (
  path: string,
  method: string,
  signal: AbortSignal,
  body?: string
): Promise<Response> =>
  fetch(getApiUrl(path), {
    body,
    cache: "no-store",
    credentials: "include",
    headers: {
      Accept: "application/json",
      ...(body && { "Content-Type": "application/json" }),
    },
    method,
    signal,
  });

export async function apiFetch(path: string, options: FetchOptions = {}): Promise<Response> {
  const { body, method = "GET", skipRefresh = false, timeoutMs = DEFAULT_TIMEOUT_MS } = options;
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs);

  try {
    const response = await executeFetch(path, method, controller.signal, body);

    const isServer = typeof window === "undefined";
    if (isTokenExpiredResponse(response) && !isServer && !skipRefresh && !shouldSkipRefresh(path)) {
      const refreshed = await refreshToken();

      if (refreshed) {
        clearTimeout(timeoutId);
        const retryController = new AbortController();
        const retryTimeoutId = setTimeout(() => retryController.abort(), timeoutMs);

        try {
          return await executeFetch(path, method, retryController.signal, body);
        } finally {
          clearTimeout(retryTimeoutId);
        }
      }
    }

    return response;
  } catch (error) {
    if (error instanceof Error && error.name === "AbortError") {
      throw new Error("Request timed out");
    }
    throw new Error(`Network error: ${error instanceof Error ? error.message : "Unknown"}`);
  } finally {
    clearTimeout(timeoutId);
  }
}

export async function parseJsonResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}));
    throw new Error(errorBody.detail || response.statusText);
  }
  return response.json();
}
