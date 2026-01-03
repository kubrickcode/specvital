"use client";

import { Building2, Globe, Lock, User } from "lucide-react";
import { useTranslations } from "next-intl";
import { useCallback, useMemo, useState } from "react";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAuth } from "@/features/auth";
import {
  EmptyStateVariant,
  InfiniteScrollLoader,
  RepositoryList,
  useAddBookmark,
  useReanalyze,
  useRemoveBookmark,
  type SortOption,
} from "@/features/dashboard";

import { useExploreRepositories } from "../hooks";
import { LoginRequiredState } from "./login-required-state";
import { MyReposTab } from "./my-repos-tab";
import { OrgReposTab } from "./org-repos-tab";
import { SearchSortControls } from "./search-sort-controls";

type ExploreTab = "community" | "my-repos" | "organizations";

export const ExploreContent = () => {
  const t = useTranslations("explore");
  const { isAuthenticated } = useAuth();

  const { addBookmark } = useAddBookmark();
  const { removeBookmark } = useRemoveBookmark();
  const { reanalyze } = useReanalyze();

  const [activeTab, setActiveTab] = useState<ExploreTab>("community");
  const [sortBy, setSortBy] = useState<SortOption>("recent");
  const [searchQuery, setSearchQuery] = useState("");

  const {
    data: communityRepositories,
    fetchNextPage,
    hasNextPage,
    isError,
    isFetchingNextPage,
    isLoading: isLoadingCommunity,
    refetch,
  } = useExploreRepositories({
    sortBy,
    sortOrder: "desc",
  });

  const filteredRepositories = useMemo(() => {
    if (!searchQuery.trim()) {
      return communityRepositories;
    }

    const query = searchQuery.toLowerCase();
    return communityRepositories.filter(
      (repo) =>
        repo.owner.toLowerCase().includes(query) ||
        repo.name.toLowerCase().includes(query) ||
        `${repo.owner}/${repo.name}`.toLowerCase().includes(query)
    );
  }, [communityRepositories, searchQuery]);

  const handleBookmarkToggle = useCallback(
    (owner: string, repo: string, isBookmarked: boolean) => {
      if (isBookmarked) {
        removeBookmark(owner, repo);
      } else {
        addBookmark(owner, repo);
      }
    },
    [addBookmark, removeBookmark]
  );

  const handleReanalyze = useCallback(
    (owner: string, repo: string) => {
      reanalyze(owner, repo);
    },
    [reanalyze]
  );

  const handleLoadMore = useCallback(() => {
    fetchNextPage();
  }, [fetchNextPage]);

  const handleRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  const handleTabChange = (value: string) => {
    setActiveTab(value as ExploreTab);
    setSearchQuery("");
  };

  const hasNoRepositories = !isLoadingCommunity && communityRepositories.length === 0 && !isError;
  const hasNoSearchResults =
    searchQuery.trim() !== "" &&
    filteredRepositories.length === 0 &&
    communityRepositories.length > 0;

  return (
    <Tabs defaultValue="community" onValueChange={handleTabChange} value={activeTab}>
      <TabsList className="mb-6">
        <TabsTrigger className="gap-2" value="community">
          <Globe className="size-4" />
          {t("tabs.community")}
        </TabsTrigger>
        <TabsTrigger className="gap-2" value="my-repos">
          {isAuthenticated ? <User className="size-4" /> : <Lock className="size-4" />}
          {t("tabs.myRepos")}
        </TabsTrigger>
        <TabsTrigger className="gap-2" value="organizations">
          {isAuthenticated ? <Building2 className="size-4" /> : <Lock className="size-4" />}
          {t("tabs.organizations")}
        </TabsTrigger>
      </TabsList>

      <TabsContent value="community">
        <div className="space-y-8">
          <p className="text-sm text-muted-foreground">{t("community.visibilityDisclosure")}</p>
          <SearchSortControls
            hasNextPage={hasNextPage}
            isLoading={isLoadingCommunity}
            onSearchChange={setSearchQuery}
            onSortChange={setSortBy}
            searchQuery={searchQuery}
            sortBy={sortBy}
            totalLoaded={communityRepositories.length}
          />

          {isLoadingCommunity ? (
            <RepositoryList
              isLoading
              onBookmarkToggle={handleBookmarkToggle}
              onReanalyze={handleReanalyze}
              repositories={[]}
            />
          ) : hasNoRepositories ? (
            <EmptyStateVariant variant="no-repos" />
          ) : hasNoSearchResults ? (
            <EmptyStateVariant searchQuery={searchQuery} variant="no-search-results" />
          ) : (
            <>
              <RepositoryList
                onBookmarkToggle={handleBookmarkToggle}
                onReanalyze={handleReanalyze}
                repositories={filteredRepositories}
              />

              {!searchQuery.trim() && (
                <InfiniteScrollLoader
                  hasError={isError}
                  hasNextPage={hasNextPage}
                  isFetchingNextPage={isFetchingNextPage}
                  onLoadMore={handleLoadMore}
                  onRetry={handleRetry}
                />
              )}
            </>
          )}
        </div>
      </TabsContent>

      <TabsContent value="my-repos">
        {isAuthenticated ? (
          <MyReposTab />
        ) : (
          <LoginRequiredState descriptionKey="myReposDescription" titleKey="myReposTitle" />
        )}
      </TabsContent>

      <TabsContent value="organizations">
        {isAuthenticated ? (
          <OrgReposTab />
        ) : (
          <LoginRequiredState
            descriptionKey="organizationsDescription"
            titleKey="organizationsTitle"
          />
        )}
      </TabsContent>
    </Tabs>
  );
};
