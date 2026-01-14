"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRef } from "react";
import { toast } from "sonner";

import type { RepositoryCard } from "@/lib/api/types";
import { invalidationEvents, useInvalidationTrigger } from "@/lib/query";
import { validateRepositoryIdentifiers } from "@/lib/validations/github";

import { checkUpdateStatus, triggerReanalyze } from "../api";
import { paginatedRepositoriesKeys } from "./use-paginated-repositories";
import type { DisplayUpdateStatus } from "../components/update-status-badge";

const POLL_INTERVAL_MS = 3000;
const POLL_TIMEOUT_MS = 60000;

type MutationVariables = { owner: string; repo: string };
type MutationContext = { previousData: InfiniteData | undefined };
type PaginatedPage = { data: RepositoryCard[]; hasNext: boolean; nextCursor?: string | null };
type InfiniteData = { pageParams: unknown[]; pages: PaginatedPage[] };

type UseReanalyzeReturn = {
  isPending: boolean;
  reanalyze: (owner: string, repo: string) => void;
};

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
  const queryClient = useQueryClient();
  const triggerInvalidation = useInvalidationTrigger();
  const pollIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const pollTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const stopPolling = () => {
    if (pollIntervalRef.current) {
      clearInterval(pollIntervalRef.current);
      pollIntervalRef.current = null;
    }
    if (pollTimeoutRef.current) {
      clearTimeout(pollTimeoutRef.current);
      pollTimeoutRef.current = null;
    }
  };

  const startPolling = (owner: string, repo: string) => {
    stopPolling();

    const poll = async () => {
      try {
        const result = await checkUpdateStatus(owner, repo);
        if (result.status === "up-to-date") {
          stopPolling();
          triggerInvalidation(invalidationEvents.ANALYSIS_COMPLETED);
        }
      } catch {
        // Ignore polling errors and continue retrying
      }
    };

    pollIntervalRef.current = setInterval(poll, POLL_INTERVAL_MS);
    pollTimeoutRef.current = setTimeout(stopPolling, POLL_TIMEOUT_MS);
  };

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
      toast.error("Failed to trigger reanalysis", {
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
      toast.success("Reanalysis queued");
      startPolling(owner, repo);
    },
  });

  return {
    isPending: mutation.isPending,
    reanalyze: (owner: string, repo: string) => mutation.mutate({ owner, repo }),
  };
};
