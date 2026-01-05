"use client";

import { ChevronDown, ChevronRight, FileText } from "lucide-react";
import { useMemo, useState } from "react";

import { cn } from "@/lib/utils";

import { SpecSuiteGroup } from "./spec-suite-group";
import type { ConvertedTestFile } from "../../types";
import { FrameworkBadge } from "../framework-badge";
import { StatusMiniBar } from "../status-mini-bar";

type SpecFileGroupProps = {
  file: ConvertedTestFile;
};

export const SpecFileGroup = ({ file }: SpecFileGroupProps) => {
  const [isExpanded, setIsExpanded] = useState(true);

  const { statusCounts, testCount } = useMemo(() => {
    const allTests = file.suites.flatMap((suite) => suite.tests);
    const counts = {
      active: 0,
      focused: 0,
      skipped: 0,
      todo: 0,
      xfail: 0,
    };
    for (const test of allTests) {
      counts[test.status]++;
    }
    return { statusCounts: counts, testCount: allTests.length };
  }, [file.suites]);

  return (
    <div
      className={cn(
        "rounded-lg border bg-card transition-shadow duration-200",
        isExpanded ? "shadow-md" : "shadow-sm"
      )}
    >
      <button
        aria-expanded={isExpanded}
        aria-label={isExpanded ? `Collapse ${file.filePath}` : `Expand ${file.filePath}`}
        className={cn(
          "flex w-full items-center gap-3 px-4 py-3 rounded-t-lg",
          "transition-all duration-200 ease-in-out",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
          isExpanded ? "bg-accent/40 hover:bg-accent/60" : "hover:bg-muted/70"
        )}
        onClick={() => setIsExpanded((prev) => !prev)}
      >
        {isExpanded ? (
          <ChevronDown className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
        ) : (
          <ChevronRight className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
        )}
        <FileText className="h-5 w-5 flex-shrink-0 text-muted-foreground" />
        <span className="flex-1 text-left text-sm font-medium truncate">{file.filePath}</span>
        <FrameworkBadge framework={file.framework} />
        <StatusMiniBar counts={statusCounts} />
        <span className="text-xs text-muted-foreground flex-shrink-0">
          {testCount} {testCount === 1 ? "test" : "tests"}
        </span>
      </button>

      {isExpanded && (
        <div className="border-t border-accent/20 bg-accent/10 px-4 py-3">
          <div className="space-y-3 pl-6">
            {file.suites.map((suite) => (
              <SpecSuiteGroup key={suite.suiteHierarchy} suite={suite} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
