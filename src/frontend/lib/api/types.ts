/**
 * Re-export types from OpenAPI generated types.
 * This provides backward-compatible type aliases for existing imports.
 *
 * @see generated-types.ts (auto-generated from backend/api/openapi.yaml)
 */

import type { components } from "./generated-types";

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

// Convenience type for analysis status
export type AnalysisStatus = AnalysisResponse["status"];
