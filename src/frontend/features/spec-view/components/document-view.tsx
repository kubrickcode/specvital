"use client";

import { FilterEmptyState } from "@/components/feedback";

import { ExecutiveSummary } from "./executive-summary";
import { ReadingProgressBar } from "./reading-progress-bar";
import { TocSidebar } from "./toc-sidebar";
import { VirtualizedDocumentView } from "./virtualized-document-view";
import { DocumentNavigationProvider } from "../contexts";
import { useDocumentFilter } from "../hooks/use-document-filter";
import type { SpecDocument, SpecLanguage } from "../types";

type DocumentViewProps = {
  document: SpecDocument;
  isGeneratingOtherLanguage?: boolean;
  isRegenerating?: boolean;
  onGenerateNewLanguage?: (language: SpecLanguage) => void;
  onLanguageSwitch?: (language: SpecLanguage) => void;
  onRegenerate?: () => void;
};

export const DocumentView = ({
  document,
  isGeneratingOtherLanguage,
  isRegenerating,
  onGenerateNewLanguage,
  onLanguageSwitch,
  onRegenerate,
}: DocumentViewProps) => {
  const { clearFilters, filteredDocument, filterInfo, hasFilter, matchCount } =
    useDocumentFilter(document);

  const showEmptyState = hasFilter && matchCount === 0;

  return (
    <DocumentNavigationProvider document={document}>
      <ReadingProgressBar />
      <div className="flex flex-col lg:flex-row lg:gap-6">
        <TocSidebar document={document} filteredDocument={filteredDocument} hasFilter={hasFilter} />

        <div className="flex-1 space-y-6 min-w-0">
          <ExecutiveSummary
            document={document}
            isGeneratingOtherLanguage={isGeneratingOtherLanguage}
            isRegenerating={isRegenerating}
            onGenerateNewLanguage={onGenerateNewLanguage}
            onLanguageSwitch={onLanguageSwitch}
            onRegenerate={onRegenerate}
          />

          {showEmptyState ? (
            <FilterEmptyState filterInfo={filterInfo} onReset={clearFilters} />
          ) : filteredDocument ? (
            <VirtualizedDocumentView document={filteredDocument} hasFilter={hasFilter} />
          ) : null}
        </div>
      </div>
    </DocumentNavigationProvider>
  );
};
