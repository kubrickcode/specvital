// Task store hooks
export { useActiveTaskCount, useBackgroundTask, useBackgroundTasks } from "./hooks";

// Polling manager and hook
export { pollingManager } from "./polling-manager";
export type { FetchStatusFn, PollingManagerConfig, TaskStatusResponse } from "./polling-manager";
export { useTaskPolling } from "./use-task-polling";
export type { UseTaskPollingOptions, UseTaskPollingResult } from "./use-task-polling";

// Task store
export {
  addTask,
  type BackgroundTask,
  type BackgroundTaskMetadata,
  type BackgroundTaskStatus,
  type BackgroundTaskType,
  clearCompletedTasks,
  getAllTasks,
  getTask,
  removeTask,
  updateTask,
} from "./task-store";
