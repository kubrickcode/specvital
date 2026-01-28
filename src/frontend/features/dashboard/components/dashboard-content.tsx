"use client";

import { Compass } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { AnalyzeDialog } from "@/features/home";
import { Link } from "@/i18n/navigation";

import {
  useAddBookmark,
  usePaginatedRepositories,
  useReanalyze,
  useRemoveBookmark,
} from "../hooks";
import { useBookmarkFilter } from "../hooks/use-bookmark-filter";
import { useOwnershipFilter } from "../hooks/use-ownership-filter";
import type { SortOption } from "../types";
import { ActiveTasksSection } from "./active-tasks-section";
import { EmptyStateVariant } from "./empty-state-variant";
import { FilterBar } from "./filter-bar";
import { InfiniteScrollLoader } from "./infinite-scroll-loader";
import { RepositoryList } from "./repository-list";
import { SummarySection } from "./summary-section";

export const DashboardContent = () => {
  const t = useTranslations("dashboard");

  const { addBookmark } = useAddBookmark();
  const { removeBookmark } = useRemoveBookmark();
  const { reanalyze } = useReanalyze();
  const { ownershipFilter } = useOwnershipFilter();
  const { bookmarkOnly } = useBookmarkFilter();

  const [sortBy, setSortBy] = useState<SortOption>("recent");
  const [searchQuery, setSearchQuery] = useState("");

  const {
    data: repositories,
    fetchNextPage,
    hasNextPage,
    isError,
    isFetchingNextPage,
    isLoading,
    refetch,
  } = usePaginatedRepositories({
    ownership: ownershipFilter,
    sortBy,
    sortOrder: "desc",
    view: "my",
  });

  const filteredRepositories = (() => {
    let result = repositories;

    if (bookmarkOnly) {
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
  })();

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

  const handleLoadMore = () => {
    fetchNextPage();
  };

  const handleRetry = () => {
    refetch();
  };

  const hasNoRepositories = !isLoading && repositories.length === 0 && !isError;
  const hasNoFilterResults =
    bookmarkOnly && filteredRepositories.length === 0 && repositories.length > 0;
  const hasNoSearchResults =
    searchQuery.trim() !== "" && filteredRepositories.length === 0 && repositories.length > 0;

  return (
    <div className="space-y-8">
      <SummarySection />

      <ActiveTasksSection />

      <FilterBar
        hasNextPage={hasNextPage && !bookmarkOnly}
        isLoading={isLoading}
        onSearchChange={setSearchQuery}
        onSortChange={setSortBy}
        searchQuery={searchQuery}
        sortBy={sortBy}
        totalLoaded={bookmarkOnly ? filteredRepositories.length : repositories.length}
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

          {!bookmarkOnly && !searchQuery.trim() && (
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
