"use client";

import { RefreshCw } from "lucide-react";
import { useTranslations } from "next-intl";

import { ErrorFallback } from "@/components/feedback";
import { Button } from "@/components/ui/button";
import { AnalysisContent, AnalysisSkeleton, useAnalysis } from "@/features/analysis";

type AnalysisPageProps = {
  owner: string;
  repo: string;
};

type SkeletonStatus = "loading" | "queued" | "analyzing";

const getSkeletonProps = (
  status: string,
  t: ReturnType<typeof useTranslations<"analyze">>
): { description: string; status: SkeletonStatus; title: string } => {
  switch (status) {
    case "queued":
      return {
        description: t("status.queuedDescription"),
        status: "queued",
        title: t("status.queuedTitle"),
      };
    case "analyzing":
      return {
        description: t("status.analyzingDescription"),
        status: "analyzing",
        title: t("status.analyzingTitle"),
      };
    default:
      return {
        description: t("status.loadingDescription"),
        status: "loading",
        title: t("status.loadingTitle"),
      };
  }
};

const getDisplayErrorMessage = (
  error: Error | null,
  status: string,
  t: ReturnType<typeof useTranslations<"analyze">>
): string => {
  if (error) {
    return error.message;
  }
  if (status === "failed") {
    return t("status.failed");
  }
  return t("status.error");
};

export const AnalysisPage = ({ owner, repo }: AnalysisPageProps) => {
  const t = useTranslations("analyze");
  const { data, error, isLoading, refetch, status } = useAnalysis(owner, repo);

  if (isLoading) {
    return <AnalysisSkeleton {...getSkeletonProps(status, t)} />;
  }

  if (status === "error" || status === "failed") {
    return (
      <ErrorFallback
        action={
          <Button className="gap-2" onClick={refetch} variant="outline">
            <RefreshCw className="h-4 w-4" />
            {t("retry")}
          </Button>
        }
        description={getDisplayErrorMessage(error, status, t)}
        title={t("status.error")}
      />
    );
  }

  if (data) {
    return <AnalysisContent result={data} />;
  }

  return null;
};
