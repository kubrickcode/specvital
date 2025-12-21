"use client";

import { useTranslations } from "next-intl";

import type { Summary } from "@/lib/api/types";
import { getFrameworkColor } from "@/lib/styles";
import { cn } from "@/lib/utils";

type StatsCardProps = {
  summary: Summary;
};

export const StatsCard = ({ summary }: StatsCardProps) => {
  const t = useTranslations("stats");
  const { active, frameworks, skipped, todo, total } = summary;

  const showFrameworkBreakdown = frameworks.length > 1;

  return (
    <div className={cn("rounded-lg border bg-card p-6 shadow-xs")}>
      <h3 className="text-lg font-semibold mb-4">{t("label")}</h3>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        <div className="flex flex-col">
          <span className="text-4xl font-extrabold text-foreground">{total}</span>
          <span className="text-sm text-muted-foreground">{t("total")}</span>
        </div>
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-status-active">{active}</span>
          <span className="text-sm text-muted-foreground">{t("active")}</span>
        </div>
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-status-skipped">{skipped}</span>
          <span className="text-sm text-muted-foreground">{t("skipped")}</span>
        </div>
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-status-todo">{todo}</span>
          <span className="text-sm text-muted-foreground">{t("todo")}</span>
        </div>
      </div>

      {showFrameworkBreakdown ? (
        <div className="mt-6 pt-6 border-t">
          <h4 className="text-sm font-medium text-muted-foreground mb-4">{t("byFramework")}</h4>
          <div className="space-y-4">
            {frameworks.map((fw, index) => {
              const percentage = total > 0 ? Math.round((fw.total / total) * 100) : 0;
              return (
                <div className="space-y-2" key={fw.framework}>
                  <div className="flex items-center justify-between">
                    <span className="flex items-center gap-2 font-medium">
                      <span
                        className="h-3 w-3 rounded-full"
                        style={{ backgroundColor: getFrameworkColor(index, fw.framework) }}
                      />
                      {fw.framework}
                      <span className="text-muted-foreground font-normal">({percentage}%)</span>
                    </span>
                    <span className="text-sm text-muted-foreground">
                      {fw.total} {t("tests")}
                    </span>
                  </div>
                  <div
                    aria-valuemax={100}
                    aria-valuemin={0}
                    aria-valuenow={percentage}
                    className="h-2 w-full rounded-full bg-muted overflow-hidden"
                    role="progressbar"
                  >
                    <div
                      className="h-full rounded-full transition-all duration-300"
                      style={{
                        backgroundColor: getFrameworkColor(index, fw.framework),
                        width: `${percentage}%`,
                      }}
                    />
                  </div>
                  <div className="text-sm text-muted-foreground">
                    <span className="text-status-active">{fw.active}</span> {t("active")}
                    {" · "}
                    <span className="text-status-skipped">{fw.skipped}</span> {t("skipped")}
                    {" · "}
                    <span className="text-status-todo">{fw.todo}</span> {t("todo")}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      ) : frameworks.length === 1 ? (
        <div className="mt-4 pt-4 border-t">
          <span className="text-sm text-muted-foreground">
            {t("framework")}:{" "}
            <span className="font-medium text-foreground">{frameworks[0].framework}</span>
          </span>
        </div>
      ) : null}
    </div>
  );
};
