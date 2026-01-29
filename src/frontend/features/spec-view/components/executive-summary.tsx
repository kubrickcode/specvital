"use client";

import {
  Check,
  ChevronDown,
  FileText,
  Globe,
  History,
  Languages,
  Plus,
  RefreshCw,
} from "lucide-react";
import { useLocale, useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

import { CacheStatsIndicator } from "./cache-stats-indicator";
import { SpecExportButton } from "./spec-export-button";
import { StatusLegend } from "./status-legend";
import { SPEC_LANGUAGES } from "../constants/spec-languages";
import type {
  BehaviorCacheStats,
  RepoVersionInfo,
  SpecDocument,
  SpecLanguage,
  VersionInfo,
} from "../types";
import { calculateDocumentStats } from "../utils/stats";

type VersionInfoWithCommit = VersionInfo | RepoVersionInfo;

type ExecutiveSummaryProps = {
  behaviorCacheStats?: BehaviorCacheStats;
  commitSha?: string;
  document: SpecDocument;
  isGeneratingOtherLanguage?: boolean;
  isLoadingVersions?: boolean;
  isRegenerating?: boolean;
  latestVersion?: number;
  onGenerateNewLanguage?: (language: SpecLanguage) => void;
  onLanguageSwitch?: (language: SpecLanguage) => void;
  onRegenerate?: () => void;
  onVersionSwitch?: (version: number) => void;
  owner?: string;
  repo?: string;
  versions?: VersionInfoWithCommit[];
};

const formatShortDate = (dateString: string, locale: string): string => {
  return new Date(dateString).toLocaleDateString(locale, {
    day: "numeric",
    month: "short",
  });
};

const formatCommitSha = (sha: string): string => (sha.length < 7 ? sha : sha.slice(0, 7));

export const ExecutiveSummary = ({
  behaviorCacheStats,
  commitSha,
  document,
  isGeneratingOtherLanguage = false,
  isLoadingVersions = false,
  isRegenerating = false,
  latestVersion,
  onGenerateNewLanguage,
  onLanguageSwitch,
  onRegenerate,
  onVersionSwitch,
  owner,
  repo,
  versions = [],
}: ExecutiveSummaryProps) => {
  const locale = useLocale();
  const t = useTranslations("specView");

  const hasExecutiveSummary = !!document.executiveSummary;
  const { behaviorCount, domainCount, featureCount } = calculateDocumentStats(document);
  const currentLanguage = document.language;
  const currentVersion = document.version;
  const availableLanguages = document.availableLanguages ?? [];
  const isDisabled = isRegenerating || isGeneratingOtherLanguage;
  const isLatestVersion = latestVersion === undefined || currentVersion === latestVersion;

  // Build sets for available vs new languages
  const availableLanguageSet = new Set(availableLanguages.map((l) => l.language));
  const newLanguages = SPEC_LANGUAGES.filter((lang) => !availableLanguageSet.has(lang));

  const handleLanguageSelect = (language: SpecLanguage) => {
    if (language === currentLanguage) return;

    if (availableLanguageSet.has(language)) {
      // Available language → instant switch (free)
      onLanguageSwitch?.(language);
    } else {
      // New language → generate (costs quota)
      onGenerateNewLanguage?.(language);
    }
  };

  const canSwitchLanguage = onLanguageSwitch || onGenerateNewLanguage;

  return (
    <Card className={cn(isGeneratingOtherLanguage && "opacity-60 pointer-events-none")}>
      <CardHeader>
        <div className="flex items-center justify-between flex-wrap gap-2">
          <div className="flex items-center gap-2">
            <FileText className="h-5 w-5 text-primary" />
            <CardTitle className="text-lg">{t("executiveSummary.title")}</CardTitle>
          </div>
          <div className="flex items-center gap-2 flex-wrap">
            {behaviorCacheStats && <CacheStatsIndicator stats={behaviorCacheStats} />}
            {currentLanguage && canSwitchLanguage ? (
              <DropdownMenu>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <DropdownMenuTrigger asChild>
                      <Button
                        className="h-auto px-2.5 py-1 gap-1.5 text-xs font-normal"
                        disabled={isDisabled}
                        variant="outline"
                      >
                        <Languages className="h-3 w-3" />
                        {currentLanguage}
                        <ChevronDown className="h-3 w-3 opacity-50" />
                      </Button>
                    </DropdownMenuTrigger>
                  </TooltipTrigger>
                  <TooltipContent>{t("executiveSummary.switchLanguageTooltip")}</TooltipContent>
                </Tooltip>
                <DropdownMenuContent align="end" className="w-64 max-h-80 overflow-y-auto">
                  {/* Available Languages Section */}
                  {availableLanguages.length > 0 && (
                    <>
                      <DropdownMenuLabel className="flex items-center gap-2 px-2 py-2 text-[11px] font-medium uppercase tracking-wider text-muted-foreground/70 bg-muted/40 -mx-1 mb-1 border-b border-border/50">
                        <Check className="h-3 w-3" />
                        {t("executiveSummary.availableLanguages")}
                      </DropdownMenuLabel>
                      {availableLanguages.map((langInfo) => {
                        const isCurrentLang = langInfo.language === currentLanguage;
                        return (
                          <DropdownMenuItem
                            className={cn(
                              "flex items-center justify-between",
                              isCurrentLang && "bg-muted font-medium"
                            )}
                            key={langInfo.language}
                            onClick={() => handleLanguageSelect(langInfo.language)}
                          >
                            <div className="flex items-center gap-2">
                              {isCurrentLang ? (
                                <Check className="h-3.5 w-3.5 text-primary" />
                              ) : (
                                <span className="w-3.5" />
                              )}
                              <span>{langInfo.language}</span>
                            </div>
                            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                              <span>
                                {t("executiveSummary.versionLabel", {
                                  version: langInfo.latestVersion,
                                })}
                              </span>
                              <span className="text-muted-foreground/60">·</span>
                              <span>{formatShortDate(langInfo.createdAt, locale)}</span>
                            </div>
                          </DropdownMenuItem>
                        );
                      })}
                    </>
                  )}

                  {/* Separator between sections */}
                  {availableLanguages.length > 0 &&
                    newLanguages.length > 0 &&
                    onGenerateNewLanguage && <DropdownMenuSeparator />}

                  {/* Generate New Section */}
                  {newLanguages.length > 0 && onGenerateNewLanguage && (
                    <>
                      <DropdownMenuLabel className="flex items-center gap-2 px-2 py-2 text-[11px] font-medium uppercase tracking-wider text-muted-foreground/70 bg-muted/40 -mx-1 mt-1 mb-1 border-b border-border/50">
                        <Plus className="h-3 w-3" />
                        {t("executiveSummary.generateNew")}
                      </DropdownMenuLabel>
                      {newLanguages.map((language) => (
                        <DropdownMenuItem
                          className="flex items-center justify-between"
                          key={language}
                          onClick={() => handleLanguageSelect(language)}
                        >
                          <div className="flex items-center gap-2">
                            <Plus className="h-3.5 w-3.5 text-muted-foreground" />
                            <span>{language}</span>
                          </div>
                        </DropdownMenuItem>
                      ))}
                    </>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              currentLanguage && (
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Badge className="text-xs gap-1" variant="outline">
                      <Globe className="h-3 w-3" />
                      {currentLanguage}
                    </Badge>
                  </TooltipTrigger>
                  <TooltipContent>{t("executiveSummary.languageTooltip")}</TooltipContent>
                </Tooltip>
              )
            )}
            {onVersionSwitch && versions.length > 1 ? (
              <DropdownMenu>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <DropdownMenuTrigger asChild>
                      <Button
                        className="h-auto px-2.5 py-1 gap-1.5 text-xs font-normal"
                        disabled={isDisabled || isLoadingVersions}
                        variant="outline"
                      >
                        <History className="h-3 w-3" />
                        {formatShortDate(document.createdAt, locale)}
                        {commitSha && (
                          <code className="text-[10px] text-muted-foreground/70 font-mono">
                            {formatCommitSha(commitSha)}
                          </code>
                        )}
                        {isLatestVersion && (
                          <span className="text-muted-foreground">
                            ({t("executiveSummary.latestLabel")})
                          </span>
                        )}
                        <ChevronDown className="h-3 w-3 opacity-50" />
                      </Button>
                    </DropdownMenuTrigger>
                  </TooltipTrigger>
                  <TooltipContent>{t("executiveSummary.switchVersionTooltip")}</TooltipContent>
                </Tooltip>
                <DropdownMenuContent align="end" className="w-56 max-h-80 overflow-y-auto">
                  <DropdownMenuLabel className="flex items-center gap-2 px-2 py-2 text-[11px] font-medium uppercase tracking-wider text-muted-foreground/70 bg-muted/40 -mx-1 mb-1 border-b border-border/50">
                    <History className="h-3 w-3" />
                    {t("executiveSummary.versionHistory")}
                  </DropdownMenuLabel>
                  {versions.map((versionInfo) => {
                    const isCurrentVersion = versionInfo.version === currentVersion;
                    const isLatest = versionInfo.version === latestVersion;
                    const versionCommitSha =
                      "commitSha" in versionInfo ? versionInfo.commitSha : undefined;
                    return (
                      <DropdownMenuItem
                        className={cn(
                          "flex items-center justify-between",
                          isCurrentVersion && "bg-muted font-medium"
                        )}
                        key={versionInfo.version}
                        onClick={() => onVersionSwitch(versionInfo.version)}
                      >
                        <div className="flex items-center gap-2">
                          {isCurrentVersion ? (
                            <Check className="h-3.5 w-3.5 text-primary" />
                          ) : (
                            <span className="w-3.5" />
                          )}
                          <span>{formatShortDate(versionInfo.createdAt, locale)}</span>
                          {versionCommitSha && (
                            <code className="text-[10px] text-muted-foreground/70 font-mono">
                              {formatCommitSha(versionCommitSha)}
                            </code>
                          )}
                          {isLatest && (
                            <span className="text-muted-foreground">
                              ({t("executiveSummary.latestLabel")})
                            </span>
                          )}
                        </div>
                        <span className="text-xs text-muted-foreground">
                          {t("executiveSummary.versionLabel", { version: versionInfo.version })}
                        </span>
                      </DropdownMenuItem>
                    );
                  })}
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    className="h-auto px-2.5 py-1 gap-1.5 text-xs font-normal cursor-default"
                    variant="outline"
                  >
                    <History className="h-3 w-3" />
                    {formatShortDate(document.createdAt, locale)}
                    {currentVersion !== undefined && (
                      <span className="text-muted-foreground">
                        ({t("executiveSummary.latestLabel")})
                      </span>
                    )}
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t("executiveSummary.dateTooltip")}</TooltipContent>
              </Tooltip>
            )}
            {owner && repo && (
              <SpecExportButton
                disabled={isDisabled}
                document={document}
                owner={owner}
                repo={repo}
              />
            )}
            {onRegenerate && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    aria-label={t("executiveSummary.regenerateAriaLabel")}
                    disabled={isDisabled}
                    onClick={onRegenerate}
                    size="icon"
                    variant="ghost"
                  >
                    <RefreshCw className={isRegenerating ? "h-4 w-4 animate-spin" : "h-4 w-4"} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t("executiveSummary.regenerateTooltip")}</TooltipContent>
              </Tooltip>
            )}
          </div>
        </div>
        <CardDescription>
          {t("stats.domains", { count: domainCount })} /{" "}
          {t("stats.features", { count: featureCount })} /{" "}
          {t("stats.behaviors", { count: behaviorCount })}
        </CardDescription>
      </CardHeader>
      {hasExecutiveSummary ? (
        <CardContent className="space-y-4">
          <p className="text-sm leading-relaxed whitespace-pre-line">{document.executiveSummary}</p>
          <StatusLegend />
        </CardContent>
      ) : (
        <CardContent>
          <p className="text-sm text-muted-foreground">{t("executiveSummary.noSummary")}</p>
        </CardContent>
      )}
    </Card>
  );
};
