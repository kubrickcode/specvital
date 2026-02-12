import { getApiUrl } from "./client";

const SKIP_REFRESH_PATHS = [
  "/api/auth/callback",
  "/api/auth/dev-login",
  "/api/auth/login",
  "/api/auth/logout",
  "/api/auth/refresh",
];

let refreshPromise: Promise<boolean> | null = null;

type RefreshResult = {
  success: boolean;
};

const attemptRefresh = async (): Promise<boolean> => {
  try {
    const response = await fetch(getApiUrl("/api/auth/refresh"), {
      credentials: "include",
      headers: { Accept: "application/json" },
      method: "POST",
    });

    if (!response.ok) {
      return false;
    }

    const data: RefreshResult = await response.json();
    return data.success === true;
  } catch {
    return false;
  }
};

export const refreshToken = async (): Promise<boolean> => {
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = attemptRefresh();

  try {
    return await refreshPromise;
  } finally {
    refreshPromise = null;
  }
};

export const shouldSkipRefresh = (path: string): boolean => {
  return SKIP_REFRESH_PATHS.some(
    (skipPath) => path === skipPath || path.startsWith(`${skipPath}?`)
  );
};

export const isTokenExpiredResponse = (response: Response): boolean => {
  return response.status === 401;
};
