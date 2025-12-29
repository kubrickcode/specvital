"use client";

import { SearchX } from "lucide-react";
import { useTranslations } from "next-intl";

export const FilterEmptyState = () => {
  const t = useTranslations("analyze.filter");

  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      <SearchX className="h-12 w-12 text-muted-foreground mb-4" />
      <h3 className="text-lg font-semibold">{t("noResults")}</h3>
      <p className="text-sm text-muted-foreground mt-1">{t("noResultsDescription")}</p>
    </div>
  );
};
