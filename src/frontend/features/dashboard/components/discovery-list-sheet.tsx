"use client";

import { ExternalLink, Lock, Search } from "lucide-react";
import Link from "next/link";
import { useFormatter, useNow, useTranslations } from "next-intl";
import { useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import type { GitHubRepository, OrganizationAccessStatus } from "@/lib/api/types";

import { OrgConnectionBanner } from "./org-connection-banner";

type DiscoveryListSheetProps = {
  accessStatus?: OrganizationAccessStatus;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  repositories: GitHubRepository[];
  title: string;
};

export const DiscoveryListSheet = ({
  accessStatus,
  isOpen,
  onOpenChange,
  repositories,
  title,
}: DiscoveryListSheetProps) => {
  const t = useTranslations("dashboard.discovery");
  const format = useFormatter();
  const now = useNow({ updateInterval: 1000 * 60 });
  const [searchQuery, setSearchQuery] = useState("");

  const filteredRepos = useMemo(() => {
    if (!searchQuery.trim()) return repositories;
    const query = searchQuery.toLowerCase();
    return repositories.filter(
      (repo) =>
        repo.name.toLowerCase().includes(query) ||
        repo.fullName.toLowerCase().includes(query) ||
        repo.description?.toLowerCase().includes(query)
    );
  }, [repositories, searchQuery]);

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      setSearchQuery("");
    }
    onOpenChange(open);
  };

  return (
    <Sheet onOpenChange={handleOpenChange} open={isOpen}>
      <SheetContent className="w-full sm:max-w-lg overflow-y-auto" side="right">
        <SheetHeader className="mb-4">
          <SheetTitle>{title}</SheetTitle>
          <SheetDescription>
            {t("sheetDescription", { count: repositories.length })}
          </SheetDescription>
        </SheetHeader>

        <div className="space-y-4 px-4 pb-4">
          {accessStatus && accessStatus !== "accessible" && (
            <OrgConnectionBanner accessStatus={accessStatus} orgLogin={title} />
          )}

          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
            <Input
              className="pl-9"
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder={t("searchPlaceholder")}
              type="search"
              value={searchQuery}
            />
          </div>

          <ul className="space-y-2">
            {filteredRepos.map((repo) => (
              <li key={repo.id}>
                <div className="flex items-start gap-3 p-3 rounded-lg border bg-card shadow-sm hover:shadow-md hover:bg-accent/50 transition-all">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <h4 className="font-medium text-sm truncate">{repo.fullName}</h4>
                      {repo.isPrivate && <Lock className="size-3 text-muted-foreground shrink-0" />}
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
                    <Button asChild size="sm" variant="cta">
                      <Link href={`/analyze/${repo.owner}/${repo.name}`}>{t("analyze")}</Link>
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
            ))}

            {filteredRepos.length === 0 && (
              <li className="text-center py-8 text-sm text-muted-foreground">
                {searchQuery ? t("noSearchResults") : t("noRepos")}
              </li>
            )}
          </ul>
        </div>
      </SheetContent>
    </Sheet>
  );
};
