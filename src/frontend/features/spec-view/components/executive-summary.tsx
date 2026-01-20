"use client";

import { Bot, Calendar, ChevronDown, FileText, Globe, Languages, RefreshCw } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

import { StatusLegend } from "./status-legend";
import { SPEC_LANGUAGES } from "../constants/spec-languages";
import type { SpecDocument, SpecLanguage } from "../types";
import { calculateDocumentStats } from "../utils/stats";

type ExecutiveSummaryProps = {
  document: SpecDocument;
  isGeneratingOtherLanguage?: boolean;
  isRegenerating?: boolean;
  onLanguageSwitch?: (language: SpecLanguage) => void;
  onRegenerate?: () => void;
};

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString(undefined, {
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    month: "short",
    year: "numeric",
  });
};

export const ExecutiveSummary = ({
  document,
  isGeneratingOtherLanguage = false,
  isRegenerating = false,
  onLanguageSwitch,
  onRegenerate,
}: ExecutiveSummaryProps) => {
  const t = useTranslations("specView");

  if (!document.executiveSummary) {
    return null;
  }

  const { behaviorCount, domainCount, featureCount } = calculateDocumentStats(document);
  const currentLanguage = document.language;
  const isDisabled = isRegenerating || isGeneratingOtherLanguage;

  return (
    <Card className={cn(isGeneratingOtherLanguage && "opacity-60 pointer-events-none")}>
      <CardHeader>
        <div className="flex items-center justify-between flex-wrap gap-2">
          <div className="flex items-center gap-2">
            <FileText className="h-5 w-5 text-primary" />
            <CardTitle className="text-lg">{t("executiveSummary.title")}</CardTitle>
          </div>
          <div className="flex items-center gap-2 flex-wrap">
            {currentLanguage && onLanguageSwitch ? (
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
                <DropdownMenuContent align="end" className="max-h-64 overflow-y-auto">
                  {SPEC_LANGUAGES.map((language) => (
                    <DropdownMenuItem
                      className={cn(language === currentLanguage && "bg-muted font-medium")}
                      key={language}
                      onClick={() => {
                        if (language !== currentLanguage) {
                          onLanguageSwitch(language);
                        }
                      }}
                    >
                      {language}
                      {language === currentLanguage && (
                        <span className="ml-auto text-xs text-muted-foreground">
                          {t("executiveSummary.currentLanguage")}
                        </span>
                      )}
                    </DropdownMenuItem>
                  ))}
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
            {document.modelId && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Badge className="text-xs gap-1" variant="outline">
                    <Bot className="h-3 w-3" />
                    {document.modelId}
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>{t("executiveSummary.modelTooltip")}</TooltipContent>
              </Tooltip>
            )}
            <Tooltip>
              <TooltipTrigger asChild>
                <Badge className="text-xs gap-1" variant="secondary">
                  <Calendar className="h-3 w-3" />
                  {formatDate(document.createdAt)}
                </Badge>
              </TooltipTrigger>
              <TooltipContent>{t("executiveSummary.dateTooltip")}</TooltipContent>
            </Tooltip>
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
      <CardContent className="space-y-4">
        <p className="text-sm leading-relaxed whitespace-pre-line">{document.executiveSummary}</p>
        <StatusLegend />
      </CardContent>
    </Card>
  );
};
