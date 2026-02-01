"use client";

import { useQuery } from "@tanstack/react-query";
import { Globe, Infinity as InfinityIcon, RefreshCw, Sparkles, Zap } from "lucide-react";
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

const LEVEL_CONFIG: Record<QuotaLevel, { progressColor: string }> = {
  danger: { progressColor: "bg-destructive" },
  normal: { progressColor: "bg-primary" },
  unlimited: { progressColor: "bg-primary" },
  warning: { progressColor: "bg-amber-500" },
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

  const tQuota = useTranslations("specView.quota");
  const specview = usage?.specview;
  const percentage = specview?.percentage ?? null;
  const reserved = specview?.reserved ?? 0;
  const level = getQuotaLevel(percentage);
  const isExceeded = isQuotaExceeded(percentage);
  const isUnlimited = specview?.limit === null || specview?.limit === undefined;
  const config = LEVEL_CONFIG[level];

  // Calculate if generation would exceed limit (including reserved)
  const afterUsage = specview ? specview.used + reserved + (estimatedCost ?? 0) : 0;
  const wouldExceed =
    !isUnlimited && specview?.limit && estimatedCost ? afterUsage > specview.limit : false;
  const reservedPercentage = specview?.limit ? (reserved / specview.limit) * 100 : 0;

  return (
    <Dialog onOpenChange={onOpenChange} open={isOpen}>
      <DialogContent className="overflow-hidden sm:max-w-md">
        {/* AI indicator gradient line */}
        <div className="absolute left-0 right-0 top-0 h-0.5 bg-gradient-to-r from-transparent via-violet-500/50 to-transparent" />

        <DialogHeader className="text-center sm:text-center">
          <DialogTitle className="flex items-center justify-center gap-2">
            {isRegenerate ? t("regenerateTitle") : t("title")}
            <span className="inline-flex items-center gap-1 rounded-full border border-violet-500/20 bg-violet-500/10 px-1.5 py-0.5 text-[10px] font-medium text-violet-600 dark:bg-violet-400/10 dark:text-violet-400">
              <Sparkles className="size-3" />
              AI
            </span>
          </DialogTitle>
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
              <div className="mb-3 flex items-center justify-between">
                <span className="text-sm font-medium">{t("specviewUsage")}</span>
                {!isUnlimited && estimatedCost !== null && (
                  <span
                    className={cn(
                      "text-xs font-medium",
                      wouldExceed ? "text-destructive" : "text-muted-foreground"
                    )}
                  >
                    {t("afterUsageShort", {
                      after: formatQuotaNumber(afterUsage),
                      limit: formatQuotaNumber(specview.limit ?? 0),
                    })}
                  </span>
                )}
              </div>

              {isUnlimited ? (
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <InfinityIcon className="h-4 w-4" />
                  <span>
                    {formatQuotaNumber(specview.used)} {t("unit")} {t("used")} · {t("unlimited")}
                    {reserved > 0 && (
                      <span className="ml-1 opacity-70">
                        · {reserved} {tQuota("processing")}
                      </span>
                    )}
                  </span>
                </div>
              ) : (
                <>
                  {/* Unified progress bar with prediction overlay */}
                  <div className="mb-3 h-2.5 overflow-hidden rounded-full bg-muted">
                    <div className="relative h-full w-full">
                      {/* Current usage */}
                      <div
                        className={cn(
                          "absolute left-0 top-0 h-full transition-all",
                          config.progressColor
                        )}
                        style={{ width: `${Math.min(percentage ?? 0, 100)}%` }}
                      />
                      {/* Reserved usage */}
                      {reserved > 0 && (
                        <div
                          className="absolute top-0 h-full bg-primary/40 transition-all"
                          style={{
                            left: `${Math.min(percentage ?? 0, 100)}%`,
                            width: `${Math.min(reservedPercentage, 100 - (percentage ?? 0))}%`,
                          }}
                        />
                      )}
                      {/* Predicted usage (striped pattern) */}
                      {estimatedCost !== null && estimatedCost > 0 && (
                        <div
                          className={cn(
                            "absolute top-0 h-full transition-all",
                            wouldExceed ? "bg-destructive/60" : "bg-primary/40"
                          )}
                          style={{
                            backgroundImage: wouldExceed
                              ? "repeating-linear-gradient(45deg, transparent, transparent 2px, rgba(255,255,255,0.3) 2px, rgba(255,255,255,0.3) 4px)"
                              : "repeating-linear-gradient(45deg, transparent, transparent 2px, rgba(255,255,255,0.4) 2px, rgba(255,255,255,0.4) 4px)",
                            left: `${Math.min((percentage ?? 0) + reservedPercentage, 100)}%`,
                            width: `${Math.min((estimatedCost / (specview.limit ?? 1)) * 100, 100 - (percentage ?? 0) - reservedPercentage)}%`,
                          }}
                        />
                      )}
                    </div>
                  </div>

                  {/* Labels under progress bar */}
                  <div className="flex items-center justify-between text-xs">
                    <div className="flex items-center gap-3">
                      <span className="flex items-center gap-1.5">
                        <span
                          className={cn("inline-block h-2 w-2 rounded-full", config.progressColor)}
                        />
                        <span className="text-muted-foreground">
                          {t("current")}: {formatQuotaNumber(specview.used)}
                        </span>
                      </span>
                      {reserved > 0 && (
                        <span className="flex items-center gap-1.5">
                          <span className="inline-block h-2 w-2 rounded-full bg-primary/40" />
                          <span className="text-muted-foreground">
                            {reserved} {tQuota("processing")}
                          </span>
                        </span>
                      )}
                      {estimatedCost !== null && estimatedCost > 0 && (
                        <span className="flex items-center gap-1.5">
                          <span
                            className={cn(
                              "inline-block h-2 w-2 rounded-full",
                              wouldExceed ? "bg-destructive/60" : "bg-primary/40"
                            )}
                            style={{
                              backgroundImage:
                                "repeating-linear-gradient(45deg, transparent, transparent 1px, rgba(255,255,255,0.5) 1px, rgba(255,255,255,0.5) 2px)",
                            }}
                          />
                          <span
                            className={cn(
                              wouldExceed ? "text-destructive" : "text-muted-foreground"
                            )}
                          >
                            ~+{formatQuotaNumber(estimatedCost)}
                          </span>
                        </span>
                      )}
                    </div>
                    <span className="text-muted-foreground">
                      {t("limit")}: {formatQuotaNumber(specview.limit ?? 0)}
                    </span>
                  </div>
                </>
              )}
            </div>

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
            {isRegenerate ? t("regenerate") : t("generate")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
