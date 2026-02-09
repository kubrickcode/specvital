"use client";

import { useTranslations } from "next-intl";
import { useTheme } from "next-themes";
import { useState } from "react";

import { useAuth } from "@/features/auth";
import { useRouter } from "@/i18n/navigation";

import { COMMAND_ACTIONS, NAVIGATION_ACTIONS } from "../lib/static-actions";
import type { SearchableItem, StaticAction } from "../types";
import { useGlobalSearchStore } from "./use-global-search-store";

export const useStaticActions = () => {
  const t = useTranslations("globalSearch");
  const router = useRouter();
  const { close } = useGlobalSearchStore();
  const { isAuthenticated } = useAuth();
  const { resolvedTheme, setTheme } = useTheme();
  const [analyzeDialogOpen, setAnalyzeDialogOpen] = useState(false);

  const executeAction = (action: StaticAction) => {
    const keepOpen =
      action.id === "action-toggle-theme-light" || action.id === "action-toggle-theme-dark";

    if (!keepOpen) {
      close();
    }

    switch (action.id) {
      case "nav-dashboard":
        router.push("/dashboard");
        break;
      case "nav-explore":
        router.push("/explore");
        break;
      case "nav-account":
        router.push("/account");
        break;
      case "nav-pricing":
        router.push("/pricing");
        break;
      case "nav-home":
        router.push("/");
        break;
      case "action-new-analysis":
        setAnalyzeDialogOpen(true);
        break;
      case "action-toggle-theme-light":
        setTheme("light");
        break;
      case "action-toggle-theme-dark":
        setTheme("dark");
        break;
    }
  };

  const getFilteredNavigationActions = (): StaticAction[] =>
    NAVIGATION_ACTIONS.filter((action) => !action.requiresAuth || isAuthenticated);

  const getFilteredCommandActions = (): StaticAction[] =>
    COMMAND_ACTIONS.filter((action) => {
      if (action.id === "action-toggle-theme-light") return resolvedTheme !== "light";
      if (action.id === "action-toggle-theme-dark") return resolvedTheme !== "dark";
      return true;
    });

  const mapActionToSearchableItem = (action: StaticAction): SearchableItem => ({
    category: "actions",
    icon: action.icon,
    id: action.id,
    keywords: action.keywords,
    label: t(action.labelKey),
    labelKey: action.labelKey,
    onSelect: () => executeAction(action),
    shortcut: action.shortcut,
  });

  const navigationItems = getFilteredNavigationActions().map(mapActionToSearchableItem);
  const commandItems = getFilteredCommandActions().map(mapActionToSearchableItem);

  return {
    analyzeDialogOpen,
    commandItems,
    navigationItems,
    setAnalyzeDialogOpen,
  };
};
