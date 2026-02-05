"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { toast } from "sonner";

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
  const t = useTranslations("dashboard.toast");
  const queryClient = useQueryClient();

  const query = useQuery({
    queryFn: () => fetchUserGitHubRepositories(),
    queryKey: myRepositoriesKeys.list(),
  });

  const refresh = async () => {
    try {
      const freshData = await fetchUserGitHubRepositories({ refresh: true });
      queryClient.setQueryData(myRepositoriesKeys.list(), freshData);
      toast.success(t("refreshed"));
    } catch (error) {
      toast.error(t("refreshFailed"), {
        description: error instanceof Error ? error.message : String(error),
      });
    }
  };

  return {
    data: query.data?.data ?? [],
    error: query.error,
    isLoading: query.isPending,
    isRefreshing: query.isFetching && !query.isPending,
    refresh,
  };
};
