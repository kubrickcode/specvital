"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";

import { AuthErrorBoundary } from "@/components/feedback";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AnalyzeDialog } from "@/features/home";

import {
  useAddBookmark,
  useBookmarkedRepositories,
  useMyAnalyses,
  useOwnershipFilter,
  useReanalyze,
  useRecentRepositories,
  useRemoveBookmark,
  useRepositorySearch,
  useTabState,
} from "../hooks";
import type { TabValue } from "../hooks";
import { DashboardHeader } from "./dashboard-header";
import { DiscoveryErrorFallback } from "./discovery-error-fallback";
import { DiscoverySection } from "./discovery-section";
import { EmptyStateVariant } from "./empty-state-variant";
import { OwnershipFilter } from "./ownership-filter";
import { RepositoryList } from "./repository-list";

export const DashboardContent = () => {
  const t = useTranslations("dashboard");
  const queryClient = useQueryClient();

  const { setTab, tab } = useTabState();
  const { ownership, setOwnership } = useOwnershipFilter();

  const { data: bookmarked = [], isLoading: isLoadingBookmarked } = useBookmarkedRepositories();
  const { data: recent = [], isLoading: isLoadingRecent } = useRecentRepositories();
  const { data: myAnalyses = [], isLoading: isLoadingMyAnalyses } = useMyAnalyses({ ownership });

  const { addBookmark } = useAddBookmark();
  const { removeBookmark } = useRemoveBookmark();
  const { reanalyze } = useReanalyze();

  const { filteredRepositories, searchQuery, setSearchQuery, setSortBy, sortBy } =
    useRepositorySearch(recent);

  const handleTabChange = (value: string) => {
    setTab(value as TabValue);
  };

  const handleBookmarkToggle = (owner: string, repo: string, isBookmarked: boolean) => {
    if (isBookmarked) {
      removeBookmark(owner, repo);
    } else {
      addBookmark(owner, repo);
    }
  };

  const handleReanalyze = (owner: string, repo: string) => {
    reanalyze(owner, repo);
  };

  const handleDiscoveryReset = () => {
    queryClient.resetQueries({ exact: false, queryKey: ["dashboard"] });
  };

  const isLoading = isLoadingBookmarked || isLoadingRecent;
  const hasNoRepositories = !isLoading && recent.length === 0;
  const hasNoMyAnalyses = !isLoadingMyAnalyses && myAnalyses.length === 0;

  return (
    <div className="space-y-8">
      <DashboardHeader
        onSearchChange={setSearchQuery}
        onSortChange={setSortBy}
        searchQuery={searchQuery}
        sortBy={sortBy}
      />

      <Tabs onValueChange={handleTabChange} value={tab}>
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <TabsList>
            <TabsTrigger value="bookmarked">{t("tabs.bookmarked")}</TabsTrigger>
            <TabsTrigger value="my-analyses">{t("tabs.myAnalyses")}</TabsTrigger>
          </TabsList>

          {tab === "my-analyses" && <OwnershipFilter onChange={setOwnership} value={ownership} />}
        </div>

        <TabsContent className="mt-6" value="bookmarked">
          {isLoadingBookmarked ? (
            <RepositoryList
              isLoading
              onBookmarkToggle={handleBookmarkToggle}
              onReanalyze={handleReanalyze}
              repositories={[]}
            />
          ) : bookmarked.length === 0 ? (
            <EmptyStateVariant variant="no-bookmarks" />
          ) : (
            <RepositoryList
              onBookmarkToggle={handleBookmarkToggle}
              onReanalyze={handleReanalyze}
              repositories={bookmarked}
            />
          )}
        </TabsContent>

        <TabsContent className="mt-6" value="my-analyses">
          {isLoadingMyAnalyses ? (
            <RepositoryList
              isLoading
              onBookmarkToggle={handleBookmarkToggle}
              onReanalyze={handleReanalyze}
              repositories={[]}
            />
          ) : hasNoMyAnalyses ? (
            <EmptyStateVariant
              action={<AnalyzeDialog variant="empty-state" />}
              variant="no-repos"
            />
          ) : (
            <RepositoryList
              onBookmarkToggle={handleBookmarkToggle}
              onReanalyze={handleReanalyze}
              repositories={myAnalyses}
            />
          )}
        </TabsContent>
      </Tabs>

      <section aria-labelledby="all-repos-heading">
        <h2 className="mb-4 text-xl font-semibold" id="all-repos-heading">
          {t("allRepositories")}
        </h2>

        {isLoading ? (
          <RepositoryList
            isLoading
            onBookmarkToggle={handleBookmarkToggle}
            onReanalyze={handleReanalyze}
            repositories={[]}
          />
        ) : hasNoRepositories ? (
          <EmptyStateVariant action={<AnalyzeDialog variant="empty-state" />} variant="no-repos" />
        ) : filteredRepositories.length === 0 ? (
          <EmptyStateVariant searchQuery={searchQuery} variant="no-search-results" />
        ) : (
          <RepositoryList
            onBookmarkToggle={handleBookmarkToggle}
            onReanalyze={handleReanalyze}
            repositories={filteredRepositories}
          />
        )}
      </section>

      <AuthErrorBoundary
        fallback={<DiscoveryErrorFallback resetErrorBoundary={handleDiscoveryReset} />}
        onReset={handleDiscoveryReset}
      >
        <DiscoverySection analyzedRepositories={recent} />
      </AuthErrorBoundary>
    </div>
  );
};
