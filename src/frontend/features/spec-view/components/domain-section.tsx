"use client";

import { ChevronDown, ChevronRight, Layers } from "lucide-react";
import { useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

import { FeatureGroup } from "./feature-group";
import type { FilteredDomain } from "../hooks/use-document-filter";

type DomainSectionProps = {
  defaultOpen?: boolean;
  domain: FilteredDomain;
  hasFilter?: boolean;
  query?: string;
};

const formatConfidence = (confidence: number): string => {
  return `${Math.round(confidence * 100)}%`;
};

const getConfidenceColor = (confidence: number): string => {
  if (confidence >= 0.9) return "bg-green-500/10 text-green-700 border-green-500/30";
  if (confidence >= 0.7) return "bg-yellow-500/10 text-yellow-700 border-yellow-500/30";
  return "bg-red-500/10 text-red-700 border-red-500/30";
};

export const DomainSection = ({
  defaultOpen = true,
  domain,
  hasFilter = false,
  query = "",
}: DomainSectionProps) => {
  const [isOpen, setIsOpen] = useState(defaultOpen);

  const featureCount = domain.features.length;
  const behaviorCount = domain.features.reduce((sum, f) => sum + f.behaviors.length, 0);

  // Show match count when filtering
  const displayCount = hasFilter
    ? `${domain.matchCount} / ${behaviorCount}`
    : `${featureCount} features / ${behaviorCount} behaviors`;

  return (
    <Card
      aria-labelledby={`domain-title-${domain.id}`}
      className="overflow-hidden"
      id={`domain-${domain.id}`}
      role="region"
      tabIndex={-1}
    >
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <Button
            aria-controls={isOpen ? `domain-content-${domain.id}` : undefined}
            aria-expanded={isOpen}
            className="flex-1 justify-start gap-3 px-0 h-auto hover:bg-transparent"
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
            <Layers aria-hidden="true" className="h-5 w-5 flex-shrink-0 text-primary" />
            <CardTitle className="text-lg" id={`domain-title-${domain.id}`}>
              {domain.name}
            </CardTitle>
          </Button>

          <div className="flex items-center gap-2">
            {domain.classificationConfidence !== undefined && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Badge
                    className={cn("text-xs", getConfidenceColor(domain.classificationConfidence))}
                    variant="outline"
                  >
                    {formatConfidence(domain.classificationConfidence)}
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>AI classification confidence</TooltipContent>
              </Tooltip>
            )}
            <Badge variant="secondary">{displayCount}</Badge>
          </div>
        </div>

        {domain.description && (
          <p className="text-sm text-muted-foreground mt-2 ml-12">{domain.description}</p>
        )}
      </CardHeader>

      {isOpen && (
        <CardContent className="pt-0 space-y-3" id={`domain-content-${domain.id}`}>
          {domain.features.map((feature) => (
            <FeatureGroup feature={feature} hasFilter={hasFilter} key={feature.id} query={query} />
          ))}
        </CardContent>
      )}
    </Card>
  );
};
