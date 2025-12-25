"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import type { GitHubOrganization } from "@/lib/api/types";

import { fetchUserGitHubOrganizations } from "../api";

export const organizationsKeys = {
  all: ["organizations"] as const,
  list: () => [...organizationsKeys.all, "list"] as const,
};

type UseOrganizationsReturn = {
  data: GitHubOrganization[];
  error: Error | null;
  isLoading: boolean;
  isRefreshing: boolean;
  refresh: () => Promise<void>;
};

export const useOrganizations = (): UseOrganizationsReturn => {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryFn: () => fetchUserGitHubOrganizations(),
    queryKey: organizationsKeys.list(),
    staleTime: 5 * 60 * 1000,
  });

  const refresh = async () => {
    try {
      const freshData = await fetchUserGitHubOrganizations({ refresh: true });
      queryClient.setQueryData(organizationsKeys.list(), freshData);
      toast.success("Refreshed from GitHub");
    } catch (error) {
      toast.error("Failed to refresh", {
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
