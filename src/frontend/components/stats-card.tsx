"use client";

import { useTranslations } from "next-intl";
import { cn } from "@/lib/utils";

type StatsCardProps = {
  active: number;
  skipped: number;
  todo: number;
  total: number;
};

export const StatsCard = ({ active, skipped, todo, total }: StatsCardProps) => {
  const t = useTranslations("stats");

  return (
    <div className={cn("rounded-lg border bg-card p-6 shadow-xs")}>
      <h3 className="text-lg font-semibold mb-4">{t("label")}</h3>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-foreground">{total}</span>
          <span className="text-sm text-muted-foreground">{t("total")}</span>
        </div>
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-green-600">{active}</span>
          <span className="text-sm text-muted-foreground">{t("active")}</span>
        </div>
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-amber-600">{skipped}</span>
          <span className="text-sm text-muted-foreground">{t("skipped")}</span>
        </div>
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-blue-600">{todo}</span>
          <span className="text-sm text-muted-foreground">{t("todo")}</span>
        </div>
      </div>
    </div>
  );
};
