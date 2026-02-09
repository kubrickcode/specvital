"use client";

import { useEffect } from "react";
import { useInView } from "react-intersection-observer";

type UseLoadMoreOptions = {
  fetchNextPage: () => void;
  hasError?: boolean;
  hasNextPage: boolean;
  isFetchingNextPage: boolean;
  rootMargin?: string;
  threshold?: number;
};

type UseLoadMoreReturn = {
  loadMoreRef: (node?: Element | null) => void;
};

export const useLoadMore = ({
  fetchNextPage,
  hasError = false,
  hasNextPage,
  isFetchingNextPage,
  rootMargin = "200px",
  threshold = 0.1,
}: UseLoadMoreOptions): UseLoadMoreReturn => {
  const { inView, ref } = useInView({
    rootMargin,
    threshold,
  });

  useEffect(() => {
    if (inView && hasNextPage && !isFetchingNextPage && !hasError) {
      fetchNextPage();
    }
  }, [fetchNextPage, hasError, hasNextPage, inView, isFetchingNextPage]);

  return {
    loadMoreRef: ref,
  };
};
