"use client";

import { ChevronDown, Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import type {
  GitHubOrganization,
  GitHubRepository,
  OrganizationAccessStatus,
  RepositoryCard as RepositoryCardType,
} from "@/lib/api/types";
import { cn } from "@/lib/utils";

import { useAllOrgRepos, useMyRepositories, useOrganizations, useUnanalyzedRepos } from "../hooks";
import { DiscoveryCard } from "./discovery-card";
import { DiscoveryListSheet } from "./discovery-list-sheet";
import { OrganizationPicker } from "./organization-picker";

type SheetState = {
  accessStatus?: OrganizationAccessStatus;
  isOpen: boolean;
  repositories: GitHubRepository[];
  title: string;
};

type OrgPickerState = {
  isOpen: boolean;
};

type DiscoverySectionProps = {
  analyzedRepositories: RepositoryCardType[];
};

export const DiscoverySection = ({ analyzedRepositories }: DiscoverySectionProps) => {
  const t = useTranslations("dashboard.discovery");

  const [isExpanded, setIsExpanded] = useState(false);

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

  const {
    isLoading: isLoadingOrgRepos,
    orgData,
    refreshOrg,
    totalUnanalyzedCount,
  } = useAllOrgRepos({
    analyzedRepositories,
    organizations: orgs,
  });

  const [sheetState, setSheetState] = useState<SheetState>({
    accessStatus: undefined,
    isOpen: false,
    repositories: [],
    title: "",
  });

  const [orgPickerState, setOrgPickerState] = useState<OrgPickerState>({
    isOpen: false,
  });

  const { count: unanalyzedPersonalCount, data: unanalyzedPersonalRepos } = useUnanalyzedRepos({
    analyzedRepositories,
    githubRepositories: myRepos,
  });

  const unanalyzedCounts = (() => {
    const counts: Record<string, number> = {};
    Object.entries(orgData).forEach(([orgLogin, data]) => {
      counts[orgLogin] = data.unanalyzedCount;
    });
    return counts;
  })();

  const handlePersonalClick = () => {
    setSheetState({
      accessStatus: undefined,
      isOpen: true,
      repositories: unanalyzedPersonalRepos,
      title: t("myRepos"),
    });
  };

  const handleOrgClick = () => {
    if (orgs.length === 0) {
      setSheetState({
        accessStatus: "restricted",
        isOpen: true,
        repositories: [],
        title: t("orgRepos"),
      });
    } else if (orgs.length === 1) {
      const org = orgs[0];
      const data = orgData[org.login];
      setSheetState({
        accessStatus: org.accessStatus,
        isOpen: true,
        repositories: data?.unanalyzedRepos ?? [],
        title: org.login,
      });
    } else {
      setOrgPickerState({ isOpen: true });
    }
  };

  const handleSelectOrg = (org: GitHubOrganization) => {
    const data = orgData[org.login];
    setOrgPickerState({ isOpen: false });
    if (data) {
      setSheetState({
        accessStatus: org.accessStatus,
        isOpen: true,
        repositories: data.unanalyzedRepos,
        title: org.login,
      });
    }
  };

  const handleOrgPickerOpenChange = (open: boolean) => {
    setOrgPickerState({ isOpen: open });
  };

  const handleSheetOpenChange = (open: boolean) => {
    setSheetState((prev) => ({ ...prev, isOpen: open }));
  };

  const handleRefreshOrgs = () => {
    refreshOrgs();
  };

  const isOrgLoading = isLoadingOrgs || isLoadingOrgRepos;
  const isOrgRefreshing = isRefreshingOrgs;
  const totalCount = unanalyzedPersonalCount + totalUnanalyzedCount;
  const isLoading = isLoadingMyRepos || isOrgLoading;

  const handleToggle = () => {
    setIsExpanded((prev) => !prev);
  };

  if (!isLoading && totalCount === 0) {
    return null;
  }

  return (
    <section aria-labelledby="discovery-heading" className="mt-8 rounded-lg border bg-muted/50">
      <Button
        aria-controls="discovery-content"
        aria-expanded={isExpanded}
        className="flex w-full items-center justify-between p-4 text-left hover:bg-muted/80"
        onClick={handleToggle}
        variant="ghost"
      >
        <div className="flex items-center gap-2">
          <Sparkles className="size-5 text-amber-500" />
          <h2 className="text-base font-semibold" id="discovery-heading">
            {t("title")}
          </h2>
          <span className="ml-1 rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700 dark:bg-amber-900 dark:text-amber-300">
            {isLoading ? "..." : totalCount}
          </span>
        </div>
        <ChevronDown
          className={cn(
            "size-5 text-muted-foreground transition-transform",
            isExpanded && "rotate-180"
          )}
        />
      </Button>

      <div
        className={cn(
          "grid transition-all duration-200 ease-in-out",
          isExpanded ? "grid-rows-[1fr]" : "grid-rows-[0fr]"
        )}
        id="discovery-content"
      >
        <div className="overflow-hidden">
          <div className="space-y-4 px-4 pb-4">
            <p className="text-sm text-muted-foreground">{t("description")}</p>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <DiscoveryCard
                count={unanalyzedPersonalCount}
                isLoading={isLoadingMyRepos}
                isRefreshing={isRefreshingMyRepos}
                onClick={handlePersonalClick}
                onRefresh={refreshMyRepos}
                type="personal"
              />
              <DiscoveryCard
                count={totalUnanalyzedCount}
                isLoading={isOrgLoading}
                isRefreshing={isOrgRefreshing}
                onClick={handleOrgClick}
                onRefresh={handleRefreshOrgs}
                type="organization"
              />
            </div>
          </div>
        </div>
      </div>

      <DiscoveryListSheet
        accessStatus={sheetState.accessStatus}
        isOpen={sheetState.isOpen}
        onOpenChange={handleSheetOpenChange}
        repositories={sheetState.repositories}
        title={sheetState.title}
      />

      <OrganizationPicker
        isLoading={isOrgLoading}
        isOpen={orgPickerState.isOpen}
        isRefreshing={isOrgRefreshing}
        onOpenChange={handleOrgPickerOpenChange}
        onRefreshOrg={refreshOrg}
        onSelectOrg={handleSelectOrg}
        organizations={orgs}
        unanalyzedCounts={unanalyzedCounts}
      />
    </section>
  );
};
