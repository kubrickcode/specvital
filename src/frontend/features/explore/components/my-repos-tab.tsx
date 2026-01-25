"use client";

import { RefreshCw, Search } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { RepositoryCard, useMyRepositories, usePaginatedRepositories } from "@/features/dashboard";
import type { GitHubRepository } from "@/lib/api/types";

type MyReposTabProps = {
  className?: string;
};

export const MyReposTab = ({ className }: MyReposTabProps) => {
  const t = useTranslations("explore.myRepos");
  const [searchQuery, setSearchQuery] = useState("");

  const { data: myRepos, error, isLoading, isRefreshing, refresh } = useMyRepositories();

  const { data: analyzedRepos, isLoading: isLoadingAnalyzed } = usePaginatedRepositories({
    view: "my",
  });

  const analyzedRepoSet = (() => {
    const set = new Set<string>();
    for (const repo of analyzedRepos) {
      set.add(`${repo.owner}/${repo.name}`);
    }
    return set;
  })();

  const filteredRepos = (() => {
    if (!searchQuery.trim()) {
      return myRepos;
    }
    const query = searchQuery.toLowerCase();
    return myRepos.filter(
      (repo) =>
        repo.name.toLowerCase().includes(query) ||
        repo.fullName.toLowerCase().includes(query) ||
        repo.description?.toLowerCase().includes(query)
    );
  })();

  const isRepoAnalyzed = (repo: GitHubRepository): boolean => {
    return analyzedRepoSet.has(repo.fullName);
  };

  const handleRefresh = () => {
    refresh();
  };

  if (isLoading) {
    return (
      <div className={className}>
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Skeleton className="h-10 flex-1" />
            <Skeleton className="h-10 w-10" />
          </div>
          <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, index) => (
              <li key={index}>
                <Skeleton className="h-48 w-full rounded-xl" />
              </li>
            ))}
          </ul>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={className}>
        <div className="text-center py-12">
          <p className="text-sm text-muted-foreground mb-4">{error.message}</p>
          <Button onClick={handleRefresh} size="sm" variant="outline">
            <RefreshCw className="mr-2 size-4" />
            {t("refresh")}
          </Button>
        </div>
      </div>
    );
  }

  if (myRepos.length === 0) {
    return (
      <div className={className}>
        <div className="text-center py-12">
          <p className="text-sm font-medium mb-1">{t("noRepos")}</p>
          <p className="text-sm text-muted-foreground">{t("noReposDescription")}</p>
        </div>
      </div>
    );
  }

  return (
    <div className={className}>
      <div className="space-y-4">
        <div className="flex items-center gap-2">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
            <Input
              className="pl-9"
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder={t("searchPlaceholder")}
              type="search"
              value={searchQuery}
            />
          </div>
          <Button disabled={isRefreshing} onClick={handleRefresh} size="icon" variant="outline">
            <RefreshCw className={`size-4 ${isRefreshing ? "animate-spin" : ""}`} />
            <span className="sr-only">{t("refresh")}</span>
          </Button>
        </div>

        <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {filteredRepos.map((repo) => {
            const analyzed = !isLoadingAnalyzed && isRepoAnalyzed(repo);

            return (
              <li key={repo.id}>
                <RepositoryCard isAnalyzed={analyzed} repo={repo} variant="unanalyzed" />
              </li>
            );
          })}

          {filteredRepos.length === 0 && searchQuery && (
            <li className="col-span-full text-center py-8 text-sm text-muted-foreground">
              {t("noSearchResults")}
            </li>
          )}
        </ul>
      </div>
    </div>
  );
};
