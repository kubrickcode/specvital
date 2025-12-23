"use client";

import { useQuery } from "@tanstack/react-query";

import type { RepositoryCard } from "@/lib/api/types";

import { fetchRecentRepositories } from "../api";

export const recentRepositoriesKeys = {
  all: ["recentRepositories"] as const,
  list: (limit?: number) => [...recentRepositoriesKeys.all, "list", limit] as const,
};

type UseRecentRepositoriesReturn = {
  data: RepositoryCard[];
  error: Error | null;
  isLoading: boolean;
};

export const useRecentRepositories = (limit?: number): UseRecentRepositoriesReturn => {
  const query = useQuery({
    queryFn: () => fetchRecentRepositories(limit),
    queryKey: recentRepositoriesKeys.list(limit),
    staleTime: 30 * 1000,
  });

  return {
    data: query.data?.data ?? [],
    error: query.error,
    isLoading: query.isPending,
  };
};
