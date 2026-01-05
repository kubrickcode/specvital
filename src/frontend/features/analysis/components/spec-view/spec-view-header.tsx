"use client";

import { useTranslations } from "next-intl";

import { LanguageSelector } from "./language-selector";
import type { ConversionLanguage, ConversionSummary } from "../../types";

type SpecViewHeaderProps = {
  language: ConversionLanguage;
  onLanguageChange: (language: ConversionLanguage) => void;
  summary: ConversionSummary | null;
};

export const SpecViewHeader = ({ language, onLanguageChange, summary }: SpecViewHeaderProps) => {
  const t = useTranslations("analyze.specView");

  return (
    <div className="flex items-center justify-between mb-4">
      <div className="flex items-center gap-4">
        {summary && (
          <p className="text-sm text-muted-foreground">
            {t("summary", {
              cached: summary.cachedCount,
              converted: summary.convertedCount,
              total: summary.totalTests,
            })}
          </p>
        )}
      </div>
      <LanguageSelector onChange={onLanguageChange} value={language} />
    </div>
  );
};
