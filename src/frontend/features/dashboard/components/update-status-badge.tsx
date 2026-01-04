"use client";

import { CheckCircle2, CircleDashed, GitCommit } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import type { UpdateStatus } from "@/lib/api/types";
import { cn } from "@/lib/utils";

type UpdateStatusBadgeProps = {
  status: UpdateStatus;
};

const STATUS_CONFIG = {
  "new-commits": {
    className: "bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300",
    Icon: GitCommit,
    labelKey: "newCommits",
  },
  unknown: {
    className: "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-300",
    Icon: CircleDashed,
    labelKey: "unknown",
  },
  "up-to-date": {
    className: "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300",
    Icon: CheckCircle2,
    labelKey: "upToDate",
  },
} as const;

export const UpdateStatusBadge = ({ status }: UpdateStatusBadgeProps) => {
  const t = useTranslations("dashboard.status");
  const config = STATUS_CONFIG[status];
  const { className, Icon, labelKey } = config;

  return (
    <Badge className={cn("gap-1", className)} variant="outline">
      <Icon aria-hidden="true" className="size-3" />
      <span className="text-xs">{t(labelKey)}</span>
    </Badge>
  );
};
