/**
 * Data view modes for Tests tab (list or tree view)
 * Separated from PrimaryTab to clarify conceptual model
 */
export const DATA_VIEW_MODES = ["list", "tree"] as const;

export type DataViewMode = (typeof DATA_VIEW_MODES)[number];

export const DEFAULT_DATA_VIEW_MODE: DataViewMode = "list";
