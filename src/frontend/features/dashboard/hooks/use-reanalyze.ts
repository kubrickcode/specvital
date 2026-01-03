"use client";

import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

import { invalidationEvents, useInvalidationTrigger } from "@/lib/query";
import { validateRepositoryIdentifiers } from "@/lib/validations/github";

import { triggerReanalyze } from "../api";

type UseReanalyzeReturn = {
  isPending: boolean;
  reanalyze: (owner: string, repo: string) => void;
};

export const useReanalyze = (): UseReanalyzeReturn => {
  const triggerInvalidation = useInvalidationTrigger();

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
      triggerInvalidation(invalidationEvents.ANALYSIS_COMPLETED);
      toast.success("Reanalysis queued");
    },
  });

  return {
    isPending: mutation.isPending,
    reanalyze: (owner: string, repo: string) => mutation.mutate({ owner, repo }),
  };
};
