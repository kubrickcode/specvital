/**
 * Type definitions for mock API responses
 * Subset of generated-types.ts for E2E testing
 */

// AI Spec Summary
export interface AiSpecSummary {
  hasSpec: boolean;
  languageCount?: number;
  latestGeneratedAt?: string;
}

// Repository Card
export interface RepositoryCard {
  id: string;
  owner: string;
  name: string;
  fullName: string;
  isBookmarked: boolean;
  isAnalyzedByMe: boolean;
  updateStatus: "up-to-date" | "new-commits" | "unknown";
  latestAnalysis?: AnalysisSummary;
  aiSpecSummary?: AiSpecSummary;
}

export interface AnalysisSummary {
  testCount: number;
  change: number;
  analyzedAt: string;
  commitSha: string;
  testSummary?: TestStatusSummary;
}

export interface TestStatusSummary {
  active: number;
  focused: number;
  skipped: number;
  todo: number;
  xfail: number;
}

// User Analyzed Repositories Response
export interface UserAnalyzedRepositoriesResponse {
  data: RepositoryCard[];
  nextCursor: string | null;
  hasNext: boolean;
}

// Repository Stats Response
export interface RepositoryStatsResponse {
  totalRepositories: number;
  totalTests: number;
}

// GitHub Repository
export interface GitHubRepository {
  id: number;
  name: string;
  fullName: string;
  owner: string;
  description?: string;
  defaultBranch: string;
  isPrivate: boolean;
  pushedAt?: string;
}

export interface GitHubRepositoriesResponse {
  data: GitHubRepository[];
}

// GitHub Organization
export interface GitHubOrganization {
  id: number;
  login: string;
  avatarUrl?: string;
  description?: string;
  accessStatus: "accessible" | "restricted" | "pending";
}

export interface GitHubOrganizationsResponse {
  data: GitHubOrganization[];
}

// Subscription
export interface PlanInfo {
  tier: "free" | "pro" | "pro_plus" | "enterprise";
  specviewMonthlyLimit?: number | null;
  analysisMonthlyLimit?: number | null;
  retentionDays?: number | null;
}

export interface UserSubscriptionResponse {
  plan: PlanInfo;
  currentPeriodStart: string;
  currentPeriodEnd: string;
}

// Usage
export interface UsageMetric {
  used: number;
  limit?: number | null;
  percentage?: number | null;
}

export interface UsageStatusResponse {
  specview: UsageMetric;
  analysis: UsageMetric;
  resetAt: string;
  plan: PlanInfo;
}

// Bookmark Response
export interface BookmarkResponse {
  success: boolean;
  isBookmarked: boolean;
}

// Reanalyze Response
export interface ReanalyzeResponse {
  status: "queued" | "analyzing";
}

// Analysis Status Response (for polling)
export interface AnalysisStatusResponse {
  owner: string;
  repo: string;
  status: "queued" | "analyzing" | "completed" | "failed";
  commitSha?: string;
  analysisId?: string;
}

// Update Status Response (for update banner)
export interface UpdateStatusResponse {
  status: "up-to-date" | "new-commits" | "unknown";
  parserOutdated: boolean;
}
