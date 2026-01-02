"use client";

import { Star } from "lucide-react";
import { useTranslations } from "next-intl";

import { Toggle } from "@/components/ui/toggle";

import { useStarredFilter } from "../hooks/use-starred-filter";

export const StarredToggle = () => {
  const t = useTranslations("dashboard.filter");
  const { setStarredOnly, starredOnly } = useStarredFilter();

  const handleToggle = (pressed: boolean) => {
    setStarredOnly(pressed ? true : null);
  };

  return (
    <Toggle
      aria-label={t("starredLabel")}
      className="h-11 gap-2 px-3 sm:h-9"
      onPressedChange={handleToggle}
      pressed={starredOnly}
      variant="outline"
    >
      <Star aria-hidden="true" className={starredOnly ? "fill-yellow-400 text-yellow-400" : ""} />
      <span className="hidden sm:inline">{t("starred")}</span>
    </Toggle>
  );
};
