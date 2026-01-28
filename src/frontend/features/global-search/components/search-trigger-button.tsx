"use client";

import { Search } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";

import { useGlobalSearchStore } from "../hooks";

const isMac = () => {
  if (typeof window === "undefined") return false;
  return /Mac|iPhone|iPad|iPod/.test(navigator.userAgent);
};

export const SearchTriggerButton = () => {
  const { open } = useGlobalSearchStore();
  const t = useTranslations("globalSearch");
  const [shortcutKey, setShortcutKey] = useState("Ctrl");
  const [ariaShortcut, setAriaShortcut] = useState("Control+K");

  useEffect(() => {
    const isMacOS = isMac();
    setShortcutKey(isMacOS ? "⌘" : "Ctrl");
    setAriaShortcut(isMacOS ? "Meta+K" : "Control+K");
  }, []);

  return (
    <>
      {/* Desktop: Text button with shortcut hint */}
      <Button
        aria-keyshortcuts={ariaShortcut}
        className="hidden w-64 justify-between gap-2 px-3 text-muted-foreground md:inline-flex"
        onClick={open}
        size="sm"
        variant="header-action"
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
      <ResponsiveTooltip content={t("title")} side="bottom" sideOffset={8}>
        <Button
          aria-label={t("title")}
          className="md:hidden"
          onClick={open}
          size="header-icon"
          variant="header-action"
        >
          <Search className="size-4" />
        </Button>
      </ResponsiveTooltip>
    </>
  );
};
