"use client";

/**
 * React hook for task polling integration.
 *
 * Provides a React-friendly interface to the polling manager,
 * with automatic cleanup on unmount (optional) and reactive state.
 */

import { useEffect, useSyncExternalStore } from "react";

import { pollingManager, type FetchStatusFn } from "./polling-manager";

// Subscription management for useSyncExternalStore
const listeners = new Set<() => void>();
let pollingSnapshot = new Map<string, boolean>();

const notifyListeners = () => {
  // Create new snapshot to trigger re-render
  pollingSnapshot = new Map(pollingSnapshot);
  listeners.forEach((listener) => listener());
};

const subscribe = (listener: () => void) => {
  listeners.add(listener);
  return () => listeners.delete(listener);
};

const getSnapshot = () => pollingSnapshot;
const getServerSnapshot = () => new Map<string, boolean>();

// Configure polling manager with session lifecycle callbacks for React state synchronization.
// This approach avoids module-level mutation of the pollingManager object.
pollingManager.configure({
  onSessionStart: (taskId: string) => {
    pollingSnapshot.set(taskId, true);
    notifyListeners();
  },
  onSessionStop: (taskId: string) => {
    pollingSnapshot.delete(taskId);
    notifyListeners();
  },
});

/**
 * Options for the useTaskPolling hook.
 */
export type UseTaskPollingOptions = {
  /**
   * Whether to automatically stop polling when the component unmounts.
   * Set to false if you want polling to continue in the background.
   * @default false
   */
  stopOnUnmount?: boolean;
};

/**
 * Return type for the useTaskPolling hook.
 */
export type UseTaskPollingResult = {
  /** Whether the task is currently being polled */
  isPolling: boolean;
  /** Start polling for the task */
  start: () => void;
  /** Stop polling for the task */
  stop: () => void;
};

/**
 * React hook for managing task polling.
 *
 * Provides a component-friendly interface to start/stop polling
 * with reactive state updates.
 *
 * @param taskId - Unique identifier for the task
 * @param fetchStatusFn - Function to fetch the current task status
 * @param options - Hook configuration options
 *
 * @example
 * ```tsx
 * const { isPolling, start, stop } = useTaskPolling(
 *   analysisId,
 *   (signal) => fetchAnalysisStatus(analysisId, signal)
 * );
 *
 * // Start polling when needed
 * start();
 *
 * // Polling continues even if component unmounts (default behavior)
 * ```
 */
export function useTaskPolling(
  taskId: string,
  fetchStatusFn: FetchStatusFn,
  options: UseTaskPollingOptions = {}
): UseTaskPollingResult {
  const { stopOnUnmount = false } = options;

  const snapshot = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
  const isCurrentlyPolling = snapshot.get(taskId) ?? pollingManager.isPolling(taskId);

  const start = () => {
    pollingManager.start(taskId, fetchStatusFn);
  };

  const stop = () => {
    pollingManager.stop(taskId);
  };

  // Handle cleanup on unmount if stopOnUnmount is true
  // Side effect: Stops polling when component unmounts (only if stopOnUnmount is enabled)
  useEffect(() => {
    if (!stopOnUnmount) {
      return;
    }

    return () => {
      if (pollingManager.isPolling(taskId)) {
        pollingManager.stop(taskId);
      }
    };
  }, [taskId, stopOnUnmount]);

  return {
    isPolling: isCurrentlyPolling,
    start,
    stop,
  };
}
