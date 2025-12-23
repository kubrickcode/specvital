"use client";

import { useQuery } from "@tanstack/react-query";

import type { RepositoryStatsResponse } from "@/lib/api/types";

import { fetchRepositoryStats } from "../api";

export const repositoryStatsKeys = {
  all: ["repositoryStats"] as const,
  stats: () => [...repositoryStatsKeys.all, "stats"] as const,
};

type UseRepositoryStatsReturn = {
  data: RepositoryStatsResponse | null;
  error: Error | null;
  isLoading: boolean;
};

export const useRepositoryStats = (): UseRepositoryStatsReturn => {
  const query = useQuery({
    queryFn: fetchRepositoryStats,
    queryKey: repositoryStatsKeys.stats(),
    staleTime: 60 * 1000,
  });

  return {
    data: query.data ?? null,
    error: query.error,
    isLoading: query.isPending,
  };
};
