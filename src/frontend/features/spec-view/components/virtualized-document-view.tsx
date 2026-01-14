"use no memo";
"use client";

import { useWindowVirtualizer } from "@tanstack/react-virtual";
import { useEffect, useLayoutEffect, useRef, useState } from "react";

import { BehaviorRow, DomainHeaderRow, FeatureHeaderRow } from "./rows";
import type { FilteredDocument } from "../hooks/use-document-filter";
import type { FlatSpecItem } from "../types";
import { flattenSpecDocument, getItemHeight } from "../utils/flatten-spec-document";

type VirtualizedDocumentViewProps = {
  document: FilteredDocument;
  hasFilter: boolean;
};

/**
 * Window-level virtualized document view
 * Optimized for large documents with many behaviors
 */
export const VirtualizedDocumentView = ({ document, hasFilter }: VirtualizedDocumentViewProps) => {
  const listRef = useRef<HTMLDivElement>(null);
  const [scrollMargin, setScrollMargin] = useState(0);

  const [expandedDomains, setExpandedDomains] = useState<Set<string>>(() => {
    return new Set(document.domains.map((d) => d.id));
  });
  const [expandedFeatures, setExpandedFeatures] = useState<Set<string>>(() => {
    return new Set(document.domains.flatMap((d) => d.features.map((f) => f.id)));
  });

  // Reset expanded state when document structure changes (e.g., filter applied, navigation)
  useEffect(() => {
    setExpandedDomains(new Set(document.domains.map((d) => d.id)));
    setExpandedFeatures(new Set(document.domains.flatMap((d) => d.features.map((f) => f.id))));
  }, [document.domains]);

  const flatItems = flattenSpecDocument(document, expandedDomains, expandedFeatures, hasFilter);

  useLayoutEffect(() => {
    setScrollMargin(listRef.current?.offsetTop ?? 0);
  }, []);

  const handleDomainToggle = (domainId: string) => {
    setExpandedDomains((prev) => {
      const next = new Set(prev);
      if (next.has(domainId)) {
        next.delete(domainId);
      } else {
        next.add(domainId);
      }
      return next;
    });
  };

  const handleFeatureToggle = (featureId: string) => {
    setExpandedFeatures((prev) => {
      const next = new Set(prev);
      if (next.has(featureId)) {
        next.delete(featureId);
      } else {
        next.add(featureId);
      }
      return next;
    });
  };

  const estimateSize = (index: number): number => {
    const item = flatItems[index];
    return item ? getItemHeight(item) : 72;
  };

  const virtualizer = useWindowVirtualizer({
    count: flatItems.length,
    estimateSize,
    overscan: 10,
    scrollMargin,
  });

  const virtualItems = virtualizer.getVirtualItems();

  const renderItem = (item: FlatSpecItem) => {
    switch (item.type) {
      case "domain-header":
        return (
          <DomainHeaderRow
            domain={item.domain}
            hasFilter={hasFilter}
            isExpanded={item.isExpanded}
            onToggle={() => handleDomainToggle(item.domainId)}
          />
        );
      case "feature-header":
        return (
          <FeatureHeaderRow
            feature={item.feature}
            hasFilter={hasFilter}
            isExpanded={item.isExpanded}
            onToggle={() => handleFeatureToggle(item.featureId)}
          />
        );
      case "behavior":
        return <BehaviorRow behavior={item.behavior} />;
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
