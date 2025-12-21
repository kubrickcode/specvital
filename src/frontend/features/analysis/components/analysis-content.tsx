"use client";

import { ExternalLink, GitCommit } from "lucide-react";
import { useTranslations } from "next-intl";

import type { AnalysisResult } from "@/lib/api";
import { formatAnalysisDate, SHORT_SHA_LENGTH } from "@/lib/utils";

import { ShareButton } from "./share-button";
import { StatsCard } from "./stats-card";
import { TestList } from "./test-list";

type AnalysisContentProps = {
  result: AnalysisResult;
};

export const AnalysisContent = ({ result }: AnalysisContentProps) => {
  const t = useTranslations("analyze");

  return (
    <main className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        <header className="space-y-2">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold">
              {result.owner}/{result.repo}
            </h1>
            <div className="flex items-center gap-2">
              <ShareButton />
              <a
                className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
                href={`https://github.com/${result.owner}/${result.repo}`}
                rel="noopener noreferrer"
                target="_blank"
              >
                {t("viewOnGitHub")}
                <ExternalLink className="h-4 w-4" />
              </a>
            </div>
          </div>
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
            <span className="flex items-center gap-1">
              <GitCommit className="h-4 w-4" />
              {result.commitSha.slice(0, SHORT_SHA_LENGTH)}
            </span>
            <span>{t("analyzedAt", { date: formatAnalysisDate(result.analyzedAt) })}</span>
          </div>
        </header>

        <StatsCard summary={result.summary} />

        <section className="space-y-4">
          <h2 className="text-xl font-semibold">{t("testSuites")}</h2>
          <TestList suites={result.suites} />
        </section>
      </div>
    </main>
  );
};
