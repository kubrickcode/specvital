"use client";

import { useTranslations } from "next-intl";

import type { TestStatus } from "@/lib/api";

import { FrameworkFilter } from "./framework-filter";
import { SearchInput } from "./search-input";
import { StatusFilter } from "./status-filter";

type SpecToolbarProps = {
  availableFrameworks: string[];
  frameworks: string[];
  hasDocument: boolean;
  hasFilter: boolean;
  matchCount: number;
  onFrameworksChange: (value: string[]) => void;
  onQueryChange: (value: string) => void;
  onStatusesChange: (value: TestStatus[]) => void;
  query: string;
  statuses: TestStatus[];
  totalCount: number;
};

export const SpecToolbar = ({
  availableFrameworks,
  frameworks,
  hasDocument,
  hasFilter,
  matchCount,
  onFrameworksChange,
  onQueryChange,
  onStatusesChange,
  query,
  statuses,
  totalCount,
}: SpecToolbarProps) => {
  const t = useTranslations("analyze.filter");

  return (
    <div className="rounded-lg border bg-muted/30 p-4 space-y-3">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
        <div className="flex-1 min-w-0">
          <SearchInput
            onChange={onQueryChange}
            placeholder={t("searchBehaviorsPlaceholder")}
            value={query}
          />
        </div>
        <div className="flex items-center gap-2">
          <div className="-mx-4 px-4 sm:mx-0 sm:px-0 overflow-x-auto scrollbar-hide">
            <div className="flex items-center gap-2 w-max">
              <StatusFilter onChange={onStatusesChange} value={statuses} />
              <FrameworkFilter
                availableFrameworks={availableFrameworks}
                onChange={onFrameworksChange}
                value={frameworks}
              />
            </div>
          </div>
        </div>
      </div>

      {hasFilter && hasDocument && (
        <p className="text-sm text-muted-foreground">
          {t("resultSummaryWithType", {
            filtered: matchCount,
            total: totalCount,
            type: "behaviors",
          })}
        </p>
      )}
    </div>
  );
};
