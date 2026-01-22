"use no memo";
"use client";

import { defaultRangeExtractor, useWindowVirtualizer } from "@tanstack/react-virtual";
import { useLayoutEffect, useRef, useState } from "react";

import { BehaviorRow, DomainHeaderRow, FeatureHeaderRow } from "./rows";
import { useDocumentNavigation } from "../contexts";
import type { FilteredDocument } from "../hooks/use-document-filter";
import type { FlatSpecItem } from "../types";
import { flattenSpecDocument, getItemHeight } from "../utils/flatten-spec-document";

type VirtualizedDocumentViewProps = {
  document: FilteredDocument;
  hasFilter: boolean;
};

/**
 * Get section ID from flat item
 */
const getSectionIdFromItem = (item: FlatSpecItem): string | null => {
  if (item.type === "domain-header") return `domain-${item.domainId}`;
  if (item.type === "feature-header") return `feature-${item.featureId}`;
  return null;
};

/**
 * Find the first header item in the visible range
 */
const findFirstVisibleHeader = (
  flatItems: FlatSpecItem[],
  startIndex: number,
  endIndex: number
): string | null => {
  for (let i = startIndex; i <= endIndex && i < flatItems.length; i++) {
    const item = flatItems[i];
    if (item && (item.type === "domain-header" || item.type === "feature-header")) {
      return getSectionIdFromItem(item);
    }
  }
  return null;
};

/**
 * Window-level virtualized document view
 * Optimized for large documents with many behaviors
 */
export const VirtualizedDocumentView = ({ document, hasFilter }: VirtualizedDocumentViewProps) => {
  const listRef = useRef<HTMLDivElement>(null);
  const [scrollMargin, setScrollMargin] = useState(0);

  // Track visible range for scroll spy
  const visibleRangeRef = useRef({ endIndex: 0, startIndex: 0 });

  const {
    expandedDomains,
    expandedFeatures,
    registerVirtualizer,
    setActiveSection,
    toggleDomain,
    toggleFeature,
  } = useDocumentNavigation();

  const flatItems = flattenSpecDocument(document, expandedDomains, expandedFeatures, hasFilter);

  useLayoutEffect(() => {
    setScrollMargin(listRef.current?.offsetTop ?? 0);
  }, []);

  const estimateSize = (index: number): number => {
    const item = flatItems[index];
    return item ? getItemHeight(item) : 72;
  };

  const virtualizer = useWindowVirtualizer({
    count: flatItems.length,
    estimateSize,
    overscan: 10,
    scrollMargin,
    // Track actual visible range (excluding overscan)
    rangeExtractor: (range) => {
      visibleRangeRef.current = {
        endIndex: range.endIndex,
        startIndex: range.startIndex,
      };
      return defaultRangeExtractor(range);
    },
    // Scroll spy: update active section when visible items change
    onChange: () => {
      const { endIndex, startIndex } = visibleRangeRef.current;
      const sectionId = findFirstVisibleHeader(flatItems, startIndex, endIndex);
      if (sectionId) {
        setActiveSection(sectionId);
      }
    },
  });

  // Register virtualizer with context for navigation
  // Using useLayoutEffect to ensure it's registered before any scroll attempts
  useLayoutEffect(() => {
    registerVirtualizer(virtualizer, flatItems);
  }, [registerVirtualizer, virtualizer, flatItems]);

  const virtualItems = virtualizer.getVirtualItems();

  const renderItem = (item: FlatSpecItem) => {
    switch (item.type) {
      case "domain-header":
        return (
          <DomainHeaderRow
            domain={item.domain}
            hasFilter={hasFilter}
            isExpanded={item.isExpanded}
            isLastInDomain={item.isLastInDomain}
            onToggle={() => toggleDomain(item.domainId)}
          />
        );
      case "feature-header":
        return (
          <FeatureHeaderRow
            feature={item.feature}
            hasFilter={hasFilter}
            isExpanded={item.isExpanded}
            isLastInDomain={item.isLastInDomain}
            onToggle={() => toggleFeature(item.featureId)}
          />
        );
      case "behavior":
        return <BehaviorRow behavior={item.behavior} isLastInDomain={item.isLastInDomain} />;
    }
  };

  return (
    <div
      ref={listRef}
      role="list"
      style={{
        height: `${virtualizer.getTotalSize()}px`,
        position: "relative",
        width: "100%",
      }}
    >
      {virtualItems.map((virtualItem) => {
        const item = flatItems[virtualItem.index];
        if (!item) return null;

        const key =
          item.type === "domain-header"
            ? `domain-${item.domainId}`
            : item.type === "feature-header"
              ? `feature-${item.featureId}`
              : `behavior-${item.behaviorId}`;

        return (
          <div
            data-index={virtualItem.index}
            key={key}
            ref={virtualizer.measureElement}
            style={{
              left: 0,
              position: "absolute",
              top: 0,
              transform: `translateY(${virtualItem.start - scrollMargin}px)`,
              width: "100%",
            }}
          >
            {renderItem(item)}
          </div>
        );
      })}
    </div>
  );
};
