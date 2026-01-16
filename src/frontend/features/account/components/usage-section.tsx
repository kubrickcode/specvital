"use client";

import { AlertCircle } from "lucide-react";
import { useLocale, useTranslations } from "next-intl";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

import type { UsageStatusResponse } from "../api";
import { getResetInfo } from "../utils";
import { UsageProgress } from "./usage-progress";

type UsageSectionProps = {
  error?: Error | null;
  isLoading: boolean;
  usage?: UsageStatusResponse;
};

export const UsageSection = ({ error, isLoading, usage }: UsageSectionProps) => {
  const t = useTranslations("account");
  const tFeatures = useTranslations("pricing.features");
  const locale = useLocale();

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">{t("usage.title")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <Skeleton className="h-16 w-full" />
          <Skeleton className="h-16 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (error || !usage) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">{t("usage.title")}</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center gap-2 py-8 text-muted-foreground">
          <AlertCircle className="h-4 w-4" />
          <span>{t("usage.unavailable")}</span>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">{t("usage.title")}</CardTitle>
          <CardDescription>
            {t("usage.resetIn", getResetInfo(usage.resetAt, locale))}
          </CardDescription>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        <UsageProgress
          label={t("usage.specview")}
          limit={usage.specview.limit ?? null}
          percentage={usage.specview.percentage ?? null}
          unit={tFeatures("specview.unit")}
          used={usage.specview.used}
        />
        <UsageProgress
          label={t("usage.analysis")}
          limit={usage.analysis.limit ?? null}
          percentage={usage.analysis.percentage ?? null}
          unit={tFeatures("analysis.unit")}
          used={usage.analysis.used}
        />
      </CardContent>
    </Card>
  );
};
