"use client";

import { GitFork } from "lucide-react";
import { useTranslations } from "next-intl";

import { Card } from "@/components/ui/card";
import { AnalyzeDialog } from "@/features/home";

export const EmptyState = () => {
  const t = useTranslations("dashboard.emptyState");

  return (
    <Card className="flex flex-col items-center justify-center py-16 px-6 text-center">
      <div className="rounded-full bg-muted p-4 mb-6">
        <GitFork aria-hidden="true" className="size-8 text-muted-foreground" />
      </div>

      <h2 className="text-xl font-semibold mb-2">{t("title")}</h2>
      <p className="text-muted-foreground mb-6 max-w-md">{t("description")}</p>

      <AnalyzeDialog variant="empty-state" />
    </Card>
  );
};
