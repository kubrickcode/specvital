"use client";

import { useTranslations } from "next-intl";

type PaginationStatusProps = {
  hasNextPage: boolean;
  isLoading?: boolean;
  totalLoaded: number;
};

export const PaginationStatus = ({
  hasNextPage,
  isLoading = false,
  totalLoaded,
}: PaginationStatusProps) => {
  const t = useTranslations("dashboard.pagination");

  if (isLoading && totalLoaded === 0) {
    return null;
  }

  return (
    <p
      aria-atomic="true"
      aria-live="polite"
      className="text-sm text-muted-foreground"
      role="status"
    >
      {hasNextPage
        ? t("showingWithMore", { count: totalLoaded })
        : t("showingAll", { count: totalLoaded })}
    </p>
  );
};
