"use client";

import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

import { addToHistory } from "@/features/analysis/api/history";
import { invalidationEvents, useInvalidationTrigger } from "@/lib/query";
import { validateRepositoryIdentifiers } from "@/lib/validations/github";

type UseAddToDashboardReturn = {
  addToDashboard: (owner: string, repo: string) => void;
  isPending: boolean;
};

export const useAddToDashboard = (): UseAddToDashboardReturn => {
  const triggerInvalidation = useInvalidationTrigger();

  const mutation = useMutation({
    mutationFn: ({ owner, repo }: { owner: string; repo: string }) => {
      validateRepositoryIdentifiers(owner, repo);
      return addToHistory(owner, repo);
    },
    onError: (error) =>
      toast.error("Failed to add to dashboard", {
        description: error instanceof Error ? error.message : String(error),
      }),
    onSuccess: () => {
      triggerInvalidation(invalidationEvents.HISTORY_CHANGED);
      toast.success("Added to dashboard");
    },
  });

  return {
    addToDashboard: (owner: string, repo: string) => mutation.mutate({ owner, repo }),
    isPending: mutation.isPending,
  };
};
