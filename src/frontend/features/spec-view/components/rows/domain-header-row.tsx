"use client";

import { ChevronDown, ChevronRight } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

import type { FilteredDomain } from "../../hooks/use-document-filter";

type DomainHeaderRowProps = {
  domain: FilteredDomain;
  hasFilter: boolean;
  isExpanded: boolean;
  isLastInDomain: boolean;
  onToggle: () => void;
};

const formatConfidence = (confidence: number): string => {
  return `${Math.round(confidence * 100)}%`;
};

const getConfidenceVariant = (confidence: number) => {
  if (confidence >= 0.9)
    return {
      bg: "bg-emerald-50 dark:bg-emerald-950/30",
      border: "border-emerald-200 dark:border-emerald-800",
      text: "text-emerald-700 dark:text-emerald-400",
    };
  if (confidence >= 0.7)
    return {
      bg: "bg-amber-50 dark:bg-amber-950/30",
      border: "border-amber-200 dark:border-amber-800",
      text: "text-amber-700 dark:text-amber-400",
    };
  return {
    bg: "bg-red-50 dark:bg-red-950/30",
    border: "border-red-200 dark:border-red-800",
    text: "text-red-700 dark:text-red-400",
  };
};

export const DomainHeaderRow = ({
  domain,
  hasFilter,
  isExpanded,
  isLastInDomain,
  onToggle,
}: DomainHeaderRowProps) => {
  const featureCount = domain.features.length;
  const behaviorCount = domain.features.reduce((sum, f) => sum + f.behaviors.length, 0);

  const displayCount = hasFilter
    ? `${domain.matchCount} / ${behaviorCount}`
    : `${featureCount} features, ${behaviorCount} behaviors`;

  const confidenceVariant =
    domain.classificationConfidence !== undefined
      ? getConfidenceVariant(domain.classificationConfidence)
      : null;

  return (
    <div
      className={cn(
        // Card styles - always top border and rounded top
        "bg-card border-t border-x border-border/60 rounded-t-lg",
        // If last in domain (collapsed), add bottom border and margin
        isLastInDomain && "border-b rounded-b-lg mb-6"
      )}
      id={`domain-${domain.id}`}
      role="region"
    >
      <div className="px-4 py-3 sm:px-6 sm:py-4">
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <Button
            aria-expanded={isExpanded}
            className="justify-start gap-2 px-0 h-auto hover:bg-transparent min-w-0"
            onClick={onToggle}
            variant="ghost"
          >
            {isExpanded ? (
              <ChevronDown className="h-5 w-5 flex-shrink-0 text-muted-foreground" />
            ) : (
              <ChevronRight className="h-5 w-5 flex-shrink-0 text-muted-foreground" />
            )}
            <h3 className="text-lg font-semibold tracking-tight">{domain.name}</h3>
          </Button>

          <div className="flex items-center gap-2 flex-wrap ml-7 sm:ml-0 sm:flex-shrink-0">
            {confidenceVariant && domain.classificationConfidence !== undefined && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Badge
                    className={cn(
                      "text-xs font-medium tabular-nums",
                      confidenceVariant.bg,
                      confidenceVariant.text,
                      confidenceVariant.border
                    )}
                    variant="outline"
                  >
                    {formatConfidence(domain.classificationConfidence)}
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>AI classification confidence</TooltipContent>
              </Tooltip>
            )}
            <span className="text-xs text-muted-foreground">{displayCount}</span>
          </div>
        </div>

        {domain.description && (
          <p className="text-sm text-muted-foreground mt-2 ml-7 leading-relaxed">
            {domain.description}
          </p>
        )}
      </div>
    </div>
  );
};
