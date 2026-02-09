"use client";

import { useQuery } from "@tanstack/react-query";

import type { AnalysisHistoryItem } from "../api";
import { fetchAnalysisHistory } from "../api";

export const analysisHistoryKeys = {
  all: ["analysisHistory"] as const,
  detail: (owner: string, repo: string) => [...analysisHistoryKeys.all, owner, repo] as const,
};

type UseAnalysisHistoryOptions = {
  enabled?: boolean;
};

type UseAnalysisHistoryReturn = {
  data: AnalysisHistoryItem[] | undefined;
  error: Error | null;
  isLoading: boolean;
};

export const useAnalysisHistory = (
  owner: string,
  repo: string,
  options: UseAnalysisHistoryOptions = {}
): UseAnalysisHistoryReturn => {
  const { enabled = true } = options;

  const query = useQuery({
    enabled,
    queryFn: async () => {
      const response = await fetchAnalysisHistory(owner, repo);
      return response.data;
    },
    queryKey: analysisHistoryKeys.detail(owner, repo),
  });

  return {
    data: query.data,
    error: query.error,
    isLoading: query.isLoading,
  };
};
