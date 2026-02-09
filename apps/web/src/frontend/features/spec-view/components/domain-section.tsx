"use client";

import { ChevronDown, ChevronRight } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { collapseTransition, expandCollapse, useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

import { FeatureGroup } from "./feature-group";
import type { FilteredDomain } from "../hooks/use-document-filter";

type DomainSectionProps = {
  defaultOpen?: boolean;
  domain: FilteredDomain;
  hasFilter?: boolean;
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

export const DomainSection = ({
  defaultOpen = true,
  domain,
  hasFilter = false,
}: DomainSectionProps) => {
  const [isOpen, setIsOpen] = useState(defaultOpen);
  const shouldReduceMotion = useReducedMotion();
  const t = useTranslations("specView");

  const featureCount = domain.features.length;
  const behaviorCount = domain.features.reduce((sum, f) => sum + f.behaviors.length, 0);

  const displayCount = hasFilter
    ? `${domain.matchCount} / ${behaviorCount}`
    : t("domain.summary", { behaviorCount, featureCount });

  const confidenceVariant =
    domain.classificationConfidence !== undefined
      ? getConfidenceVariant(domain.classificationConfidence)
      : null;

  return (
    <Card
      aria-labelledby={`domain-title-${domain.id}`}
      className="overflow-hidden border-border/60 py-0"
      id={`domain-${domain.id}`}
      role="region"
      tabIndex={-1}
    >
      <CardHeader className="!px-4 !py-3 sm:!px-6 sm:!py-4">
        {/* Header: desktop=horizontal, mobile=vertical */}
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <Button
            aria-controls={isOpen ? `domain-content-${domain.id}` : undefined}
            aria-expanded={isOpen}
            className="justify-start gap-2 px-0 h-auto hover:bg-transparent min-w-0"
            onClick={() => setIsOpen(!isOpen)}
            variant="ghost"
          >
            {isOpen ? (
              <ChevronDown
                aria-hidden="true"
                className="h-5 w-5 flex-shrink-0 text-muted-foreground"
              />
            ) : (
              <ChevronRight
                aria-hidden="true"
                className="h-5 w-5 flex-shrink-0 text-muted-foreground"
              />
            )}
            <CardTitle
              className="text-lg font-semibold tracking-tight"
              id={`domain-title-${domain.id}`}
            >
              {domain.name}
            </CardTitle>
          </Button>

          {/* Metadata: mobile=aligned with title, desktop=right aligned */}
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
                <TooltipContent>{t("domain.aiConfidence")}</TooltipContent>
              </Tooltip>
            )}
            <span className="text-xs text-muted-foreground">{displayCount}</span>
          </div>
        </div>

        {/* Description: aligned with title (chevron width + gap = ml-7) */}
        {domain.description && (
          <p className="text-sm text-muted-foreground mt-2 ml-7 leading-relaxed">
            {domain.description}
          </p>
        )}
      </CardHeader>

      <AnimatePresence initial={false}>
        {isOpen && (
          <motion.div
            animate={shouldReduceMotion ? {} : "expanded"}
            exit={shouldReduceMotion ? {} : "collapsed"}
            initial={shouldReduceMotion ? false : "collapsed"}
            transition={shouldReduceMotion ? undefined : collapseTransition}
            variants={shouldReduceMotion ? undefined : expandCollapse}
          >
            <CardContent
              className="!px-4 sm:!px-6 pt-0 pb-4 space-y-3"
              id={`domain-content-${domain.id}`}
            >
              {domain.features.map((feature) => (
                <FeatureGroup feature={feature} hasFilter={hasFilter} key={feature.id} />
              ))}
            </CardContent>
          </motion.div>
        )}
      </AnimatePresence>
    </Card>
  );
};
