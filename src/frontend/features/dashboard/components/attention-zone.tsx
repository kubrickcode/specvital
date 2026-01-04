"use client";

import { AlertCircle, ChevronRight, RefreshCw } from "lucide-react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useMemo } from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import type { RepositoryCard } from "@/lib/api/types";
import { cn } from "@/lib/utils";

import { SeverityBadge, type AttentionLabelKey, type SeverityLevel } from "./severity-badge";

type AttentionItem = {
  action: "reanalyze" | "view";
  labelKey: AttentionLabelKey;
  repo: RepositoryCard;
  severity: SeverityLevel;
};

type AttentionZoneProps = {
  onReanalyze?: (owner: string, repo: string) => void;
  repositories: RepositoryCard[];
};

const SKIPPED_RATIO_THRESHOLD = 0.3;
const MAX_DISPLAY_ITEMS = 3;

const getAttentionItems = (repositories: RepositoryCard[]): AttentionItem[] => {
  const items: AttentionItem[] = [];

  for (const repo of repositories) {
    if (!repo.latestAnalysis?.testSummary) continue;

    const { testSummary } = repo.latestAnalysis;

    if (testSummary.focused > 0) {
      items.push({
        action: "view",
        labelKey: "focused",
        repo,
        severity: "critical",
      });
      continue;
    }

    if (repo.updateStatus === "new-commits") {
      items.push({
        action: "reanalyze",
        labelKey: "newCommits",
        repo,
        severity: "warning",
      });
      continue;
    }

    const total = repo.latestAnalysis.testCount ?? 0;
    if (total > 0 && testSummary.skipped / total > SKIPPED_RATIO_THRESHOLD) {
      items.push({
        action: "view",
        labelKey: "skippedHigh",
        repo,
        severity: "info",
      });
    }
  }

  return items.sort((a, b) => {
    const severityOrder: Record<SeverityLevel, number> = { critical: 0, info: 2, warning: 1 };
    return severityOrder[a.severity] - severityOrder[b.severity];
  });
};

type AttentionCardProps = {
  item: AttentionItem;
  onReanalyze?: (owner: string, repo: string) => void;
};

const AttentionCard = ({ item, onReanalyze }: AttentionCardProps) => {
  const t = useTranslations("dashboard.attention");
  const { action, labelKey, repo, severity } = item;

  const handleReanalyzeClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    onReanalyze?.(repo.owner, repo.name);
  };

  return (
    <Link
      className="group block outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 rounded-lg"
      href={`/analyze/${repo.owner}/${repo.name}`}
    >
      <Card
        className={cn(
          "h-full",
          "hover:border-primary/20 hover:-translate-y-0.5",
          "group-focus-visible:shadow-md group-focus-visible:border-primary/20"
        )}
        depth="raised"
      >
        <CardContent className="p-4 space-y-3">
          <div className="flex items-start justify-between gap-2">
            <h4 className="font-medium text-sm truncate flex-1" title={repo.fullName}>
              {repo.fullName}
            </h4>
            <SeverityBadge labelKey={labelKey} level={severity} />
          </div>

          <p className="text-xs text-muted-foreground">{t(`${labelKey}Description`)}</p>

          {action === "reanalyze" && onReanalyze && (
            <Button
              className="h-7 px-2 text-xs w-full"
              onClick={handleReanalyzeClick}
              size="sm"
              variant="outline"
            >
              <RefreshCw aria-hidden="true" className="size-3 mr-1" />
              {t("reanalyze")}
            </Button>
          )}

          {action === "view" && (
            <Button asChild className="h-7 px-2 text-xs w-full" size="sm" variant="outline">
              <span>
                {t("viewDetails")}
                <ChevronRight aria-hidden="true" className="size-3 ml-1" />
              </span>
            </Button>
          )}
        </CardContent>
      </Card>
    </Link>
  );
};

export const AttentionZone = ({ onReanalyze, repositories }: AttentionZoneProps) => {
  const t = useTranslations("dashboard.attention");

  const attentionItems = useMemo(() => getAttentionItems(repositories), [repositories]);

  if (attentionItems.length === 0) {
    return null;
  }

  const displayItems = attentionItems.slice(0, MAX_DISPLAY_ITEMS);
  const remainingCount = attentionItems.length - MAX_DISPLAY_ITEMS;

  return (
    <section aria-labelledby="attention-zone-title" className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <AlertCircle aria-hidden="true" className="size-5 text-amber-500" />
          <h2 className="text-lg font-semibold" id="attention-zone-title">
            {t("title")}
          </h2>
          <span className="text-sm text-muted-foreground" role="status">
            <span className="sr-only">{t("itemCount", { count: attentionItems.length })}</span>
            <span aria-hidden="true">({attentionItems.length})</span>
          </span>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {displayItems.map((item) => (
          <AttentionCard
            item={item}
            key={`${item.repo.id}-${item.labelKey}`}
            onReanalyze={onReanalyze}
          />
        ))}

        {remainingCount > 0 && (
          <Card className="flex items-center justify-center bg-muted/50">
            <CardContent className="p-4 text-center">
              <p className="text-sm text-muted-foreground">
                {t("moreItems", { count: remainingCount })}
              </p>
            </CardContent>
          </Card>
        )}
      </div>
    </section>
  );
};
