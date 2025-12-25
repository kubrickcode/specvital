"use client";

import { useInfiniteQuery } from "@tanstack/react-query";

import type {
  OwnershipFilter,
  RepositoryCard,
  UserAnalyzedRepositoriesResponse,
} from "@/lib/api/types";

import { fetchUserAnalyzedRepositories } from "../api";

export const myAnalysesKeys = {
  all: ["my-analyses"] as const,
  list: (ownership: OwnershipFilter) => [...myAnalysesKeys.all, "list", ownership] as const,
};

type UseMyAnalysesParams = {
  limit?: number;
  ownership?: OwnershipFilter;
};

type UseMyAnalysesReturn = {
  data: RepositoryCard[];
  error: Error | null;
  fetchNextPage: () => void;
  hasNextPage: boolean;
  isFetchingNextPage: boolean;
  isLoading: boolean;
  refetch: () => void;
};

export const useMyAnalyses = (params?: UseMyAnalysesParams): UseMyAnalysesReturn => {
  const ownership = params?.ownership ?? "all";
  const limit = params?.limit ?? 20;

  const query = useInfiniteQuery({
    getNextPageParam: (lastPage: UserAnalyzedRepositoriesResponse) =>
      lastPage.hasNext ? lastPage.nextCursor : undefined,
    initialPageParam: undefined as string | undefined,
    queryFn: ({ pageParam }) =>
      fetchUserAnalyzedRepositories({
        cursor: pageParam,
        limit,
        ownership,
      }),
    queryKey: myAnalysesKeys.list(ownership),
    staleTime: 30 * 1000,
  });

  const data = query.data?.pages.flatMap((page) => page.data) ?? [];

  return {
    data,
    error: query.error,
    fetchNextPage: query.fetchNextPage,
    hasNextPage: query.hasNextPage,
    isFetchingNextPage: query.isFetchingNextPage,
    isLoading: query.isPending,
    refetch: query.refetch,
  };
};
