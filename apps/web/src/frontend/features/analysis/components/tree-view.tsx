"use client";
"use no memo";

import { useWindowVirtualizer } from "@tanstack/react-virtual";
import { useLayoutEffect, useRef, useState } from "react";

import type { TestSuite } from "@/lib/api";

import { EmptyState } from "./empty-state";
import { TreeNode } from "./tree-node";
import type { FlatTreeItem } from "../types";
import { buildFileTree, flattenTree } from "../utils";

type TreeViewProps = {
  suites: TestSuite[];
};

const ESTIMATED_ITEM_HEIGHT = 40;
const TEST_ITEM_HEIGHT = 44;

export const TreeView = ({ suites }: TreeViewProps) => {
  const listRef = useRef<HTMLDivElement>(null);
  const [scrollMargin, setScrollMargin] = useState(0);

  const tree = buildFileTree(suites);

  const [expandedPaths, setExpandedPaths] = useState<Set<string>>(() => new Set());

  const flatItems: FlatTreeItem[] = flattenTree(tree, expandedPaths);

  useLayoutEffect(() => {
    setScrollMargin(listRef.current?.offsetTop ?? 0);
  }, []);

  const handleToggle = (path: string) => {
    setExpandedPaths((prev) => {
      const next = new Set(prev);
      if (next.has(path)) {
        next.delete(path);
      } else {
        next.add(path);
      }
      return next;
    });
  };

  const getItemHeight = (index: number) => {
    const item = flatItems[index];
    if (!item) return ESTIMATED_ITEM_HEIGHT;

    if (item.node.type === "file" && expandedPaths.has(item.node.path)) {
      return ESTIMATED_ITEM_HEIGHT + item.node.testCount * TEST_ITEM_HEIGHT;
    }

    return ESTIMATED_ITEM_HEIGHT;
  };

  const virtualizer = useWindowVirtualizer({
    count: flatItems.length,
    estimateSize: getItemHeight,
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
        const item = flatItems[virtualItem.index];
        if (!item) return null;

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
            <TreeNode
              isExpanded={expandedPaths.has(item.node.path)}
              item={item}
              onToggle={handleToggle}
            />
          </div>
        );
      })}
    </div>
  );
};
