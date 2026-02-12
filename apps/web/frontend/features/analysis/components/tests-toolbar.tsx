"use client";

import { useTranslations } from "next-intl";

import type { AnalysisResult, TestStatus } from "@/lib/api";

import { DataViewToggle } from "./data-view-toggle";
import { ExportButton } from "./export-button";
import { FrameworkFilter } from "./framework-filter";
import { SearchInput } from "./search-input";
import { StatusFilter } from "./status-filter";
import type { DataViewMode } from "../types/data-view-mode";

type TestsToolbarProps = {
  availableFrameworks: string[];
  data: AnalysisResult;
  dataViewMode: DataViewMode;
  filteredCount: number;
  frameworks: string[];
  hasFilter: boolean;
  onDataViewModeChange: (value: DataViewMode) => void;
  onFrameworksChange: (value: string[]) => void;
  onQueryChange: (value: string) => void;
  onStatusesChange: (value: TestStatus[]) => void;
  query: string;
  statuses: TestStatus[];
  totalCount: number;
};

export const TestsToolbar = ({
  availableFrameworks,
  data,
  dataViewMode,
  filteredCount,
  frameworks,
  hasFilter,
  onDataViewModeChange,
  onFrameworksChange,
  onQueryChange,
  onStatusesChange,
  query,
  statuses,
  totalCount,
}: TestsToolbarProps) => {
  const t = useTranslations("analyze.filter");

  return (
    <div className="py-4 space-y-3">
      {/* Desktop: Single row */}
      <div className="hidden sm:flex sm:items-center sm:gap-2">
        <div className="flex-1 min-w-0">
          <SearchInput
            onChange={onQueryChange}
            placeholder={t("searchPlaceholder")}
            value={query}
          />
        </div>
        <StatusFilter onChange={onStatusesChange} value={statuses} />
        <FrameworkFilter
          availableFrameworks={availableFrameworks}
          onChange={onFrameworksChange}
          value={frameworks}
        />
        <div className="w-px h-5 bg-border" />
        <DataViewToggle onChange={onDataViewModeChange} value={dataViewMode} />
        <div className="w-px h-5 bg-border" />
        <ExportButton data={data} />
      </div>

      {/* Mobile: Two rows */}
      <div className="flex flex-col gap-3 sm:hidden">
        {/* Row 1: Search + Export icon */}
        <div className="flex items-center gap-2">
          <div className="flex-1 min-w-0">
            <SearchInput
              onChange={onQueryChange}
              placeholder={t("searchPlaceholder")}
              value={query}
            />
          </div>
          <ExportButton data={data} size="icon" />
        </div>
        {/* Row 2: Filters + View Toggle */}
        <div className="flex items-center gap-2">
          <StatusFilter onChange={onStatusesChange} value={statuses} />
          <FrameworkFilter
            availableFrameworks={availableFrameworks}
            onChange={onFrameworksChange}
            value={frameworks}
          />
          <div className="w-px h-5 bg-border" />
          <DataViewToggle onChange={onDataViewModeChange} value={dataViewMode} />
        </div>
      </div>

      {hasFilter && (
        <p className="text-xs text-muted-foreground tabular-nums">
          {filteredCount.toLocaleString()} / {totalCount.toLocaleString()}
        </p>
      )}
    </div>
  );
};
