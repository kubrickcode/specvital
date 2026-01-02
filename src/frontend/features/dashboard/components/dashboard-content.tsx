"use client";

import { useQueryClient } from "@tanstack/react-query";
import { Compass } from "lucide-react";
import { useTranslations } from "next-intl";
import { useCallback, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { AnalyzeDialog } from "@/features/home";
import { Link } from "@/i18n/navigation";
import type { ViewFilterParam } from "@/lib/api/types";

import { fetchPaginatedRepositories } from "../api";
import {
  paginatedRepositoriesKeys,
  useAddBookmark,
  usePaginatedRepositories,
  useReanalyze,
  useRemoveBookmark,
} from "../hooks";
import { type OwnershipFilter, useOwnershipFilter } from "../hooks/use-ownership-filter";
import { useStarredFilter } from "../hooks/use-starred-filter";
import type { SortOption } from "../types";
import { AttentionZone } from "./attention-zone";
import { EmptyStateVariant } from "./empty-state-variant";
import { FilterBar } from "./filter-bar";
import { LoadMoreButton } from "./load-more-button";
import { RepositoryList } from "./repository-list";
import { SummarySection } from "./summary-section";

const mapOwnershipToViewParam = (ownership: OwnershipFilter): ViewFilterParam | undefined => {
  switch (ownership) {
    case "mine":
      return "my";
    case "all":
    default:
      return undefined;
  }
};

export const DashboardContent = () => {
  const t = useTranslations("dashboard");
  const queryClient = useQueryClient();

  const { addBookmark } = useAddBookmark();
  const { removeBookmark } = useRemoveBookmark();
  const { reanalyze } = useReanalyze();
  const { ownershipFilter } = useOwnershipFilter();
  const { starredOnly } = useStarredFilter();

  const [sortBy, setSortBy] = useState<SortOption>("recent");
  const [searchQuery, setSearchQuery] = useState("");

  const viewParam = mapOwnershipToViewParam(ownershipFilter);

  const {
    data: repositories,
    fetchNextPage,
    hasNextPage,
    isError,
    isFetchingNextPage,
    isLoading,
    refetch,
  } = usePaginatedRepositories({
    sortBy,
    sortOrder: "desc",
    view: viewParam,
  });

  const filteredRepositories = useMemo(() => {
    let result = repositories;

    if (starredOnly) {
      result = result.filter((repo) => repo.isBookmarked);
    }

    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      result = result.filter(
        (repo) =>
          repo.owner.toLowerCase().includes(query) ||
          repo.name.toLowerCase().includes(query) ||
          `${repo.owner}/${repo.name}`.toLowerCase().includes(query)
      );
    }

    return result;
  }, [repositories, searchQuery, starredOnly]);

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

    const lastPage = queryClient.getQueryData<{
      pageParams: (string | undefined)[];
      pages: { hasNext: boolean; nextCursor?: string | null }[];
    }>(paginatedRepositoriesKeys.list({ limit: 10, sortBy, sortOrder: "desc", view: viewParam }));

    const nextCursor = lastPage?.pages.at(-1)?.nextCursor;
    if (!nextCursor) return;

    queryClient.prefetchInfiniteQuery({
      initialPageParam: undefined as string | undefined,
      queryFn: () =>
        fetchPaginatedRepositories({
          cursor: nextCursor,
          limit: 10,
          sortBy,
          sortOrder: "desc",
          view: viewParam,
        }),
      queryKey: paginatedRepositoriesKeys.list({
        limit: 10,
        sortBy,
        sortOrder: "desc",
        view: viewParam,
      }),
      staleTime: 30 * 1000,
    });
  }, [hasNextPage, isFetchingNextPage, queryClient, sortBy, viewParam]);

  const handleRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  const hasNoRepositories = !isLoading && repositories.length === 0 && !isError;
  const hasNoFilterResults =
    starredOnly && filteredRepositories.length === 0 && repositories.length > 0;
  const hasNoSearchResults =
    searchQuery.trim() !== "" && filteredRepositories.length === 0 && repositories.length > 0;

  return (
    <div className="space-y-8">
      <SummarySection />

      <AttentionZone onReanalyze={handleReanalyze} repositories={repositories} />

      <FilterBar
        hasNextPage={hasNextPage && !starredOnly}
        isLoading={isLoading}
        onSearchChange={setSearchQuery}
        onSortChange={setSortBy}
        searchQuery={searchQuery}
        sortBy={sortBy}
        totalLoaded={starredOnly ? filteredRepositories.length : repositories.length}
      />

      {isLoading ? (
        <RepositoryList
          isLoading
          onBookmarkToggle={handleBookmarkToggle}
          onReanalyze={handleReanalyze}
          repositories={[]}
        />
      ) : hasNoRepositories ? (
        <EmptyStateVariant action={<AnalyzeDialog variant="empty-state" />} variant="no-repos" />
      ) : hasNoFilterResults ? (
        <EmptyStateVariant variant="no-bookmarks" />
      ) : hasNoSearchResults ? (
        <EmptyStateVariant searchQuery={searchQuery} variant="no-search-results" />
      ) : (
        <>
          <RepositoryList
            onBookmarkToggle={handleBookmarkToggle}
            onReanalyze={handleReanalyze}
            repositories={filteredRepositories}
          />

          {!starredOnly && (
            <LoadMoreButton
              hasError={isError}
              hasNextPage={hasNextPage && !searchQuery.trim()}
              isFetchingNextPage={isFetchingNextPage}
              onLoadMore={handleLoadMore}
              onPrefetch={handlePrefetchNextPage}
              onRetry={handleRetry}
            />
          )}
        </>
      )}

      <div className="flex justify-center pt-4 border-t">
        <Link href="/explore">
          <Button className="gap-2" variant="ghost">
            <Compass aria-hidden="true" className="size-4" />
            {t("exploreCta")}
          </Button>
        </Link>
      </div>
    </div>
  );
};
