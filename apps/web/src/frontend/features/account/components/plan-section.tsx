"use client";

import { AlertCircle, ArrowUpRight, Mail } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Link } from "@/i18n/navigation";
import type { components } from "@/lib/api/generated-types";

import { formatNumber } from "../utils";

type PlanTier = components["schemas"]["PlanTier"];

type PlanSectionProps = {
  error?: Error | null;
  isLoading: boolean;
  plan?: components["schemas"]["PlanInfo"];
};

const getTierVariant = (tier: PlanTier): "default" | "destructive" | "outline" | "secondary" => {
  switch (tier) {
    case "enterprise":
      return "destructive";
    case "pro_plus":
    case "pro":
      return "default";
    default:
      return "secondary";
  }
};

export const PlanSection = ({ error, isLoading, plan }: PlanSectionProps) => {
  const t = useTranslations("account");

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">{t("plan.title")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <Skeleton className="h-6 w-32" />
          <Skeleton className="h-20 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (error || !plan) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">{t("plan.title")}</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center gap-2 py-8 text-muted-foreground">
          <AlertCircle className="h-4 w-4" />
          <span>{t("plan.unavailable")}</span>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">{t("plan.title")}</CardTitle>
          <Badge variant={getTierVariant(plan.tier)}>{t(`plan.tier.${plan.tier}`)}</Badge>
        </div>
      </CardHeader>
      <CardContent>
        <div className="grid gap-3 text-sm">
          <div className="flex justify-between">
            <span className="text-muted-foreground">{t("plan.specviewLimit")}</span>
            <span className="font-medium">
              {formatNumber(plan.specviewMonthlyLimit)}/{t("plan.month")}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">{t("plan.analysisLimit")}</span>
            <span className="font-medium">
              {formatNumber(plan.analysisMonthlyLimit)}/{t("plan.month")}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">{t("plan.retention")}</span>
            <span className="font-medium">
              {plan.retentionDays === null || plan.retentionDays === undefined
                ? t("plan.unlimited")
                : t("plan.days", { days: plan.retentionDays })}
            </span>
          </div>
        </div>
      </CardContent>
      <CardFooter>
        {plan.tier === "enterprise" ? (
          <Button asChild className="w-full" variant="outline">
            <a href="mailto:support@specvital.com">
              <Mail className="size-4" />
              {t("plan.contactUs")}
            </a>
          </Button>
        ) : (
          <Button asChild className="w-full" variant="cta">
            <Link href="/pricing">
              <ArrowUpRight className="size-4" />
              {t("plan.upgrade")}
            </Link>
          </Button>
        )}
      </CardFooter>
    </Card>
  );
};
