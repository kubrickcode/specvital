"use client";

import { Bot, Calendar, FileText, Globe, RefreshCw } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";

import { StatusLegend } from "./status-legend";
import type { SpecDocument } from "../types";
import { calculateDocumentStats } from "../utils/stats";

type ExecutiveSummaryProps = {
  document: SpecDocument;
  isRegenerating?: boolean;
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
  isRegenerating = false,
  onRegenerate,
}: ExecutiveSummaryProps) => {
  const t = useTranslations("specView");

  if (!document.executiveSummary) {
    return null;
  }

  const { behaviorCount, domainCount, featureCount } = calculateDocumentStats(document);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between flex-wrap gap-2">
          <div className="flex items-center gap-2">
            <FileText className="h-5 w-5 text-primary" />
            <CardTitle className="text-lg">{t("executiveSummary.title")}</CardTitle>
          </div>
          <div className="flex items-center gap-2 flex-wrap">
            {document.language && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Badge className="text-xs gap-1" variant="outline">
                    <Globe className="h-3 w-3" />
                    {document.language}
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>{t("executiveSummary.languageTooltip")}</TooltipContent>
              </Tooltip>
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
                    disabled={isRegenerating}
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
