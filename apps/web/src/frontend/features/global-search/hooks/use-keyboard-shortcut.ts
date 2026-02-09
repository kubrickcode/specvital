"use client";

import { useEffect } from "react";

import { globalSearchStore } from "./use-global-search-store";

export const useKeyboardShortcut = () => {
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      // Skip if user is typing in an input field
      const target = event.target as HTMLElement;
      const isTyping =
        target.tagName === "INPUT" || target.tagName === "TEXTAREA" || target.isContentEditable;

      if (isTyping) {
        return;
      }

      // Use userAgentData with fallback to userAgent (navigator.platform is deprecated)
      const isMac =
        // @ts-expect-error -- userAgentData not in all TS definitions
        navigator.userAgentData?.platform?.toLowerCase().includes("mac") ??
        /mac/i.test(navigator.userAgent);

      const isModifierPressed = isMac ? event.metaKey : event.ctrlKey;
      const isKKey = event.key.toLowerCase() === "k";

      if (isModifierPressed && isKKey) {
        event.preventDefault();
        globalSearchStore.toggle();
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, []);
};
