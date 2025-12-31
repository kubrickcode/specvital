"use client";

import { useInfiniteQuery } from "@tanstack/react-query";

import type { RepositoryCard } from "@/lib/api/types";

import { fetchPaginatedRepositories, type PaginatedRepositoriesParams } from "../api";

const DEFAULT_LIMIT = 20;

export type PaginatedRepositoriesOptions = Omit<PaginatedRepositoriesParams, "cursor">;

export const paginatedRepositoriesKeys = {
  all: ["paginatedRepositories"] as const,
  list: (options: PaginatedRepositoriesOptions) =>
    [...paginatedRepositoriesKeys.all, "list", options] as const,
};

type UsePaginatedRepositoriesReturn = {
  data: RepositoryCard[];
  error: Error | null;
  fetchNextPage: () => void;
  hasNextPage: boolean;
  isError: boolean;
  isFetchingNextPage: boolean;
  isLoading: boolean;
};

export const usePaginatedRepositories = (
  options: PaginatedRepositoriesOptions = {}
): UsePaginatedRepositoriesReturn => {
  const { limit = DEFAULT_LIMIT, sortBy, sortOrder, view } = options;

  const query = useInfiniteQuery({
    getNextPageParam: (lastPage: { hasNext: boolean; nextCursor?: string | null }) =>
      lastPage.hasNext ? lastPage.nextCursor : undefined,
    initialPageParam: undefined as string | undefined,
    queryFn: ({ pageParam }) =>
      fetchPaginatedRepositories({
        cursor: pageParam,
        limit,
        sortBy,
        sortOrder,
        view,
      }),
    queryKey: paginatedRepositoriesKeys.list({ limit, sortBy, sortOrder, view }),
    staleTime: 30 * 1000,
  });

  const data = query.data?.pages.flatMap((page) => page.data) ?? [];

  return {
    data,
    error: query.error,
    fetchNextPage: query.fetchNextPage,
    hasNextPage: query.hasNextPage,
    isError: query.isError,
    isFetchingNextPage: query.isFetchingNextPage,
    isLoading: query.isPending,
  };
};
