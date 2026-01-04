"use client";

import { motion } from "motion/react";
import { useTranslations } from "next-intl";

import { Card, CardContent } from "@/components/ui/card";
import type { Summary } from "@/lib/api/types";
import { easeOutTransition, useReducedMotion } from "@/lib/motion";
import { getFrameworkSolidColor } from "@/lib/styles";

type StatsCardProps = {
  summary: Summary;
};

export const StatsCard = ({ summary }: StatsCardProps) => {
  const t = useTranslations("stats");
  const { active, frameworks, skipped, todo, total } = summary;
  const shouldReduceMotion = useReducedMotion();

  const showFrameworkBreakdown = frameworks.length > 1;

  return (
    <Card depth="raised">
      <CardContent className="space-y-6">
        <h3 className="text-lg font-semibold">{t("label")}</h3>
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
          <div className="flex flex-col">
            <span className="text-5xl font-extrabold text-primary">{total}</span>
            <span className="text-sm text-muted-foreground">{t("total")}</span>
          </div>
          <div className="flex flex-col border-l-4 border-status-active pl-3">
            <span className="text-3xl font-bold text-status-active">{active}</span>
            <span className="text-sm text-muted-foreground">{t("active")}</span>
          </div>
          <div className="flex flex-col border-l-4 border-status-skipped pl-3">
            <span className="text-3xl font-bold text-status-skipped">{skipped}</span>
            <span className="text-sm text-muted-foreground">{t("skipped")}</span>
          </div>
          <div className="flex flex-col border-l-4 border-status-todo pl-3">
            <span className="text-3xl font-bold text-status-todo">{todo}</span>
            <span className="text-sm text-muted-foreground">{t("todo")}</span>
          </div>
        </div>

        {showFrameworkBreakdown ? (
          <div className="pt-6 border-t border-border/50">
            <h4 className="text-sm font-medium text-muted-foreground mb-4">{t("byFramework")}</h4>
            <div className="space-y-4">
              {frameworks.map((fw, index) => {
                const percentage = total > 0 ? Math.round((fw.total / total) * 100) : 0;
                return (
                  <div className="space-y-2" key={fw.framework}>
                    <div className="flex items-center justify-between">
                      <span className="flex items-center gap-2 font-medium">
                        <span
                          aria-hidden="true"
                          className="h-3 w-3 rounded-full"
                          style={{ backgroundColor: getFrameworkSolidColor(fw.framework) }}
                        />
                        {fw.framework}
                        <span className="text-muted-foreground font-normal">({percentage}%)</span>
                      </span>
                      <span className="text-sm text-muted-foreground">
                        {fw.total} {t("tests")}
                      </span>
                    </div>
                    <div
                      aria-label={t("frameworkProgress", { framework: fw.framework })}
                      aria-valuemax={100}
                      aria-valuemin={0}
                      aria-valuenow={percentage}
                      aria-valuetext={t("frameworkProgressDetail", {
                        active: fw.active,
                        framework: fw.framework,
                        percentage,
                        skipped: fw.skipped,
                        todo: fw.todo,
                        total: fw.total,
                      })}
                      className="h-2 w-full rounded-full bg-muted overflow-hidden"
                      role="progressbar"
                    >
                      <motion.div
                        animate={{ width: `${percentage}%` }}
                        className="h-full rounded-full"
                        initial={shouldReduceMotion ? false : { width: 0 }}
                        style={{ backgroundColor: getFrameworkSolidColor(fw.framework) }}
                        transition={
                          shouldReduceMotion
                            ? undefined
                            : { ...easeOutTransition, delay: 0.2 + index * 0.1 }
                        }
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
          <div className="pt-4 border-t border-border/50">
            <span className="text-sm text-muted-foreground">
              {t("framework")}:{" "}
              <span className="font-medium text-foreground">{frameworks[0].framework}</span>
            </span>
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
};
