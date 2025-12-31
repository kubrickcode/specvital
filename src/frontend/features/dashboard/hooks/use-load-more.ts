"use client";

import { useEffect } from "react";
import { useInView } from "react-intersection-observer";

type UseLoadMoreOptions = {
  fetchNextPage: () => void;
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
  hasNextPage,
  isFetchingNextPage,
  rootMargin = "100px",
  threshold = 0.1,
}: UseLoadMoreOptions): UseLoadMoreReturn => {
  const { inView, ref } = useInView({
    rootMargin,
    threshold,
  });

  useEffect(() => {
    if (inView && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [fetchNextPage, hasNextPage, inView, isFetchingNextPage]);

  return {
    loadMoreRef: ref,
  };
};
