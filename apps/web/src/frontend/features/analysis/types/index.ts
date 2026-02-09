export type { FileTreeNode, FileTreeNodeType, FlatTreeItem } from "./file-tree";
export type { FilterState } from "./filter";

// Analysis Status
export type { AnalysisStatus, WaitingStatus } from "./analysis-status";

// Primary Tab (Tests vs AI Spec)
export { DEFAULT_PRIMARY_TAB, PRIMARY_TABS } from "./primary-tab";
export type { PrimaryTab } from "./primary-tab";

// Data View Mode (list vs tree) - for Tests tab
export { DATA_VIEW_MODES, DEFAULT_DATA_VIEW_MODE } from "./data-view-mode";
export type { DataViewMode } from "./data-view-mode";

// Legacy exports - kept for backward compatibility
// @deprecated Use PrimaryTab and DataViewMode instead
export { DEFAULT_VIEW_MODE, VIEW_MODES } from "./view-mode";
export type { ViewMode } from "./view-mode";
