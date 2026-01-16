"use client";

import { AlertTriangle, Gauge, Infinity as InfinityIcon, Sparkles } from "lucide-react";
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
import { Link } from "@/i18n/navigation";
import { cn } from "@/lib/utils";

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
  const { close, confirm, estimatedCost, isOpen, onOpenChange, usage } = useQuotaConfirmDialog();

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
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>

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
            {t("generate")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
