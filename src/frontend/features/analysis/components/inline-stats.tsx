"use client";

import { useTranslations } from "next-intl";

import type { Summary } from "@/lib/api/types";
import { cn } from "@/lib/utils";

type InlineStatsProps = {
  summary: Summary;
};

export const InlineStats = ({ summary }: InlineStatsProps) => {
  const t = useTranslations("stats");
  const { active, frameworks, skipped, todo, total } = summary;

  const sortedFrameworks = [...frameworks].sort((a, b) => b.total - a.total);

  return (
    <div className="rounded-lg border border-border/60 bg-card/50 p-5">
      <div className="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
        {/* Left: Total + Status */}
        <div className="flex flex-wrap items-center gap-x-6 gap-y-3">
          {/* Total - Hero metric */}
          <div className="flex items-baseline gap-2">
            <span className="text-4xl font-bold tabular-nums tracking-tight">
              {total.toLocaleString()}
            </span>
            <span className="text-sm font-medium text-muted-foreground">{t("total")}</span>
          </div>

          {/* Separator */}
          <div className="h-8 w-px bg-border/60 hidden sm:block" />

          {/* Status breakdown */}
          <div className="flex items-center gap-5">
            <StatusItem color="active" label={t("active")} value={active} />
            <StatusItem color="skipped" label={t("skipped")} value={skipped} />
            <StatusItem color="todo" label={t("todo")} value={todo} />
          </div>
        </div>

        {/* Right: Frameworks */}
        {sortedFrameworks.length > 1 && (
          <div className="flex flex-wrap items-center gap-3">
            {sortedFrameworks.map((fw) => {
              const percent = total > 0 ? Math.round((fw.total / total) * 100) : 0;
              return (
                <span
                  className="inline-flex items-center gap-2 rounded-md bg-muted/50 px-2.5 py-1.5 text-sm"
                  key={fw.framework}
                >
                  <span className="font-medium">{fw.framework}</span>
                  <span className="tabular-nums text-muted-foreground">
                    {fw.total.toLocaleString()}
                  </span>
                  <span className="text-muted-foreground/60">({percent}%)</span>
                </span>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};

type StatusItemProps = {
  color: "active" | "skipped" | "todo";
  label: string;
  value: number;
};

const StatusItem = ({ color, label, value }: StatusItemProps) => {
  const colorClasses = {
    active: "bg-status-active",
    skipped: "bg-status-skipped",
    todo: "bg-status-todo",
  };

  return (
    <span className="inline-flex items-center gap-2">
      <span className={cn("w-2.5 h-2.5 rounded-full", colorClasses[color])} />
      <span className="text-sm text-muted-foreground">{label}</span>
      <span className="text-lg font-semibold tabular-nums">{value.toLocaleString()}</span>
    </span>
  );
};
