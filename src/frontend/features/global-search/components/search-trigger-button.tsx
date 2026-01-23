"use client";

import { Search } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";

import { useGlobalSearchStore } from "../hooks";

const isMac = () => {
  if (typeof window === "undefined") return false;
  return navigator.platform.toUpperCase().indexOf("MAC") >= 0;
};

export const SearchTriggerButton = () => {
  const { open } = useGlobalSearchStore();
  const t = useTranslations("globalSearch");
  const [shortcutKey, setShortcutKey] = useState("Ctrl");

  useEffect(() => {
    setShortcutKey(isMac() ? "⌘" : "Ctrl");
  }, []);

  return (
    <>
      {/* Desktop: Text button with shortcut hint */}
      <Button
        aria-keyshortcuts={isMac() ? "Meta+K" : "Control+K"}
        className="hidden h-9 w-64 justify-between gap-2 px-3 text-sm text-muted-foreground md:inline-flex"
        onClick={open}
        variant="outline"
      >
        <span className="flex items-center gap-2">
          <Search className="size-4" />
          <span>{t("placeholder")}</span>
        </span>
        <kbd className="pointer-events-none hidden h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
          {shortcutKey}K
        </kbd>
      </Button>

      {/* Mobile: Icon-only button */}
      <Button
        aria-label={t("title")}
        className="size-9 md:hidden"
        onClick={open}
        size="icon"
        variant="ghost"
      >
        <Search className="size-5" />
      </Button>
    </>
  );
};
