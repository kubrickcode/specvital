"use client";

import { ChevronDown, ChevronRight } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";
import { useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { collapseTransition, expandCollapse, useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

import { BehaviorItem } from "./behavior-item";
import { VirtualizedBehaviorList, VIRTUALIZATION_THRESHOLD } from "./virtualized-behavior-list";
import type { FilteredFeature } from "../hooks/use-document-filter";

type FeatureGroupProps = {
  defaultOpen?: boolean;
  feature: FilteredFeature;
  hasFilter?: boolean;
};

export const FeatureGroup = ({
  defaultOpen = true,
  feature,
  hasFilter = false,
}: FeatureGroupProps) => {
  const [isOpen, setIsOpen] = useState(defaultOpen);
  const shouldReduceMotion = useReducedMotion();

  const visibleBehaviors = hasFilter
    ? feature.behaviors.filter((b) => b.hasMatch)
    : feature.behaviors;

  const behaviorCount = visibleBehaviors.length;
  const totalCount = feature.behaviors.length;

  const displayCount = hasFilter ? `${feature.matchCount}/${totalCount}` : behaviorCount;

  return (
    <div
      aria-labelledby={`feature-title-${feature.id}`}
      className={cn(
        "border-l-2 ml-3 pl-3",
        isOpen ? "border-primary/40" : "border-muted-foreground/20",
        "transition-colors"
      )}
      id={`feature-${feature.id}`}
      role="region"
      tabIndex={-1}
    >
      <Button
        aria-controls={isOpen ? `feature-content-${feature.id}` : undefined}
        aria-expanded={isOpen}
        className={cn(
          "w-full justify-start gap-2 px-3 py-2 h-auto",
          "text-left font-medium rounded-md",
          "hover:bg-muted/50"
        )}
        onClick={() => setIsOpen(!isOpen)}
        variant="ghost"
      >
        {isOpen ? (
          <ChevronDown aria-hidden="true" className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
        ) : (
          <ChevronRight
            aria-hidden="true"
            className="h-4 w-4 flex-shrink-0 text-muted-foreground"
          />
        )}

        <span className="flex-1 truncate text-sm" id={`feature-title-${feature.id}`}>
          {feature.name}
        </span>

        <Badge className="text-xs tabular-nums" variant="secondary">
          {displayCount}
        </Badge>
      </Button>

      <AnimatePresence initial={false}>
        {isOpen && (
          <motion.div
            animate={shouldReduceMotion ? {} : "expanded"}
            className="mt-1 space-y-0.5 pb-1"
            exit={shouldReduceMotion ? {} : "collapsed"}
            id={`feature-content-${feature.id}`}
            initial={shouldReduceMotion ? false : "collapsed"}
            transition={shouldReduceMotion ? undefined : collapseTransition}
            variants={shouldReduceMotion ? undefined : expandCollapse}
          >
            {feature.description && (
              <p className="px-3 py-1.5 text-sm text-muted-foreground">{feature.description}</p>
            )}
            {visibleBehaviors.length >= VIRTUALIZATION_THRESHOLD ? (
              <VirtualizedBehaviorList behaviors={visibleBehaviors} />
            ) : (
              <div className="space-y-0.5" role="list">
                {visibleBehaviors.map((behavior) => (
                  <BehaviorItem behavior={behavior} key={behavior.id} />
                ))}
              </div>
            )}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
