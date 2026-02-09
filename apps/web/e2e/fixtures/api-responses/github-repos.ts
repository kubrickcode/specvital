/**
 * Mock data for GitHub repositories API
 * Endpoint: /api/user/github/repositories
 */

import type { GitHubRepository, GitHubRepositoriesResponse } from "./types";

const createGitHubRepo = (
  index: number,
  overrides: Partial<GitHubRepository> = {}
): GitHubRepository => ({
  id: 100000 + index,
  name: `github-repo-${index}`,
  fullName: `e2e-user/github-repo-${index}`,
  owner: "e2e-user",
  description: `Test repository ${index} for E2E testing`,
  defaultBranch: "main",
  isPrivate: index % 4 === 0,
  pushedAt: new Date(Date.now() - index * 86400000).toISOString(),
  ...overrides,
});

// User's GitHub repositories (15 repos)
export const mockGitHubReposList: GitHubRepository[] = Array.from(
  { length: 15 },
  (_, i) => createGitHubRepo(i + 1)
);

export const mockGitHubReposResponse: GitHubRepositoriesResponse = {
  data: mockGitHubReposList,
};

// Empty state
export const mockGitHubReposEmpty: GitHubRepositoriesResponse = {
  data: [],
};

// Organization repositories
const createOrgRepo = (
  orgLogin: string,
  index: number
): GitHubRepository => ({
  id: 200000 + index,
  name: `org-repo-${index}`,
  fullName: `${orgLogin}/org-repo-${index}`,
  owner: orgLogin,
  description: `Organization repository ${index}`,
  defaultBranch: "main",
  isPrivate: false,
  pushedAt: new Date(Date.now() - index * 86400000).toISOString(),
});

export const mockOrgRepos: Record<string, GitHubRepositoriesResponse> = {
  "test-org": {
    data: Array.from({ length: 10 }, (_, i) => createOrgRepo("test-org", i + 1)),
  },
  "another-org": {
    data: Array.from({ length: 5 }, (_, i) => createOrgRepo("another-org", i + 1)),
  },
};
