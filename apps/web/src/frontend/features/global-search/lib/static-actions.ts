import { BookOpen, Home, LayoutDashboard, Moon, Plus, Search, Sun, Wallet } from "lucide-react";

import type { StaticAction } from "../types";

export const NAVIGATION_ACTIONS: StaticAction[] = [
  {
    icon: LayoutDashboard,
    id: "nav-dashboard",
    keywords: ["home", "main", "repos"],
    labelKey: "actions.goToDashboard",
    requiresAuth: true,
    type: "navigation",
  },
  {
    icon: Search,
    id: "nav-explore",
    keywords: ["browse", "community"],
    labelKey: "actions.goToExplore",
    type: "navigation",
  },
  {
    icon: Wallet,
    id: "nav-account",
    keywords: ["settings", "profile", "plan"],
    labelKey: "actions.goToAccount",
    requiresAuth: true,
    type: "navigation",
  },
  {
    icon: BookOpen,
    id: "nav-pricing",
    keywords: ["plans", "billing"],
    labelKey: "actions.goToPricing",
    type: "navigation",
  },
  {
    icon: Home,
    id: "nav-home",
    keywords: ["landing", "start"],
    labelKey: "actions.goToHome",
    type: "navigation",
  },
];

export const COMMAND_ACTIONS: StaticAction[] = [
  {
    icon: Plus,
    id: "action-new-analysis",
    keywords: ["analyze", "repository", "github", "new"],
    labelKey: "actions.newAnalysis",
    type: "action",
  },
  {
    icon: Sun,
    id: "action-toggle-theme-light",
    keywords: ["light", "day", "bright"],
    labelKey: "actions.switchToLight",
    type: "action",
  },
  {
    icon: Moon,
    id: "action-toggle-theme-dark",
    keywords: ["dark", "night"],
    labelKey: "actions.switchToDark",
    type: "action",
  },
];

export const ALL_STATIC_ACTIONS: StaticAction[] = [...NAVIGATION_ACTIONS, ...COMMAND_ACTIONS];
