"use client";

import { useLocale, useTranslations } from "next-intl";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";

import { FilterEmptyState } from "@/components/feedback";
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
  useQuotaConfirmDialog,
  useRepoSpecView,
  useRepoVersionHistory,
  useSpecView,
} from "@/features/spec-view";
import type { GenerationState, SpecLanguage } from "@/features/spec-view";
import { addTask, getTask, removeTask, updateTask } from "@/lib/background-tasks";

import { SpecToolbar } from "./spec-toolbar";
import { useFilterState } from "../hooks/use-filter-state";

type SpecPanelProps = {
  analysisId: string;
  availableFrameworks: string[];
  owner: string;
  repo: string;
  totalTests: number;
};

export const SpecPanel = ({
  analysisId,
  availableFrameworks,
  owner,
  repo,
  totalTests,
}: SpecPanelProps) => {
  const t = useTranslations("specView.toast");
  const locale = useLocale();

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
    updateStatus: updateProgressStatus,
  } = useGenerationProgress();

  // Local state: tracking whether user initiated a generation
  const [isWaitingForGeneration, setIsWaitingForGeneration] = useState(false);
  const [pendingLanguage, setPendingLanguage] = useState<SpecLanguage | null>(null);
  const [selectedVersion, setSelectedVersion] = useState<number | undefined>(undefined);

  // Spec document - use repository-based API for cross-analysis version access
  const {
    accessError: repoAccessError,
    behaviorCacheStats,
    commitSha,
    data: specDocument,
    isFetching,
  } = useRepoSpecView(owner, repo, {
    language: pendingLanguage ?? undefined,
    version: selectedVersion,
  });

  // Generation functionality - analysisId based (generation targets specific analysis)
  const {
    accessError: genAccessError,
    generationState,
    requestGenerate,
    serverStatus,
  } = useSpecView(analysisId, {
    language: pendingLanguage ?? undefined,
  });

  // Combine access errors (repo error takes precedence)
  const accessError = repoAccessError ?? genAccessError;

  // Version history - use repository-based API for cross-analysis version access
  const { data: repoVersionHistory, isLoading: isLoadingVersions } = useRepoVersionHistory(
    owner,
    repo,
    specDocument?.language,
    {
      enabled: Boolean(specDocument?.language),
    }
  );

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

  // Determine effective state: combine server state with local waiting state
  const effectiveState: GenerationState = (() => {
    if (!isWaitingForGeneration) {
      return generationState;
    }
    // While waiting for generation, only trust server-provided status
    if (serverStatus === "pending") return "pending";
    if (serverStatus === "running") return "running";
    if (serverStatus === "completed") return "completed";
    if (serverStatus === "failed") return "failed";
    // Server returned document instead of status → generation completed
    // But only if fetch completed (to avoid treating cached data as new completion)
    if (generationState === "completed" && !isFetching) return "completed";
    // Still fetching or no server status yet
    return "pending";
  })();

  // Single effect to handle generation state changes
  const prevEffectiveStateRef = useRef<GenerationState | null>(null);
  useEffect(() => {
    if (!isWaitingForGeneration) {
      prevEffectiveStateRef.current = null;
      return;
    }

    // Update modal status for in-progress states
    if (
      effectiveState === "requesting" ||
      effectiveState === "pending" ||
      effectiveState === "running"
    ) {
      updateProgressStatus(effectiveState === "requesting" ? "pending" : effectiveState);

      // Update TaskStore status
      const taskId = `spec-generation-${analysisId}`;
      const existingTask = getTask(taskId);
      if (existingTask) {
        if (effectiveState === "pending" && existingTask.status !== "queued") {
          updateTask(taskId, { status: "queued" });
        } else if (effectiveState === "running" && existingTask.status !== "processing") {
          updateTask(taskId, { startedAt: new Date().toISOString(), status: "processing" });
        }
      }
    }

    // Handle completion
    if (effectiveState === "completed" && prevEffectiveStateRef.current !== "completed") {
      updateProgressStatus("completed");
      setIsWaitingForGeneration(false);
      setPendingLanguage(null);

      // Remove from TaskStore on completion
      removeTask(`spec-generation-${analysisId}`);

      if (isProgressInBackground) {
        toast.success(t("generationComplete.title"), {
          description: t("generationComplete.description"),
        });
      }
      closeProgressModal();
    }

    // Handle failure
    if (effectiveState === "failed") {
      setIsWaitingForGeneration(false);
      setPendingLanguage(null);

      // Remove from TaskStore on failure
      removeTask(`spec-generation-${analysisId}`);

      if (isProgressInBackground) {
        toast.error(t("generateFailed.title"));
      }
      closeProgressModal();
    }

    prevEffectiveStateRef.current = effectiveState;
  }, [
    analysisId,
    effectiveState,
    isWaitingForGeneration,
    isProgressInBackground,
    updateProgressStatus,
    closeProgressModal,
    t,
  ]);

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

  const startGeneration = (language: SpecLanguage, forceRegenerate = false) => {
    setPendingLanguage(language);
    setIsWaitingForGeneration(true);
    requestGenerate(language, forceRegenerate);
    openProgressModal({
      analysisId,
      onViewDocument: () => {},
      status: "pending",
    });

    // Register task to global TaskStore
    addTask({
      id: `spec-generation-${analysisId}`,
      metadata: {
        analysisId,
        language,
        owner,
        repo,
      },
      startedAt: null,
      status: "queued",
      type: "spec-generation",
    });
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
      onConfirm: (selectedLanguage, forceRegenerate) =>
        startGeneration(selectedLanguage, forceRegenerate ?? false),
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
      onConfirm: (selectedLanguage, isForceRegenerate) =>
        startGeneration(selectedLanguage, isForceRegenerate ?? false),
      usage: usageData ?? null,
    });
  };

  /**
   * Switch to an already-generated language (free, instant).
   * No authentication required for viewing existing documents.
   */
  const handleExistingLanguageSwitch = (language: SpecLanguage) => {
    setPendingLanguage(language);
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
      onConfirm: () => startGeneration(language, false),
      usage: usageData ?? null,
    });
  };

  // Derive UI states from effectiveState
  const isGenerating = effectiveState === "pending" || effectiveState === "running";
  const isGeneratingOtherLanguage = isWaitingForGeneration && specDocument !== null;

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
      // Calculate latest version from repository-based version history
      const latestVersion =
        repoVersionHistory?.data && repoVersionHistory.data.length > 0
          ? Math.max(...repoVersionHistory.data.map((v) => v.version))
          : undefined;

      return (
        <DocumentView
          behaviorCacheStats={behaviorCacheStats}
          commitSha={commitSha}
          document={specDocument}
          isGeneratingOtherLanguage={isGeneratingOtherLanguage}
          isLoadingVersions={isLoadingVersions}
          isRegenerating={isGenerating && !isGeneratingOtherLanguage}
          latestVersion={latestVersion}
          onGenerateNewLanguage={handleGenerateNewLanguage}
          onLanguageSwitch={handleExistingLanguageSwitch}
          onRegenerate={handleRegenerate}
          onVersionSwitch={handleVersionSwitch}
          onViewLatest={handleViewLatest}
          versions={repoVersionHistory?.data}
        />
      );
    }

    if (effectiveState === "pending" || effectiveState === "running") {
      return <GenerationStatus onRetry={() => requestGenerate()} status={effectiveState} />;
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
      <GenerationProgressModal />
      {renderContent()}
    </div>
  );
};
