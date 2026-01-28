"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useLocale, useTranslations } from "next-intl";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";

import { FilterEmptyState, LoadingFallback } from "@/components/feedback";
import { useUsage } from "@/features/account";
import { useAuth, useSpecLoginDialog } from "@/features/auth";
import {
  DocumentView,
  EmptyDocument,
  GenerationProgressModal,
  GenerationStatus,
  QuotaConfirmDialog,
  SpecAccessError,
  useDocumentFilter,
  useGenerationProgress,
  specGenerationStatusKeys,
  useQuotaConfirmDialog,
  useRepoSpecView,
  useRepoVersionHistory,
  useSpecGenerationStatus,
  useSpecView,
} from "@/features/spec-view";
import type {
  SpecGenerationMode,
  SpecGenerationStatusEnum,
  SpecLanguage,
} from "@/features/spec-view";
import { getSpecGenerationTaskId } from "@/features/spec-view/utils/task-ids";
import { addTask, getTask, removeTask, updateTask, useTaskStartedAt } from "@/lib/background-tasks";

import { SpecToolbar } from "./spec-toolbar";
import { useFilterState } from "../hooks/use-filter-state";

type SpecPanelProps = {
  analysisCommitSha: string;
  analysisId: string;
  availableFrameworks: string[];
  owner: string;
  repo: string;
  totalTests: number;
};

export const SpecPanel = ({
  analysisCommitSha,
  analysisId,
  availableFrameworks,
  owner,
  repo,
  totalTests,
}: SpecPanelProps) => {
  const t = useTranslations("specView.toast");
  const locale = useLocale();
  const queryClient = useQueryClient();

  // Filter state (shared URL params)
  const { frameworks, query, setFrameworks, setQuery, setStatuses, statuses } = useFilterState();

  // Auth & quota
  const { isAuthenticated } = useAuth();
  const { open: openSpecLoginDialog } = useSpecLoginDialog();
  const { data: usageData } = useUsage(isAuthenticated);
  const { open: openQuotaConfirmDialog } = useQuotaConfirmDialog();

  // Progress modal
  const {
    bringToForeground: bringProgressToForeground,
    close: closeProgressModal,
    isInBackground: isProgressInBackground,
    open: openProgressModal,
  } = useGenerationProgress();

  // Local state
  // generatingLanguage: tracks language currently being generated (enables polling when non-null)
  // selectedLanguage: tracks language for document viewing (does not trigger polling)
  const [generatingLanguage, setGeneratingLanguage] = useState<SpecLanguage | null>(null);
  const [selectedLanguage, setSelectedLanguage] = useState<SpecLanguage | undefined>(undefined);
  const [selectedVersion, setSelectedVersion] = useState<number | undefined>(undefined);

  // Determine which language to use for document fetching
  // Priority: generatingLanguage (during generation) > selectedLanguage (user selection)
  const documentLanguage = generatingLanguage ?? selectedLanguage;

  // Spec document - use repository-based API for cross-analysis version access
  const {
    accessError: repoAccessError,
    behaviorCacheStats,
    commitSha,
    data: specDocument,
    isFetching,
    isLoading: isLoadingDocument,
  } = useRepoSpecView(owner, repo, {
    language: documentLanguage,
    version: selectedVersion,
  });

  // Generation mutation - analysisId based (generation targets specific analysis)
  const {
    accessError: genAccessError,
    isRequesting,
    requestGenerate,
  } = useSpecView(analysisId, {
    language: documentLanguage,
  });

  // Generation status polling - only enabled when generating
  const { status: serverStatus } = useSpecGenerationStatus(analysisId, {
    enabled: generatingLanguage !== null,
    language: generatingLanguage ?? undefined,
    onCompleted: () => {
      // Deduplicate with SpecGenerationMonitor (global):
      // getTask returns null if already removed by another handler
      const taskId = `spec-generation-${analysisId}`;
      const task = getTask(taskId);
      if (task) {
        removeTask(taskId);
        if (isProgressInBackground) {
          toast.success(t("generationComplete.title"), {
            description: t("generationComplete.description"),
          });
        }
      }
      closeProgressModal();
      setSelectedLanguage(generatingLanguage ?? undefined);
      setGeneratingLanguage(null);
    },
    onFailed: () => {
      const taskId = `spec-generation-${analysisId}`;
      const task = getTask(taskId);
      if (task) {
        removeTask(taskId);
        if (isProgressInBackground) {
          toast.error(t("generateFailed.title"));
        }
      }
      closeProgressModal();
      setGeneratingLanguage(null);
    },
    owner,
    repo,
  });

  // Combine access errors (repo error takes precedence)
  const accessError = repoAccessError ?? genAccessError;

  // Version history - use repository-based API for cross-analysis version access
  const {
    data: repoVersionHistory,
    isFetching: isFetchingVersions,
    isLoading: isLoadingVersions,
  } = useRepoVersionHistory(owner, repo, specDocument?.language, {
    enabled: Boolean(specDocument?.language),
  });

  // Document filter
  const { matchCount } = useDocumentFilter(specDocument);

  // Calculate total behaviors
  const totalBehaviors =
    specDocument?.domains.reduce(
      (acc, domain) =>
        acc + domain.features.reduce((fAcc, feature) => fAcc + feature.behaviors.length, 0),
      0
    ) ?? totalTests;

  const hasFilter = query.trim().length > 0 || frameworks.length > 0 || statuses.length > 0;

  // Update TaskStore based on server status
  useEffect(() => {
    if (generatingLanguage === null || !serverStatus) return;

    const taskId = `spec-generation-${analysisId}`;
    const existingTask = getTask(taskId);
    if (existingTask) {
      if (serverStatus === "pending" && existingTask.status !== "queued") {
        updateTask(taskId, { status: "queued" });
      } else if (serverStatus === "running" && existingTask.status !== "processing") {
        updateTask(taskId, { startedAt: new Date().toISOString(), status: "processing" });
      }
    }
  }, [analysisId, generatingLanguage, serverStatus]);

  // Show toast when switching to background
  const prevIsInBackgroundRef = useRef(false);
  useEffect(() => {
    if (isProgressInBackground && !prevIsInBackgroundRef.current) {
      toast.info(t("generationInProgress.title"), {
        action: {
          label: t("generationInProgress.viewProgress"),
          onClick: () => bringProgressToForeground(),
        },
        description: t("generationInProgress.description"),
      });
    }
    prevIsInBackgroundRef.current = isProgressInBackground;
  }, [isProgressInBackground, t, bringProgressToForeground]);

  const startGeneration = async (language: SpecLanguage, mode: SpecGenerationMode = "initial") => {
    // Clear stale generation status cache to prevent previous "completed"
    // from being misinterpreted as new generation completion
    queryClient.removeQueries({
      queryKey: specGenerationStatusKeys.status(analysisId, language),
    });

    // Open modal immediately for instant UX feedback (status=null shows "pending" state)
    openProgressModal({
      analysisId,
      onViewDocument: () => {},
    });

    try {
      // Await mutation so the backend River job is guaranteed to exist
      // before polling starts — prevents race condition where polling
      // returns stale "completed" from a previous generation
      await requestGenerate(language, mode);

      // Register background task AFTER mutation succeeds.
      // SpecGenerationMonitor subscribes to the task store and immediately
      // starts polling when a task appears. If the task is added before the
      // mutation, the monitor polls the backend before the new job exists,
      // receives the previous job's "completed" status, and fires
      // onCompleted prematurely — closing the modal and showing a stale toast.
      addTask({
        id: `spec-generation-${analysisId}`,
        metadata: {
          analysisId,
          language,
          owner,
          repo,
        },
        startedAt: new Date().toISOString(),
        status: "queued",
        type: "spec-generation",
      });

      setGeneratingLanguage(language);
    } catch (error) {
      // Toast already shown by useSpecView.onError handler
      console.error("Spec generation request failed:", error);
      closeProgressModal();
    }
  };

  const handleGenerate = () => {
    if (!isAuthenticated) {
      openSpecLoginDialog();
      return;
    }

    openQuotaConfirmDialog({
      analysisId,
      estimatedCost: totalTests,
      locale,
      onConfirm: (confirmedLanguage, mode) => startGeneration(confirmedLanguage, mode),
      usage: usageData ?? null,
    });
  };

  const handleRegenerate = () => {
    if (!isAuthenticated) {
      openSpecLoginDialog();
      return;
    }

    openQuotaConfirmDialog({
      analysisId,
      estimatedCost: totalTests,
      initialLanguage: specDocument?.language,
      isRegenerate: true,
      locale,
      onConfirm: (confirmedLanguage, mode) => startGeneration(confirmedLanguage, mode),
      usage: usageData ?? null,
    });
  };

  /**
   * Switch to an already-generated language (free, instant).
   * No authentication required for viewing existing documents.
   * Does NOT trigger polling - only updates document query.
   */
  const handleExistingLanguageSwitch = (language: SpecLanguage) => {
    setSelectedLanguage(language);
    setSelectedVersion(undefined); // Reset to latest version when switching language
  };

  /**
   * Switch to a specific version (free, instant).
   */
  const handleVersionSwitch = (version: number) => {
    setSelectedVersion(version);
  };

  /**
   * Switch to the latest version (triggered from old version banner).
   */
  const handleViewLatest = () => {
    setSelectedVersion(undefined);
  };

  /**
   * Generate a new language (costs quota, requires auth).
   */
  const handleGenerateNewLanguage = (language: SpecLanguage) => {
    if (!isAuthenticated) {
      openSpecLoginDialog();
      return;
    }

    openQuotaConfirmDialog({
      analysisId,
      estimatedCost: totalTests,
      initialLanguage: language,
      locale,
      onConfirm: (_confirmedLanguage, mode) => startGeneration(language, mode),
      usage: usageData ?? null,
    });
  };

  // Derive UI states from server status (single source of truth)
  const isGenerating = serverStatus === "pending" || serverStatus === "running";
  const isGeneratingOtherLanguage = generatingLanguage !== null && specDocument !== null;

  // Map server status to GenerationStatus component status
  const displayStatus: SpecGenerationStatusEnum | null = serverStatus;

  const resetFilters = () => {
    setQuery(null);
    setFrameworks(null);
    setStatuses(null);
  };

  const filterInfo = {
    frameworks,
    query,
    statuses,
  };

  const renderContent = () => {
    // Handle access errors first
    if (accessError) {
      return <SpecAccessError type={accessError} />;
    }

    if (specDocument) {
      if (hasFilter && matchCount === 0) {
        return <FilterEmptyState filterInfo={filterInfo} onReset={resetFilters} />;
      }
      // Suppress latestVersion while either query is refetching to avoid
      // a race condition where version history updates before the document,
      // causing a transient isViewingOldVersion=true flash (banner flicker).
      const isRefetchingEither = isFetching || isFetchingVersions;
      const latestVersion =
        !isRefetchingEither && repoVersionHistory?.data && repoVersionHistory.data.length > 0
          ? Math.max(...repoVersionHistory.data.map((v) => v.version))
          : undefined;

      return (
        <DocumentView
          analysisCommitSha={analysisCommitSha}
          behaviorCacheStats={behaviorCacheStats}
          commitSha={commitSha}
          document={specDocument}
          isGeneratingOtherLanguage={isGeneratingOtherLanguage}
          isLoadingVersions={isLoadingVersions}
          isRegenerating={isGenerating && !isGeneratingOtherLanguage}
          latestVersion={latestVersion}
          onGenerateForCurrentCommit={handleGenerate}
          onGenerateNewLanguage={handleGenerateNewLanguage}
          onLanguageSwitch={handleExistingLanguageSwitch}
          onRegenerate={handleRegenerate}
          onVersionSwitch={handleVersionSwitch}
          onViewLatest={handleViewLatest}
          owner={owner}
          repo={repo}
          versions={repoVersionHistory?.data}
        />
      );
    }

    // Show generation status when actively generating without existing document
    if (displayStatus === "pending" || displayStatus === "running") {
      return <GenerationStatus startedAt={taskStartedAt} status={displayStatus} />;
    }

    // Show loading while initially fetching spec document
    if (isLoadingDocument || isRequesting) {
      return <LoadingFallback className="py-16" fullScreen={false} />;
    }

    const specviewQuota = usageData?.specview
      ? {
          limit: usageData.specview.limit ?? null,
          percentage: usageData.specview.percentage ?? null,
          used: usageData.specview.used,
        }
      : null;

    return <EmptyDocument isLoading={false} onGenerate={handleGenerate} quota={specviewQuota} />;
  };

  const taskStartedAt = useTaskStartedAt(getSpecGenerationTaskId(analysisId));

  return (
    <div aria-labelledby="tab-spec" className="p-5 space-y-4" id="tabpanel-spec" role="tabpanel">
      <SpecToolbar
        availableFrameworks={availableFrameworks}
        frameworks={frameworks}
        hasDocument={Boolean(specDocument)}
        hasFilter={hasFilter}
        matchCount={matchCount}
        onFrameworksChange={setFrameworks}
        onQueryChange={setQuery}
        onStatusesChange={setStatuses}
        query={query}
        statuses={statuses}
        totalCount={totalBehaviors}
      />

      <QuotaConfirmDialog />
      <GenerationProgressModal startedAt={taskStartedAt} status={serverStatus} />
      {renderContent()}
    </div>
  );
};
