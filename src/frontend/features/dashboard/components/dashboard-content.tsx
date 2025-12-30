"use client";

import { useQueryClient } from "@tanstack/react-query";
import { ArrowUpDown, Search } from "lucide-react";
import { useTranslations } from "next-intl";

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

import {
  useAddBookmark,
  useReanalyze,
  useRecentRepositories,
  useRemoveBookmark,
  useRepositorySearch,
} from "../hooks";
import type { SortOption } from "../types";
import { DiscoveryErrorFallback } from "./discovery-error-fallback";
import { DiscoverySection } from "./discovery-section";
import { EmptyStateVariant } from "./empty-state-variant";
import { RepositoryList } from "./repository-list";

const SORT_OPTIONS: SortOption[] = ["name", "recent", "tests"];

const isSortOption = (value: string): value is SortOption =>
  SORT_OPTIONS.includes(value as SortOption);

type SearchSortControlsProps = {
  onSearchChange: (query: string) => void;
  onSortChange: (sort: SortOption) => void;
  searchQuery: string;
  sortBy: SortOption;
};

const SearchSortControls = ({
  onSearchChange,
  onSortChange,
  searchQuery,
  sortBy,
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
    <div className="mb-4 flex flex-col gap-3 sm:flex-row">
      <div className="relative flex-1">
        <Search
          aria-hidden="true"
          className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
        />
        <Input
          aria-label={t("searchPlaceholder")}
          className="pl-9"
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder={t("searchPlaceholder")}
          type="search"
          value={searchQuery}
        />
      </div>

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button className="w-full sm:w-auto" variant="outline">
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
  );
};

export const DashboardContent = () => {
  const queryClient = useQueryClient();

  const { data: repositories = [], isLoading } = useRecentRepositories();

  const { addBookmark } = useAddBookmark();
  const { removeBookmark } = useRemoveBookmark();
  const { reanalyze } = useReanalyze();

  const { filteredRepositories, searchQuery, setSearchQuery, setSortBy, sortBy } =
    useRepositorySearch(repositories);

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

  const hasNoRepositories = !isLoading && repositories.length === 0;

  return (
    <div className="space-y-8">
      <SearchSortControls
        onSearchChange={setSearchQuery}
        onSortChange={setSortBy}
        searchQuery={searchQuery}
        sortBy={sortBy}
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
      ) : filteredRepositories.length === 0 ? (
        <EmptyStateVariant searchQuery={searchQuery} variant="no-search-results" />
      ) : (
        <RepositoryList
          onBookmarkToggle={handleBookmarkToggle}
          onReanalyze={handleReanalyze}
          repositories={filteredRepositories}
        />
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
