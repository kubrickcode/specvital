"use client";

import { ExternalLink, GitBranch, GitCommit, Info } from "lucide-react";
import { motion } from "motion/react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";
import type { components } from "@/lib/api/generated-types";
import { useTruncateDetection } from "@/lib/hooks";
import { fadeInUp } from "@/lib/motion";
import { formatAnalysisDate, SHORT_SHA_LENGTH } from "@/lib/utils";

import { ExportButton } from "./export-button";
import { ShareButton } from "./share-button";

type AnalysisResult = components["schemas"]["AnalysisResult"];

type AnalysisHeaderProps = {
  analyzedAt: string;
  branchName?: string;
  commitSha: string;
  committedAt?: string;
  data?: AnalysisResult;
  owner: string;
  parserVersion?: string;
  repo: string;
};

export const AnalysisHeader = ({
  analyzedAt,
  branchName,
  commitSha,
  committedAt,
  data,
  owner,
  parserVersion,
  repo,
}: AnalysisHeaderProps) => {
  const t = useTranslations("analyze");
  const { isTruncated, ref } = useTruncateDetection<HTMLHeadingElement>();
  const fullName = `${owner}/${repo}`;

  const titleElement = (
    <h1 className="text-xl font-bold sm:text-2xl truncate" ref={ref}>
      {fullName}
    </h1>
  );

  return (
    <motion.header variants={fadeInUp}>
      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        {/* Repository info section */}
        <div className="space-y-1 min-w-0">
          {isTruncated ? (
            <ResponsiveTooltip content={fullName}>{titleElement}</ResponsiveTooltip>
          ) : (
            titleElement
          )}
          <div className="flex items-center gap-4 text-sm text-muted-foreground [&_svg]:translate-y-[2px]">
            {branchName && (
              <span className="inline-flex items-center gap-1.5">
                <GitBranch className="size-4" />
                <span className="font-medium text-foreground">{branchName}</span>
              </span>
            )}
            <ResponsiveTooltip
              content={
                <div className="space-y-1">
                  {committedAt && (
                    <div>{t("committedAt", { date: formatAnalysisDate(committedAt) })}</div>
                  )}
                  <div>{t("analyzedAt", { date: formatAnalysisDate(analyzedAt) })}</div>
                  {parserVersion && <div>{t("parserVersion", { version: parserVersion })}</div>}
                </div>
              }
              side="bottom"
            >
              <button
                className="inline-flex items-center gap-1.5 hover:text-foreground transition-colors"
                type="button"
              >
                <GitCommit className="size-4" />
                <span>{commitSha.slice(0, SHORT_SHA_LENGTH)}</span>
                <Info className="size-3.5 opacity-50" />
              </button>
            </ResponsiveTooltip>
          </div>
        </div>

        {/* Action buttons section */}
        <div className="flex items-center gap-2 shrink-0">
          {data && <ExportButton data={data} />}
          <ShareButton />
          <Button asChild size="sm" variant="outline">
            <a
              href={`https://github.com/${owner}/${repo}`}
              rel="noopener noreferrer"
              target="_blank"
            >
              {t("viewOnGitHub")}
              <ExternalLink className="h-4 w-4" />
            </a>
          </Button>
        </div>
      </div>
    </motion.header>
  );
};
