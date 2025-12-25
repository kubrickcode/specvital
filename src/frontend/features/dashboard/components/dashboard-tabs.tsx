"use client";

import { History, Star } from "lucide-react";
import { useTranslations } from "next-intl";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import type { TabValue } from "../hooks";
import { useTabState } from "../hooks";

type DashboardTabsProps = {
  bookmarkedContent: React.ReactNode;
  myAnalysesContent: React.ReactNode;
};

export const DashboardTabs = ({ bookmarkedContent, myAnalysesContent }: DashboardTabsProps) => {
  const t = useTranslations("dashboard.tabs");
  const { setTab, tab } = useTabState();

  const handleTabChange = (value: string) => {
    setTab(value as TabValue);
  };

  return (
    <Tabs onValueChange={handleTabChange} value={tab}>
      <TabsList>
        <TabsTrigger value="bookmarked">
          <Star aria-hidden="true" className="size-4" />
          {t("bookmarked")}
        </TabsTrigger>
        <TabsTrigger value="my-analyses">
          <History aria-hidden="true" className="size-4" />
          {t("myAnalyses")}
        </TabsTrigger>
      </TabsList>
      <TabsContent value="bookmarked">{bookmarkedContent}</TabsContent>
      <TabsContent value="my-analyses">{myAnalysesContent}</TabsContent>
    </Tabs>
  );
};
