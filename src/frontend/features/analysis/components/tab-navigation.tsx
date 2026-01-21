"use client";

import { useTranslations } from "next-intl";

import { cn } from "@/lib/utils";

import type { PrimaryTab } from "../types/primary-tab";

type TabNavigationProps = {
  activeTab: PrimaryTab;
  onTabChange: (tab: PrimaryTab) => void;
};

export const TabNavigation = ({ activeTab, onTabChange }: TabNavigationProps) => {
  const t = useTranslations("analyze.tabs");

  const tabs: { id: string; label: string; panelId: string; value: PrimaryTab }[] = [
    { id: "tab-tests", label: t("tests"), panelId: "tabpanel-tests", value: "tests" },
    { id: "tab-spec", label: t("spec"), panelId: "tabpanel-spec", value: "spec" },
  ];

  return (
    <div className="flex gap-1" role="tablist">
      {tabs.map((tab) => (
        <button
          aria-controls={tab.panelId}
          aria-selected={activeTab === tab.value}
          className={cn(
            "relative px-3 py-1.5 text-sm font-medium rounded-md transition-colors",
            "hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
            activeTab === tab.value
              ? "text-foreground bg-muted"
              : "text-muted-foreground hover:text-foreground"
          )}
          data-state={activeTab === tab.value ? "active" : "inactive"}
          id={tab.id}
          key={tab.value}
          onClick={() => onTabChange(tab.value)}
          role="tab"
          type="button"
        >
          {tab.label}
        </button>
      ))}
    </div>
  );
};
