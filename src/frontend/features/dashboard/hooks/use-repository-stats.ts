"use client";

import { useQuery } from "@tanstack/react-query";

import { useAuth } from "@/features/auth";
import type { RepositoryStatsResponse } from "@/lib/api/types";

import { fetchRepositoryStats } from "../api";

export const repositoryStatsKeys = {
  all: ["repositoryStats"] as const,
};

type UseRepositoryStatsReturn = {
  data: RepositoryStatsResponse | undefined;
  error: Error | null;
  isError: boolean;
  isLoading: boolean;
  refetch: () => void;
};

export const useRepositoryStats = (): UseRepositoryStatsReturn => {
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const isEnabled = isAuthenticated && !isAuthLoading;

  const query = useQuery({
    enabled: isEnabled,
    queryFn: fetchRepositoryStats,
    queryKey: repositoryStatsKeys.all,
    staleTime: 60 * 1000,
  });

  return {
    data: query.data,
    error: query.error,
    isError: query.isError,
    isLoading: query.status === "pending",
    refetch: query.refetch,
  };
};
