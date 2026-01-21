export {
  BehaviorItem,
  DocumentView,
  DomainSection,
  DomainStatsBadge,
  EmptyDocument,
  ExecutiveSummary,
  FeatureGroup,
  GenerationProgressModal,
  GenerationStatus,
  QuotaConfirmDialog,
  QuotaIndicator,
  TocSidebar,
} from "./components";

export {
  specViewKeys,
  useDocumentFilter,
  useGenerationProgress,
  useQuotaConfirmDialog,
  useScrollSync,
  useSpecView,
} from "./hooks";
export type { GenerationState } from "./hooks";

export { calculateDocumentStats, calculateDomainStats, isQuotaExceeded } from "./utils";

export type {
  SpecBehavior,
  SpecDocument,
  SpecDomain,
  SpecFeature,
  SpecGenerationStatusEnum,
  SpecLanguage,
} from "./types";
