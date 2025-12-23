"use client";

import { useQuery } from "@tanstack/react-query";

import type { RepositoryCard } from "@/lib/api/types";

import { fetchBookmarkedRepositories } from "../api";

export const bookmarkKeys = {
  all: ["bookmarks"] as const,
  list: () => [...bookmarkKeys.all, "list"] as const,
};

type UseBookmarkedRepositoriesReturn = {
  data: RepositoryCard[];
  error: Error | null;
  isLoading: boolean;
};

export const useBookmarkedRepositories = (): UseBookmarkedRepositoriesReturn => {
  const query = useQuery({
    queryFn: fetchBookmarkedRepositories,
    queryKey: bookmarkKeys.list(),
    staleTime: 30 * 1000,
  });

  return {
    data: query.data?.data ?? [],
    error: query.error,
    isLoading: query.isPending,
  };
};
