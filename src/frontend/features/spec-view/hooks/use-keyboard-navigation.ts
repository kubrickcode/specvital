"use client";

import { useCallback, useEffect, useRef } from "react";

type NavigationItem = {
  element: HTMLElement;
  id: string;
  type: "behavior" | "domain" | "feature";
};

type UseKeyboardNavigationOptions = {
  enabled?: boolean;
  onNavigate?: (item: NavigationItem) => void;
};

const FOCUSABLE_SELECTOR = '[id^="domain-"], [id^="feature-"], [id^="behavior-"]';

/**
 * Keyboard navigation for document view
 * - j/k or ArrowDown/ArrowUp: Navigate between sections
 * - Enter/Space: Expand/collapse current section
 * - Escape: Reset focus to document container
 */
export const useKeyboardNavigation = ({
  enabled = true,
  onNavigate,
}: UseKeyboardNavigationOptions = {}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const currentIndexRef = useRef(-1);

  const getNavigableElements = useCallback((): NavigationItem[] => {
    if (!containerRef.current) return [];

    const elements = containerRef.current.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR);

    return Array.from(elements).map((element) => {
      const id = element.id;
      let type: NavigationItem["type"] = "domain";
      if (id.startsWith("feature-")) type = "feature";
      if (id.startsWith("behavior-")) type = "behavior";
      return { element, id, type };
    });
  }, []);

  const navigateTo = useCallback(
    (index: number) => {
      const items = getNavigableElements();
      if (items.length === 0) return;

      const clampedIndex = Math.max(0, Math.min(index, items.length - 1));
      const item = items[clampedIndex];

      if (item) {
        currentIndexRef.current = clampedIndex;
        item.element.scrollIntoView({ behavior: "smooth", block: "center" });
        item.element.focus({ preventScroll: true });
        onNavigate?.(item);
      }
    },
    [getNavigableElements, onNavigate]
  );

  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      if (!enabled) return;

      // Ignore if typing in input
      const target = event.target as HTMLElement;
      if (target.tagName === "INPUT" || target.tagName === "TEXTAREA" || target.isContentEditable) {
        return;
      }

      const items = getNavigableElements();
      if (items.length === 0) return;

      switch (event.key) {
        case "j":
        case "ArrowDown": {
          event.preventDefault();
          navigateTo(currentIndexRef.current + 1);
          break;
        }
        case "k":
        case "ArrowUp": {
          event.preventDefault();
          navigateTo(currentIndexRef.current - 1);
          break;
        }
        case "Home": {
          event.preventDefault();
          navigateTo(0);
          break;
        }
        case "End": {
          event.preventDefault();
          navigateTo(items.length - 1);
          break;
        }
        case "Escape": {
          event.preventDefault();
          currentIndexRef.current = -1;
          containerRef.current?.focus();
          break;
        }
      }
    },
    [enabled, getNavigableElements, navigateTo]
  );

  useEffect(() => {
    if (!enabled) return;

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [enabled, handleKeyDown]);

  return {
    containerProps: {
      "aria-label": "Specification document",
      ref: containerRef,
      role: "region" as const,
      tabIndex: -1,
    },
    navigateTo,
  };
};
