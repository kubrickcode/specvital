"use client";

import { ExternalLink, GitBranch, GitCommit, Info } from "lucide-react";
import { motion } from "motion/react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { fadeInUp } from "@/lib/motion";
import { formatAnalysisDate, SHORT_SHA_LENGTH } from "@/lib/utils";

import { ShareButton } from "./share-button";

type AnalysisHeaderProps = {
  analyzedAt: string;
  branchName?: string;
  commitSha: string;
  committedAt?: string;
  owner: string;
  repo: string;
};

export const AnalysisHeader = ({
  analyzedAt,
  branchName,
  commitSha,
  committedAt,
  owner,
  repo,
}: AnalysisHeaderProps) => {
  const t = useTranslations("analyze");

  return (
    <motion.header variants={fadeInUp}>
      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        {/* 저장소 정보 세트 */}
        <div className="space-y-1 min-w-0">
          <h1 className="text-xl font-bold sm:text-2xl truncate">
            {owner}/{repo}
          </h1>
          <div className="flex items-center gap-4 text-sm text-muted-foreground [&_svg]:translate-y-[2px]">
            {branchName && (
              <span className="inline-flex items-center gap-1.5">
                <GitBranch className="size-4" />
                <span className="font-medium text-foreground">{branchName}</span>
              </span>
            )}
            <Popover>
              <PopoverTrigger asChild>
                <button
                  className="inline-flex items-center gap-1.5 hover:text-foreground transition-colors"
                  type="button"
                >
                  <GitCommit className="size-4" />
                  <span>{commitSha.slice(0, SHORT_SHA_LENGTH)}</span>
                  <Info className="size-3.5 opacity-50" />
                </button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-3 text-sm" side="bottom">
                <div className="space-y-1.5">
                  {committedAt && (
                    <div>{t("committedAt", { date: formatAnalysisDate(committedAt) })}</div>
                  )}
                  <div>{t("analyzedAt", { date: formatAnalysisDate(analyzedAt) })}</div>
                </div>
              </PopoverContent>
            </Popover>
          </div>
        </div>

        {/* 액션 버튼 세트 */}
        <div className="flex items-center gap-2 shrink-0">
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
