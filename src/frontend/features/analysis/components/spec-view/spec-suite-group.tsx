"use client";

import { ChevronDown, ChevronRight } from "lucide-react";
import { useMemo, useState } from "react";

import { cn } from "@/lib/utils";

import { SpecItem } from "./spec-item";
import type { ConvertedTestSuite } from "../../types";
import { calculateStatusCounts } from "../../utils/calculate-status-counts";
import { StatusMiniBar } from "../status-mini-bar";

type SpecSuiteGroupProps = {
  suite: ConvertedTestSuite;
};

export const SpecSuiteGroup = ({ suite }: SpecSuiteGroupProps) => {
  const [isExpanded, setIsExpanded] = useState(true);

  const testCount = suite.tests.length;
  const statusCounts = useMemo(() => calculateStatusCounts(suite.tests), [suite.tests]);

  return (
    <div className="border-l-2 border-muted pl-3">
      <button
        aria-expanded={isExpanded}
        className={cn(
          "flex w-full items-center gap-2 py-2 text-left",
          "hover:bg-muted/30 rounded-md px-2 -ml-2",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        )}
        onClick={() => setIsExpanded((prev) => !prev)}
      >
        {isExpanded ? (
          <ChevronDown className="h-4 w-4 text-muted-foreground" />
        ) : (
          <ChevronRight className="h-4 w-4 text-muted-foreground" />
        )}
        <span className="flex-1 text-sm font-medium truncate">{suite.suiteName}</span>
        <StatusMiniBar counts={statusCounts} />
        <span className="text-xs text-muted-foreground">
          {testCount} {testCount === 1 ? "test" : "tests"}
        </span>
      </button>

      {isExpanded && (
        <div className="mt-1 space-y-0.5">
          {suite.tests.map((item) => (
            <SpecItem item={item} key={`${item.line}-${item.originalName}`} />
          ))}
        </div>
      )}
    </div>
  );
};
