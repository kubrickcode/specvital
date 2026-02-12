"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { toast } from "sonner";

import type { RepositoryCard } from "@/lib/api/types";
import { validateRepositoryIdentifiers } from "@/lib/validations/github";

import { addBookmark, removeBookmark } from "../api";
import { paginatedRepositoriesKeys } from "./use-paginated-repositories";

type MutationVariables = { owner: string; repo: string };
type MutationContext = { previousData: InfiniteData | undefined };
type PaginatedPage = { data: RepositoryCard[]; hasNext: boolean; nextCursor?: string | null };
type InfiniteData = { pageParams: unknown[]; pages: PaginatedPage[] };

type UseAddBookmarkReturn = {
  addBookmark: (owner: string, repo: string) => void;
  isPending: boolean;
};

type UseRemoveBookmarkReturn = {
  isPending: boolean;
  removeBookmark: (owner: string, repo: string) => void;
};

const toggleBookmarkInPages = (
  pages: PaginatedPage[],
  owner: string,
  repo: string,
  isBookmarked: boolean
): PaginatedPage[] =>
  pages.map((page) => ({
    ...page,
    data: page.data.map((r) => (r.owner === owner && r.name === repo ? { ...r, isBookmarked } : r)),
  }));

export const useAddBookmark = (): UseAddBookmarkReturn => {
  const t = useTranslations("dashboard.toast");
  const queryClient = useQueryClient();

  const mutation = useMutation<unknown, Error, MutationVariables, MutationContext>({
    mutationFn: ({ owner, repo }) => {
      validateRepositoryIdentifiers(owner, repo);
      return addBookmark(owner, repo);
    },
    onError: (error, _variables, context) => {
      if (context?.previousData) {
        queryClient.setQueriesData<InfiniteData>(
          { queryKey: paginatedRepositoriesKeys.all },
          () => context.previousData
        );
      }
      toast.error(t("bookmarkAddFailed"), {
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
            pages: toggleBookmarkInPages(old.pages, owner, repo, true),
          };
        }
      );

      return { previousData };
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: paginatedRepositoriesKeys.all });
      toast.success(t("bookmarkAdded"));
    },
  });

  return {
    addBookmark: (owner: string, repo: string) => mutation.mutate({ owner, repo }),
    isPending: mutation.isPending,
  };
};

export const useRemoveBookmark = (): UseRemoveBookmarkReturn => {
  const t = useTranslations("dashboard.toast");
  const queryClient = useQueryClient();

  const mutation = useMutation<unknown, Error, MutationVariables, MutationContext>({
    mutationFn: ({ owner, repo }) => {
      validateRepositoryIdentifiers(owner, repo);
      return removeBookmark(owner, repo);
    },
    onError: (error, _variables, context) => {
      if (context?.previousData) {
        queryClient.setQueriesData<InfiniteData>(
          { queryKey: paginatedRepositoriesKeys.all },
          () => context.previousData
        );
      }
      toast.error(t("bookmarkRemoveFailed"), {
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
            pages: toggleBookmarkInPages(old.pages, owner, repo, false),
          };
        }
      );

      return { previousData };
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: paginatedRepositoriesKeys.all });
      toast.success(t("bookmarkRemoved"));
    },
  });

  return {
    isPending: mutation.isPending,
    removeBookmark: (owner: string, repo: string) => mutation.mutate({ owner, repo }),
  };
};
