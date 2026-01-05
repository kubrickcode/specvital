"use client";

import { useWindowVirtualizer } from "@tanstack/react-virtual";
import { useLayoutEffect, useRef, useState } from "react";

import { SpecFileGroup } from "./spec-file-group";
import type { ConvertedTestFile } from "../../types";
import { EmptyState } from "../empty-state";

type SpecViewContentProps = {
  files: ConvertedTestFile[];
};

const ITEM_GAP = 12;
const ESTIMATED_ITEM_HEIGHT = 60;

export const SpecViewContent = ({ files }: SpecViewContentProps) => {
  const listRef = useRef<HTMLDivElement>(null);
  const [scrollMargin, setScrollMargin] = useState(0);

  useLayoutEffect(() => {
    setScrollMargin(listRef.current?.offsetTop ?? 0);
  }, []);

  const virtualizer = useWindowVirtualizer({
    count: files.length,
    estimateSize: () => ESTIMATED_ITEM_HEIGHT + ITEM_GAP,
    overscan: 5,
    scrollMargin,
  });

  if (files.length === 0) {
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
        const file = files[virtualItem.index];
        if (!file) {
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
              <SpecFileGroup file={file} />
            </div>
          </div>
        );
      })}
    </div>
  );
};
