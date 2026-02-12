"use client";

import { useTranslations } from "next-intl";

export const SearchKeyboardHints = () => {
  const t = useTranslations("globalSearch");

  return (
    <div className="hidden border-t px-3 py-2 text-xs text-muted-foreground md:flex md:items-center md:justify-between">
      <div className="flex items-center gap-4">
        <span className="flex items-center gap-1">
          <kbd className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px]">↑</kbd>
          <kbd className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px]">↓</kbd>
          <span className="ml-1">{t("hints.navigate")}</span>
        </span>
        <span className="flex items-center gap-1">
          <kbd className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px]">↵</kbd>
          <span className="ml-1">{t("hints.select")}</span>
        </span>
      </div>
      <span className="flex items-center gap-1">
        <kbd className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px]">Esc</kbd>
        <span className="ml-1">{t("hints.close")}</span>
      </span>
    </div>
  );
};
