import type { BackgroundTaskStatus } from "../task-store";

export const isTaskActive = (status: BackgroundTaskStatus): boolean =>
  status === "queued" || status === "processing";

export { formatElapsedTime } from "./format-elapsed-time";
export { getTaskPageUrl } from "./get-task-page-url";
