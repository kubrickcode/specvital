"use client";

import { useQueries, useQueryClient } from "@tanstack/react-query";
import { useMemo } from "react";
import { toast } from "sonner";

import type { GitHubOrganization, GitHubRepository, RepositoryCard } from "@/lib/api/types";

import { fetchOrganizationRepositories } from "../api";
import { organizationReposKeys } from "./use-organization-repos";

type UseAllOrgReposParams = {
  analyzedRepositories: RepositoryCard[];
  organizations: GitHubOrganization[];
};

type OrgRepoData = {
  isLoading: boolean;
  repositories: GitHubRepository[];
  unanalyzedCount: number;
  unanalyzedRepos: GitHubRepository[];
};

type UseAllOrgReposReturn = {
  isLoading: boolean;
  orgData: Record<string, OrgRepoData>;
  refreshOrg: (org: string) => Promise<void>;
  totalUnanalyzedCount: number;
};

export const useAllOrgRepos = ({
  analyzedRepositories,
  organizations,
}: UseAllOrgReposParams): UseAllOrgReposReturn => {
  const queryClient = useQueryClient();

  const analyzedFullNames = useMemo(
    () => new Set(analyzedRepositories.map((r) => r.fullName.toLowerCase())),
    [analyzedRepositories]
  );

  const queries = useQueries({
    queries: organizations.map((org) => ({
      enabled: Boolean(org.login),
      queryFn: () => fetchOrganizationRepositories(org.login),
      queryKey: organizationReposKeys.byOrg(org.login),
      staleTime: 5 * 60 * 1000,
    })),
  });

  const isLoading = queries.some((q) => q.isPending);

  const orgData = useMemo(() => {
    const data: Record<string, OrgRepoData> = {};

    organizations.forEach((org, index) => {
      const query = queries[index];
      const repositories = query?.data?.data ?? [];
      const unanalyzedRepos = repositories.filter(
        (repo) => !analyzedFullNames.has(repo.fullName.toLowerCase())
      );

      data[org.login] = {
        isLoading: query?.isPending ?? true,
        repositories,
        unanalyzedCount: unanalyzedRepos.length,
        unanalyzedRepos,
      };
    });

    return data;
  }, [analyzedFullNames, organizations, queries]);

  const totalUnanalyzedCount = useMemo(
    () => Object.values(orgData).reduce((sum, data) => sum + data.unanalyzedCount, 0),
    [orgData]
  );

  const refreshOrg = async (org: string) => {
    try {
      const freshData = await fetchOrganizationRepositories(org, { refresh: true });
      queryClient.setQueryData(organizationReposKeys.byOrg(org), freshData);
      toast.success(`Refreshed ${org} from GitHub`);
    } catch (error) {
      toast.error(`Failed to refresh ${org}`, {
        description: error instanceof Error ? error.message : String(error),
      });
    }
  };

  return {
    isLoading,
    orgData,
    refreshOrg,
    totalUnanalyzedCount,
  };
};
