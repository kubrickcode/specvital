"use client";

import { RefreshCw, Star } from "lucide-react";
import Link from "next/link";
import { useFormatter, useNow, useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";
import type { RepositoryCard as RepositoryCardType } from "@/lib/api/types";
import { cn } from "@/lib/utils";

import { TestDeltaBadge } from "./test-delta-badge";
import { UpdateStatusBadge } from "./update-status-badge";

type RepositoryCardProps = {
  onBookmarkToggle?: (owner: string, repo: string, isBookmarked: boolean) => void;
  onReanalyze?: (owner: string, repo: string) => void;
  repo: RepositoryCardType;
};

export const RepositoryCard = ({ onBookmarkToggle, onReanalyze, repo }: RepositoryCardProps) => {
  const format = useFormatter();
  const now = useNow({ updateInterval: 60_000 });
  const t = useTranslations("dashboard.card");
  const { fullName, isBookmarked, latestAnalysis, name, owner, updateStatus } = repo;
  const hasNewCommits = updateStatus === "new-commits";
  const hasAnalysis = Boolean(latestAnalysis);

  const handleBookmarkClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    onBookmarkToggle?.(owner, name, isBookmarked);
  };

  const handleReanalyzeClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    onReanalyze?.(owner, name);
  };

  return (
    <Link
      className="group block outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 rounded-xl"
      href={`/analyze/${owner}/${name}`}
    >
      <Card
        className={cn(
          "relative h-full p-4 transition-all duration-200",
          "hover:shadow-md hover:border-primary/20",
          "group-focus-visible:shadow-md group-focus-visible:border-primary/20"
        )}
      >
        <div className="flex items-start justify-between gap-2 mb-3">
          <div className="min-w-0 flex-1">
            <h3 className="font-semibold text-sm truncate" title={fullName}>
              {fullName}
            </h3>
          </div>
          <ResponsiveTooltip content={isBookmarked ? t("removeBookmark") : t("addBookmark")}>
            <Button
              aria-label={isBookmarked ? t("removeBookmark") : t("addBookmark")}
              aria-pressed={isBookmarked}
              className={cn(
                "size-8 shrink-0",
                isBookmarked && "text-amber-500 hover:text-amber-600"
              )}
              onClick={handleBookmarkClick}
              size="icon"
              variant="ghost"
            >
              <Star aria-hidden="true" className={cn("size-4", isBookmarked && "fill-current")} />
            </Button>
          </ResponsiveTooltip>
        </div>

        {hasAnalysis && latestAnalysis ? (
          <div className="space-y-3">
            <div className="flex items-baseline gap-2">
              <span className="text-2xl font-bold">{latestAnalysis.testCount}</span>
              <span className="text-sm text-muted-foreground">{t("tests")}</span>
              {latestAnalysis.change !== 0 && (
                <TestDeltaBadge delta={latestAnalysis.change} isCompact />
              )}
            </div>

            <div className="flex items-center gap-2">
              <UpdateStatusBadge status={updateStatus} />
            </div>

            <div className="flex items-center justify-between pt-2 border-t">
              <ResponsiveTooltip
                content={format.dateTime(new Date(latestAnalysis.analyzedAt), {
                  dateStyle: "medium",
                  timeStyle: "short",
                })}
              >
                <time
                  className="text-xs text-muted-foreground"
                  dateTime={latestAnalysis.analyzedAt}
                >
                  {format.relativeTime(new Date(latestAnalysis.analyzedAt), now)}
                </time>
              </ResponsiveTooltip>

              {hasNewCommits && (
                <ResponsiveTooltip content={t("analyzeLatest")}>
                  <Button
                    aria-label={t("reanalyze")}
                    className="h-7 px-2 text-xs"
                    onClick={handleReanalyzeClick}
                    size="sm"
                    variant="outline"
                  >
                    <RefreshCw aria-hidden="true" className="size-3 mr-1" />
                    {t("update")}
                  </Button>
                </ResponsiveTooltip>
              )}
            </div>
          </div>
        ) : (
          <div className="flex items-center justify-center h-20 text-sm text-muted-foreground">
            {t("noAnalysis")}
          </div>
        )}
      </Card>
    </Link>
  );
};
