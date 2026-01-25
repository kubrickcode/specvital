// Components
export { TaskBadge, TasksDropdownSection, TaskStatusBadge } from "./components";

// Task store hooks
export { useActiveTaskCount, useBackgroundTask, useBackgroundTasks } from "./hooks";

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
