import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type {
  GitHubAppInstallationsResponse,
  GitHubAppInstallUrlResponse,
  GitHubOrganizationsResponse,
  GitHubRepositoriesResponse,
} from "@/lib/api/types";

type FetchOptions = {
  refresh?: boolean;
};

export const fetchUserGitHubRepositories = async (
  options?: FetchOptions
): Promise<GitHubRepositoriesResponse> => {
  const params = new URLSearchParams();
  if (options?.refresh) {
    params.set("refresh", "true");
  }
  const query = params.toString();
  const path = `/api/user/github/repositories${query ? `?${query}` : ""}`;
  const response = await apiFetch(path);
  return parseJsonResponse<GitHubRepositoriesResponse>(response);
};

export const fetchUserGitHubOrganizations = async (
  options?: FetchOptions
): Promise<GitHubOrganizationsResponse> => {
  const params = new URLSearchParams();
  if (options?.refresh) {
    params.set("refresh", "true");
  }
  const query = params.toString();
  const path = `/api/user/github/organizations${query ? `?${query}` : ""}`;
  const response = await apiFetch(path);
  return parseJsonResponse<GitHubOrganizationsResponse>(response);
};

export const fetchOrganizationRepositories = async (
  org: string,
  options?: FetchOptions
): Promise<GitHubRepositoriesResponse> => {
  const params = new URLSearchParams();
  if (options?.refresh) {
    params.set("refresh", "true");
  }
  const query = params.toString();
  const path = `/api/user/github/organizations/${org}/repositories${query ? `?${query}` : ""}`;
  const response = await apiFetch(path);
  return parseJsonResponse<GitHubRepositoriesResponse>(response);
};

export const fetchGitHubAppInstallations = async (): Promise<GitHubAppInstallationsResponse> => {
  const response = await apiFetch("/api/user/github-app/installations");
  return parseJsonResponse<GitHubAppInstallationsResponse>(response);
};

export const fetchGitHubAppInstallUrl = async (): Promise<GitHubAppInstallUrlResponse> => {
  const response = await apiFetch("/api/user/github-app/install-url");
  return parseJsonResponse<GitHubAppInstallUrlResponse>(response);
};
