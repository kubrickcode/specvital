"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useRef } from "react";

import { paginatedRepositoriesKeys } from "@/features/dashboard";
import type { AnalysisResponse, AnalysisResult } from "@/lib/api/types";
import { userActiveTasksKeys } from "@/lib/background-tasks";

import { fetchAnalysis } from "../api";

class AnalysisTimeoutError extends Error {
  constructor() {
    super("Analysis timeout: exceeded 5 minute limit");
    this.name = "AnalysisTimeoutError";
  }
}

const INITIAL_INTERVAL_MS = 300;
const MAX_INTERVAL_MS = 1000;
const BACKOFF_MULTIPLIER = 1.5;
const MAX_WAIT_MS = 5 * 60 * 1000; // 5 minutes

export const analysisKeys = {
  all: ["analysis"] as const,
  detail: (owner: string, repo: string, commit?: string | null) =>
    commit
      ? ([...analysisKeys.all, owner, repo, commit] as const)
      : ([...analysisKeys.all, owner, repo] as const),
};

const isTerminalStatus = (response: AnalysisResponse): boolean =>
  response.status === "completed" || response.status === "failed";

type UseAnalysisOptions = {
  commit?: string | null;
  enabled?: boolean;
};

type UseAnalysisReturn = {
  data: AnalysisResult | null;
  error: Error | null;
  isLoading: boolean;
  refetch: () => void;
  startedAt: string | null;
  status: AnalysisResponse["status"] | "error" | "pending";
};

export const useAnalysis = (
  owner: string,
  repo: string,
  options: UseAnalysisOptions = {}
): UseAnalysisReturn => {
  const { commit, enabled = true } = options;
  const queryClient = useQueryClient();
  const intervalRef = useRef(INITIAL_INTERVAL_MS);
  const startTimeRef = useRef(Date.now());
  const hasTriggeredInvalidation = useRef(false);
  // Track when waiting started (queued or analyzing) for elapsed time display
  const waitingStartedAtRef = useRef<string | null>(null);

  useEffect(() => {
    intervalRef.current = INITIAL_INTERVAL_MS;
    startTimeRef.current = Date.now();
    hasTriggeredInvalidation.current = false;
    waitingStartedAtRef.current = null;
  }, [owner, repo, commit]);

  // Specific commit query: no polling (already completed)
  const isSpecificCommitQuery = !!commit;

  const query = useQuery({
    enabled,
    queryFn: async () => {
      if (!isSpecificCommitQuery && Date.now() - startTimeRef.current > MAX_WAIT_MS) {
        throw new AnalysisTimeoutError();
      }
      return fetchAnalysis(owner, repo, { commit });
    },
    queryKey: analysisKeys.detail(owner, repo, commit),
    refetchInterval: (query) => {
      // Specific commit queries don't need polling
      if (isSpecificCommitQuery) {
        return false;
      }

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
  });

  const response = query.data;

  const refetch = () => {
    intervalRef.current = INITIAL_INTERVAL_MS;
    startTimeRef.current = Date.now();
    waitingStartedAtRef.current = null;
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

  // Track when analysis enters waiting state (queued or analyzing)
  // Only used as fallback when server doesn't provide startedAt
  useEffect(() => {
    const isWaitingStatus = status === "queued" || status === "analyzing";
    if (isWaitingStatus && !waitingStartedAtRef.current) {
      waitingStartedAtRef.current = new Date().toISOString();
    }
  }, [status]);

  useEffect(() => {
    if (response?.status === "completed" && !hasTriggeredInvalidation.current) {
      hasTriggeredInvalidation.current = true;
      queryClient.removeQueries({ queryKey: paginatedRepositoriesKeys.all });
      // Invalidate usage for fresh quota display
      queryClient.invalidateQueries({ queryKey: ["user", "usage"] });
      // Refresh active tasks list from server
      queryClient.invalidateQueries({ queryKey: userActiveTasksKeys.all });
    }
  }, [response?.status, queryClient]);

  const isLoading =
    query.isPending || status === "queued" || status === "analyzing" || status === "pending";

  // Prefer server startedAt (analyzing state), fallback to client tracking (queued state)
  const serverStartedAt = response?.status === "analyzing" ? response.startedAt : undefined;
  const startedAt = serverStartedAt ?? waitingStartedAtRef.current;

  return {
    data,
    error: query.error,
    isLoading,
    refetch,
    startedAt,
    status,
  };
};
