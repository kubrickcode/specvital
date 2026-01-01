"use client";

import { ExternalLink, Lock, RefreshCw, Search } from "lucide-react";
import { useFormatter, useNow, useTranslations } from "next-intl";
import { useCallback, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { useMyRepositories, usePaginatedRepositories } from "@/features/dashboard";
import { Link } from "@/i18n/navigation";
import type { GitHubRepository } from "@/lib/api/types";

type MyReposTabProps = {
  className?: string;
};

export const MyReposTab = ({ className }: MyReposTabProps) => {
  const t = useTranslations("explore.myRepos");
  const format = useFormatter();
  const now = useNow({ updateInterval: 1000 * 60 });
  const [searchQuery, setSearchQuery] = useState("");

  const { data: myRepos, error, isLoading, isRefreshing, refresh } = useMyRepositories();

  const { data: analyzedRepos, isLoading: isLoadingAnalyzed } = usePaginatedRepositories({
    view: "my",
  });

  const analyzedRepoSet = useMemo(() => {
    const set = new Set<string>();
    for (const repo of analyzedRepos) {
      set.add(`${repo.owner}/${repo.name}`);
    }
    return set;
  }, [analyzedRepos]);

  const filteredRepos = useMemo(() => {
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
  }, [myRepos, searchQuery]);

  const isRepoAnalyzed = (repo: GitHubRepository): boolean => {
    return analyzedRepoSet.has(repo.fullName);
  };

  const handleRefresh = useCallback(() => {
    refresh();
  }, [refresh]);

  if (isLoading) {
    return (
      <div className={className}>
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Skeleton className="h-10 flex-1" />
            <Skeleton className="h-10 w-10" />
          </div>
          <ul className="space-y-2">
            {Array.from({ length: 5 }).map((_, index) => (
              <li key={index}>
                <Skeleton className="h-24 w-full rounded-lg" />
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

        <ul className="space-y-2">
          {filteredRepos.map((repo) => {
            const analyzed = !isLoadingAnalyzed && isRepoAnalyzed(repo);

            return (
              <li key={repo.id}>
                <div className="flex items-start gap-3 p-3 rounded-lg border bg-card shadow-sm hover:shadow-md hover:bg-accent/50 transition-all">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <h4 className="font-medium text-sm truncate">{repo.fullName}</h4>
                      {repo.isPrivate && <Lock className="size-3 text-muted-foreground shrink-0" />}
                      {!isLoadingAnalyzed && (
                        <span
                          className={`text-xs px-1.5 py-0.5 rounded ${
                            analyzed
                              ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400"
                              : "bg-muted text-muted-foreground"
                          }`}
                        >
                          {analyzed ? t("analyzed") : t("notAnalyzed")}
                        </span>
                      )}
                    </div>
                    {repo.description && (
                      <p className="text-xs text-muted-foreground line-clamp-2 mt-1">
                        {repo.description}
                      </p>
                    )}
                    {repo.pushedAt && (
                      <p className="text-xs text-muted-foreground mt-1">
                        {t("lastPush")}: {format.relativeTime(new Date(repo.pushedAt), now)}
                      </p>
                    )}
                  </div>
                  <div className="flex flex-col gap-1 shrink-0">
                    <Button asChild size="sm" variant={analyzed ? "outline" : "cta"}>
                      <Link href={`/analyze/${repo.owner}/${repo.name}`}>
                        {analyzed ? t("view") : t("analyze")}
                      </Link>
                    </Button>
                    <Button asChild className="h-7" size="sm" variant="ghost">
                      <a
                        href={`https://github.com/${repo.fullName}`}
                        rel="noopener noreferrer"
                        target="_blank"
                      >
                        <ExternalLink className="size-3 mr-1" />
                        GitHub
                      </a>
                    </Button>
                  </div>
                </div>
              </li>
            );
          })}

          {filteredRepos.length === 0 && searchQuery && (
            <li className="text-center py-8 text-sm text-muted-foreground">
              {t("noSearchResults")}
            </li>
          )}
        </ul>
      </div>
    </div>
  );
};
