"use client";

import { useEffect, useMemo } from "react";

import { usePaginatedRepositories } from "@/features/dashboard";

export const useDashboardRepoIds = (): Set<string> => {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } = usePaginatedRepositories({
    limit: 50,
    view: "my",
  });

  useEffect(() => {
    if (hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  return useMemo(() => new Set(data.map((repo) => repo.id)), [data]);
};
