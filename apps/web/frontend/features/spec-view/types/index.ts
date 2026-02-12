import type { components } from "@/lib/api/generated-types";

// Flat spec item types for virtualization
export type {
  FlatSpecBehaviorItem,
  FlatSpecDomainItem,
  FlatSpecFeatureItem,
  FlatSpecItem,
  FlatSpecItemType,
} from "./flat-spec-item";

// Spec View domain types
export type SpecDocument = components["schemas"]["SpecDocument"];
export type SpecDomain = components["schemas"]["SpecDomain"];
export type SpecFeature = components["schemas"]["SpecFeature"];
export type SpecBehavior = components["schemas"]["SpecBehavior"];
export type SpecBehaviorSourceInfo = components["schemas"]["SpecBehaviorSourceInfo"];
export type SpecLanguage = components["schemas"]["SpecLanguage"];
export type BehaviorCacheStats = components["schemas"]["BehaviorCacheStats"];

// Response types
export type SpecDocumentResponse = components["schemas"]["SpecDocumentResponse"];
export type SpecDocumentCompleted = components["schemas"]["SpecDocumentCompleted"];
export type SpecDocumentGenerating = components["schemas"]["SpecDocumentGenerating"];
export type AvailableLanguageInfo = components["schemas"]["AvailableLanguageInfo"];
export type VersionInfo = components["schemas"]["VersionInfo"];
export type VersionHistoryResponse = components["schemas"]["VersionHistoryResponse"];
export type CacheAvailabilityResponse = components["schemas"]["CacheAvailabilityResponse"];

// Repository-based response types (cross-analysis version access)
export type RepoSpecDocumentResponse = components["schemas"]["RepoSpecDocumentResponse"];
export type RepoSpecDocumentCompleted = components["schemas"]["RepoSpecDocumentCompleted"];
export type RepoSpecDocumentEmpty = components["schemas"]["RepoSpecDocumentEmpty"];
export type RepoSpecDocument = components["schemas"]["RepoSpecDocument"];
export type RepoVersionHistoryResponse = components["schemas"]["RepoVersionHistoryResponse"];
export type RepoVersionInfo = components["schemas"]["RepoVersionInfo"];

// Generation types
export type SpecGenerationStatus = components["schemas"]["SpecGenerationStatus"];
export type SpecGenerationStatusEnum = components["schemas"]["SpecGenerationStatusEnum"];
export type SpecGenerationStatusResponse = components["schemas"]["SpecGenerationStatusResponse"];
export type RequestSpecGenerationRequest = components["schemas"]["RequestSpecGenerationRequest"];
export type RequestSpecGenerationResponse = components["schemas"]["RequestSpecGenerationResponse"];
export type SpecGenerationMode = components["schemas"]["SpecGenerationMode"];
export type CachePredictionResponse = components["schemas"]["CachePredictionResponse"];

// Type guards
export const isDocumentCompleted = (
  response: SpecDocumentResponse
): response is SpecDocumentCompleted => response.status === "completed";

export const isDocumentGenerating = (
  response: SpecDocumentResponse
): response is SpecDocumentGenerating => response.status === "generating";

// Type guards for repository-based responses
export const isRepoDocumentCompleted = (
  response: RepoSpecDocumentResponse
): response is RepoSpecDocumentCompleted => response.status === "completed";

export const isRepoDocumentEmpty = (
  response: RepoSpecDocumentResponse
): response is RepoSpecDocumentEmpty => response.status === "empty";
