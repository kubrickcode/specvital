"use client";

import { ChevronDown, ChevronRight, FileText, Folder, FolderOpen } from "lucide-react";
import { useMemo } from "react";

import { FrameworkBadge } from "@/components/ui/framework-badge";
import { cn } from "@/lib/utils";

import { StatusMiniBar } from "./status-mini-bar";
import { TestItem } from "./test-item";
import type { FlatTreeItem } from "../types";
import { calculateStatusCounts } from "../utils/calculate-status-counts";

type TreeNodeProps = {
  isExpanded: boolean;
  item: FlatTreeItem;
  onToggle: (path: string) => void;
};

const INDENT_SIZE = 20;

export const TreeNode = ({ isExpanded, item, onToggle }: TreeNodeProps) => {
  const { depth, hasChildren, node } = item;
  const isDirectory = node.type === "directory";

  const statusCounts = useMemo(() => {
    if (isDirectory) {
      return { active: 0, skipped: 0, todo: 0 };
    }
    const allTests = node.suites.flatMap((suite) => suite.tests);
    return calculateStatusCounts(allTests);
  }, [isDirectory, node.suites]);

  const handleClick = () => {
    onToggle(node.path);
  };

  const paddingLeft = depth * INDENT_SIZE;

  if (isDirectory) {
    return (
      <button
        aria-expanded={isExpanded}
        aria-label={isExpanded ? `Collapse ${node.name}` : `Expand ${node.name}`}
        className={cn(
          "flex w-full items-center gap-2 px-3 py-2 text-left rounded-md",
          "hover:bg-muted/50 transition-colors",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        )}
        onClick={handleClick}
        style={{ paddingLeft: `${paddingLeft + 12}px` }}
      >
        {hasChildren &&
          (isExpanded ? (
            <ChevronDown className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
          ))}
        {!hasChildren && <span className="w-4" />}
        {isExpanded ? (
          <FolderOpen className="h-4 w-4 flex-shrink-0 text-amber-500" />
        ) : (
          <Folder className="h-4 w-4 flex-shrink-0 text-amber-500" />
        )}
        <span className={cn("flex-1 text-sm font-medium", !isExpanded && "truncate")}>
          {node.name}
        </span>
        <span className="text-xs text-muted-foreground flex-shrink-0">{node.testCount}</span>
      </button>
    );
  }

  const { suites } = node;
  const firstSuite = suites[0];

  return (
    <div className="flex flex-col">
      <button
        aria-expanded={isExpanded}
        aria-label={isExpanded ? `Collapse ${node.name}` : `Expand ${node.name}`}
        className={cn(
          "flex w-full items-center gap-2 px-3 py-2 text-left rounded-md",
          "hover:bg-muted/50 transition-colors",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
          isExpanded && "bg-accent/40"
        )}
        onClick={handleClick}
        style={{ paddingLeft: `${paddingLeft + 12}px` }}
      >
        {isExpanded ? (
          <ChevronDown className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
        ) : (
          <ChevronRight className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
        )}
        <FileText className="h-4 w-4 flex-shrink-0 text-blue-500" />
        <span className={cn("flex-1 text-sm", !isExpanded && "truncate")}>{node.name}</span>
        {firstSuite && <FrameworkBadge framework={firstSuite.framework} />}
        <StatusMiniBar counts={statusCounts} />
        <span className="text-xs text-muted-foreground flex-shrink-0">{node.testCount}</span>
      </button>
      {isExpanded && suites.length > 0 && (
        <div
          className="border-l border-border/50 ml-4"
          style={{ marginLeft: `${paddingLeft + 24}px` }}
        >
          {suites.flatMap((suite) =>
            suite.tests.map((test) => (
              <TestItem
                key={`${suite.filePath}-${suite.suiteName}-${test.line}-${test.name}`}
                test={test}
              />
            ))
          )}
        </div>
      )}
    </div>
  );
};
