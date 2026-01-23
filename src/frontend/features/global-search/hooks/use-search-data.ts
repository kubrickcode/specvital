"use client";

import { useQuery } from "@tanstack/react-query";

import { useAuth } from "@/features/auth";

import {
  fetchCommunityRepositories,
  fetchUserAnalyzedRepositories,
  fetchUserBookmarks,
} from "../api";
import type { RepositorySearchItem } from "../lib/fuse-config";

const SEARCH_DATA_STALE_TIME = 30 * 1000;
const SEARCH_DATA_GC_TIME = 5 * 60 * 1000;

type UseSearchDataReturn = {
  allItems: RepositorySearchItem[];
  hasError: boolean;
  isAuthenticated: boolean;
  isLoading: boolean;
};

export const useSearchData = (): UseSearchDataReturn => {
  const { isAuthenticated } = useAuth();

  const userReposQuery = useQuery({
    enabled: isAuthenticated,
    gcTime: SEARCH_DATA_GC_TIME,
    queryFn: fetchUserAnalyzedRepositories,
    queryKey: ["global-search", "user-repos"],
    staleTime: SEARCH_DATA_STALE_TIME,
  });

  const bookmarksQuery = useQuery({
    enabled: isAuthenticated,
    gcTime: SEARCH_DATA_GC_TIME,
    queryFn: fetchUserBookmarks,
    queryKey: ["global-search", "bookmarks"],
    staleTime: SEARCH_DATA_STALE_TIME,
  });

  const communityQuery = useQuery({
    gcTime: SEARCH_DATA_GC_TIME,
    queryFn: fetchCommunityRepositories,
    queryKey: ["global-search", "community"],
    staleTime: SEARCH_DATA_STALE_TIME,
  });

  const isLoading =
    (isAuthenticated && (userReposQuery.isLoading || bookmarksQuery.isLoading)) ||
    communityQuery.isLoading;

  const hasError =
    (isAuthenticated && (userReposQuery.isError || bookmarksQuery.isError)) ||
    communityQuery.isError;

  const allItems: RepositorySearchItem[] = [];
  const seenIds = new Set<string>();

  if (userReposQuery.data) {
    for (const repo of userReposQuery.data.data) {
      if (!seenIds.has(repo.id)) {
        seenIds.add(repo.id);
        allItems.push({
          category: "repositories",
          fullName: repo.fullName,
          id: repo.id,
          isAnalyzedByMe: repo.isAnalyzedByMe,
          isBookmarked: repo.isBookmarked,
          name: repo.name,
          owner: repo.owner,
        });
      }
    }
  }

  if (bookmarksQuery.data) {
    for (const repo of bookmarksQuery.data.data) {
      if (!seenIds.has(repo.id)) {
        seenIds.add(repo.id);
        allItems.push({
          category: "bookmarks",
          fullName: repo.fullName,
          id: repo.id,
          isAnalyzedByMe: repo.isAnalyzedByMe,
          isBookmarked: repo.isBookmarked,
          name: repo.name,
          owner: repo.owner,
        });
      }
    }
  }

  if (communityQuery.data) {
    for (const repo of communityQuery.data.data) {
      if (!seenIds.has(repo.id)) {
        seenIds.add(repo.id);
        allItems.push({
          category: "community",
          fullName: repo.fullName,
          id: repo.id,
          isAnalyzedByMe: repo.isAnalyzedByMe,
          isBookmarked: repo.isBookmarked,
          name: repo.name,
          owner: repo.owner,
        });
      }
    }
  }

  return {
    allItems,
    hasError,
    isAuthenticated,
    isLoading,
  };
};
