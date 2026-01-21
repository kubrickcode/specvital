"use client";

import { FilterEmptyState } from "@/components/feedback";
import type { TestSuite } from "@/lib/api";

import { TestList } from "./test-list";
import { TestsToolbar } from "./tests-toolbar";
import { TreeView } from "./tree-view";
import { useDataViewMode } from "../hooks/use-data-view-mode";
import { useFilterState } from "../hooks/use-filter-state";
import { filterSuites } from "../utils/filter-suites";

type TestsPanelProps = {
  availableFrameworks: string[];
  suites: TestSuite[];
  totalCount: number;
};

export const TestsPanel = ({ availableFrameworks, suites, totalCount }: TestsPanelProps) => {
  const { frameworks, query, setFrameworks, setQuery, setStatuses, statuses } = useFilterState();
  const { dataViewMode, setDataViewMode } = useDataViewMode();

  const filteredSuites = filterSuites(suites, { frameworks, query, statuses });
  const filteredCount = filteredSuites.reduce((acc, suite) => acc + suite.tests.length, 0);

  const hasFilter = query.trim().length > 0 || frameworks.length > 0 || statuses.length > 0;
  const hasResults = filteredSuites.length > 0;

  const resetFilters = () => {
    setQuery(null);
    setFrameworks(null);
    setStatuses(null);
  };

  return (
    <div aria-labelledby="tab-tests" className="p-5 space-y-4" id="tabpanel-tests" role="tabpanel">
      <TestsToolbar
        availableFrameworks={availableFrameworks}
        dataViewMode={dataViewMode}
        filteredCount={filteredCount}
        frameworks={frameworks}
        hasFilter={hasFilter}
        onDataViewModeChange={setDataViewMode}
        onFrameworksChange={setFrameworks}
        onQueryChange={setQuery}
        onStatusesChange={setStatuses}
        query={query}
        statuses={statuses}
        totalCount={totalCount}
      />

      {hasFilter && !hasResults ? (
        <FilterEmptyState filterInfo={{ frameworks, query, statuses }} onReset={resetFilters} />
      ) : dataViewMode === "tree" ? (
        <TreeView suites={filteredSuites} />
      ) : (
        <TestList suites={filteredSuites} />
      )}
    </div>
  );
};
