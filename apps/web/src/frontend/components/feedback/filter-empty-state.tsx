"use client";

import { RotateCcw, SearchX } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { TestStatus } from "@/lib/api";

export type FilterInfo = {
  frameworks: string[];
  query: string;
  statuses: TestStatus[];
};

type FilterEmptyStateProps = {
  filterInfo: FilterInfo;
  onReset: () => void;
};

export const FilterEmptyState = ({ filterInfo, onReset }: FilterEmptyStateProps) => {
  const t = useTranslations("analyze.filter");

  const hasActiveFilters =
    filterInfo.query.trim().length > 0 ||
    filterInfo.frameworks.length > 0 ||
    filterInfo.statuses.length > 0;

  return (
    <div className="rounded-lg border bg-card p-8 shadow-xs">
      <div className="flex flex-col items-center justify-center text-center">
        <div className="rounded-full bg-muted p-4 mb-4">
          <SearchX aria-hidden="true" className="h-8 w-8 text-muted-foreground" />
        </div>
        <h3 className="text-lg font-semibold">{t("noResults")}</h3>
        <p className="text-sm text-muted-foreground mt-1">{t("noResultsDescription")}</p>

        {hasActiveFilters && (
          <div className="mt-6 space-y-4">
            <div className="flex flex-wrap justify-center gap-2">
              {filterInfo.query.trim() && (
                <Badge className="text-xs" variant="secondary">
                  {t("searchLabel")}: &quot;{filterInfo.query}&quot;
                </Badge>
              )}
              {filterInfo.frameworks.map((framework) => (
                <Badge className="text-xs" key={framework} variant="secondary">
                  {framework}
                </Badge>
              ))}
              {filterInfo.statuses.map((status) => (
                <Badge className="text-xs" key={status} variant="secondary">
                  {t(
                    `status${status.charAt(0).toUpperCase()}${status.slice(1)}` as Parameters<
                      typeof t
                    >[0]
                  )}
                </Badge>
              ))}
            </div>

            <Button onClick={onReset} size="sm" variant="outline">
              <RotateCcw className="mr-2 h-4 w-4" />
              {t("resetFilters")}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};
