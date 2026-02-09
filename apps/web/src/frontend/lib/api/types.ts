/**
 * Re-export types from OpenAPI generated types.
 * This provides backward-compatible type aliases for existing imports.
 *
 * @see generated-types.ts (auto-generated from backend/api/openapi.yaml)
 */

import type { components, operations } from "./generated-types";

// Core domain types
export type AnalysisResult = components["schemas"]["AnalysisResult"];
export type TestSuite = components["schemas"]["TestSuite"];
export type TestCase = components["schemas"]["TestCase"];
export type TestStatus = components["schemas"]["TestStatus"];
export type Framework = components["schemas"]["Framework"];
export type Summary = components["schemas"]["Summary"];
export type FrameworkSummary = components["schemas"]["FrameworkSummary"];

// API response types
export type AnalysisResponse = components["schemas"]["AnalysisResponse"];
export type CompletedResponse = components["schemas"]["CompletedResponse"];
export type AnalyzingResponse = components["schemas"]["AnalyzingResponse"];
export type QueuedResponse = components["schemas"]["QueuedResponse"];
export type FailedResponse = components["schemas"]["FailedResponse"];

// Error types (RFC 7807)
export type ProblemDetail = components["schemas"]["ProblemDetail"];
export type RateLimitInfo = components["schemas"]["RateLimitInfo"];

// Auth types
export type LoginResponse = components["schemas"]["LoginResponse"];
export type LogoutResponse = components["schemas"]["LogoutResponse"];
export type RefreshResponse = components["schemas"]["RefreshResponse"];
export type UserInfo = components["schemas"]["UserInfo"];

// Convenience type for analysis status
export type AnalysisStatus = AnalysisResponse["status"];

// Dashboard types
export type RepositoryCard = components["schemas"]["RepositoryCard"];
export type AnalysisSummary = components["schemas"]["AnalysisSummary"];
export type UpdateStatus = components["schemas"]["UpdateStatus"];
export type UpdateStatusResponse = components["schemas"]["UpdateStatusResponse"];
export type BookmarkResponse = components["schemas"]["BookmarkResponse"];
export type BookmarkedRepositoriesResponse =
  components["schemas"]["BookmarkedRepositoriesResponse"];
export type PaginatedRepositoriesResponse = components["schemas"]["PaginatedRepositoriesResponse"];
export type SortByParam = components["schemas"]["SortByParam"];
export type SortOrderParam = components["schemas"]["SortOrderParam"];
export type ViewFilterParam = components["schemas"]["ViewFilterParam"];
export type OwnershipFilterParam = components["schemas"]["OwnershipFilterParam"];
export type RepositoryStatsResponse = components["schemas"]["RepositoryStatsResponse"];
export type UserAnalyzedRepositoriesResponse =
  components["schemas"]["UserAnalyzedRepositoriesResponse"];

// Ownership filter type - derived from OpenAPI spec
export const OWNERSHIP_FILTER_VALUES = ["all", "mine", "organization"] as const;
export type OwnershipFilter = NonNullable<
  NonNullable<operations["getUserAnalyzedRepositories"]["parameters"]["query"]>["ownership"]
>;

// GitHub types
export type GitHubOrganization = components["schemas"]["GitHubOrganization"];
export type GitHubOrganizationsResponse = components["schemas"]["GitHubOrganizationsResponse"];
export type GitHubRepository = components["schemas"]["GitHubRepository"];
export type GitHubRepositoriesResponse = components["schemas"]["GitHubRepositoriesResponse"];
export type OrganizationAccessStatus = components["schemas"]["OrganizationAccessStatus"];

// GitHub App types
export type GitHubAppInstallation = components["schemas"]["GitHubAppInstallation"];
export type GitHubAppInstallationsResponse =
  components["schemas"]["GitHubAppInstallationsResponse"];
export type GitHubAppInstallUrlResponse = components["schemas"]["GitHubAppInstallUrlResponse"];

// AI Spec types
export type AiSpecSummary = components["schemas"]["AiSpecSummary"];

// Active Tasks types
export type ActiveTask = components["schemas"]["ActiveTask"];
export type UserActiveTasksResponse = components["schemas"]["UserActiveTasksResponse"];
