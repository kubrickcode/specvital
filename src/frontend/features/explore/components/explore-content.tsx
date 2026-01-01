"use client";

import { useQueryClient } from "@tanstack/react-query";
import { Building2, Globe, User } from "lucide-react";
import { useTranslations } from "next-intl";
import { useCallback, useMemo, useState } from "react";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  EmptyStateVariant,
  fetchPaginatedRepositories,
  LoadMoreButton,
  paginatedRepositoriesKeys,
  RepositoryList,
  useAddBookmark,
  useReanalyze,
  useRemoveBookmark,
  type SortOption,
} from "@/features/dashboard";

import { useExploreRepositories } from "../hooks";
import { MyReposTab } from "./my-repos-tab";
import { OrgReposTab } from "./org-repos-tab";
import { SearchSortControls } from "./search-sort-controls";

const DEFAULT_PAGE_LIMIT = 10;

type ExploreTab = "community" | "my-repos" | "organizations";

export const ExploreContent = () => {
  const t = useTranslations("explore");
  const queryClient = useQueryClient();

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

  const handlePrefetchNextPage = useCallback(() => {
    if (!hasNextPage || isFetchingNextPage) return;

    const queryKey = paginatedRepositoriesKeys.list({
      limit: DEFAULT_PAGE_LIMIT,
      sortBy,
      sortOrder: "desc",
      view: "community",
    });

    const lastPage = queryClient.getQueryData<{
      pageParams: (string | undefined)[];
      pages: { hasNext: boolean; nextCursor?: string | null }[];
    }>(queryKey);

    const nextCursor = lastPage?.pages.at(-1)?.nextCursor;
    if (!nextCursor) return;

    queryClient.prefetchInfiniteQuery({
      initialPageParam: undefined as string | undefined,
      queryFn: () =>
        fetchPaginatedRepositories({
          cursor: nextCursor,
          limit: DEFAULT_PAGE_LIMIT,
          sortBy,
          sortOrder: "desc",
          view: "community",
        }),
      queryKey,
      staleTime: 30 * 1000,
    });
  }, [hasNextPage, isFetchingNextPage, queryClient, sortBy]);

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
          <User className="size-4" />
          {t("tabs.myRepos")}
        </TabsTrigger>
        <TabsTrigger className="gap-2" value="organizations">
          <Building2 className="size-4" />
          {t("tabs.organizations")}
        </TabsTrigger>
      </TabsList>

      <TabsContent value="community">
        <div className="space-y-8">
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

              <LoadMoreButton
                hasError={isError}
                hasNextPage={hasNextPage && !searchQuery.trim()}
                isFetchingNextPage={isFetchingNextPage}
                onLoadMore={handleLoadMore}
                onPrefetch={handlePrefetchNextPage}
                onRetry={handleRetry}
              />
            </>
          )}
        </div>
      </TabsContent>

      <TabsContent value="my-repos">
        <MyReposTab />
      </TabsContent>

      <TabsContent value="organizations">
        <OrgReposTab />
      </TabsContent>
    </Tabs>
  );
};
