"use client";

import { useVirtualizer } from "@tanstack/react-virtual";
import { useRef } from "react";

import { BehaviorItem } from "./behavior-item";
import type { FilteredBehavior } from "../hooks/use-document-filter";

type VirtualizedBehaviorListProps = {
  behaviors: FilteredBehavior[];
};

const ITEM_HEIGHT_ESTIMATE = 72;
const OVERSCAN = 5;

/**
 * Virtualized behavior list for large documents (100+ behaviors)
 * Only renders visible items for performance
 */
export const VirtualizedBehaviorList = ({ behaviors }: VirtualizedBehaviorListProps) => {
  const parentRef = useRef<HTMLDivElement>(null);

  const virtualizer = useVirtualizer({
    count: behaviors.length,
    estimateSize: () => ITEM_HEIGHT_ESTIMATE,
    getScrollElement: () => parentRef.current,
    overscan: OVERSCAN,
  });

  const items = virtualizer.getVirtualItems();

  return (
    <div
      aria-label={`${behaviors.length} behaviors`}
      className="h-[400px] overflow-auto"
      ref={parentRef}
      role="list"
    >
      <div className="relative w-full" style={{ height: `${virtualizer.getTotalSize()}px` }}>
        {items.map((virtualItem) => {
          const behavior = behaviors[virtualItem.index];
          if (!behavior) return null;

          return (
            <div
              className="absolute top-0 left-0 w-full"
              data-index={virtualItem.index}
              key={virtualItem.key}
              ref={virtualizer.measureElement}
              style={{
                transform: `translateY(${virtualItem.start}px)`,
              }}
            >
              <BehaviorItem behavior={behavior} />
            </div>
          );
        })}
      </div>
    </div>
  );
};

/**
 * Threshold for switching to virtualized list
 */
export const VIRTUALIZATION_THRESHOLD = 30;
