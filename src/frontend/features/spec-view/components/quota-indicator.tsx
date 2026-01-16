"use client";

import { AlertTriangle, Gauge, Infinity as InfinityIcon, TrendingUp } from "lucide-react";
import { useTranslations } from "next-intl";

import { Link } from "@/i18n/navigation";
import { cn } from "@/lib/utils";

import { formatQuotaNumber, getQuotaLevel, isQuotaExceeded, type QuotaLevel } from "../utils/quota";

type QuotaIndicatorProps = {
  limit: number | null;
  percentage: number | null;
  used: number;
};

const LEVEL_CONFIG: Record<
  QuotaLevel,
  {
    bgClass: string;
    icon: typeof Gauge;
    textClass: string;
  }
> = {
  danger: {
    bgClass: "bg-destructive/10 border-destructive/20",
    icon: AlertTriangle,
    textClass: "text-destructive",
  },
  normal: {
    bgClass: "bg-muted/50 border-border",
    icon: Gauge,
    textClass: "text-muted-foreground",
  },
  unlimited: {
    bgClass: "bg-primary/5 border-primary/20",
    icon: InfinityIcon,
    textClass: "text-primary",
  },
  warning: {
    bgClass: "bg-amber-500/10 border-amber-500/20",
    icon: TrendingUp,
    textClass: "text-amber-600 dark:text-amber-500",
  },
};

export const QuotaIndicator = ({ limit, percentage, used }: QuotaIndicatorProps) => {
  const t = useTranslations("specView.quota");
  const level = getQuotaLevel(percentage);
  const config = LEVEL_CONFIG[level];
  const Icon = config.icon;

  const isAtLimit = isQuotaExceeded(percentage);

  return (
    <div
      className={cn(
        "flex items-center justify-between rounded-lg border px-3 py-2",
        config.bgClass
      )}
    >
      <div className="flex items-center gap-2">
        <Icon className={cn("h-4 w-4", config.textClass)} />
        <span className={cn("text-sm font-medium", config.textClass)}>
          {level === "unlimited" ? (
            <>
              {formatQuotaNumber(used)} {t("unlimited")}
            </>
          ) : (
            <>
              {formatQuotaNumber(used)} / {formatQuotaNumber(limit ?? 0)} {t("thisMonth")}
            </>
          )}
        </span>
      </div>

      {level === "warning" && (
        <span className={cn("text-xs", config.textClass)}>{t("warning")}</span>
      )}

      {level === "danger" && !isAtLimit && (
        <span className={cn("text-xs", config.textClass)}>{t("almostFull")}</span>
      )}

      {isAtLimit && (
        <Link
          className={cn("text-xs underline underline-offset-2", config.textClass)}
          href="/account"
        >
          {t("viewAccount")}
        </Link>
      )}
    </div>
  );
};
