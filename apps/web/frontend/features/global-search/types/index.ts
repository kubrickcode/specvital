import type { LucideIcon } from "lucide-react";

export type SearchCategory = "recent" | "repositories" | "bookmarks" | "community" | "actions";

export type SearchableItem = {
  category: SearchCategory;
  icon?: LucideIcon;
  id: string;
  keywords?: string[];
  label: string;
  labelKey?: string;
  metadata?: Record<string, unknown>;
  onSelect: () => void;
  shortcut?: string;
};

export type StaticAction = {
  icon: LucideIcon;
  id: string;
  keywords?: string[];
  labelKey: string;
  requiresAuth?: boolean;
  shortcut?: string;
  type: "navigation" | "action";
};

export type GlobalSearchState = {
  isOpen: boolean;
};
