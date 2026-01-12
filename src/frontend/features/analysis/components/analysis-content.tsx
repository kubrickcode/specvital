"use client";

import { motion } from "motion/react";
import { useTranslations } from "next-intl";
import { useMemo } from "react";

import type { AnalysisResult } from "@/lib/api";
import { createStaggerContainer, fadeInUp, useReducedMotion } from "@/lib/motion";

import { AnalysisHeader } from "./analysis-header";
import { FilterBar, FilterSummary } from "./filter-bar";
import { FilterEmptyState } from "./filter-empty-state";
import { StatsCard } from "./stats-card";
import { TestList } from "./test-list";
import { TreeView } from "./tree-view";
import { useFilterState } from "../hooks/use-filter-state";
import { useViewMode } from "../hooks/use-view-mode";
import type { ViewMode } from "../types";
import { filterSuites } from "../utils/filter-suites";

type AnalysisContentProps = {
  result: AnalysisResult;
};

const pageStaggerContainer = createStaggerContainer(0.1, 0);

export const AnalysisContent = ({ result }: AnalysisContentProps) => {
  const t = useTranslations("analyze");
  const { frameworks, query, setFrameworks, setQuery, setStatuses, statuses } = useFilterState();
  const { setViewMode, viewMode } = useViewMode();
  const shouldReduceMotion = useReducedMotion();

  const containerVariants = shouldReduceMotion ? {} : pageStaggerContainer;
  const itemVariants = shouldReduceMotion ? {} : fadeInUp;

  const availableFrameworks = useMemo(
    () => result.summary.frameworks.map((f) => f.framework),
    [result.summary.frameworks]
  );

  const filteredSuites = useMemo(
    () => filterSuites(result.suites, { frameworks, query, statuses }),
    [result.suites, frameworks, query, statuses]
  );

  const filteredTestCount = useMemo(
    () => filteredSuites.reduce((acc, suite) => acc + suite.tests.length, 0),
    [filteredSuites]
  );

  const hasFilter = query.trim().length > 0 || frameworks.length > 0 || statuses.length > 0;
  const hasResults = filteredSuites.length > 0;

  const handleViewModeChange = (mode: ViewMode) => {
    setViewMode(mode);
  };

  const renderContent = () => {
    if (hasFilter && !hasResults) {
      return <FilterEmptyState />;
    }

    if (viewMode === "tree") {
      return <TreeView suites={filteredSuites} />;
    }

    return <TestList suites={filteredSuites} />;
  };

  return (
    <motion.main
      animate="visible"
      className="container mx-auto px-4 py-8"
      initial={shouldReduceMotion ? false : "hidden"}
      variants={containerVariants}
    >
      <div className="space-y-6">
        <AnalysisHeader
          analyzedAt={result.analyzedAt}
          branchName={result.branchName}
          commitSha={result.commitSha}
          committedAt={result.committedAt}
          data={result}
          owner={result.owner}
          parserVersion={result.parserVersion}
          repo={result.repo}
        />

        <motion.div variants={itemVariants}>
          <StatsCard summary={result.summary} />
        </motion.div>

        <motion.section className="space-y-4" variants={itemVariants}>
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <h2 className="text-xl font-semibold">{t("testSuites")}</h2>
              <FilterSummary
                filteredCount={filteredTestCount}
                hasFilter={hasFilter}
                totalCount={result.summary.total}
              />
            </div>
          </div>
          <FilterBar
            availableFrameworks={availableFrameworks}
            frameworks={frameworks}
            onFrameworksChange={setFrameworks}
            onQueryChange={setQuery}
            onStatusesChange={setStatuses}
            onViewModeChange={handleViewModeChange}
            query={query}
            statuses={statuses}
            viewMode={viewMode}
          />
          {renderContent()}
        </motion.section>
      </div>
    </motion.main>
  );
};
