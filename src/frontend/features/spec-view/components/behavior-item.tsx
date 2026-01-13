"use client";

import { Check, Circle, CircleDashed, Crosshair, RefreshCw, XCircle } from "lucide-react";
import { useMemo, useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import type { TestStatus } from "@/lib/api/types";
import { cn } from "@/lib/utils";

import { HighlightedText } from "./highlighted-text";
import type { FilteredBehavior } from "../hooks/use-document-filter";
import { findHighlightRanges } from "../utils/highlight";

type BehaviorItemProps = {
  behavior: FilteredBehavior;
  query?: string;
};

const STATUS_CONFIG: Record<TestStatus, { color: string; icon: typeof Check; label: string }> = {
  active: {
    color: "text-green-600",
    icon: Check,
    label: "Active test",
  },
  focused: {
    color: "text-purple-500",
    icon: Crosshair,
    label: "Focused test",
  },
  skipped: {
    color: "text-amber-500",
    icon: CircleDashed,
    label: "Skipped test",
  },
  todo: {
    color: "text-blue-500",
    icon: Circle,
    label: "Todo test",
  },
  xfail: {
    color: "text-red-400",
    icon: XCircle,
    label: "Expected failure",
  },
} as const;

export const BehaviorItem = ({ behavior, query = "" }: BehaviorItemProps) => {
  const [isShowingOriginal, setIsShowingOriginal] = useState(false);
  const sourceInfo = behavior.sourceInfo;
  const config = sourceInfo ? STATUS_CONFIG[sourceInfo.status] : null;
  const Icon = config?.icon ?? Check;

  // Calculate highlight ranges for the displayed text
  const displayText = isShowingOriginal ? behavior.originalName : behavior.convertedDescription;
  const highlightRanges = useMemo(
    () => (query ? findHighlightRanges(displayText, query) : behavior.highlightRanges),
    [displayText, query, behavior.highlightRanges]
  );

  return (
    <div
      className={cn(
        "flex items-start gap-3 px-3 py-2.5 rounded-md",
        "hover:bg-muted/50 transition-colors group",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      )}
      id={`behavior-${behavior.id}`}
      role="listitem"
      tabIndex={-1}
    >
      <Tooltip>
        <TooltipTrigger asChild>
          <span className="mt-0.5">
            <Icon
              aria-label={config?.label ?? "Test"}
              className={cn("h-4 w-4 flex-shrink-0", config?.color ?? "text-muted-foreground")}
            />
          </span>
        </TooltipTrigger>
        <TooltipContent>{config?.label ?? "Test"}</TooltipContent>
      </Tooltip>

      <div className="flex-1 min-w-0 space-y-1">
        <p className="text-sm leading-relaxed">
          <HighlightedText ranges={highlightRanges} text={displayText} />
        </p>

        {sourceInfo && (
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <span className="font-mono truncate max-w-[200px]">{sourceInfo.filePath}</span>
            <span className="font-mono">L:{sourceInfo.lineNumber}</span>
            <Badge className="text-[10px] px-1.5" variant="outline">
              {sourceInfo.framework}
            </Badge>
          </div>
        )}
      </div>

      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            className="opacity-0 group-hover:opacity-100 transition-opacity"
            onClick={() => setIsShowingOriginal(!isShowingOriginal)}
            size="icon"
            variant="ghost"
          >
            <RefreshCw className="h-3.5 w-3.5" />
            <span className="sr-only">Toggle original/converted</span>
          </Button>
        </TooltipTrigger>
        <TooltipContent>{isShowingOriginal ? "Show converted" : "Show original"}</TooltipContent>
      </Tooltip>
    </div>
  );
};
