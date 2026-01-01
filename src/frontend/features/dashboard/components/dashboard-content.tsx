"use client";

import { useQueryClient } from "@tanstack/react-query";
import { ArrowUpDown, Search } from "lucide-react";
import { useTranslations } from "next-intl";
import { useCallback, useMemo, useState } from "react";

import { AuthErrorBoundary } from "@/components/feedback";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { AnalyzeDialog } from "@/features/home";
import type { ViewFilterParam } from "@/lib/api/types";

import { fetchPaginatedRepositories } from "../api";
import {
  paginatedRepositoriesKeys,
  useAddBookmark,
  usePaginatedRepositories,
  useReanalyze,
  useRemoveBookmark,
  useViewFilter,
} from "../hooks";
import type { SortOption } from "../types";
import { DiscoveryErrorFallback } from "./discovery-error-fallback";
import { DiscoverySection } from "./discovery-section";
import { EmptyStateVariant } from "./empty-state-variant";
import { LoadMoreButton } from "./load-more-button";
import { PaginationStatus } from "./pagination-status";
import { RepositoryList } from "./repository-list";
import { ViewFilterDropdown } from "./view-filter-dropdown";

const SORT_OPTIONS: SortOption[] = ["name", "recent", "tests"];

const isSortOption = (value: string): value is SortOption =>
  SORT_OPTIONS.includes(value as SortOption);

const mapViewFilterToParam = (view: string): ViewFilterParam | undefined => {
  switch (view) {
    case "mine":
      return "my";
    case "community":
      return "community";
    case "starred":
    case "all":
    default:
      return undefined;
  }
};

type SearchSortControlsProps = {
  hasNextPage: boolean;
  isLoading: boolean;
  onSearchChange: (query: string) => void;
  onSortChange: (sort: SortOption) => void;
  searchQuery: string;
  sortBy: SortOption;
  totalLoaded: number;
};

const SearchSortControls = ({
  hasNextPage,
  isLoading,
  onSearchChange,
  onSortChange,
  searchQuery,
  sortBy,
  totalLoaded,
}: SearchSortControlsProps) => {
  const t = useTranslations("dashboard");

  const sortLabels: Record<SortOption, string> = {
    name: t("sort.name"),
    recent: t("sort.recent"),
    tests: t("sort.tests"),
  };

  const handleSortChange = (value: string) => {
    if (isSortOption(value)) {
      onSortChange(value);
    }
  };

  return (
    <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:flex-1">
        <div className="relative flex-1 sm:max-w-sm">
          <Search
            aria-hidden="true"
            className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
          />
          <Input
            aria-label={t("searchPlaceholder")}
            className="h-11 pl-10 sm:h-9 sm:pl-9"
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder={t("searchPlaceholder")}
            type="search"
            value={searchQuery}
          />
        </div>

        <div className="flex w-full gap-2 sm:w-auto">
          <ViewFilterDropdown />

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button className="h-11 flex-1 sm:h-9 sm:flex-none" variant="outline">
                <ArrowUpDown aria-hidden="true" />
                <span>
                  {t("sort.label")}: {sortLabels[sortBy]}
                </span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-40">
              <DropdownMenuRadioGroup onValueChange={handleSortChange} value={sortBy}>
                <DropdownMenuRadioItem value="recent">{sortLabels.recent}</DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="name">{sortLabels.name}</DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="tests">{sortLabels.tests}</DropdownMenuRadioItem>
              </DropdownMenuRadioGroup>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <PaginationStatus hasNextPage={hasNextPage} isLoading={isLoading} totalLoaded={totalLoaded} />
    </div>
  );
};

export const DashboardContent = () => {
  const queryClient = useQueryClient();

  const { addBookmark } = useAddBookmark();
  const { removeBookmark } = useRemoveBookmark();
  const { reanalyze } = useReanalyze();
  const { viewFilter } = useViewFilter();

  const [sortBy, setSortBy] = useState<SortOption>("recent");
  const [searchQuery, setSearchQuery] = useState("");

  const viewParam = mapViewFilterToParam(viewFilter);

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

    if (viewFilter === "starred") {
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
  }, [repositories, searchQuery, viewFilter]);

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

  const handleDiscoveryReset = useCallback(() => {
    queryClient.resetQueries({ exact: false, queryKey: ["dashboard"] });
  }, [queryClient]);

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
    viewFilter === "starred" && filteredRepositories.length === 0 && repositories.length > 0;
  const hasNoSearchResults =
    searchQuery.trim() !== "" && filteredRepositories.length === 0 && repositories.length > 0;

  return (
    <div className="space-y-8">
      <SearchSortControls
        hasNextPage={hasNextPage && viewFilter !== "starred"}
        isLoading={isLoading}
        onSearchChange={setSearchQuery}
        onSortChange={setSortBy}
        searchQuery={searchQuery}
        sortBy={sortBy}
        totalLoaded={viewFilter === "starred" ? filteredRepositories.length : repositories.length}
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

          {viewFilter !== "starred" && (
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

      <AuthErrorBoundary
        fallback={<DiscoveryErrorFallback resetErrorBoundary={handleDiscoveryReset} />}
        onReset={handleDiscoveryReset}
      >
        <DiscoverySection analyzedRepositories={repositories} />
      </AuthErrorBoundary>
    </div>
  );
};
