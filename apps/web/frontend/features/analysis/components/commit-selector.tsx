"use client";

import { Check, ChevronDown, GitCommit, History } from "lucide-react";
import { useLocale, useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { cn } from "@/lib/utils";

import type { AnalysisHistoryItem } from "../api";
import { useAnalysisHistory } from "../hooks";

type CommitSelectorProps = {
  currentCommitSha: string;
  onCommitSelect?: (commitSha: string | null) => void;
  owner: string;
  repo: string;
};

const formatShortDate = (dateString: string, locale: string): string => {
  return new Date(dateString).toLocaleDateString(locale, {
    day: "numeric",
    month: "short",
  });
};

const formatCommitSha = (sha: string): string => (sha.length < 7 ? sha : sha.slice(0, 7));

export const CommitSelector = ({
  currentCommitSha,
  onCommitSelect,
  owner,
  repo,
}: CommitSelectorProps) => {
  const locale = useLocale();
  const t = useTranslations("analyze");
  const { data: history, error, isLoading } = useAnalysisHistory(owner, repo);

  // Show static commit SHA during loading, on error, or when only one commit exists
  if (isLoading || error || !history || history.length <= 1) {
    return <span className="font-mono font-medium">{formatCommitSha(currentCommitSha)}</span>;
  }

  const currentItem = history.find((item) => item.commitSha === currentCommitSha);
  const isCurrentHead = currentItem?.isHead ?? false;

  const handleCommitSelect = (commitSha: string) => {
    if (commitSha === currentCommitSha) return;
    onCommitSelect?.(commitSha);
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          className="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded"
          type="button"
        >
          <span className="font-mono font-medium">{formatCommitSha(currentCommitSha)}</span>
          {isCurrentHead && (
            <Badge className="h-4 px-1 text-[10px] font-normal" variant="secondary">
              HEAD
            </Badge>
          )}
          <ChevronDown className="h-3 w-3 opacity-50" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" className="w-72 max-h-80 overflow-y-auto">
        <DropdownMenuLabel className="flex items-center gap-2 px-2 py-2 text-[11px] font-medium uppercase tracking-wider text-muted-foreground/70 bg-muted/40 -mx-1 mb-1 border-b border-border/50">
          <History className="h-3 w-3" />
          {t("commitHistory.title")}
        </DropdownMenuLabel>
        {history.map((item) => (
          <CommitItem
            currentCommitSha={currentCommitSha}
            item={item}
            key={item.id}
            locale={locale}
            onSelect={handleCommitSelect}
            t={t}
          />
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

type CommitItemProps = {
  currentCommitSha: string;
  item: AnalysisHistoryItem;
  locale: string;
  onSelect: (commitSha: string) => void;
  t: ReturnType<typeof useTranslations<"analyze">>;
};

const CommitItem = ({ currentCommitSha, item, locale, onSelect, t }: CommitItemProps) => {
  const isSelected = item.commitSha === currentCommitSha;

  return (
    <DropdownMenuItem
      className={cn("flex flex-col items-start gap-0.5 py-2", isSelected && "bg-muted")}
      onClick={() => onSelect(item.commitSha)}
    >
      <div className="flex items-center gap-2 w-full">
        <div className="flex items-center gap-1.5">
          {isSelected ? (
            <Check className="h-3.5 w-3.5 text-primary shrink-0" />
          ) : (
            <GitCommit className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
          )}
          <code className="font-mono text-xs font-medium">{formatCommitSha(item.commitSha)}</code>
        </div>
        <div className="flex items-center gap-1.5 ml-auto">
          {item.isHead && (
            <Badge className="h-4 px-1 text-[10px] font-normal" variant="secondary">
              HEAD
            </Badge>
          )}
          <span className="text-xs text-muted-foreground">
            {formatShortDate(item.completedAt, locale)}
          </span>
        </div>
      </div>
      <div className="flex items-center gap-2 pl-5 text-xs text-muted-foreground">
        <span>{t("commitHistory.tests", { count: item.totalTests })}</span>
        {item.branchName && (
          <>
            <span className="text-muted-foreground/50">·</span>
            <span>{item.branchName}</span>
          </>
        )}
      </div>
    </DropdownMenuItem>
  );
};
