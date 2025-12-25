"use client";

import { Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";
import { useMemo, useState } from "react";

import type { GitHubRepository, RepositoryCard as RepositoryCardType } from "@/lib/api/types";

import { useMyRepositories, useOrganizations, useUnanalyzedRepos } from "../hooks";
import { DiscoveryCard } from "./discovery-card";
import { DiscoveryListSheet } from "./discovery-list-sheet";

type SheetState = {
  isOpen: boolean;
  repositories: GitHubRepository[];
  title: string;
};

type DiscoverySectionProps = {
  analyzedRepositories: RepositoryCardType[];
};

export const DiscoverySection = ({ analyzedRepositories }: DiscoverySectionProps) => {
  const t = useTranslations("dashboard.discovery");

  const {
    data: myRepos,
    isLoading: isLoadingMyRepos,
    isRefreshing: isRefreshingMyRepos,
    refresh: refreshMyRepos,
  } = useMyRepositories();

  const {
    data: orgs,
    isLoading: isLoadingOrgs,
    isRefreshing: isRefreshingOrgs,
    refresh: refreshOrgs,
  } = useOrganizations();

  const [sheetState, setSheetState] = useState<SheetState>({
    isOpen: false,
    repositories: [],
    title: "",
  });

  const { count: unanalyzedPersonalCount, data: unanalyzedPersonalRepos } = useUnanalyzedRepos({
    analyzedRepositories,
    githubRepositories: myRepos,
  });

  const orgRepoCount = useMemo(() => orgs.length, [orgs]);

  const handlePersonalClick = () => {
    setSheetState({
      isOpen: true,
      repositories: unanalyzedPersonalRepos,
      title: t("myRepos"),
    });
  };

  const handleOrgClick = () => {
    setSheetState({
      isOpen: true,
      repositories: [],
      title: t("orgRepos"),
    });
  };

  const handleSheetOpenChange = (open: boolean) => {
    setSheetState((prev) => ({ ...prev, isOpen: open }));
  };

  return (
    <section aria-labelledby="discovery-heading" className="mt-8">
      <div className="flex items-center gap-2 mb-4">
        <Sparkles className="size-5 text-amber-500" />
        <h2 className="text-lg font-semibold" id="discovery-heading">
          {t("title")}
        </h2>
      </div>

      <p className="text-sm text-muted-foreground mb-4">{t("description")}</p>

      <div className="grid gap-4 grid-cols-1 sm:grid-cols-2">
        <DiscoveryCard
          count={unanalyzedPersonalCount}
          isLoading={isLoadingMyRepos}
          isRefreshing={isRefreshingMyRepos}
          onClick={handlePersonalClick}
          onRefresh={refreshMyRepos}
          type="personal"
        />
        <DiscoveryCard
          count={orgRepoCount}
          isLoading={isLoadingOrgs}
          isRefreshing={isRefreshingOrgs}
          onClick={handleOrgClick}
          onRefresh={refreshOrgs}
          type="organization"
        />
      </div>

      <DiscoveryListSheet
        isOpen={sheetState.isOpen}
        onOpenChange={handleSheetOpenChange}
        repositories={sheetState.repositories}
        title={sheetState.title}
      />
    </section>
  );
};
