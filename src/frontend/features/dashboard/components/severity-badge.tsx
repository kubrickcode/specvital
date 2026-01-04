"use client";

import { AlertTriangle, Focus, SkipForward } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

export type AttentionLabelKey = "focused" | "newCommits" | "skippedHigh";
export type SeverityLevel = "critical" | "warning" | "info";

type SeverityBadgeProps = {
  labelKey: AttentionLabelKey;
  level: SeverityLevel;
};

const SEVERITY_CONFIG = {
  critical: {
    className:
      "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300 border-red-200 dark:border-red-800",
    Icon: Focus,
  },
  info: {
    className:
      "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300 border-blue-200 dark:border-blue-800",
    Icon: SkipForward,
  },
  warning: {
    className:
      "bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300 border-amber-200 dark:border-amber-800",
    Icon: AlertTriangle,
  },
} as const;

export const SeverityBadge = ({ labelKey, level }: SeverityBadgeProps) => {
  const t = useTranslations("dashboard.attention");
  const { className, Icon } = SEVERITY_CONFIG[level];

  return (
    <Badge className={cn("gap-1 text-xs font-medium", className)} variant="outline">
      <Icon aria-hidden="true" className="size-3" />
      {t(labelKey)}
    </Badge>
  );
};
