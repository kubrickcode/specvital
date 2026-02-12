"use client";

import { useQuery } from "@tanstack/react-query";
import { useRef } from "react";

import { useAuth } from "@/features/auth";

import {
  fetchCommunityRepositories,
  fetchUserAnalyzedRepositories,
  fetchUserBookmarks,
} from "../api";
import type { RepositorySearchItem } from "../lib/fuse-config";

type UseSearchDataReturn = {
  allItems: RepositorySearchItem[];
  hasError: boolean;
  isAuthenticated: boolean;
  isLoading: boolean;
};

// Cache key for shallow comparison
type CacheKey = {
  bookmarksData: unknown;
  communityData: unknown;
  userReposData: unknown;
};

export const useSearchData = (): UseSearchDataReturn => {
  const { isAuthenticated } = useAuth();
  const cachedItemsRef = useRef<RepositorySearchItem[]>([]);
  const cacheKeyRef = useRef<CacheKey | null>(null);

  const userReposQuery = useQuery({
    enabled: isAuthenticated,
    queryFn: fetchUserAnalyzedRepositories,
    queryKey: ["global-search", "user-repos"],
  });

  const bookmarksQuery = useQuery({
    enabled: isAuthenticated,
    queryFn: fetchUserBookmarks,
    queryKey: ["global-search", "bookmarks"],
  });

  const communityQuery = useQuery({
    queryFn: fetchCommunityRepositories,
    queryKey: ["global-search", "community"],
  });

  const isLoading =
    (isAuthenticated && (userReposQuery.isLoading || bookmarksQuery.isLoading)) ||
    communityQuery.isLoading;

  const hasError =
    (isAuthenticated && (userReposQuery.isError || bookmarksQuery.isError)) ||
    communityQuery.isError;

  // Check if data sources changed (shallow comparison)
  const currentKey: CacheKey = {
    bookmarksData: bookmarksQuery.data,
    communityData: communityQuery.data,
    userReposData: userReposQuery.data,
  };

  const isCacheValid =
    cacheKeyRef.current !== null &&
    cacheKeyRef.current.userReposData === currentKey.userReposData &&
    cacheKeyRef.current.bookmarksData === currentKey.bookmarksData &&
    cacheKeyRef.current.communityData === currentKey.communityData;

  // Return cached items if data hasn't changed
  if (isCacheValid) {
    return {
      allItems: cachedItemsRef.current,
      hasError,
      isAuthenticated,
      isLoading,
    };
  }

  // Build new items array
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

  // Update cache
  cachedItemsRef.current = allItems;
  cacheKeyRef.current = currentKey;

  return {
    allItems,
    hasError,
    isAuthenticated,
    isLoading,
  };
};
