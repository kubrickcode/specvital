"use client";

import { RefreshCw } from "lucide-react";
import { useTranslations } from "next-intl";
import { AnalysisContent, useAnalysis } from "@/features/analysis";
import { Button } from "@/components/ui/button";
import { ErrorFallback, LoadingFallback } from "@/components/feedback";
import { getErrorMessage } from "@/lib/api";

type AnalysisPageProps = {
  owner: string;
  repo: string;
};

const getStatusMessage = (
  status: string,
  t: ReturnType<typeof useTranslations<"analyze">>
): string => {
  switch (status) {
    case "queued":
      return t("status.queued");
    case "analyzing":
      return t("status.analyzing");
    default:
      return t("status.loading");
  }
};

const getDisplayErrorMessage = (
  error: Error | null,
  status: string,
  t: ReturnType<typeof useTranslations<"analyze">>
): string => {
  if (error) {
    return getErrorMessage(error);
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
    return <LoadingFallback message={getStatusMessage(status, t)} />;
  }

  if (status === "error" || status === "failed") {
    return (
      <ErrorFallback
        title={t("status.error")}
        description={getDisplayErrorMessage(error, status, t)}
        action={
          <Button onClick={refetch} variant="outline" className="gap-2">
            <RefreshCw className="h-4 w-4" />
            {t("retry")}
          </Button>
        }
      />
    );
  }

  if (data) {
    return <AnalysisContent result={data} />;
  }

  return null;
};
