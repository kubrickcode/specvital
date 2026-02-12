import type { SpecGenerationStatusEnum } from "../types";

export const isGenerationActive = (status: SpecGenerationStatusEnum): boolean =>
  status === "pending" || status === "running";

export {
  copySpecToClipboard,
  downloadSpecMarkdown,
  exportSpecToMarkdown,
  generateSpecFilename,
} from "./export-spec-markdown";
export { flattenSpecDocument } from "./flatten-spec-document";
export { findHighlightRanges, type HighlightRange } from "./highlight";
export { formatQuotaNumber, getQuotaLevel, isQuotaExceeded, type QuotaLevel } from "./quota";
export { calculateDocumentStats, calculateDomainStats } from "./stats";
export { getSpecGenerationTaskId } from "./task-ids";
