"use client";

import { useQuery } from "@tanstack/react-query";
import {
  AlertTriangle,
  Gauge,
  Globe,
  Infinity as InfinityIcon,
  RefreshCw,
  Sparkles,
  Zap,
} from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Link } from "@/i18n/navigation";
import { cn } from "@/lib/utils";

import { fetchCacheAvailability } from "../api";
import { LanguageCombobox } from "./language-combobox";
import { useQuotaConfirmDialog } from "../hooks/use-quota-confirm-dialog";
import { formatQuotaNumber, getQuotaLevel, isQuotaExceeded, type QuotaLevel } from "../utils/quota";

const LEVEL_CONFIG: Record<
  QuotaLevel,
  {
    bgClass: string;
    iconColor: string;
    progressColor: string;
  }
> = {
  danger: {
    bgClass: "bg-destructive/10",
    iconColor: "text-destructive",
    progressColor: "bg-destructive",
  },
  normal: {
    bgClass: "bg-muted",
    iconColor: "text-muted-foreground",
    progressColor: "bg-primary",
  },
  unlimited: {
    bgClass: "bg-primary/10",
    iconColor: "text-primary",
    progressColor: "bg-primary",
  },
  warning: {
    bgClass: "bg-amber-500/10",
    iconColor: "text-amber-500",
    progressColor: "bg-amber-500",
  },
};

const getIcon = (level: QuotaLevel, isUnlimited: boolean) => {
  if (isUnlimited) return InfinityIcon;
  if (level === "danger" || level === "warning") return AlertTriangle;
  return Gauge;
};

export const QuotaConfirmDialog = () => {
  const t = useTranslations("specView.quotaConfirm");
  const tGenerate = useTranslations("specView.generate");
  const {
    analysisId,
    close,
    confirm,
    estimatedCost,
    forceRegenerate,
    isOpen,
    isRegenerate,
    isSameCommit,
    onOpenChange,
    regeneratingLanguage,
    selectedLanguage,
    setForceRegenerate,
    setSelectedLanguage,
    usage,
  } = useQuotaConfirmDialog();

  // Fetch cache availability when dialog is open
  const { data: cacheAvailability, isError: isCacheAvailabilityError } = useQuery({
    enabled: isOpen && !!analysisId,
    queryFn: () => fetchCacheAvailability(analysisId!),
    queryKey: ["cache-availability", analysisId],
    staleTime: 60000, // 1 minute
  });

  // Check if previous spec exists for selected language.
  // Cache availability API excludes the current analysis, so for regeneration
  // we supplement with the known language from the current analysis.
  const hasCacheFromCurrentAnalysis = regeneratingLanguage === selectedLanguage;
  const hasCacheFromPreviousAnalysis =
    !isCacheAvailabilityError && (cacheAvailability?.languages?.[selectedLanguage] ?? false);
  const hasPreviousSpec = hasCacheFromCurrentAnalysis || hasCacheFromPreviousAnalysis;

  const specview = usage?.specview;
  const percentage = specview?.percentage ?? null;
  const level = getQuotaLevel(percentage);
  const isExceeded = isQuotaExceeded(percentage);
  const isUnlimited = specview?.limit === null || specview?.limit === undefined;
  const config = LEVEL_CONFIG[level];

  // Calculate if generation would exceed limit
  const afterUsage = specview ? specview.used + (estimatedCost ?? 0) : 0;
  const wouldExceed =
    !isUnlimited && specview?.limit && estimatedCost ? afterUsage > specview.limit : false;

  const Icon = getIcon(level, isUnlimited);

  return (
    <Dialog onOpenChange={onOpenChange} open={isOpen}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader className="text-center sm:text-center">
          <div
            className={cn(
              "mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full",
              config.bgClass
            )}
          >
            <Icon className={cn("h-6 w-6", config.iconColor)} />
          </div>
          <DialogTitle>{isRegenerate ? t("regenerateTitle") : t("title")}</DialogTitle>
          <DialogDescription>
            {isRegenerate ? t("regenerateDescription") : t("description")}
          </DialogDescription>
        </DialogHeader>

        <div className="mt-4 space-y-4">
          {/* Language Selection */}
          <div className="space-y-2">
            <Label
              className="flex items-center gap-2 text-sm font-medium"
              htmlFor="language-select"
            >
              <Globe className="h-4 w-4 text-muted-foreground" />
              {t("outputLanguage")}
            </Label>
            <LanguageCombobox onValueChange={setSelectedLanguage} value={selectedLanguage} />
          </div>

          {/* Cache availability check error warning */}
          {isCacheAvailabilityError && (
            <p className="text-xs text-muted-foreground">{tGenerate("cacheCheckFailed")}</p>
          )}

          {/* Analysis Mode Selection - only show when cache is available */}
          {hasPreviousSpec && (
            <div className="space-y-2">
              <Label className="flex items-center gap-2 text-sm font-medium">
                <Zap className="h-4 w-4 text-muted-foreground" />
                {tGenerate("analysisMode")}
              </Label>

              <RadioGroup
                className="grid grid-cols-1 gap-2"
                onValueChange={(value) => setForceRegenerate(value === "fresh")}
                value={forceRegenerate ? "fresh" : "cache"}
              >
                <label
                  className={cn(
                    "flex items-start gap-3 rounded-lg border p-3 transition-colors",
                    isSameCommit ? "cursor-not-allowed opacity-50 border-border" : "cursor-pointer",
                    !forceRegenerate && !isSameCommit
                      ? "border-primary bg-primary/5"
                      : "border-border hover:border-muted-foreground/50"
                  )}
                  htmlFor="cache-mode"
                >
                  <RadioGroupItem
                    className="mt-0.5"
                    disabled={isSameCommit}
                    id="cache-mode"
                    value="cache"
                  />
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium">{tGenerate("withCache")}</span>
                      {!isSameCommit && (
                        <span className="rounded-full bg-primary/10 px-2 py-0.5 text-xs font-medium text-primary">
                          {tGenerate("recommended")}
                        </span>
                      )}
                    </div>
                    <p className="mt-0.5 text-xs text-muted-foreground">
                      {isSameCommit
                        ? tGenerate("sameCommitCacheDisabled")
                        : tGenerate("withCacheBenefit")}
                    </p>
                  </div>
                </label>
                <label
                  className={cn(
                    "flex cursor-pointer items-start gap-3 rounded-lg border p-3 transition-colors",
                    forceRegenerate
                      ? "border-primary bg-primary/5"
                      : "border-border hover:border-muted-foreground/50"
                  )}
                  htmlFor="fresh-mode"
                >
                  <RadioGroupItem className="mt-0.5" id="fresh-mode" value="fresh" />
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <RefreshCw className="h-3 w-3 text-muted-foreground" />
                      <span className="text-sm font-medium">{tGenerate("fresh")}</span>
                    </div>
                    <p className="mt-0.5 text-xs text-muted-foreground">
                      {tGenerate("freshWarning")}
                    </p>
                  </div>
                </label>
              </RadioGroup>
            </div>
          )}
        </div>

        {specview && (
          <div className="mt-4 space-y-4">
            <div className="rounded-lg border bg-muted/30 p-4">
              <div className="mb-2 flex items-center justify-between">
                <span className="text-sm font-medium">{t("specviewUsage")}</span>
                {!isUnlimited && (
                  <span
                    className={cn(
                      "text-sm font-semibold",
                      level === "danger" && "text-destructive",
                      level === "warning" && "text-amber-600 dark:text-amber-500"
                    )}
                  >
                    {Math.round(percentage ?? 0)}%
                  </span>
                )}
              </div>

              {isUnlimited ? (
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <InfinityIcon className="h-4 w-4" />
                  <span>
                    {formatQuotaNumber(specview.used)} {t("unit")} {t("used")} · {t("unlimited")}
                  </span>
                </div>
              ) : (
                <>
                  <div className="mb-2 h-2 overflow-hidden rounded-full bg-muted">
                    <div
                      className={cn("h-full transition-all", config.progressColor)}
                      style={{ width: `${Math.min(percentage ?? 0, 100)}%` }}
                    />
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {formatQuotaNumber(specview.used)} / {formatQuotaNumber(specview.limit ?? 0)}{" "}
                    {t("thisMonth")}
                  </p>
                </>
              )}
            </div>

            {estimatedCost !== null && (
              <div
                className={cn(
                  "rounded-lg border p-3",
                  wouldExceed
                    ? "border-destructive/20 bg-destructive/10"
                    : "border-primary/20 bg-primary/5"
                )}
              >
                <p
                  className={cn(
                    "text-sm font-medium",
                    wouldExceed ? "text-destructive" : "text-primary"
                  )}
                >
                  {t("estimatedCost", { count: formatQuotaNumber(estimatedCost) })}
                </p>
                {!isUnlimited && (
                  <p className="mt-1 text-xs text-muted-foreground">
                    {t("afterGeneration", {
                      after: formatQuotaNumber(afterUsage),
                      limit: formatQuotaNumber(specview?.limit ?? 0),
                    })}
                  </p>
                )}
              </div>
            )}

            {isRegenerate && (
              <div className="rounded-lg border border-amber-500/20 bg-amber-500/10 p-3">
                <p className="text-sm text-amber-600 dark:text-amber-500">
                  {t("regenerateWarning")}
                </p>
              </div>
            )}

            {level === "warning" && (
              <div className="rounded-lg border border-amber-500/20 bg-amber-500/10 p-3">
                <p className="text-sm text-amber-600 dark:text-amber-500">{t("warningMessage")}</p>
              </div>
            )}

            {level === "danger" && !isExceeded && (
              <div className="rounded-lg border border-destructive/20 bg-destructive/10 p-3">
                <p className="text-sm text-destructive">{t("dangerMessage")}</p>
              </div>
            )}

            {isExceeded && (
              <div className="rounded-lg border border-destructive/20 bg-destructive/10 p-3">
                <p className="text-sm text-destructive">{t("exceededMessage")}</p>
                <Link
                  className="mt-2 inline-block text-sm font-medium text-destructive underline underline-offset-2"
                  href="/account"
                  onClick={close}
                >
                  {t("viewAccount")}
                </Link>
              </div>
            )}
          </div>
        )}

        <DialogFooter className="mt-4 flex-col gap-2 sm:flex-row">
          <Button className="w-full sm:w-auto" onClick={close} variant="outline">
            {t("cancel")}
          </Button>
          <Button
            className="w-full sm:w-auto"
            disabled={isExceeded || wouldExceed}
            onClick={confirm}
            variant="default"
          >
            <Sparkles className="mr-2 h-4 w-4" />
            {isRegenerate ? t("regenerate") : t("generate")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
