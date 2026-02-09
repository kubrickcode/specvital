"use client";

import { Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

import { useLoadMore } from "../hooks";

type InfiniteScrollLoaderProps = {
  hasError?: boolean;
  hasNextPage: boolean;
  isFetchingNextPage: boolean;
  onLoadMore: () => void;
  onRetry?: () => void;
};

export const InfiniteScrollLoader = ({
  hasError = false,
  hasNextPage,
  isFetchingNextPage,
  onLoadMore,
  onRetry,
}: InfiniteScrollLoaderProps) => {
  const t = useTranslations("dashboard.pagination");

  const { loadMoreRef } = useLoadMore({
    fetchNextPage: onLoadMore,
    hasError,
    hasNextPage,
    isFetchingNextPage,
  });

  if (!hasNextPage && !hasError) {
    return (
      <p
        aria-live="polite"
        className="py-4 text-center text-sm text-muted-foreground"
        role="status"
      >
        {t("allLoaded")}
      </p>
    );
  }

  if (hasError) {
    return (
      <div className="flex flex-col items-center gap-2 py-4">
        <p className="text-sm text-destructive" role="alert">
          {t("loadError")}
        </p>
        <Button onClick={onRetry ?? onLoadMore} size="sm" variant="outline">
          {t("retry")}
        </Button>
      </div>
    );
  }

  return (
    <div className="relative flex flex-col items-center gap-4 py-4" ref={loadMoreRef}>
      {isFetchingNextPage ? (
        <div
          aria-busy="true"
          aria-live="polite"
          className="flex items-center gap-2 text-sm text-muted-foreground"
          role="status"
        >
          <Loader2 aria-hidden="true" className="size-4 animate-spin" />
          <span>{t("loading")}</span>
        </div>
      ) : (
        <Button
          aria-label={t("loadMore")}
          className="sr-only focus:not-sr-only focus:absolute focus:top-1/2 focus:left-1/2 focus:-translate-x-1/2 focus:-translate-y-1/2"
          onClick={onLoadMore}
          variant="outline"
        >
          {t("loadMore")}
        </Button>
      )}
    </div>
  );
};
