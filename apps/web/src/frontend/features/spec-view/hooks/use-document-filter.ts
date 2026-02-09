"use client";

import { parseAsArrayOf, parseAsString, useQueryState } from "nuqs";
import { useDeferredValue, useTransition } from "react";

import type { TestStatus } from "@/lib/api";

import type { SpecBehavior, SpecDocument, SpecDomain, SpecFeature } from "../types";
import { findHighlightRanges } from "../utils/highlight";
import type { HighlightRange } from "../utils/highlight";

// Re-export for convenience
export type { HighlightRange } from "../utils/highlight";

const VALID_STATUSES: TestStatus[] = ["active", "focused", "skipped", "todo", "xfail"];

const queryParser = parseAsString.withDefault("");
const arrayParser = parseAsArrayOf(parseAsString, ",").withDefault([]);

export type FilteredDocument = {
  domains: FilteredDomain[];
  matchCount: number;
  totalCount: number;
};

export type FilteredDomain = Omit<SpecDomain, "features"> & {
  features: FilteredFeature[];
  hasMatch: boolean;
  matchCount: number;
};

export type FilteredFeature = Omit<SpecFeature, "behaviors"> & {
  behaviors: FilteredBehavior[];
  hasMatch: boolean;
  matchCount: number;
};

export type FilteredBehavior = SpecBehavior & {
  hasMatch: boolean;
  highlightRanges: HighlightRange[];
};

const matchesBehavior = (
  behavior: SpecBehavior,
  query: string,
  statuses: TestStatus[],
  frameworks: string[]
): { highlightRanges: HighlightRange[]; matches: boolean } => {
  // Status filter
  if (statuses.length > 0 && behavior.sourceInfo) {
    if (!statuses.includes(behavior.sourceInfo.status)) {
      return { highlightRanges: [], matches: false };
    }
  }

  // Framework filter
  if (frameworks.length > 0 && behavior.sourceInfo) {
    if (!frameworks.includes(behavior.sourceInfo.framework)) {
      return { highlightRanges: [], matches: false };
    }
  }

  // Query search
  if (query) {
    const descriptionRanges = findHighlightRanges(behavior.convertedDescription, query);
    const originalRanges = findHighlightRanges(behavior.originalName, query);

    if (descriptionRanges.length === 0 && originalRanges.length === 0) {
      return { highlightRanges: [], matches: false };
    }

    return { highlightRanges: descriptionRanges, matches: true };
  }

  return { highlightRanges: [], matches: true };
};

export const useDocumentFilter = (document: SpecDocument | null) => {
  const [isPending, startTransition] = useTransition();
  const [query, setQuery] = useQueryState("q", queryParser.withOptions({ startTransition }));
  const [rawStatuses, setRawStatuses] = useQueryState(
    "statuses",
    arrayParser.withOptions({ startTransition })
  );
  const [frameworks, setFrameworks] = useQueryState(
    "frameworks",
    arrayParser.withOptions({ startTransition })
  );

  // Defer expensive filtering to keep UI responsive during input
  const deferredQuery = useDeferredValue(query);
  const deferredStatuses = useDeferredValue(rawStatuses);
  const deferredFrameworks = useDeferredValue(frameworks);

  const statuses = deferredStatuses.filter((s): s is TestStatus =>
    VALID_STATUSES.includes(s as TestStatus)
  );

  const hasFilter =
    deferredQuery.length > 0 || statuses.length > 0 || deferredFrameworks.length > 0;

  const computeFilteredDocument = (): FilteredDocument | null => {
    if (!document) return null;

    let totalCount = 0;
    let matchCount = 0;

    const filteredDomains: FilteredDomain[] = document.domains
      .map((domain): FilteredDomain => {
        let domainMatchCount = 0;

        const filteredFeatures: FilteredFeature[] = domain.features
          .map((feature): FilteredFeature => {
            let featureMatchCount = 0;

            const filteredBehaviors: FilteredBehavior[] = feature.behaviors.map(
              (behavior): FilteredBehavior => {
                totalCount++;
                const { highlightRanges, matches } = matchesBehavior(
                  behavior,
                  deferredQuery,
                  statuses,
                  deferredFrameworks
                );

                if (matches) {
                  matchCount++;
                  featureMatchCount++;
                }

                return {
                  ...behavior,
                  hasMatch: matches,
                  highlightRanges,
                };
              }
            );

            domainMatchCount += featureMatchCount;

            return {
              ...feature,
              behaviors: filteredBehaviors,
              hasMatch: featureMatchCount > 0,
              matchCount: featureMatchCount,
            };
          })
          .filter((feature) => !hasFilter || feature.hasMatch);

        return {
          ...domain,
          features: filteredFeatures,
          hasMatch: domainMatchCount > 0,
          matchCount: domainMatchCount,
        };
      })
      .filter((domain) => !hasFilter || domain.hasMatch);

    return {
      domains: filteredDomains,
      matchCount,
      totalCount,
    };
  };

  const filteredDocument = computeFilteredDocument();

  const clearFilters = () => {
    setQuery("");
    setRawStatuses([]);
    setFrameworks([]);
  };

  const filterInfo = {
    frameworks: deferredFrameworks,
    query: deferredQuery,
    statuses,
  };

  return {
    clearFilters,
    filteredDocument,
    filterInfo,
    hasFilter,
    isPending,
    matchCount: filteredDocument?.matchCount ?? 0,
    query,
  } as const;
};
