"use client";

import { FilterEmptyState } from "@/components/feedback";

import { ExecutiveSummary } from "./executive-summary";
import { OldVersionBanner } from "./old-version-banner";
import { ReadingProgressBar } from "./reading-progress-bar";
import { TocSidebar } from "./toc-sidebar";
import { VirtualizedDocumentView } from "./virtualized-document-view";
import { DocumentNavigationProvider } from "../contexts";
import { useDocumentFilter } from "../hooks/use-document-filter";
import type {
  BehaviorCacheStats,
  RepoVersionInfo,
  SpecDocument,
  SpecLanguage,
  VersionInfo,
} from "../types";

type DocumentViewProps = {
  analysisCommitSha?: string;
  behaviorCacheStats?: BehaviorCacheStats;
  commitSha?: string;
  document: SpecDocument;
  isGeneratingOtherLanguage?: boolean;
  isLoadingVersions?: boolean;
  isRegenerating?: boolean;
  latestVersion?: number;
  onGenerateForCurrentCommit?: () => void;
  onGenerateNewLanguage?: (language: SpecLanguage) => void;
  onLanguageSwitch?: (language: SpecLanguage) => void;
  onRegenerate?: () => void;
  onVersionSwitch?: (version: number) => void;
  onViewLatest?: () => void;
  owner?: string;
  repo?: string;
  versions?: (VersionInfo | RepoVersionInfo)[];
};

export const DocumentView = ({
  analysisCommitSha,
  behaviorCacheStats,
  commitSha,
  document,
  isGeneratingOtherLanguage,
  isLoadingVersions,
  isRegenerating,
  latestVersion,
  onGenerateForCurrentCommit,
  onGenerateNewLanguage,
  onLanguageSwitch,
  onRegenerate,
  onVersionSwitch,
  onViewLatest,
  owner,
  repo,
  versions,
}: DocumentViewProps) => {
  const { clearFilters, filteredDocument, filterInfo, hasFilter, matchCount } =
    useDocumentFilter(document);

  const showEmptyState = hasFilter && matchCount === 0;
  const currentVersion = document.version;
  const isViewingOldVersion =
    latestVersion !== undefined && currentVersion !== undefined && currentVersion < latestVersion;
  // Check if spec was generated from a different commit than current analysis
  const isOutdatedCommit =
    analysisCommitSha !== undefined && commitSha !== undefined && analysisCommitSha !== commitSha;

  return (
    <DocumentNavigationProvider document={document}>
      <ReadingProgressBar />
      {isOutdatedCommit && onGenerateForCurrentCommit && (
        <OldVersionBanner onViewLatest={onGenerateForCurrentCommit} outdatedCommit />
      )}
      {!isOutdatedCommit && isViewingOldVersion && onViewLatest && (
        <OldVersionBanner
          currentVersion={currentVersion!}
          latestVersion={latestVersion!}
          onViewLatest={onViewLatest}
        />
      )}
      <div className="flex flex-col lg:flex-row lg:gap-6">
        <TocSidebar document={document} filteredDocument={filteredDocument} hasFilter={hasFilter} />

        <div className="flex-1 space-y-6 min-w-0">
          <ExecutiveSummary
            behaviorCacheStats={behaviorCacheStats}
            commitSha={commitSha}
            document={document}
            isGeneratingOtherLanguage={isGeneratingOtherLanguage}
            isLoadingVersions={isLoadingVersions}
            isRegenerating={isRegenerating}
            latestVersion={latestVersion}
            onGenerateNewLanguage={onGenerateNewLanguage}
            onLanguageSwitch={onLanguageSwitch}
            onRegenerate={isOutdatedCommit ? onGenerateForCurrentCommit : onRegenerate}
            onVersionSwitch={onVersionSwitch}
            owner={owner}
            repo={repo}
            versions={versions}
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
