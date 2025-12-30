"use client";

import { useTranslations } from "next-intl";

import type { TestStatus } from "@/lib/api";

import { FrameworkFilter } from "./framework-filter";
import { SearchInput } from "./search-input";
import { StatusFilter } from "./status-filter";
import { ViewModeToggle } from "./view-mode-toggle";
import type { ViewMode } from "../types";

type FilterBarProps = {
  availableFrameworks: string[];
  frameworks: string[];
  onFrameworksChange: (value: string[]) => void;
  onQueryChange: (value: string) => void;
  onStatusesChange: (value: TestStatus[]) => void;
  onViewModeChange: (value: ViewMode) => void;
  query: string;
  statuses: TestStatus[];
  viewMode: ViewMode;
};

type FilterSummaryProps = {
  filteredCount: number;
  hasFilter: boolean;
  totalCount: number;
};

export const FilterSummary = ({ filteredCount, hasFilter, totalCount }: FilterSummaryProps) => {
  const t = useTranslations("analyze.filter");

  if (!hasFilter) {
    return null;
  }

  return (
    <span className="text-sm text-muted-foreground">
      {t("resultSummary", { filtered: filteredCount, total: totalCount })}
    </span>
  );
};

export const FilterBar = ({
  availableFrameworks,
  frameworks,
  onFrameworksChange,
  onQueryChange,
  onStatusesChange,
  onViewModeChange,
  query,
  statuses,
  viewMode,
}: FilterBarProps) => {
  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
      <div className="flex-1">
        <SearchInput onChange={onQueryChange} value={query} />
      </div>
      {/* Mobile: full-bleed horizontal scroll area */}
      <div className="-mx-4 px-4 sm:mx-0 sm:px-0 overflow-x-auto scrollbar-hide">
        <div className="flex items-center gap-2 w-max sm:w-auto">
          <StatusFilter onChange={onStatusesChange} value={statuses} />
          <FrameworkFilter
            availableFrameworks={availableFrameworks}
            onChange={onFrameworksChange}
            value={frameworks}
          />
          <ViewModeToggle onChange={onViewModeChange} value={viewMode} />
        </div>
      </div>
    </div>
  );
};
