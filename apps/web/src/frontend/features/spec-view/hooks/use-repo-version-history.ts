"use client";

import { useQuery } from "@tanstack/react-query";

import { fetchRepoVersionHistory } from "../api";
import type { RepoVersionHistoryResponse, SpecLanguage } from "../types";
import { repoSpecViewKeys } from "./use-repo-spec-view";

type UseRepoVersionHistoryOptions = {
  enabled?: boolean;
};

type UseRepoVersionHistoryReturn = {
  data: RepoVersionHistoryResponse | undefined;
  error: Error | null;
  isFetching: boolean;
  isLoading: boolean;
};

export const useRepoVersionHistory = (
  owner: string,
  repo: string,
  language: SpecLanguage | undefined,
  options: UseRepoVersionHistoryOptions = {}
): UseRepoVersionHistoryReturn => {
  const { enabled = true } = options;

  const query = useQuery({
    enabled: enabled && Boolean(owner) && Boolean(repo) && Boolean(language),
    queryFn: () => fetchRepoVersionHistory(owner, repo, language!),
    queryKey: repoSpecViewKeys.versions(owner, repo, language!),
  });

  return {
    data: query.data,
    error: query.error,
    isFetching: query.isFetching,
    isLoading: query.isPending,
  };
};
