"use client";

import { useLocale, useTranslations } from "next-intl";
import { useRef, useEffect, useState } from "react";
import { toast } from "sonner";

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
import type { SpecLanguage } from "@/features/spec-view";

import { FilterEmptyState } from "./filter-empty-state";
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

  // Spec view
  const {
    data: specDocument,
    generationStatus,
    isGenerating,
    requestGenerate,
  } = useSpecView(analysisId);

  // Language switch state
  const [isGeneratingOtherLanguage, setIsGeneratingOtherLanguage] = useState(false);
  const [pendingLanguage, setPendingLanguage] = useState<SpecLanguage | null>(null);

  // Document filter
  const { matchCount } = useDocumentFilter(specDocument);

  // Calculate total behaviors
  const totalBehaviors =
    specDocument?.domains.reduce(
      (acc, d) => acc + d.features.reduce((facc, f) => facc + f.behaviors.length, 0),
      0
    ) ?? totalTests;

  const hasFilter = query.trim().length > 0 || frameworks.length > 0 || statuses.length > 0;

  // Sync generation status with progress modal
  useEffect(() => {
    if (generationStatus) {
      updateProgressStatus(generationStatus);

      if (isProgressInBackground) {
        if (generationStatus === "completed") {
          closeProgressModal();
          toast.success(t("generationComplete.title"), {
            description: t("generationComplete.description"),
          });
        } else if (generationStatus === "failed") {
          closeProgressModal();
          toast.error(t("generateFailed.title"));
        }
      } else {
        if (generationStatus === "completed" || generationStatus === "failed") {
          closeProgressModal();
        }
      }
    }
  }, [generationStatus, updateProgressStatus, closeProgressModal, isProgressInBackground, t]);

  // Handle completion when document becomes available
  const prevSpecDocumentRef = useRef(specDocument);
  useEffect(() => {
    if (specDocument && !prevSpecDocumentRef.current) {
      updateProgressStatus("completed");
      if (isProgressInBackground) {
        closeProgressModal();
        toast.success(t("generationComplete.title"), {
          description: t("generationComplete.description"),
        });
      } else {
        closeProgressModal();
      }
    }
    prevSpecDocumentRef.current = specDocument;
  }, [specDocument, updateProgressStatus, closeProgressModal, isProgressInBackground, t]);

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

  const handleGenerate = () => {
    if (!isAuthenticated) {
      openSpecLoginDialog();
      return;
    }

    openQuotaConfirmDialog({
      analysisId,
      estimatedCost: totalTests,
      locale,
      onConfirm: (selectedLanguage) => {
        requestGenerate(selectedLanguage, false);
        openProgressModal({
          analysisId,
          onViewDocument: () => {},
          status: "pending",
        });
      },
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
      onConfirm: (selectedLanguage, isForceRegenerate) => {
        requestGenerate(selectedLanguage, isForceRegenerate);
        openProgressModal({
          analysisId,
          onViewDocument: () => {},
          status: "pending",
        });
      },
      usage: usageData ?? null,
    });
  };

  const handleLanguageSwitch = (language: SpecLanguage) => {
    if (!isAuthenticated) {
      openSpecLoginDialog();
      return;
    }

    setPendingLanguage(language);
    openQuotaConfirmDialog({
      analysisId,
      estimatedCost: totalTests,
      locale,
      onConfirm: () => {
        setIsGeneratingOtherLanguage(true);
        requestGenerate(language, false);
        openProgressModal({
          analysisId,
          onViewDocument: () => {},
          status: "pending",
        });
      },
      usage: usageData ?? null,
    });
  };

  // Reset language switch state when generation completes or language matches
  useEffect(() => {
    const shouldReset =
      (pendingLanguage && specDocument?.language === pendingLanguage) ||
      generationStatus === "completed" ||
      generationStatus === "failed";

    if (shouldReset) {
      setIsGeneratingOtherLanguage(false);
      setPendingLanguage(null);
    }
  }, [specDocument?.language, pendingLanguage, generationStatus]);

  const renderContent = () => {
    if (specDocument) {
      if (hasFilter && matchCount === 0) {
        return <FilterEmptyState />;
      }
      return (
        <DocumentView
          document={specDocument}
          isGeneratingOtherLanguage={isGeneratingOtherLanguage}
          isRegenerating={isGenerating && !isGeneratingOtherLanguage}
          onLanguageSwitch={handleLanguageSwitch}
          onRegenerate={handleRegenerate}
        />
      );
    }

    if (isGenerating && generationStatus) {
      return <GenerationStatus onRetry={() => requestGenerate()} status={generationStatus} />;
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
    <div aria-labelledby="tab-spec" className="space-y-4" id="tabpanel-spec" role="tabpanel">
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
