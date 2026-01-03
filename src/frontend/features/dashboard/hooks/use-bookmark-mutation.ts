"use client";

import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

import { invalidationEvents, useInvalidationTrigger } from "@/lib/query";
import { validateRepositoryIdentifiers } from "@/lib/validations/github";

import { addBookmark, removeBookmark } from "../api";

type UseAddBookmarkReturn = {
  addBookmark: (owner: string, repo: string) => void;
  isPending: boolean;
};

type UseRemoveBookmarkReturn = {
  isPending: boolean;
  removeBookmark: (owner: string, repo: string) => void;
};

export const useAddBookmark = (): UseAddBookmarkReturn => {
  const triggerInvalidation = useInvalidationTrigger();

  const mutation = useMutation({
    mutationFn: ({ owner, repo }: { owner: string; repo: string }) => {
      validateRepositoryIdentifiers(owner, repo);
      return addBookmark(owner, repo);
    },
    onError: (error) =>
      toast.error("Failed to add bookmark", {
        description: error instanceof Error ? error.message : String(error),
      }),
    onSuccess: () => {
      triggerInvalidation(invalidationEvents.BOOKMARK_CHANGED);
      toast.success("Bookmark added");
    },
  });

  return {
    addBookmark: (owner: string, repo: string) => mutation.mutate({ owner, repo }),
    isPending: mutation.isPending,
  };
};

export const useRemoveBookmark = (): UseRemoveBookmarkReturn => {
  const triggerInvalidation = useInvalidationTrigger();

  const mutation = useMutation({
    mutationFn: ({ owner, repo }: { owner: string; repo: string }) => {
      validateRepositoryIdentifiers(owner, repo);
      return removeBookmark(owner, repo);
    },
    onError: (error) =>
      toast.error("Failed to remove bookmark", {
        description: error instanceof Error ? error.message : String(error),
      }),
    onSuccess: () => {
      triggerInvalidation(invalidationEvents.BOOKMARK_CHANGED);
      toast.success("Bookmark removed");
    },
  });

  return {
    isPending: mutation.isPending,
    removeBookmark: (owner: string, repo: string) => mutation.mutate({ owner, repo }),
  };
};
