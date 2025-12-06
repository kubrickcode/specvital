"use client";

import { useWindowVirtualizer } from "@tanstack/react-virtual";
import { useLayoutEffect, useRef, useState } from "react";
import { useTranslations } from "next-intl";
import type { TestSuite } from "@/lib/api";
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

  const t = useTranslations("analyze");

  if (suites.length === 0) {
    return (
      <div className="rounded-lg border bg-card p-6 text-center text-muted-foreground">
        {t("noTestSuites")}
      </div>
    );
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
