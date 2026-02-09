/**
 * Mock data for Dashboard repository list API
 * Endpoint: /api/user/analyzed-repositories
 */

import type {
  UserAnalyzedRepositoriesResponse,
  RepositoryCard,
} from "./types";

const createRepository = (
  index: number,
  overrides: Partial<RepositoryCard> = {}
): RepositoryCard => ({
  id: `repo-${index}`,
  owner: index % 3 === 0 ? "e2e-user" : index % 3 === 1 ? "test-org" : "other-user",
  name: `test-repo-${index}`,
  fullName: `${index % 3 === 0 ? "e2e-user" : index % 3 === 1 ? "test-org" : "other-user"}/test-repo-${index}`,
  isBookmarked: index < 3,
  isAnalyzedByMe: true,
  updateStatus: index === 0 ? "new-commits" : "up-to-date",
  latestAnalysis: {
    testCount: 50 + index * 10,
    change: index % 5 === 0 ? 5 : 0,
    analyzedAt: new Date(Date.now() - index * 86400000).toISOString(),
    commitSha: `abc${index}def${index}`,
    testSummary: {
      active: 45 + index * 9,
      focused: index === 1 ? 2 : 0,
      skipped: 3,
      todo: 2,
      xfail: 0,
    },
  },
  aiSpecSummary: index <= 5 ? {
    hasSpec: true,
    languageCount: index % 3 + 1,
    latestGeneratedAt: new Date(Date.now() - index * 3600000).toISOString(),
  } : undefined,
  ...overrides,
});

// 25 repositories for infinite scroll testing
export const mockRepositoriesList: RepositoryCard[] = Array.from(
  { length: 25 },
  (_, i) => createRepository(i + 1)
);

// First page response (20 items)
export const mockRepositoriesPage1: UserAnalyzedRepositoriesResponse = {
  data: mockRepositoriesList.slice(0, 20),
  nextCursor: "cursor-page-2",
  hasNext: true,
};

// Second page response (5 items)
export const mockRepositoriesPage2: UserAnalyzedRepositoriesResponse = {
  data: mockRepositoriesList.slice(20, 25),
  nextCursor: null,
  hasNext: false,
};

// Empty state
export const mockRepositoriesEmpty: UserAnalyzedRepositoriesResponse = {
  data: [],
  nextCursor: null,
  hasNext: false,
};

// Filtered: Mine only
export const mockRepositoriesMine: UserAnalyzedRepositoriesResponse = {
  data: mockRepositoriesList.filter((r) => r.owner === "e2e-user"),
  nextCursor: null,
  hasNext: false,
};

// Filtered: Organization only
export const mockRepositoriesOrg: UserAnalyzedRepositoriesResponse = {
  data: mockRepositoriesList.filter((r) => r.owner === "test-org"),
  nextCursor: null,
  hasNext: false,
};

// Filtered: Bookmarked only
export const mockRepositoriesBookmarked: UserAnalyzedRepositoriesResponse = {
  data: mockRepositoriesList.filter((r) => r.isBookmarked),
  nextCursor: null,
  hasNext: false,
};

// Sorted by name (A-Z)
export const mockRepositoriesSortedByName: UserAnalyzedRepositoriesResponse = {
  data: [...mockRepositoriesList].sort((a, b) => a.name.localeCompare(b.name)),
  nextCursor: null,
  hasNext: false,
};

// Sorted by test count
export const mockRepositoriesSortedByTests: UserAnalyzedRepositoriesResponse = {
  data: [...mockRepositoriesList].sort(
    (a, b) => (b.latestAnalysis?.testCount ?? 0) - (a.latestAnalysis?.testCount ?? 0)
  ),
  nextCursor: null,
  hasNext: false,
};

// Search result
export const mockRepositoriesSearchResult: UserAnalyzedRepositoriesResponse = {
  data: mockRepositoriesList.filter((r) => r.name.includes("repo-1")),
  nextCursor: null,
  hasNext: false,
};

// No search results
export const mockRepositoriesNoResults: UserAnalyzedRepositoriesResponse = {
  data: [],
  nextCursor: null,
  hasNext: false,
};
