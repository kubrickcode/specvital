"use client";

import { cn } from "@/lib/utils";

import { formatNumber } from "../utils";

type UsageProgressProps = {
  label: string;
  limit: number | null;
  percentage: number | null;
  unit: string;
  used: number;
};

const getColorClass = (percentage: number | null): string => {
  if (percentage === null) return "bg-muted-foreground";
  if (percentage >= 90) return "bg-destructive";
  if (percentage >= 70) return "bg-amber-500";
  return "bg-primary";
};

export const UsageProgress = ({ label, limit, percentage, unit, used }: UsageProgressProps) => {
  const isUnlimited = limit === null;
  const displayPercentage = percentage ?? 0;

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
        <div className="h-2 overflow-hidden rounded-full bg-muted">
          <div
            className={cn("h-full transition-all duration-300", getColorClass(percentage))}
            style={{ width: `${Math.min(displayPercentage, 100)}%` }}
          />
        </div>
      )}

      <p className="text-sm text-muted-foreground">
        {formatNumber(used)} / {isUnlimited ? "∞" : formatNumber(limit)} {unit}
      </p>
    </div>
  );
};
