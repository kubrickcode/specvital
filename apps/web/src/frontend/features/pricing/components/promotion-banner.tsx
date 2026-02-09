"use client";

import { Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";

export const PromotionBanner = () => {
  const t = useTranslations("pricing.promotion");

  return (
    <div className="mb-8 rounded-lg border border-amber-500/30 bg-gradient-to-r from-amber-500/10 via-yellow-500/10 to-amber-500/10 p-4">
      <div className="flex items-center justify-center gap-2 text-center">
        <Sparkles className="size-5 shrink-0 text-amber-500" />
        <p className="text-sm font-medium">
          <span className="font-semibold text-amber-600 dark:text-amber-400">{t("title")}</span>
          <span className="mx-2 text-muted-foreground">—</span>
          <span className="text-muted-foreground">{t("description")}</span>
        </p>
      </div>
    </div>
  );
};
