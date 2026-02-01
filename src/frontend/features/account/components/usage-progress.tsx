"use client";

import { useTranslations } from "next-intl";

import { cn } from "@/lib/utils";

import { formatNumber } from "../utils";

type UsageProgressProps = {
  label: string;
  limit: number | null;
  percentage: number | null;
  reserved: number;
  unit: string;
  used: number;
};

const getColorClass = (percentage: number | null): string => {
  if (percentage === null) return "bg-muted-foreground";
  if (percentage >= 90) return "bg-destructive";
  if (percentage >= 70) return "bg-amber-500";
  return "bg-primary";
};

export const UsageProgress = ({
  label,
  limit,
  percentage,
  reserved,
  unit,
  used,
}: UsageProgressProps) => {
  const t = useTranslations("specView.quota");
  const isUnlimited = limit === null;
  const displayPercentage = percentage ?? 0;
  const reservedPercentage = limit ? (reserved / limit) * 100 : 0;

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">{label}</span>
        {!isUnlimited && (
          <span
            className={cn("text-sm font-medium", displayPercentage >= 90 && "text-destructive")}
          >
            {Math.round(displayPercentage)}%
          </span>
        )}
      </div>

      {isUnlimited ? (
        <div className="flex h-2 items-center rounded-full bg-muted px-2">
          <span className="text-xs text-muted-foreground">∞</span>
        </div>
      ) : (
        <div className="relative h-2 overflow-hidden rounded-full bg-muted">
          {/* Used portion */}
          <div
            className={cn(
              "absolute left-0 top-0 h-full transition-all duration-300",
              getColorClass(percentage)
            )}
            style={{ width: `${Math.min(displayPercentage, 100)}%` }}
          />
          {/* Reserved portion */}
          {reserved > 0 && (
            <div
              className="absolute top-0 h-full bg-primary/40 transition-all duration-300"
              style={{
                left: `${Math.min(displayPercentage, 100)}%`,
                width: `${Math.min(reservedPercentage, 100 - displayPercentage)}%`,
              }}
            />
          )}
        </div>
      )}

      <div className="flex items-center justify-between text-sm text-muted-foreground">
        <span>
          {formatNumber(used)} / {isUnlimited ? "∞" : formatNumber(limit)} {unit}
        </span>
        {reserved > 0 && (
          <span className="text-xs">
            {reserved} {t("processing")}
          </span>
        )}
      </div>
    </div>
  );
};
