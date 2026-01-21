"use client";

import { ExternalLink } from "lucide-react";
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
    <h1 className="text-2xl font-semibold tracking-tight sm:text-3xl truncate" ref={ref}>
      {fullName}
    </h1>
  );

  return (
    <motion.header className="py-4" variants={fadeInUp}>
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        {/* Repository info */}
        <div className="space-y-1.5 min-w-0">
          {isTruncated ? (
            <ResponsiveTooltip content={fullName}>{titleElement}</ResponsiveTooltip>
          ) : (
            titleElement
          )}

          {/* Metadata line - simplified */}
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            {branchName && (
              <>
                <span className="font-medium text-foreground">{branchName}</span>
                <span className="text-border">·</span>
              </>
            )}
            <ResponsiveTooltip
              content={
                <div className="space-y-1 text-xs">
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
                className="font-mono text-xs hover:text-foreground transition-colors"
                type="button"
              >
                {commitSha.slice(0, SHORT_SHA_LENGTH)}
              </button>
            </ResponsiveTooltip>
          </div>
        </div>

        {/* Action buttons - more subtle */}
        <div className="flex items-center gap-1.5 shrink-0">
          {data && <ExportButton data={data} />}
          <ShareButton />
          <Button asChild size="sm" variant="ghost">
            <a
              href={`https://github.com/${owner}/${repo}`}
              rel="noopener noreferrer"
              target="_blank"
            >
              <span className="sr-only sm:not-sr-only sm:mr-1.5">{t("viewOnGitHub")}</span>
              <ExternalLink className="h-4 w-4" />
            </a>
          </Button>
        </div>
      </div>
    </motion.header>
  );
};
