"use client";

import { Minus, TrendingDown, TrendingUp } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

type TestDeltaBadgeProps = {
  delta: number;
  isCompact?: boolean;
};

const getDeltaInfo = (delta: number) => {
  if (delta > 0) {
    return {
      className: "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400",
      Icon: TrendingUp,
      labelKey: "increased",
      prefix: "+",
    } as const;
  }

  if (delta < 0) {
    return {
      className: "bg-rose-100 text-rose-800 dark:bg-rose-900/30 dark:text-rose-400",
      Icon: TrendingDown,
      labelKey: "decreased",
      prefix: "",
    } as const;
  }

  return {
    className: "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400",
    Icon: Minus,
    labelKey: "unchanged",
    prefix: "",
  } as const;
};

export const TestDeltaBadge = ({ delta, isCompact = false }: TestDeltaBadgeProps) => {
  const t = useTranslations("dashboard.delta");
  const { className, Icon, labelKey, prefix } = getDeltaInfo(delta);

  return (
    <Badge className={cn("gap-1", className)} variant="outline">
      <Icon aria-hidden="true" className={cn("size-3", isCompact && "size-2.5")} />
      <span className={cn(isCompact && "text-[10px]")}>
        {prefix}
        {delta}
      </span>
      <span className="sr-only">{t(labelKey, { count: Math.abs(delta) })}</span>
    </Badge>
  );
};
