"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";

import type { GitHubRepository } from "@/lib/api/types";

import { fetchOrganizationRepositories } from "../api";

export const organizationReposKeys = {
  all: ["organization-repos"] as const,
  byOrg: (org: string) => [...organizationReposKeys.all, org] as const,
};

type UseOrganizationReposParams = {
  org: string;
};

type UseOrganizationReposReturn = {
  data: GitHubRepository[];
  error: Error | null;
  isLoading: boolean;
  isRefreshing: boolean;
  refresh: () => Promise<void>;
};

export const useOrganizationRepos = ({
  org,
}: UseOrganizationReposParams): UseOrganizationReposReturn => {
  const queryClient = useQueryClient();

  const query = useQuery({
    enabled: Boolean(org),
    queryFn: () => fetchOrganizationRepositories(org),
    queryKey: organizationReposKeys.byOrg(org),
  });

  const refresh = async () => {
    const freshData = await fetchOrganizationRepositories(org, { refresh: true });
    queryClient.setQueryData(organizationReposKeys.byOrg(org), freshData);
  };

  return {
    data: query.data?.data ?? [],
    error: query.error,
    isLoading: query.isPending,
    isRefreshing: query.isFetching && !query.isPending,
    refresh,
  };
};
