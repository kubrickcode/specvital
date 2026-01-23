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
  SpecAccessError,
  TocSidebar,
} from "./components";

export {
  specViewKeys,
  useDocumentFilter,
  useGenerationProgress,
  useQuotaConfirmDialog,
  useSpecView,
  useVersionHistory,
} from "./hooks";
export type { AccessErrorType, GenerationState } from "./hooks";

export { calculateDocumentStats, calculateDomainStats, isQuotaExceeded } from "./utils";

export type {
  SpecBehavior,
  SpecDocument,
  SpecDomain,
  SpecFeature,
  SpecGenerationStatusEnum,
  SpecLanguage,
} from "./types";
