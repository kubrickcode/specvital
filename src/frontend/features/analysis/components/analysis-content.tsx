"use client";

import { ExternalLink, GitCommit } from "lucide-react";
import { useTranslations } from "next-intl";
import { useMemo } from "react";

import type { AnalysisResult } from "@/lib/api";
import { formatAnalysisDate, SHORT_SHA_LENGTH } from "@/lib/utils";

import { FilterBar, FilterSummary } from "./filter-bar";
import { FilterEmptyState } from "./filter-empty-state";
import { ShareButton } from "./share-button";
import { StatsCard } from "./stats-card";
import { TestList } from "./test-list";
import { useFilterState } from "../hooks/use-filter-state";
import { filterSuites } from "../utils/filter-suites";

type AnalysisContentProps = {
  result: AnalysisResult;
};

export const AnalysisContent = ({ result }: AnalysisContentProps) => {
  const t = useTranslations("analyze");
  const { frameworks, query, setFrameworks, setQuery, setStatuses, statuses } = useFilterState();

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
    <main className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        <header className="space-y-2">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold">
              {result.owner}/{result.repo}
            </h1>
            <div className="flex items-center gap-2">
              <ShareButton />
              <a
                className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
                href={`https://github.com/${result.owner}/${result.repo}`}
                rel="noopener noreferrer"
                target="_blank"
              >
                {t("viewOnGitHub")}
                <ExternalLink className="h-4 w-4" />
              </a>
            </div>
          </div>
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
            <span className="flex items-center gap-1">
              <GitCommit className="h-4 w-4" />
              {result.commitSha.slice(0, SHORT_SHA_LENGTH)}
            </span>
            <span>{t("analyzedAt", { date: formatAnalysisDate(result.analyzedAt) })}</span>
          </div>
        </header>

        <StatsCard summary={result.summary} />

        <section className="space-y-4">
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
            query={query}
            statuses={statuses}
          />
          {hasFilter && !hasResults ? <FilterEmptyState /> : <TestList suites={filteredSuites} />}
        </section>
      </div>
    </main>
  );
};
