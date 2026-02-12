"use client";

import { debounce } from "es-toolkit";
import { useEffect, useRef, useState } from "react";

import { useRouter } from "@/i18n/navigation";

import { useGlobalSearchStore } from "./use-global-search-store";
import { useSearchData } from "./use-search-data";
import { createRepositoryFuse } from "../lib/fuse-config";
import {
  groupResultsByCategory,
  scoreAndSortResults,
  type ScoredResult,
} from "../lib/score-results";

const SEARCH_DEBOUNCE_MS = 150;
const MAX_RESULTS_PER_CATEGORY = 5;

type GroupedResults = {
  bookmarks: ScoredResult[];
  community: ScoredResult[];
  repositories: ScoredResult[];
};

type UseDebouncedSearchReturn = {
  groupedResults: GroupedResults;
  hasError: boolean;
  hasResults: boolean;
  isAuthenticated: boolean;
  isLoading: boolean;
  navigateToRepository: (owner: string, repo: string) => void;
  query: string;
  setQuery: (query: string) => void;
  totalResults: number;
};

const EMPTY_RESULTS: GroupedResults = {
  bookmarks: [],
  community: [],
  repositories: [],
};

export const useDebouncedSearch = (): UseDebouncedSearchReturn => {
  const router = useRouter();
  const { close } = useGlobalSearchStore();
  const { allItems, hasError, isAuthenticated, isLoading } = useSearchData();

  const [query, setQuery] = useState("");
  const [debouncedQuery, setDebouncedQuery] = useState("");
  const [groupedResults, setGroupedResults] = useState<GroupedResults>(EMPTY_RESULTS);

  const fuseRef = useRef(createRepositoryFuse(allItems));
  const prevItemsRef = useRef(allItems);

  const debouncedSetQueryRef = useRef(
    debounce((newQuery: string) => {
      setDebouncedQuery(newQuery);
    }, SEARCH_DEBOUNCE_MS)
  );

  useEffect(() => {
    const debouncedFn = debouncedSetQueryRef.current;

    if (!query.trim()) {
      setDebouncedQuery("");
      setGroupedResults(EMPTY_RESULTS);
      debouncedFn.cancel();
      return;
    }

    debouncedFn(query);

    return () => debouncedFn.cancel();
  }, [query]);

  useEffect(() => {
    if (allItems !== prevItemsRef.current) {
      fuseRef.current = createRepositoryFuse(allItems);
      prevItemsRef.current = allItems;
    }

    if (!debouncedQuery.trim() || allItems.length === 0) {
      setGroupedResults(EMPTY_RESULTS);
      return;
    }

    const results = fuseRef.current.search(debouncedQuery);
    const scored = scoreAndSortResults(results);
    const grouped = groupResultsByCategory(scored);

    setGroupedResults({
      bookmarks: (grouped.get("bookmarks") ?? []).slice(0, MAX_RESULTS_PER_CATEGORY),
      community: (grouped.get("community") ?? []).slice(0, MAX_RESULTS_PER_CATEGORY),
      repositories: (grouped.get("repositories") ?? []).slice(0, MAX_RESULTS_PER_CATEGORY),
    });
  }, [debouncedQuery, allItems]);

  const navigateToRepository = (owner: string, repo: string) => {
    close();
    router.push(`/analyze/${owner}/${repo}`);
  };

  const totalResults =
    groupedResults.repositories.length +
    groupedResults.bookmarks.length +
    groupedResults.community.length;

  const hasResults = totalResults > 0;

  return {
    groupedResults,
    hasError,
    hasResults,
    isAuthenticated,
    isLoading,
    navigateToRepository,
    query,
    setQuery,
    totalResults,
  };
};
