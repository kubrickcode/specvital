"use client";

import { Check } from "lucide-react";
import { useTranslations } from "next-intl";

import { FrameworkBadge } from "@/components/ui/framework-badge";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

import { HighlightedText } from "./highlighted-text";
import { STATUS_CONFIG } from "../constants/status-config";
import type { FilteredBehavior } from "../hooks/use-document-filter";

type BehaviorItemProps = {
  behavior: FilteredBehavior;
};

export const BehaviorItem = ({ behavior }: BehaviorItemProps) => {
  const t = useTranslations("specView.behavior");
  const tStatus = useTranslations("specView.statusLegend");
  const sourceInfo = behavior.sourceInfo;
  const config = sourceInfo ? STATUS_CONFIG[sourceInfo.status] : null;
  const Icon = config?.icon ?? Check;
  const statusLabel = config ? tStatus(config.labelKey) : "Test";

  const hasDifferentOriginal = behavior.originalName !== behavior.convertedDescription;

  const content = (
    <div
      className={cn(
        "group flex items-start gap-3 px-3 py-2.5 mx-1 rounded-lg",
        "hover:bg-muted/50 active:bg-muted/70",
        "transition-colors duration-150",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1"
      )}
      id={`behavior-${behavior.id}`}
      role="listitem"
      tabIndex={0}
    >
      <Tooltip>
        <TooltipTrigger asChild>
          <span
            className={cn(
              "flex items-center justify-center w-6 h-6 rounded-full mt-0.5 flex-shrink-0",
              "ring-1 transition-all duration-150",
              config?.bgColor ?? "bg-muted",
              config?.ringColor ?? "ring-border",
              "group-hover:ring-2"
            )}
          >
            <Icon
              aria-label={statusLabel}
              className={cn("h-3.5 w-3.5", config?.color ?? "text-muted-foreground")}
            />
          </span>
        </TooltipTrigger>
        <TooltipContent side="left">{statusLabel}</TooltipContent>
      </Tooltip>

      <div className="flex-1 min-w-0 space-y-1.5">
        {/* Behavior description */}
        <p className="text-sm leading-relaxed text-foreground/90">
          <HighlightedText ranges={behavior.highlightRanges} text={behavior.convertedDescription} />
        </p>

        {/* Source info as inline chips */}
        {sourceInfo && (
          <div className="flex flex-wrap items-center gap-1.5">
            <span
              className={cn(
                "inline-flex items-center gap-1 px-2 py-0.5 rounded-md",
                "bg-muted/60 text-xs font-mono text-muted-foreground",
                "max-w-[200px] md:max-w-[280px] truncate"
              )}
            >
              {sourceInfo.filePath}
            </span>
            <span
              className={cn(
                "inline-flex items-center px-1.5 py-0.5 rounded-md",
                "bg-muted/60 text-xs font-mono text-muted-foreground tabular-nums"
              )}
            >
              L{sourceInfo.lineNumber}
            </span>
            <FrameworkBadge framework={sourceInfo.framework} />
          </div>
        )}
      </div>
    </div>
  );

  if (!hasDifferentOriginal) {
    return content;
  }

  return (
    <ResponsiveTooltip
      content={
        <div className="flex flex-col gap-1.5 max-w-xs">
          <span className="text-xs font-medium opacity-70">{t("original")}</span>
          <span className="font-mono text-xs leading-relaxed">{behavior.originalName}</span>
        </div>
      }
      side="bottom"
    >
      {content}
    </ResponsiveTooltip>
  );
};
