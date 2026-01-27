"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useRef } from "react";

import { paginatedRepositoriesKeys } from "@/features/dashboard";
import type { AnalysisResponse, AnalysisResult } from "@/lib/api/types";
import { addTask, getTask, removeTask, updateTask } from "@/lib/background-tasks";

import { fetchAnalysis } from "../api";

class AnalysisTimeoutError extends Error {
  constructor() {
    super("Analysis timeout: exceeded 5 minute limit");
    this.name = "AnalysisTimeoutError";
  }
}

const INITIAL_INTERVAL_MS = 300;
const MAX_INTERVAL_MS = 5000;
const BACKOFF_MULTIPLIER = 1.5;
const MAX_WAIT_MS = 5 * 60 * 1000; // 5 minutes

export const analysisKeys = {
  all: ["analysis"] as const,
  detail: (owner: string, repo: string) => [...analysisKeys.all, owner, repo] as const,
};

const isTerminalStatus = (response: AnalysisResponse): boolean =>
  response.status === "completed" || response.status === "failed";

type UseAnalysisReturn = {
  data: AnalysisResult | null;
  error: Error | null;
  isLoading: boolean;
  refetch: () => void;
  startedAt: string | null;
  status: AnalysisResponse["status"] | "error" | "pending";
};

export const useAnalysis = (owner: string, repo: string): UseAnalysisReturn => {
  const queryClient = useQueryClient();
  const intervalRef = useRef(INITIAL_INTERVAL_MS);
  const startTimeRef = useRef(Date.now());
  const hasTriggeredInvalidation = useRef(false);

  useEffect(() => {
    intervalRef.current = INITIAL_INTERVAL_MS;
    startTimeRef.current = Date.now();
    hasTriggeredInvalidation.current = false;
  }, [owner, repo]);

  const query = useQuery({
    queryFn: async () => {
      if (Date.now() - startTimeRef.current > MAX_WAIT_MS) {
        throw new AnalysisTimeoutError();
      }
      return fetchAnalysis(owner, repo);
    },
    queryKey: analysisKeys.detail(owner, repo),
    refetchInterval: (query) => {
      const response = query.state.data;

      if (response && isTerminalStatus(response)) {
        return false;
      }

      if (Date.now() - startTimeRef.current > MAX_WAIT_MS) {
        return false;
      }

      const interval = intervalRef.current;
      intervalRef.current = Math.min(interval * BACKOFF_MULTIPLIER, MAX_INTERVAL_MS);
      return interval;
    },
    retry: false,
    staleTime: 0,
  });

  const response = query.data;

  const refetch = () => {
    intervalRef.current = INITIAL_INTERVAL_MS;
    startTimeRef.current = Date.now();
    query.refetch();
  };

  let status: AnalysisResponse["status"] | "error" | "pending";
  if (query.error) {
    status = "error";
  } else if (response) {
    status = response.status;
  } else {
    status = "pending";
  }

  const data = response?.status === "completed" ? response.data : null;

  useEffect(() => {
    if (response?.status === "completed" && !hasTriggeredInvalidation.current) {
      hasTriggeredInvalidation.current = true;
      queryClient.invalidateQueries({ queryKey: paginatedRepositoriesKeys.all });
    }
  }, [response?.status, queryClient]);

  // Sync with global TaskStore for background progress tracking
  const taskId = `analysis-${owner}-${repo}`;
  useEffect(() => {
    if (status === "queued") {
      const existingTask = getTask(taskId);
      if (!existingTask) {
        addTask({
          id: taskId,
          metadata: { owner, repo },
          startedAt: null,
          status: "queued",
          type: "analysis",
        });
      }
    } else if (status === "analyzing") {
      const existingTask = getTask(taskId);
      if (existingTask && existingTask.status !== "processing") {
        updateTask(taskId, { startedAt: new Date().toISOString(), status: "processing" });
      } else if (!existingTask) {
        addTask({
          id: taskId,
          metadata: { owner, repo },
          startedAt: new Date().toISOString(),
          status: "processing",
          type: "analysis",
        });
      }
    } else if (status === "completed" || status === "failed" || status === "error") {
      // Deduplication with AnalysisMonitor (global)
      const existingTask = getTask(taskId);
      if (existingTask) {
        removeTask(taskId);
      }
    }
  }, [status, owner, repo, taskId]);

  const isLoading =
    query.isPending || status === "queued" || status === "analyzing" || status === "pending";

  // Get startedAt from TaskStore with fallback for race condition
  const taskFromStore = getTask(taskId);
  let startedAt = taskFromStore?.startedAt ?? null;

  // Fallback: If analyzing but no startedAt, use current time as approximation
  if (status === "analyzing" && !startedAt) {
    startedAt = new Date().toISOString();
    if (process.env.NODE_ENV === "development") {
      console.warn(`[useAnalysis] TaskStore missing startedAt for ${taskId}, using fallback`);
    }
  }

  return {
    data,
    error: query.error,
    isLoading,
    refetch,
    startedAt,
    status,
  };
};
