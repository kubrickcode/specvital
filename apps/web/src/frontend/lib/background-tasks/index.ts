// Components
export { TaskBadge, TasksDropdownSection, TaskStatusBadge } from "./components";

// Task store hooks (sessionStorage-based - for local page state)
export { useActiveTaskCount, useActiveTasks, useBackgroundTasks, useTaskStartedAt } from "./hooks";

// Task store (sessionStorage-based - for local page state)
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
  useTaskStore,
} from "./task-store";

// Server-based active tasks (for authenticated users)
export { useUserActiveTasks, userActiveTasksKeys } from "./use-user-active-tasks";
