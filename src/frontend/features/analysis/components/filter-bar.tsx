"use client";

import { useTranslations } from "next-intl";

import type { TestStatus } from "@/lib/api";

import { AIAnalysisButton } from "./ai-analysis";
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

  return (
    <span
      aria-atomic="true"
      aria-live="polite"
      className="text-sm text-muted-foreground"
      role="status"
    >
      {hasFilter ? t("resultSummary", { filtered: filteredCount, total: totalCount }) : null}
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
    <>
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
        <div className="flex-1">
          <SearchInput onChange={onQueryChange} value={query} />
        </div>
        {/* Mobile: full-bleed horizontal scroll area with fixed ViewModeToggle */}
        <div className="flex items-center gap-2">
          <div className="-mx-4 px-4 sm:mx-0 sm:px-0 overflow-x-auto scrollbar-hide flex-1 min-w-0">
            <div className="flex items-center gap-2 w-max">
              <StatusFilter onChange={onStatusesChange} value={statuses} />
              <FrameworkFilter
                availableFrameworks={availableFrameworks}
                onChange={onFrameworksChange}
                value={frameworks}
              />
            </div>
          </div>
          <ViewModeToggle onChange={onViewModeChange} value={viewMode} />
          <div className="hidden sm:block">
            <AIAnalysisButton />
          </div>
        </div>
      </div>
      {/* AI Analysis button: full width on mobile only */}
      <div className="block sm:hidden">
        <AIAnalysisButton className="w-full" />
      </div>
    </>
  );
};
