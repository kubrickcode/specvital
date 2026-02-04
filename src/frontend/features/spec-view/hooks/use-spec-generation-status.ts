"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useRef } from "react";

import { paginatedRepositoriesKeys, repositoryStatsKeys } from "@/features/dashboard";
import { getTask, updateTask } from "@/lib/background-tasks";

import { fetchGenerationStatus, ForbiddenError, UnauthorizedError } from "../api";
import type { SpecGenerationStatusEnum, SpecLanguage } from "../types";
import { repoSpecViewKeys } from "./use-repo-spec-view";
import { specViewKeys } from "./use-spec-view";
import { getSpecGenerationTaskId } from "../utils/task-ids";

const POLLING_INTERVAL_MS = 1000;
const MAX_POLLING_DURATION_MS = 5 * 60 * 1000; // 5 minutes

export const specGenerationStatusKeys = {
  all: ["spec-generation-status"] as const,
  status: (analysisId: string, language?: SpecLanguage) =>
    language
      ? ([...specGenerationStatusKeys.all, analysisId, language] as const)
      : ([...specGenerationStatusKeys.all, analysisId] as const),
};

type UseSpecGenerationStatusOptions = {
  enabled?: boolean;
  language?: SpecLanguage;
  onCompleted?: () => void;
  onFailed?: () => void;
  owner?: string;
  repo?: string;
};

type UseSpecGenerationStatusReturn = {
  error: Error | null;
  isPolling: boolean;
  status: SpecGenerationStatusEnum | null;
};

/**
 * Dedicated hook for polling spec generation status.
 * - Polls GET /api/spec-view/status/{analysisId} endpoint
 * - Only polls when status is pending/running
 * - Stops polling and invalidates document queries on completed/failed
 */
export const useSpecGenerationStatus = (
  analysisId: string,
  options: UseSpecGenerationStatusOptions = {}
): UseSpecGenerationStatusReturn => {
  const { enabled = false, language, onCompleted, onFailed, owner, repo } = options;
  const queryClient = useQueryClient();
  const pollingStartTimeRef = useRef<number | null>(null);
  const previousStatusRef = useRef<SpecGenerationStatusEnum | null>(null);

  // Store callbacks in refs to avoid dependency issues
  const onCompletedRef = useRef(onCompleted);
  const onFailedRef = useRef(onFailed);
  onCompletedRef.current = onCompleted;
  onFailedRef.current = onFailed;

  // Reset refs when dependencies change or when a new polling session begins.
  // Tracking `enabled` ensures refs are cleared for same-language regeneration,
  // where analysisId and language stay constant but polling is disabled then
  // re-enabled.
  useEffect(() => {
    if (enabled) {
      pollingStartTimeRef.current = null;
      previousStatusRef.current = null;
    }
  }, [analysisId, enabled, language]);

  const query = useQuery({
    enabled: enabled && Boolean(analysisId),
    queryFn: () => fetchGenerationStatus(analysisId, language),
    queryKey: specGenerationStatusKeys.status(analysisId, language),
    refetchInterval: (query) => {
      // Initialize polling start time on first evaluation
      if (pollingStartTimeRef.current === null) {
        pollingStartTimeRef.current = Date.now();
      }

      // Stop polling after max duration regardless of status
      if (Date.now() - pollingStartTimeRef.current > MAX_POLLING_DURATION_MS) {
        pollingStartTimeRef.current = null;
        return false;
      }

      const status = query.state.data?.status;

      // Stop polling on terminal states
      if (status === "completed" || status === "failed") {
        pollingStartTimeRef.current = null;
        return false;
      }

      // Keep polling for: pending, running, not_found (job may not exist yet),
      // undefined (first fetch or error recovery)
      return POLLING_INTERVAL_MS;
    },
    retry: 2,
    staleTime: 0,
  });

  const status = query.data?.status ?? null;

  // Handle status transitions
  useEffect(() => {
    if (!enabled || !status) return;

    const prevStatus = previousStatusRef.current;

    // Sync task store with server status for accurate UI feedback
    if (status === "pending" || status === "running") {
      const taskId = getSpecGenerationTaskId(analysisId);
      const task = getTask(taskId);
      if (task) {
        const needsStartedAt = !task.startedAt;
        const needsProcessing = status === "running" && task.status !== "processing";
        if (needsStartedAt || needsProcessing) {
          updateTask(taskId, {
            ...(needsStartedAt && { startedAt: new Date().toISOString() }),
            ...(needsProcessing && { status: "processing" as const }),
          });
        }
      }
    }

    // Detect completion transition
    if (status === "completed" && prevStatus !== "completed") {
      // Invalidate document queries for fresh data
      queryClient.invalidateQueries({
        queryKey: specViewKeys.document(analysisId, language),
      });
      queryClient.invalidateQueries({
        queryKey: specViewKeys.document(analysisId),
      });
      // Invalidate repo-based queries for immediate document display
      if (owner && repo) {
        queryClient.invalidateQueries({
          queryKey: repoSpecViewKeys.document(owner, repo, language),
        });
        queryClient.invalidateQueries({
          queryKey: repoSpecViewKeys.document(owner, repo),
        });
      }
      queryClient.invalidateQueries({
        queryKey: repoSpecViewKeys.all,
      });
      // Reset to initial state: shows skeleton on mounted dashboard + fresh fetch on navigation
      queryClient.resetQueries({
        queryKey: paginatedRepositoriesKeys.all,
      });
      queryClient.invalidateQueries({
        queryKey: repositoryStatsKeys.all,
      });
      // Invalidate usage for fresh quota display in re-analyze modal
      queryClient.invalidateQueries({
        queryKey: ["user", "usage"],
      });

      onCompletedRef.current?.();
    }

    // Detect failure transition
    if (status === "failed" && prevStatus !== "failed") {
      onFailedRef.current?.();
    }

    previousStatusRef.current = status;
  }, [analysisId, enabled, language, owner, queryClient, repo, status]);

  const isPolling = status === "pending" || status === "running";

  // Determine error type
  const error = (() => {
    if (query.error instanceof UnauthorizedError) return query.error;
    if (query.error instanceof ForbiddenError) return query.error;
    return query.error ?? null;
  })();

  return {
    error,
    isPolling,
    status,
  };
};
