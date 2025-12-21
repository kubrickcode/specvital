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
        position: "relative",
        width: "100%",
      }}
    >
      {virtualItems.map((virtualItem) => {
        const suite = suites[virtualItem.index];
        if (!suite) {
          return null;
        }

        return (
          <div
            data-index={virtualItem.index}
            key={virtualItem.key}
            ref={virtualizer.measureElement}
            style={{
              left: 0,
              position: "absolute",
              top: 0,
              transform: `translateY(${virtualItem.start - scrollMargin}px)`,
              width: "100%",
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
