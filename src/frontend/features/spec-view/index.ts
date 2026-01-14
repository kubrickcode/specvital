export {
  BehaviorItem,
  DocumentView,
  DomainSection,
  DomainStatsBadge,
  EmptyDocument,
  ExecutiveSummary,
  FeatureGroup,
  GenerationStatus,
  TocSidebar,
} from "./components";

export { specViewKeys, useDocumentFilter, useScrollSync, useSpecView } from "./hooks";

export { calculateDocumentStats, calculateDomainStats } from "./utils";

export type {
  SpecBehavior,
  SpecDocument,
  SpecDomain,
  SpecFeature,
  SpecGenerationStatusEnum,
  SpecLanguage,
} from "./types";
