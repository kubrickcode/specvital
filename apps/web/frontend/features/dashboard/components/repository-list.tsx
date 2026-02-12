"use client";

import type { RepositoryCard as RepositoryCardType } from "@/lib/api/types";

import { RepositoryCard } from "./repository-card";
import { RepositoryGrid } from "./repository-grid";
import { RepositorySkeleton } from "./repository-skeleton";

type RepositoryListVariant = "dashboard" | "explore";

type RepositoryListProps = {
  isLoading?: boolean;
  onBookmarkToggle?: (owner: string, repo: string, isBookmarked: boolean) => void;
  onReanalyze?: (owner: string, repo: string) => void;
  repositories: RepositoryCardType[];
  variant?: RepositoryListVariant;
};

export const RepositoryList = ({
  isLoading = false,
  onBookmarkToggle,
  onReanalyze,
  repositories,
  variant = "dashboard",
}: RepositoryListProps) => {
  if (isLoading) {
    return (
      <RepositoryGrid ariaLabel="Loading repositories" isLoading>
        {Array.from({ length: 12 }).map((_, index) => (
          <li key={index}>
            <RepositorySkeleton />
          </li>
        ))}
      </RepositoryGrid>
    );
  }

  if (repositories.length === 0) {
    return null;
  }

  return (
    <RepositoryGrid ariaLabel="Repository list">
      {repositories.map((repo) => (
        <li key={repo.id}>
          <RepositoryCard
            onBookmarkToggle={onBookmarkToggle}
            onReanalyze={onReanalyze}
            repo={repo}
            variant={variant}
          />
        </li>
      ))}
    </RepositoryGrid>
  );
};
