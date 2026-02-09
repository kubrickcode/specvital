export const VIEW_MODES = ["list", "tree", "document"] as const;

export type ViewMode = (typeof VIEW_MODES)[number];

export const DEFAULT_VIEW_MODE: ViewMode = "list";
