"use client";

import { Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";
import { forwardRef } from "react";

import { Button } from "@/components/ui/button";

type LoadMoreButtonState = "default" | "loading" | "error" | "end";

type LoadMoreButtonProps = {
  hasError?: boolean;
  hasNextPage: boolean;
  isFetchingNextPage: boolean;
  onLoadMore: () => void;
  onRetry?: () => void;
};

const getButtonState = ({
  hasError,
  hasNextPage,
  isFetchingNextPage,
}: Pick<
  LoadMoreButtonProps,
  "hasError" | "hasNextPage" | "isFetchingNextPage"
>): LoadMoreButtonState => {
  if (hasError) return "error";
  if (isFetchingNextPage) return "loading";
  if (!hasNextPage) return "end";
  return "default";
};

export const LoadMoreButton = forwardRef<HTMLButtonElement, LoadMoreButtonProps>(
  ({ hasError = false, hasNextPage, isFetchingNextPage, onLoadMore, onRetry }, ref) => {
    const t = useTranslations("dashboard.pagination");

    const state = getButtonState({ hasError, hasNextPage, isFetchingNextPage });

    if (state === "end") {
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

    if (state === "error") {
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
      <div className="flex justify-center py-4">
        <Button
          aria-busy={state === "loading"}
          disabled={state === "loading"}
          onClick={onLoadMore}
          ref={ref}
          variant="outline"
        >
          {state === "loading" ? (
            <>
              <Loader2 aria-hidden="true" className="animate-spin" />
              <span>{t("loading")}</span>
            </>
          ) : (
            <span>{t("loadMore")}</span>
          )}
        </Button>
      </div>
    );
  }
);

LoadMoreButton.displayName = "LoadMoreButton";
