"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";

import type { GitHubRepository } from "@/lib/api/types";

import { fetchUserGitHubRepositories } from "../api";

export const myRepositoriesKeys = {
  all: ["my-repositories"] as const,
  list: () => [...myRepositoriesKeys.all, "list"] as const,
};

type UseMyRepositoriesReturn = {
  data: GitHubRepository[];
  error: Error | null;
  isLoading: boolean;
  isRefreshing: boolean;
  refresh: () => Promise<void>;
};

export const useMyRepositories = (): UseMyRepositoriesReturn => {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryFn: () => fetchUserGitHubRepositories(),
    queryKey: myRepositoriesKeys.list(),
    staleTime: 5 * 60 * 1000,
  });

  const refresh = async () => {
    const freshData = await fetchUserGitHubRepositories({ refresh: true });
    queryClient.setQueryData(myRepositoriesKeys.list(), freshData);
  };

  return {
    data: query.data?.data ?? [],
    error: query.error,
    isLoading: query.isPending,
    isRefreshing: query.isFetching && !query.isPending,
    refresh,
  };
};
