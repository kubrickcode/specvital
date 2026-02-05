"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
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
  const t = useTranslations("dashboard.toast");
  const queryClient = useQueryClient();

  const query = useQuery({
    queryFn: () => fetchUserGitHubOrganizations(),
    queryKey: organizationsKeys.list(),
  });

  const refresh = async () => {
    try {
      const freshData = await fetchUserGitHubOrganizations({ refresh: true });
      queryClient.setQueryData(organizationsKeys.list(), freshData);
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
