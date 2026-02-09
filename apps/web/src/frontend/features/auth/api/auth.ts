import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { LoginResponse, LogoutResponse, RefreshResponse, UserInfo } from "@/lib/api/types";

export async function fetchCurrentUser(): Promise<UserInfo | null> {
  const response = await apiFetch("/api/auth/me");
  if (response.status === 401) return null;
  return parseJsonResponse(response);
}

export async function fetchLogin(): Promise<LoginResponse> {
  const response = await apiFetch("/api/auth/login");
  return parseJsonResponse(response);
}

export async function fetchLogout(): Promise<LogoutResponse> {
  const response = await apiFetch("/api/auth/logout", { method: "POST" });
  if (response.status === 401) return { success: true };
  return parseJsonResponse(response);
}

export async function fetchRefresh(): Promise<RefreshResponse> {
  const response = await apiFetch("/api/auth/refresh", { method: "POST", skipRefresh: true });
  return parseJsonResponse(response);
}
