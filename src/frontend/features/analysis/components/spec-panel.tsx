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
  useDocumentFilter,
  useGenerationProgress,
  useQuotaConfirmDialog,
  useSpecView,
} from "@/features/spec-view";
import type { GenerationState, SpecLanguage } from "@/features/spec-view";

import { SpecToolbar } from "./spec-toolbar";
import { useFilterState } from "../hooks/use-filter-state";

type SpecPanelProps = {
  analysisId: string;
  availableFrameworks: string[];
  totalTests: number;
};

export const SpecPanel = ({ analysisId, availableFrameworks, totalTests }: SpecPanelProps) => {
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

  // Spec view - pass pendingLanguage to poll correct status during generation
  const {
    data: specDocument,
    generationState,
    isFetching,
    requestGenerate,
    serverStatus,
  } = useSpecView(analysisId, {
    language: pendingLanguage ?? undefined,
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
    }

    // Handle completion
    if (effectiveState === "completed" && prevEffectiveStateRef.current !== "completed") {
      updateProgressStatus("completed");
      setIsWaitingForGeneration(false);
      setPendingLanguage(null);

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

      if (isProgressInBackground) {
        toast.error(t("generateFailed.title"));
      }
      closeProgressModal();
    }

    prevEffectiveStateRef.current = effectiveState;
  }, [
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
      onConfirm: (selectedLanguage) => startGeneration(selectedLanguage, false),
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
    if (specDocument) {
      if (hasFilter && matchCount === 0) {
        return <FilterEmptyState filterInfo={filterInfo} onReset={resetFilters} />;
      }
      return (
        <DocumentView
          document={specDocument}
          isGeneratingOtherLanguage={isGeneratingOtherLanguage}
          isRegenerating={isGenerating && !isGeneratingOtherLanguage}
          onGenerateNewLanguage={handleGenerateNewLanguage}
          onLanguageSwitch={handleExistingLanguageSwitch}
          onRegenerate={handleRegenerate}
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
