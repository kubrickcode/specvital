"use client";

import { ExternalLink, GitCommit } from "lucide-react";
import { motion } from "motion/react";
import { useTranslations } from "next-intl";
import { useMemo } from "react";

import { Button } from "@/components/ui/button";
import type { AnalysisResult } from "@/lib/api";
import { createStaggerContainer, fadeInUp, useReducedMotion } from "@/lib/motion";
import { formatAnalysisDate, SHORT_SHA_LENGTH } from "@/lib/utils";

import { FilterBar, FilterSummary } from "./filter-bar";
import { FilterEmptyState } from "./filter-empty-state";
import { ShareButton } from "./share-button";
import { StatsCard } from "./stats-card";
import { TestList } from "./test-list";
import { TreeView } from "./tree-view";
import { useFilterState } from "../hooks/use-filter-state";
import { useViewMode } from "../hooks/use-view-mode";
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

  return (
    <motion.main
      animate="visible"
      className="container mx-auto px-4 py-8"
      initial={shouldReduceMotion ? false : "hidden"}
      variants={containerVariants}
    >
      <div className="space-y-6">
        <motion.header className="space-y-2" variants={itemVariants}>
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <h1 className="text-xl font-bold sm:text-2xl truncate min-w-0">
              {result.owner}/{result.repo}
            </h1>
            <div className="flex items-center gap-2 shrink-0">
              <ShareButton />
              <Button asChild size="sm" variant="ghost">
                <a
                  href={`https://github.com/${result.owner}/${result.repo}`}
                  rel="noopener noreferrer"
                  target="_blank"
                >
                  {t("viewOnGitHub")}
                  <ExternalLink className="h-4 w-4" />
                </a>
              </Button>
            </div>
          </div>
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
            <span className="flex items-center gap-1">
              <GitCommit className="h-4 w-4" />
              {result.commitSha.slice(0, SHORT_SHA_LENGTH)}
            </span>
            <span>{t("analyzedAt", { date: formatAnalysisDate(result.analyzedAt) })}</span>
          </div>
        </motion.header>

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
            onViewModeChange={setViewMode}
            query={query}
            statuses={statuses}
            viewMode={viewMode}
          />
          {hasFilter && !hasResults ? (
            <FilterEmptyState />
          ) : viewMode === "tree" ? (
            <TreeView suites={filteredSuites} />
          ) : (
            <TestList suites={filteredSuites} />
          )}
        </motion.section>
      </div>
    </motion.main>
  );
};
