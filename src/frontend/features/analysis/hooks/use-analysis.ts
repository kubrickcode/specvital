"use client";

import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef } from "react";

import type { AnalysisResponse, AnalysisResult } from "@/lib/api/types";

import { fetchAnalysis } from "../api";

class AnalysisTimeoutError extends Error {
  constructor() {
    super("Analysis timeout: exceeded 5 minute limit");
    this.name = "AnalysisTimeoutError";
  }
}

const INITIAL_INTERVAL_MS = 1000;
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
  status: AnalysisResponse["status"] | "error" | "pending";
};

export const useAnalysis = (owner: string, repo: string): UseAnalysisReturn => {
  const intervalRef = useRef(INITIAL_INTERVAL_MS);
  const startTimeRef = useRef(Date.now());

  useEffect(() => {
    intervalRef.current = INITIAL_INTERVAL_MS;
    startTimeRef.current = Date.now();
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

  const isLoading =
    query.isPending || status === "queued" || status === "analyzing" || status === "pending";

  return {
    data,
    error: query.error,
    isLoading,
    refetch,
    status,
  };
};
