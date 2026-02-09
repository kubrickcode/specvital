"use client";

import { RefreshCw } from "lucide-react";
import { useLocale, useTranslations } from "next-intl";

import { ErrorFallback } from "@/components/feedback";
import { Button } from "@/components/ui/button";
import { useUsage } from "@/features/account";
import {
  AnalysisContent,
  AnalysisSkeleton,
  useAnalysis,
  useAutoTrackHistory,
  useCommitSelect,
} from "@/features/analysis";
import type { AnalysisStatus } from "@/features/analysis";
import { useAuth } from "@/features/auth";
import { Link } from "@/i18n/navigation";

type AnalysisPageProps = {
  owner: string;
  repo: string;
};

const mapToSkeletonStatus = (status: string): AnalysisStatus => {
  if (status === "queued") return "queued";
  if (status === "analyzing") return "analyzing";
  return "loading";
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

const isQuotaExceeded = (
  analysis:
    | { limit?: number | null; percentage?: number | null; reserved: number; used: number }
    | undefined
): boolean => {
  if (!analysis) return false;
  if (analysis.limit === null || analysis.limit === undefined) return false; // Unlimited
  const totalUsed = analysis.used + analysis.reserved;
  return totalUsed >= analysis.limit;
};

const formatResetDate = (resetAt: string, locale: string): string => {
  return new Date(resetAt).toLocaleDateString(locale, {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
};

export const AnalysisPage = ({ owner, repo }: AnalysisPageProps) => {
  const t = useTranslations("analyze");
  const locale = useLocale();
  const { isAuthenticated } = useAuth();
  const { data: usageData, isLoading: isLoadingUsage } = useUsage(isAuthenticated);
  const { commitSha: selectedCommit } = useCommitSelect();

  // Check if analysis quota is exceeded for authenticated users
  const quotaExceeded = isAuthenticated && isQuotaExceeded(usageData?.analysis);

  const { data, error, isLoading, refetch, startedAt, status } = useAnalysis(owner, repo, {
    commit: selectedCommit,
    enabled: !quotaExceeded,
  });

  useAutoTrackHistory(owner, repo, status === "completed");

  // Show loading while checking quota for authenticated users
  if (isAuthenticated && isLoadingUsage) {
    return <AnalysisSkeleton owner={owner} repo={repo} startedAt={null} status="loading" />;
  }

  // Show quota exceeded error
  if (quotaExceeded && usageData) {
    const analysis = usageData.analysis;
    const description = `${t("quotaExceeded.description", {
      resetDate: formatResetDate(usageData.resetAt, locale),
    })} ${t("quotaExceeded.currentUsage", {
      limit: analysis.limit ?? 0,
      used: analysis.used + analysis.reserved,
    })}`;

    return (
      <ErrorFallback
        action={
          <Button asChild variant="outline">
            <Link href="/account">{t("quotaExceeded.viewAccount")}</Link>
          </Button>
        }
        description={description}
        title={t("quotaExceeded.title")}
      />
    );
  }

  if (isLoading) {
    return (
      <AnalysisSkeleton
        owner={owner}
        repo={repo}
        startedAt={startedAt}
        status={mapToSkeletonStatus(status)}
      />
    );
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
