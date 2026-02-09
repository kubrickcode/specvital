"use client";

import { useState } from "react";

import type { RepositoryCard } from "@/lib/api/types";

import type { SortOption } from "../types";

type UseRepositorySearchReturn = {
  filteredRepositories: RepositoryCard[];
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  setSortBy: (sort: SortOption) => void;
  sortBy: SortOption;
};

export const useRepositorySearch = (repositories: RepositoryCard[]): UseRepositorySearchReturn => {
  const [searchQuery, setSearchQuery] = useState("");
  const [sortBy, setSortBy] = useState<SortOption>("recent");

  const filteredRepositories = (() => {
    let result = [...repositories];

    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      result = result.filter(
        (repo) =>
          repo.owner.toLowerCase().includes(query) || repo.name.toLowerCase().includes(query)
      );
    }

    result.sort((a, b) => {
      // Bookmarked repos always come first
      if (a.isBookmarked !== b.isBookmarked) {
        return a.isBookmarked ? -1 : 1;
      }

      switch (sortBy) {
        case "name":
          return `${a.owner}/${a.name}`.localeCompare(`${b.owner}/${b.name}`);
        case "tests":
          return (b.latestAnalysis?.testCount ?? 0) - (a.latestAnalysis?.testCount ?? 0);
        case "recent":
        default: {
          const aTime = a.latestAnalysis?.analyzedAt
            ? new Date(a.latestAnalysis.analyzedAt).getTime()
            : 0;
          const bTime = b.latestAnalysis?.analyzedAt
            ? new Date(b.latestAnalysis.analyzedAt).getTime()
            : 0;
          return bTime - aTime;
        }
      }
    });

    return result;
  })();

  return {
    filteredRepositories,
    searchQuery,
    setSearchQuery,
    setSortBy,
    sortBy,
  };
};
