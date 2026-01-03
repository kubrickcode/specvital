export const invalidationEvents = {
  ANALYSIS_COMPLETED: "analysis:completed",
  BOOKMARK_CHANGED: "bookmark:changed",
} as const;

export type InvalidationEvent = (typeof invalidationEvents)[keyof typeof invalidationEvents];

export const invalidationRegistry: Record<InvalidationEvent, readonly unknown[][]> = {
  [invalidationEvents.ANALYSIS_COMPLETED]: [["paginatedRepositories"]],
  [invalidationEvents.BOOKMARK_CHANGED]: [["paginatedRepositories"]],
};
