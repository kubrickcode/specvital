"use client";

import { Zap } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";

import type { BehaviorCacheStats } from "../types";

type CacheStatsIndicatorProps = {
  stats: BehaviorCacheStats;
};

export const CacheStatsIndicator = ({ stats }: CacheStatsIndicatorProps) => {
  const t = useTranslations("specView.cacheStats");

  const { cachedBehaviors, generatedBehaviors, hitRate, totalBehaviors } = stats;

  // Don't show if no cache hits (first generation or forced regeneration)
  if (cachedBehaviors === 0) {
    return null;
  }

  const hitRatePercent = Math.round(hitRate * 100);
  const isFullCache = hitRatePercent === 100;

  const summaryText = isFullCache
    ? t("allFromCache")
    : t("summary", {
        cached: cachedBehaviors,
        rate: hitRatePercent,
        total: totalBehaviors,
      });

  const savedCallsText = t("costSaved", { saved: cachedBehaviors });

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Badge className="gap-1.5 text-xs font-normal" variant="secondary">
          <Zap className="h-3 w-3 text-amber-500" />
          <span>{t("badgeLabel", { rate: hitRatePercent })}</span>
        </Badge>
      </TooltipTrigger>
      <TooltipContent className="max-w-xs">
        <div className="space-y-1">
          <p className="font-medium">{summaryText}</p>
          <p className="text-muted-foreground">{savedCallsText}</p>
          {generatedBehaviors > 0 && (
            <p className="text-xs text-muted-foreground">
              {t("generatedCount", { count: generatedBehaviors })}
            </p>
          )}
        </div>
      </TooltipContent>
    </Tooltip>
  );
};
