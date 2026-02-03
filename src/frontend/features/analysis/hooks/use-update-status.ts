"use client";

import { useQuery } from "@tanstack/react-query";

import type { UpdateStatus } from "@/lib/api/types";

import { checkUpdateStatus } from "../api";

const STALE_TIME_MS = 0; // Always fetch fresh data on mount

export const updateStatusKeys = {
  all: ["updateStatus"] as const,
  detail: (owner: string, repo: string) => [...updateStatusKeys.all, owner, repo] as const,
};

type UseUpdateStatusReturn = {
  analyzedCommitSha: string | undefined;
  isLoading: boolean;
  latestCommitSha: string | undefined;
  parserOutdated: boolean;
  status: UpdateStatus | undefined;
};

export const useUpdateStatus = (
  owner: string,
  repo: string,
  enabled: boolean = true
): UseUpdateStatusReturn => {
  const query = useQuery({
    enabled,
    queryFn: () => checkUpdateStatus(owner, repo),
    queryKey: updateStatusKeys.detail(owner, repo),
    staleTime: STALE_TIME_MS,
  });

  return {
    analyzedCommitSha: query.data?.analyzedCommitSha,
    isLoading: query.isPending,
    latestCommitSha: query.data?.latestCommitSha,
    parserOutdated: query.data?.parserOutdated ?? false,
    status: query.data?.status,
  };
};
