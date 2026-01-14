"use client";

import { ChevronDown, ChevronRight } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

import type { FilteredFeature } from "../../hooks/use-document-filter";

type FeatureHeaderRowProps = {
  feature: FilteredFeature;
  hasFilter: boolean;
  isExpanded: boolean;
  isLastInDomain: boolean;
  onToggle: () => void;
};

export const FeatureHeaderRow = ({
  feature,
  hasFilter,
  isExpanded,
  isLastInDomain,
  onToggle,
}: FeatureHeaderRowProps) => {
  const visibleBehaviorCount = hasFilter
    ? feature.behaviors.filter((b) => b.hasMatch).length
    : feature.behaviors.length;
  const totalCount = feature.behaviors.length;

  const displayCount = hasFilter ? `${feature.matchCount}/${totalCount}` : visibleBehaviorCount;

  return (
    <div
      className={cn(
        // Card middle styles - side borders and background
        "bg-card border-x border-border/60 px-4 sm:px-6",
        // If last in domain, add bottom border and margin
        isLastInDomain && "border-b rounded-b-lg mb-6 pb-4"
      )}
      id={`feature-${feature.id}`}
      role="region"
    >
      <div
        className={cn(
          "border-l-2 pl-3",
          isExpanded ? "border-primary/40" : "border-muted-foreground/20",
          "transition-colors"
        )}
      >
        <Button
          aria-expanded={isExpanded}
          className={cn(
            "w-full justify-start gap-2 px-3 py-2 h-auto",
            "text-left font-medium rounded-md",
            "hover:bg-muted/50"
          )}
          onClick={onToggle}
          variant="ghost"
        >
          {isExpanded ? (
            <ChevronDown className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
          )}

          <span className="flex-1 truncate text-sm">{feature.name}</span>

          <Badge className="text-xs tabular-nums" variant="secondary">
            {displayCount}
          </Badge>
        </Button>

        {isExpanded && feature.description && (
          <p className="px-3 py-1.5 text-sm text-muted-foreground">{feature.description}</p>
        )}
      </div>
    </div>
  );
};
