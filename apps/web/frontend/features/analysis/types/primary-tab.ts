export const PRIMARY_TABS = ["tests", "spec"] as const;

export type PrimaryTab = (typeof PRIMARY_TABS)[number];

export const DEFAULT_PRIMARY_TAB: PrimaryTab = "tests";
