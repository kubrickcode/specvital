"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import type { AnalysisResponse, RepositoryCard } from "@/lib/api/types";
import { userActiveTasksKeys } from "@/lib/background-tasks";
import { validateRepositoryIdentifiers } from "@/lib/validations/github";

import { fetchAnalysisStatus } from "../../analysis/api";
import { triggerReanalyze } from "../api";
import { paginatedRepositoriesKeys } from "./use-paginated-repositories";
import { repositoryStatsKeys } from "./use-repository-stats";
import type { DisplayUpdateStatus } from "../components/update-status-badge";

const POLL_INTERVAL_MS = 1000;

type MutationVariables = { owner: string; repo: string };
type MutationContext = { previousData: InfiniteData | undefined };
type PaginatedPage = { data: RepositoryCard[]; hasNext: boolean; nextCursor?: string | null };
type InfiniteData = { pageParams: unknown[]; pages: PaginatedPage[] };

type UseReanalyzeReturn = {
  isPending: boolean;
  reanalyze: (owner: string, repo: string) => void;
};

const isTerminalStatus = (status: AnalysisResponse["status"]): boolean =>
  status === "completed" || status === "failed";

const updateRepoStatus = (
  pages: PaginatedPage[],
  owner: string,
  repo: string,
  status: DisplayUpdateStatus
): PaginatedPage[] =>
  pages.map((page) => ({
    ...page,
    data: page.data.map((r) =>
      r.owner === owner && r.name === repo
        ? { ...r, updateStatus: status as RepositoryCard["updateStatus"] }
        : r
    ),
  }));

export const useReanalyze = (): UseReanalyzeReturn => {
  const t = useTranslations("analysis.updateBanner");
  const queryClient = useQueryClient();
  const [pollingTarget, setPollingTarget] = useState<{ owner: string; repo: string } | null>(null);

  const pollingQuery = useQuery({
    enabled: pollingTarget !== null,
    queryFn: async () => {
      if (!pollingTarget) throw new Error("No polling target");
      return fetchAnalysisStatus(pollingTarget.owner, pollingTarget.repo);
    },
    queryKey: ["reanalyzePolling", pollingTarget?.owner, pollingTarget?.repo],
    refetchInterval: (query) => {
      const data = query.state.data;
      if (!data) return POLL_INTERVAL_MS;
      return isTerminalStatus(data.status) ? false : POLL_INTERVAL_MS;
    },
    // Prevent returning cached "completed" status from previous polling sessions
    gcTime: 0,
  });

  // Clean up polling state on completion
  useEffect(() => {
    if (!pollingTarget || !pollingQuery.data) return;
    if (!isTerminalStatus(pollingQuery.data.status)) return;

    queryClient.invalidateQueries({ queryKey: paginatedRepositoriesKeys.all });
    queryClient.invalidateQueries({ queryKey: repositoryStatsKeys.all });
    // Refresh active tasks list from server
    queryClient.invalidateQueries({ queryKey: userActiveTasksKeys.all });
    setPollingTarget(null);
  }, [pollingTarget, pollingQuery.data, queryClient]);

  const mutation = useMutation({
    mutationFn: ({ owner, repo }: MutationVariables) => {
      validateRepositoryIdentifiers(owner, repo);
      return triggerReanalyze(owner, repo);
    },
    onError: (error, _variables, context: MutationContext | undefined) => {
      if (context?.previousData) {
        queryClient.setQueriesData<InfiniteData>(
          { queryKey: paginatedRepositoriesKeys.all },
          () => context.previousData
        );
      }
      toast.error(t("reanalyzeFailed"), {
        description: error instanceof Error ? error.message : String(error),
      });
    },
    onMutate: async ({ owner, repo }) => {
      await queryClient.cancelQueries({ queryKey: paginatedRepositoriesKeys.all });

      const previousData = queryClient.getQueryData<InfiniteData>(paginatedRepositoriesKeys.all);

      queryClient.setQueriesData<InfiniteData>(
        { queryKey: paginatedRepositoriesKeys.all },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            pages: updateRepoStatus(old.pages, owner, repo, "analyzing"),
          };
        }
      );

      return { previousData };
    },
    onSuccess: (_data, { owner, repo }) => {
      toast.success(t("reanalyzeQueued"));

      // Clear cached polling data before starting new polling session
      // Prevents refetchInterval from seeing stale "completed" status
      queryClient.removeQueries({ queryKey: ["reanalyzePolling", owner, repo] });

      // Refresh active tasks list from server (new task will appear)
      queryClient.invalidateQueries({ queryKey: userActiveTasksKeys.all });

      setPollingTarget({ owner, repo });
    },
  });

  return {
    isPending: mutation.isPending,
    reanalyze: (owner: string, repo: string) => mutation.mutate({ owner, repo }),
  };
};
