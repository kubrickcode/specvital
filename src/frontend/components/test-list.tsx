"use client";

import { useWindowVirtualizer } from "@tanstack/react-virtual";
import { useLayoutEffect, useRef, useState } from "react";
import type { TestSuite } from "@/lib/api";
import { EmptyState } from "./empty-state";
import { TestSuiteAccordion } from "./test-suite-accordion";

type TestListProps = {
  suites: TestSuite[];
};

const ITEM_GAP = 12;
const ESTIMATED_ITEM_HEIGHT = 60;

export const TestList = ({ suites }: TestListProps) => {
  const listRef = useRef<HTMLDivElement>(null);
  const [scrollMargin, setScrollMargin] = useState(0);

  useLayoutEffect(() => {
    setScrollMargin(listRef.current?.offsetTop ?? 0);
  }, []);

  const virtualizer = useWindowVirtualizer({
    count: suites.length,
    estimateSize: () => ESTIMATED_ITEM_HEIGHT + ITEM_GAP,
    overscan: 5,
    scrollMargin,
  });

  if (suites.length === 0) {
    return <EmptyState />;
  }

  const virtualItems = virtualizer.getVirtualItems();

  return (
    <div
      ref={listRef}
      style={{
        height: `${virtualizer.getTotalSize()}px`,
        width: "100%",
        position: "relative",
      }}
    >
      {virtualItems.map((virtualItem) => {
        const suite = suites[virtualItem.index];
        if (!suite) {
          return null;
        }

        return (
          <div
            key={virtualItem.key}
            data-index={virtualItem.index}
            ref={virtualizer.measureElement}
            style={{
              position: "absolute",
              top: 0,
              left: 0,
              width: "100%",
              transform: `translateY(${virtualItem.start - scrollMargin}px)`,
            }}
          >
            <div className="pb-3">
              <TestSuiteAccordion suite={suite} />
            </div>
          </div>
        );
      })}
    </div>
  );
};
