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

// Response types
export type SpecDocumentResponse = components["schemas"]["SpecDocumentResponse"];
export type SpecDocumentCompleted = components["schemas"]["SpecDocumentCompleted"];
export type SpecDocumentGenerating = components["schemas"]["SpecDocumentGenerating"];

// Generation types
export type SpecGenerationStatus = components["schemas"]["SpecGenerationStatus"];
export type SpecGenerationStatusEnum = components["schemas"]["SpecGenerationStatusEnum"];
export type SpecGenerationStatusResponse = components["schemas"]["SpecGenerationStatusResponse"];
export type RequestSpecGenerationRequest = components["schemas"]["RequestSpecGenerationRequest"];
export type RequestSpecGenerationResponse = components["schemas"]["RequestSpecGenerationResponse"];

// Type guards
export const isDocumentCompleted = (
  response: SpecDocumentResponse
): response is SpecDocumentCompleted => response.status === "completed";

export const isDocumentGenerating = (
  response: SpecDocumentResponse
): response is SpecDocumentGenerating => response.status === "generating";
