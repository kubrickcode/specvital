"use client";

import { Sparkles } from "lucide-react";
import { useFormatter, useNow, useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Link } from "@/i18n/navigation";
import type { AiSpecSummary } from "@/lib/api/types";
import { cn } from "@/lib/utils";

type AiSpecBadgeProps = {
  owner: string;
  repo: string;
  summary: AiSpecSummary;
};

export const AiSpecBadge = ({ owner, repo, summary }: AiSpecBadgeProps) => {
  const t = useTranslations("dashboard.aiSpec");
  const format = useFormatter();
  const now = useNow({ updateInterval: 60_000 });

  const { hasSpec, languageCount, latestGeneratedAt } = summary;

  if (!hasSpec) {
    return null;
  }

  const handleWrapperClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  return (
    <div onClick={handleWrapperClick} onKeyDown={(e) => e.stopPropagation()}>
      <Popover>
        <PopoverTrigger asChild>
          <Badge
            aria-label={t("badge")}
            className={cn(
              "gap-1 cursor-pointer transition-colors",
              "bg-violet-100 text-violet-800 hover:bg-violet-200",
              "dark:bg-violet-900/30 dark:text-violet-300 dark:hover:bg-violet-900/50"
            )}
            variant="outline"
          >
            <Sparkles aria-hidden="true" className="size-3" />
            <span className="text-xs">{t("badge")}</span>
          </Badge>
        </PopoverTrigger>
        <PopoverContent align="start" className="w-56 p-3">
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <Sparkles aria-hidden="true" className="size-4 text-violet-500" />
              <span className="font-medium text-sm">{t("badge")}</span>
            </div>
            <div className="space-y-1 text-xs text-muted-foreground">
              {languageCount !== undefined && languageCount > 0 && (
                <p>{t("languageCount", { count: languageCount })}</p>
              )}
              {latestGeneratedAt && (
                <p>
                  {t("lastGenerated", {
                    time: format.relativeTime(new Date(latestGeneratedAt), now),
                  })}
                </p>
              )}
            </div>
            <Link
              className="block text-xs text-violet-600 dark:text-violet-400 font-medium pt-1 hover:underline"
              href={`/analyze/${owner}/${repo}?tab=spec`}
            >
              {t("viewDetails")}
            </Link>
          </div>
        </PopoverContent>
      </Popover>
    </div>
  );
};
