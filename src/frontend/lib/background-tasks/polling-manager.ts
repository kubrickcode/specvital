/**
 * Polling Manager for Background Tasks
 *
 * A singleton module that manages polling for background task status updates.
 * Designed to work independently of React component lifecycle, ensuring
 * polling continues even when components unmount.
 */

import type { BackgroundTaskStatus } from "./task-store";

// Default polling intervals in milliseconds
const DEFAULT_POLLING_INTERVAL_PROCESSING = 3000;
const DEFAULT_POLLING_INTERVAL_QUEUED = 5000;

/**
 * Generic response structure for task status polling.
 * Implementations should return data conforming to this interface.
 */
export type TaskStatusResponse = {
  data?: unknown;
  error?: string;
  status: BackgroundTaskStatus;
};

/**
 * Type guard to check if a task has reached a terminal state.
 */
const isTaskTerminal = (status: BackgroundTaskStatus): boolean =>
  status === "completed" || status === "failed";

/**
 * Configuration options for the polling manager.
 */
export type PollingManagerConfig = {
  /** Callback when a task completes successfully */
  onComplete?: (taskId: string, data?: unknown) => void;
  /** Callback when a task fails */
  onError?: (taskId: string, error: Error) => void;
  /** Callback when a polling session starts */
  onSessionStart?: (taskId: string) => void;
  /** Callback when a polling session stops */
  onSessionStop?: (taskId: string) => void;
  /** Callback when a task's status changes */
  onStatusChange?: (taskId: string, status: BackgroundTaskStatus) => void;
  /** Polling interval for tasks in "processing" status (default: 3000ms) */
  pollingIntervalProcessing?: number;
  /** Polling interval for tasks in "queued" status (default: 5000ms) */
  pollingIntervalQueued?: number;
};

/**
 * Function type for fetching task status.
 * Accepts an AbortSignal to allow cancellation of in-flight requests.
 */
export type FetchStatusFn = (signal: AbortSignal) => Promise<TaskStatusResponse>;

/**
 * Internal state for each active polling session.
 */
type PollingSession = {
  abortController: AbortController;
  fetchStatusFn: FetchStatusFn;
  lastStatus: BackgroundTaskStatus | null;
  timeoutId: ReturnType<typeof setTimeout> | null;
};

// Module-level state
let config: PollingManagerConfig = {};
const activeSessions = new Map<string, PollingSession>();

/**
 * Get the appropriate polling interval based on task status.
 */
const getPollingInterval = (status: BackgroundTaskStatus): number => {
  if (status === "processing") {
    return config.pollingIntervalProcessing ?? DEFAULT_POLLING_INTERVAL_PROCESSING;
  }
  return config.pollingIntervalQueued ?? DEFAULT_POLLING_INTERVAL_QUEUED;
};

/**
 * Execute a single poll for a task's status.
 */
const executePoll = async (taskId: string): Promise<void> => {
  const session = activeSessions.get(taskId);
  if (!session) {
    return;
  }

  // Check if aborted before fetching
  if (session.abortController.signal.aborted) {
    return;
  }

  try {
    const response = await session.fetchStatusFn(session.abortController.signal);

    // Check if aborted after fetching (race condition protection)
    if (session.abortController.signal.aborted) {
      return;
    }

    // Handle status change
    if (session.lastStatus !== response.status) {
      session.lastStatus = response.status;
      config.onStatusChange?.(taskId, response.status);
    }

    // Handle terminal states
    if (isTaskTerminal(response.status)) {
      if (response.status === "completed") {
        config.onComplete?.(taskId, response.data);
      } else if (response.status === "failed") {
        const errorMessage = response.error ?? "Task failed with unknown error";
        config.onError?.(taskId, new Error(errorMessage));
      }
      // Stop polling for terminal states
      stopPolling(taskId);
      return;
    }

    // Schedule next poll for non-terminal states
    scheduleNextPoll(taskId, response.status);
  } catch (error) {
    // Check if aborted (fetch was cancelled)
    if (session.abortController.signal.aborted) {
      return;
    }

    // Report error but continue polling
    config.onError?.(taskId, error instanceof Error ? error : new Error("Unknown polling error"));

    // Retry with queued interval on error
    scheduleNextPoll(taskId, "queued");
  }
};

/**
 * Schedule the next poll for a task.
 */
const scheduleNextPoll = (taskId: string, status: BackgroundTaskStatus): void => {
  const session = activeSessions.get(taskId);
  if (!session || session.abortController.signal.aborted) {
    return;
  }

  const interval = getPollingInterval(status);
  session.timeoutId = setTimeout(() => {
    executePoll(taskId);
  }, interval);
};

/**
 * Start polling for a task's status.
 * If already polling for this task, the existing session continues unchanged.
 */
const startPolling = (taskId: string, fetchStatusFn: FetchStatusFn): void => {
  // Don't start a new session if one already exists
  if (activeSessions.has(taskId)) {
    return;
  }

  const session: PollingSession = {
    abortController: new AbortController(),
    fetchStatusFn,
    lastStatus: null,
    timeoutId: null,
  };

  activeSessions.set(taskId, session);
  config.onSessionStart?.(taskId);

  // Execute first poll immediately
  executePoll(taskId);
};

/**
 * Stop polling for a specific task.
 */
const stopPolling = (taskId: string): void => {
  const session = activeSessions.get(taskId);
  if (!session) {
    return;
  }

  // Abort any in-flight requests
  session.abortController.abort();

  // Clear any scheduled polls
  if (session.timeoutId !== null) {
    clearTimeout(session.timeoutId);
  }

  activeSessions.delete(taskId);
  config.onSessionStop?.(taskId);
};

/**
 * Stop all active polling sessions.
 */
const stopAllPolling = (): void => {
  for (const taskId of activeSessions.keys()) {
    stopPolling(taskId);
  }
};

/**
 * Check if a task is currently being polled.
 */
const isPolling = (taskId: string): boolean => {
  return activeSessions.has(taskId);
};

/**
 * Configure the polling manager.
 * Merges provided config with existing config.
 */
const configure = (newConfig: PollingManagerConfig): void => {
  config = { ...config, ...newConfig };
};

/**
 * Get current polling manager configuration.
 * Useful for testing and debugging.
 */
const getConfig = (): Readonly<PollingManagerConfig> => {
  return { ...config };
};

/**
 * Get the count of active polling sessions.
 * Useful for testing and debugging.
 */
const getActiveSessionCount = (): number => {
  return activeSessions.size;
};

/**
 * Reset the polling manager to its initial state.
 * Stops all active sessions and clears configuration.
 * Primarily intended for testing.
 */
const reset = (): void => {
  stopAllPolling();
  config = {};
};

/**
 * Singleton polling manager instance.
 * Manages background task polling independently of React component lifecycle.
 */
export const pollingManager = {
  configure,
  getActiveSessionCount,
  getConfig,
  isPolling,
  reset,
  start: startPolling,
  stop: stopPolling,
  stopAll: stopAllPolling,
} as const;
