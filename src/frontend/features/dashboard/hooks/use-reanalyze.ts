"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { validateRepositoryIdentifiers } from "@/lib/validations/github";

import { triggerReanalyze } from "../api";
import { myAnalysesKeys } from "./use-my-analyses";
import { recentRepositoriesKeys } from "./use-recent-repositories";
import { repositoryStatsKeys } from "./use-repository-stats";

type UseReanalyzeReturn = {
  isPending: boolean;
  reanalyze: (owner: string, repo: string) => void;
};

export const useReanalyze = (): UseReanalyzeReturn => {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: ({ owner, repo }: { owner: string; repo: string }) => {
      validateRepositoryIdentifiers(owner, repo);
      return triggerReanalyze(owner, repo);
    },
    onError: (error) =>
      toast.error("Failed to trigger reanalysis", {
        description: error instanceof Error ? error.message : String(error),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: myAnalysesKeys.all });
      queryClient.invalidateQueries({ queryKey: recentRepositoriesKeys.all });
      queryClient.invalidateQueries({ queryKey: repositoryStatsKeys.all });
      toast.success("Reanalysis queued");
    },
  });

  return {
    isPending: mutation.isPending,
    reanalyze: (owner: string, repo: string) => mutation.mutate({ owner, repo }),
  };
};
