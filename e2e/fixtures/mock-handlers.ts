/**
 * Playwright route handlers for API mocking
 * Usage: await setupMockHandlers(page, { repositories: mockRepositoriesPage1 })
 *
 * Note: Next.js proxies /api/* requests to backend via rewrites.
 * We intercept browser requests (relative paths) not backend URLs.
 */

import type { Page, Route } from "@playwright/test";
import type {
  UserAnalyzedRepositoriesResponse,
  RepositoryStatsResponse,
  GitHubRepositoriesResponse,
  GitHubOrganizationsResponse,
  UserSubscriptionResponse,
  UsageStatusResponse,
  BookmarkResponse,
  ReanalyzeResponse,
  AnalysisCompletedResponse,
  SpecDocumentResponse,
  RequestSpecGenerationResponse,
  VersionHistoryResponse,
  RepoSpecDocumentResponse,
  RepoVersionHistoryResponse,
} from "./api-responses";

export interface MockHandlersOptions {
  repositories?: UserAnalyzedRepositoriesResponse;
  repositoriesPage2?: UserAnalyzedRepositoriesResponse;
  stats?: RepositoryStatsResponse;
  githubRepos?: GitHubRepositoriesResponse;
  githubOrgs?: GitHubOrganizationsResponse;
  orgRepos?: Record<string, GitHubRepositoriesResponse>;
  subscription?: UserSubscriptionResponse;
  usage?: UsageStatusResponse;
  bookmarkResponse?: BookmarkResponse;
  reanalyzeResponse?: ReanalyzeResponse;
  // Analysis page mocks
  analysis?: AnalysisCompletedResponse;
  specDocument?: SpecDocumentResponse;
  specDocumentByVersion?: Record<number, SpecDocumentResponse>;
  specGeneration?: RequestSpecGenerationResponse;
  versionHistory?: VersionHistoryResponse;
  // Repository-based spec document mocks (cross-analysis version access)
  repoSpecDocument?: RepoSpecDocumentResponse;
  repoSpecDocumentByVersion?: Record<number, RepoSpecDocumentResponse>;
  repoVersionHistory?: RepoVersionHistoryResponse;
  // Dynamic handlers for testing specific scenarios
  onRepositoriesRequest?: (url: URL) => UserAnalyzedRepositoriesResponse | null;
  onBookmark?: (owner: string, repo: string, method: string) => BookmarkResponse;
  onReanalyze?: (owner: string, repo: string) => ReanalyzeResponse;
  onSpecGeneration?: (analysisId: string, language: string) => RequestSpecGenerationResponse;
}

type RouteHandler = (route: Route) => Promise<void>;

/**
 * Setup mock handlers for API routes
 * Routes intercept browser requests (which go through Next.js proxy)
 */
export async function setupMockHandlers(
  page: Page,
  options: MockHandlersOptions = {}
): Promise<void> {
  const handlers: Array<{ pattern: string | RegExp; handler: RouteHandler }> = [];

  // /api/repositories/recent (Dashboard uses this endpoint)
  if (options.repositories || options.onRepositoriesRequest) {
    handlers.push({
      pattern: "**/api/repositories/recent*",
      handler: async (route) => {
        const url = new URL(route.request().url());

        // Check for dynamic handler first
        if (options.onRepositoriesRequest) {
          const response = options.onRepositoriesRequest(url);
          if (response) {
            return route.fulfill({
              status: 200,
              contentType: "application/json",
              body: JSON.stringify(response),
            });
          }
        }

        // Check for pagination
        const cursor = url.searchParams.get("cursor");
        if (cursor === "cursor-page-2" && options.repositoriesPage2) {
          return route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify(options.repositoriesPage2),
          });
        }

        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(options.repositories),
        });
      },
    });
  }

  // /api/repositories/stats
  if (options.stats) {
    handlers.push({
      pattern: "**/api/repositories/stats",
      handler: async (route) => {
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(options.stats),
        });
      },
    });
  }

  // /api/user/github/repositories
  if (options.githubRepos) {
    handlers.push({
      pattern: "**/api/user/github/repositories*",
      handler: async (route) => {
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(options.githubRepos),
        });
      },
    });
  }

  // /api/user/github/organizations (exact match, not including sub-paths)
  if (options.githubOrgs) {
    handlers.push({
      pattern: /\/api\/user\/github\/organizations(\?.*)?$/,
      handler: async (route) => {
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(options.githubOrgs),
        });
      },
    });
  }

  // /api/user/github/organizations/{org}/repositories
  if (options.orgRepos) {
    handlers.push({
      pattern: /\/api\/user\/github\/organizations\/([^/]+)\/repositories/,
      handler: async (route) => {
        const url = new URL(route.request().url());
        const orgMatch = url.pathname.match(
          /\/api\/user\/github\/organizations\/([^/]+)\/repositories/
        );
        const org = orgMatch?.[1];

        if (org && options.orgRepos?.[org]) {
          return route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify(options.orgRepos[org]),
          });
        }

        return route.fulfill({
          status: 404,
          contentType: "application/json",
          body: JSON.stringify({ error: "Organization not found" }),
        });
      },
    });
  }

  // /api/user/subscription
  if (options.subscription) {
    handlers.push({
      pattern: "**/api/user/subscription",
      handler: async (route) => {
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(options.subscription),
        });
      },
    });
  }

  // /api/usage/current
  if (options.usage) {
    handlers.push({
      pattern: "**/api/usage/current",
      handler: async (route) => {
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(options.usage),
        });
      },
    });
  }

  // /api/repositories/{owner}/{repo}/bookmark (POST/DELETE)
  if (options.bookmarkResponse || options.onBookmark) {
    handlers.push({
      pattern: /\/api\/repositories\/([^/]+)\/([^/]+)\/bookmark/,
      handler: async (route) => {
        const url = new URL(route.request().url());
        const match = url.pathname.match(
          /\/api\/repositories\/([^/]+)\/([^/]+)\/bookmark/
        );
        const owner = match?.[1] ?? "";
        const repo = match?.[2] ?? "";
        const method = route.request().method();

        let response: BookmarkResponse;
        if (options.onBookmark) {
          response = options.onBookmark(owner, repo, method);
        } else {
          response = options.bookmarkResponse ?? {
            success: true,
            isBookmarked: method === "POST",
          };
        }

        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(response),
        });
      },
    });
  }

  // /api/repositories/{owner}/{repo}/reanalyze
  if (options.reanalyzeResponse || options.onReanalyze) {
    handlers.push({
      pattern: /\/api\/repositories\/([^/]+)\/([^/]+)\/reanalyze/,
      handler: async (route) => {
        const url = new URL(route.request().url());
        const match = url.pathname.match(
          /\/api\/repositories\/([^/]+)\/([^/]+)\/reanalyze/
        );
        const owner = match?.[1] ?? "";
        const repo = match?.[2] ?? "";

        let response: ReanalyzeResponse;
        if (options.onReanalyze) {
          response = options.onReanalyze(owner, repo);
        } else {
          response = options.reanalyzeResponse ?? { status: "queued" };
        }

        return route.fulfill({
          status: 202,
          contentType: "application/json",
          body: JSON.stringify(response),
        });
      },
    });
  }

  // /api/analyze/{owner}/{repo}
  if (options.analysis) {
    handlers.push({
      pattern: /\/api\/analyze\/([^/]+)\/([^/]+)$/,
      handler: async (route) => {
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(options.analysis),
        });
      },
    });
  }

  // /api/spec-view/{analysisId} - Get spec document (supports version query param)
  if (options.specDocument || options.specDocumentByVersion) {
    handlers.push({
      pattern: /\/api\/spec-view\/[0-9a-f-]+(\?.*)?$/,
      handler: async (route) => {
        const url = new URL(route.request().url());
        const versionParam = url.searchParams.get("version");

        // If version is specified and we have version-specific documents
        if (versionParam && options.specDocumentByVersion) {
          const version = parseInt(versionParam, 10);
          const versionDoc = options.specDocumentByVersion[version];
          if (versionDoc) {
            return route.fulfill({
              status: 200,
              contentType: "application/json",
              body: JSON.stringify(versionDoc),
            });
          }
        }

        // Default response
        const response = options.specDocument;
        if (response?.status === "generating" && response.generationStatus?.status === "not_found") {
          return route.fulfill({
            status: 404,
            contentType: "application/json",
            body: JSON.stringify({ detail: "Not found" }),
          });
        }
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(response),
        });
      },
    });
  }

  // /api/spec-view/{analysisId}/versions - Get version history
  if (options.versionHistory) {
    handlers.push({
      pattern: /\/api\/spec-view\/[0-9a-f-]+\/versions/,
      handler: async (route) => {
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(options.versionHistory),
        });
      },
    });
  }

  // /api/spec-view/status/{analysisId} - Get generation status
  if (options.specGeneration) {
    handlers.push({
      pattern: /\/api\/spec-view\/status\/[0-9a-f-]+$/,
      handler: async (route) => {
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(options.specGeneration),
        });
      },
    });
  }

  // /api/spec-view/generate - Request generation
  if (options.specGeneration || options.onSpecGeneration) {
    handlers.push({
      pattern: "**/api/spec-view/generate",
      handler: async (route) => {
        const body = route.request().postDataJSON();
        const analysisId = body?.analysisId ?? "";
        const language = body?.language ?? "English";

        let response: RequestSpecGenerationResponse;
        if (options.onSpecGeneration) {
          response = options.onSpecGeneration(analysisId, language);
        } else {
          response = options.specGeneration ?? { status: "pending", analysisId };
        }

        return route.fulfill({
          status: 202,
          contentType: "application/json",
          body: JSON.stringify(response),
        });
      },
    });
  }

  // /api/spec-view/repository/{owner}/{repo} - Get spec document by repository (supports version query param)
  // Also handles backward compatibility: converts specDocument to repoSpecDocument format
  if (options.repoSpecDocument || options.repoSpecDocumentByVersion || options.specDocument || options.specDocumentByVersion) {
    handlers.push({
      pattern: /\/api\/spec-view\/repository\/([^/]+)\/([^/]+)(\?.*)?$/,
      handler: async (route) => {
        const url = new URL(route.request().url());
        const versionParam = url.searchParams.get("version");

        // If version is specified and we have version-specific documents (repo format)
        if (versionParam && options.repoSpecDocumentByVersion) {
          const version = parseInt(versionParam, 10);
          const versionDoc = options.repoSpecDocumentByVersion[version];
          if (versionDoc) {
            return route.fulfill({
              status: 200,
              contentType: "application/json",
              body: JSON.stringify(versionDoc),
            });
          }
        }

        // Backward compatibility: convert specDocumentByVersion to repoSpecDocument format
        if (versionParam && options.specDocumentByVersion) {
          const version = parseInt(versionParam, 10);
          const versionDoc = options.specDocumentByVersion[version];
          if (versionDoc && versionDoc.status === "completed" && versionDoc.data) {
            const repoDoc: RepoSpecDocumentResponse = {
              status: "completed",
              data: {
                ...versionDoc.data,
                commitSha: "abc123def",
              },
            };
            return route.fulfill({
              status: 200,
              contentType: "application/json",
              body: JSON.stringify(repoDoc),
            });
          }
        }

        // Prefer repoSpecDocument if provided
        if (options.repoSpecDocument) {
          return route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify(options.repoSpecDocument),
          });
        }

        // Backward compatibility: convert specDocument to repoSpecDocument format
        if (options.specDocument) {
          // Handle "not_found" case as "empty"
          if (
            options.specDocument.status === "generating" &&
            options.specDocument.generationStatus?.status === "not_found"
          ) {
            return route.fulfill({
              status: 200,
              contentType: "application/json",
              body: JSON.stringify({
                status: "empty",
                message: "No spec document found for this repository",
              }),
            });
          }

          // Convert completed specDocument to repoSpecDocument format
          if (options.specDocument.status === "completed" && options.specDocument.data) {
            const repoDoc: RepoSpecDocumentResponse = {
              status: "completed",
              data: {
                ...options.specDocument.data,
                commitSha: "abc123def", // Default commitSha for backward compatibility
              },
            };
            return route.fulfill({
              status: 200,
              contentType: "application/json",
              body: JSON.stringify(repoDoc),
            });
          }

          // For generating state, return empty
          return route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify({
              status: "empty",
              message: "Spec document is being generated",
            }),
          });
        }

        // Fallback: return empty
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({ status: "empty" }),
        });
      },
    });
  }

  // /api/spec-view/repository/{owner}/{repo}/versions - Get version history by repository
  // Also handles backward compatibility: converts versionHistory to repoVersionHistory format
  if (options.repoVersionHistory || options.versionHistory) {
    handlers.push({
      pattern: /\/api\/spec-view\/repository\/([^/]+)\/([^/]+)\/versions/,
      handler: async (route) => {
        // Prefer repoVersionHistory if provided
        if (options.repoVersionHistory) {
          return route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify(options.repoVersionHistory),
          });
        }

        // Backward compatibility: convert versionHistory to repoVersionHistory format
        if (options.versionHistory) {
          const repoVersionHistory: RepoVersionHistoryResponse = {
            data: options.versionHistory.data.map((v, index) => ({
              id: `doc-${v.version}`,
              analysisId: "550e8400-e29b-41d4-a716-446655440000",
              version: v.version,
              language: options.versionHistory!.language,
              modelId: v.modelId,
              createdAt: v.createdAt,
              commitSha: index === 0 ? "abc123def" : "older123",
            })),
            language: options.versionHistory.language,
          };
          return route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify(repoVersionHistory),
          });
        }

        // Fallback: empty list
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({ data: [], language: "English" }),
        });
      },
    });
  }

  // Register all handlers
  for (const { pattern, handler } of handlers) {
    await page.route(pattern, handler);
  }
}

/**
 * Clear all mock handlers
 */
export async function clearMockHandlers(page: Page): Promise<void> {
  await page.unrouteAll();
}

/**
 * Create a dynamic repositories handler for filtering/sorting scenarios
 */
export function createRepositoriesHandler(
  baseData: UserAnalyzedRepositoriesResponse,
  filters: {
    sortedByName?: UserAnalyzedRepositoriesResponse;
    sortedByTests?: UserAnalyzedRepositoriesResponse;
    filteredMine?: UserAnalyzedRepositoriesResponse;
    filteredOrg?: UserAnalyzedRepositoriesResponse;
    filteredOthers?: UserAnalyzedRepositoriesResponse;
    searchResults?: Record<string, UserAnalyzedRepositoriesResponse>;
  } = {}
): (url: URL) => UserAnalyzedRepositoriesResponse {
  return (url: URL) => {
    const sortBy = url.searchParams.get("sortBy");
    const ownership = url.searchParams.get("ownership");
    const search = url.searchParams.get("search") || url.searchParams.get("q");

    // Handle sorting
    if (sortBy === "name" && filters.sortedByName) {
      return filters.sortedByName;
    }
    if (sortBy === "tests" && filters.sortedByTests) {
      return filters.sortedByTests;
    }

    // Handle ownership filter
    if (ownership === "mine" && filters.filteredMine) {
      return filters.filteredMine;
    }
    if (ownership === "organization" && filters.filteredOrg) {
      return filters.filteredOrg;
    }
    if (ownership === "others" && filters.filteredOthers) {
      return filters.filteredOthers;
    }

    // Handle search
    if (search && filters.searchResults?.[search]) {
      return filters.searchResults[search];
    }

    return baseData;
  };
}
