"use client";

import { useTranslations } from "next-intl";

import {
  useAddBookmark,
  useBookmarkedRepositories,
  useReanalyze,
  useRecentRepositories,
  useRemoveBookmark,
  useRepositorySearch,
} from "../hooks";
import { BookmarkedSection } from "./bookmarked-section";
import { DashboardHeader } from "./dashboard-header";
import { EmptyState } from "./empty-state";
import { RepositoryList } from "./repository-list";

export const DashboardContent = () => {
  const t = useTranslations("dashboard");

  const { data: bookmarked = [], isLoading: isLoadingBookmarked } = useBookmarkedRepositories();
  const { data: recent = [], isLoading: isLoadingRecent } = useRecentRepositories();

  const { addBookmark } = useAddBookmark();
  const { removeBookmark } = useRemoveBookmark();
  const { reanalyze } = useReanalyze();

  const { filteredRepositories, searchQuery, setSearchQuery, setSortBy, sortBy } =
    useRepositorySearch(recent);

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

  const isLoading = isLoadingBookmarked || isLoadingRecent;
  const hasNoRepositories = !isLoading && recent.length === 0;

  return (
    <div className="space-y-8">
      <DashboardHeader
        onSearchChange={setSearchQuery}
        onSortChange={setSortBy}
        searchQuery={searchQuery}
        sortBy={sortBy}
      />

      {!isLoading && bookmarked.length > 0 && (
        <BookmarkedSection
          onBookmarkToggle={handleBookmarkToggle}
          onReanalyze={handleReanalyze}
          repositories={bookmarked}
        />
      )}

      <section aria-labelledby="all-repos-heading">
        <h2 className="mb-4 text-xl font-semibold" id="all-repos-heading">
          {t("allRepositories", { defaultValue: "All Repositories" })}
        </h2>

        {isLoading ? (
          <RepositoryList
            isLoading
            onBookmarkToggle={handleBookmarkToggle}
            onReanalyze={handleReanalyze}
            repositories={[]}
          />
        ) : hasNoRepositories ? (
          <EmptyState />
        ) : (
          <RepositoryList
            onBookmarkToggle={handleBookmarkToggle}
            onReanalyze={handleReanalyze}
            repositories={filteredRepositories}
          />
        )}
      </section>
    </div>
  );
};
