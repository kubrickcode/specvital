"use client";

import { ChevronDown, Github, Info } from "lucide-react";
import { motion } from "motion/react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";
import { useTruncateDetection } from "@/lib/hooks";
import { fadeInUp } from "@/lib/motion";
import { cn, formatAnalysisDate } from "@/lib/utils";

import { CommitSelector } from "./commit-selector";
import { ShareButton } from "./share-button";

type AnalysisHeaderProps = {
  analyzedAt: string;
  branchName?: string;
  commitSha: string;
  committedAt?: string;
  onCommitSelect?: (commitSha: string | null) => void;
  owner: string;
  parserVersion?: string;
  repo: string;
};

export const AnalysisHeader = ({
  analyzedAt,
  branchName,
  commitSha,
  committedAt,
  onCommitSelect,
  owner,
  parserVersion,
  repo,
}: AnalysisHeaderProps) => {
  const t = useTranslations("analyze");
  const { isTruncated, ref } = useTruncateDetection<HTMLHeadingElement>();
  const [isOpen, setIsOpen] = useState(false);
  const fullName = `${owner}/${repo}`;

  const titleElement = (
    <h1 className="text-2xl font-semibold tracking-tight sm:text-3xl truncate" ref={ref}>
      {fullName}
    </h1>
  );

  return (
    <motion.header className="py-4" variants={fadeInUp}>
      <Collapsible onOpenChange={setIsOpen} open={isOpen}>
        {/* Repository name */}
        <div className="min-w-0">
          {isTruncated ? (
            <ResponsiveTooltip content={fullName}>{titleElement}</ResponsiveTooltip>
          ) : (
            titleElement
          )}
        </div>

        {/* Details button and Action buttons in one row */}
        <div className="mt-1.5 flex items-center justify-between gap-2">
          <CollapsibleTrigger asChild>
            <button
              aria-label={t("details")}
              className="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded px-2 py-1 -ml-2"
              type="button"
            >
              <Info className="h-3.5 w-3.5" />
              {t("details")}
              <ChevronDown
                className={cn("h-3.5 w-3.5 transition-transform", isOpen && "rotate-180")}
              />
            </button>
          </CollapsibleTrigger>

          {/* Action buttons */}
          <div className="flex items-center gap-1.5 shrink-0">
            <ShareButton />
            <Button asChild size="sm" variant="outline">
              <a
                href={`https://github.com/${owner}/${repo}`}
                rel="noopener noreferrer"
                target="_blank"
              >
                <Github className="h-4 w-4" />
                <span className="sr-only sm:not-sr-only sm:ml-1.5">{t("viewOnGitHub")}</span>
              </a>
            </Button>
          </div>
        </div>

        {/* Details panel */}
        <CollapsibleContent>
          <div className="mt-2 space-y-1 text-xs text-muted-foreground pl-5 border-l-2 border-border">
            {branchName && (
              <div>
                {t("branch")}: <span className="font-medium">{branchName}</span>
              </div>
            )}
            <div className="flex items-center gap-1">
              <span>{t("commit")}:</span>
              <CommitSelector
                currentCommitSha={commitSha}
                onCommitSelect={onCommitSelect}
                owner={owner}
                repo={repo}
              />
            </div>
            {committedAt && (
              <div>{t("committedAt", { date: formatAnalysisDate(committedAt) })}</div>
            )}
            <div>{t("analyzedAt", { date: formatAnalysisDate(analyzedAt) })}</div>
            {parserVersion && <div>{t("parserVersion", { version: parserVersion })}</div>}
          </div>
        </CollapsibleContent>
      </Collapsible>
    </motion.header>
  );
};
