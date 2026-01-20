"use client";

import { useTranslations } from "next-intl";

import { ExecutiveSummary } from "./executive-summary";
import { FilterEmptyState } from "./filter-empty-state";
import { ReadingProgressBar } from "./reading-progress-bar";
import { TocSidebar } from "./toc-sidebar";
import { VirtualizedDocumentView } from "./virtualized-document-view";
import { useDocumentFilter } from "../hooks/use-document-filter";
import { useScrollSync } from "../hooks/use-scroll-sync";
import type { SpecDocument, SpecLanguage } from "../types";

type DocumentViewProps = {
  document: SpecDocument;
  isGeneratingOtherLanguage?: boolean;
  isRegenerating?: boolean;
  onLanguageSwitch?: (language: SpecLanguage) => void;
  onRegenerate?: () => void;
};

export const DocumentView = ({
  document,
  isGeneratingOtherLanguage,
  isRegenerating,
  onLanguageSwitch,
  onRegenerate,
}: DocumentViewProps) => {
  const t = useTranslations("analyze.filter");
  useScrollSync();

  const { clearFilters, filteredDocument, hasFilter, matchCount } = useDocumentFilter(document);

  const showEmptyState = hasFilter && matchCount === 0;

  return (
    <>
      <ReadingProgressBar />
      <div className="flex flex-col lg:flex-row lg:gap-6">
        <TocSidebar document={document} filteredDocument={filteredDocument} hasFilter={hasFilter} />

        <div className="flex-1 space-y-6 min-w-0">
          <ExecutiveSummary
            document={document}
            isGeneratingOtherLanguage={isGeneratingOtherLanguage}
            isRegenerating={isRegenerating}
            onLanguageSwitch={onLanguageSwitch}
            onRegenerate={onRegenerate}
          />

          {showEmptyState ? (
            <FilterEmptyState
              description={t("noResultsDescription")}
              onClearFilters={clearFilters}
              title={t("noResults")}
            />
          ) : filteredDocument ? (
            <VirtualizedDocumentView document={filteredDocument} hasFilter={hasFilter} />
          ) : null}
        </div>
      </div>
    </>
  );
};
