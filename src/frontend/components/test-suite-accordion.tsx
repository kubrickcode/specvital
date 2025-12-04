"use client";

import { ChevronDown, ChevronRight, FileText } from "lucide-react";
import { useState } from "react";
import type { TestSuite } from "@/lib/api";
import { cn } from "@/lib/utils";
import { FrameworkBadge } from "./framework-badge";
import { TestItem } from "./test-item";

type TestSuiteAccordionProps = {
  suite: TestSuite;
};

export const TestSuiteAccordion = ({ suite }: TestSuiteAccordionProps) => {
  const [isExpanded, setIsExpanded] = useState(true);

  const toggleExpanded = () => {
    setIsExpanded((prev) => !prev);
  };

  const testCount = suite.tests.length;

  return (
    <div className="rounded-lg border bg-card shadow-sm">
      <button
        className={cn(
          "flex w-full items-center gap-3 px-4 py-3",
          "hover:bg-muted/50 transition-colors",
          "focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
        )}
        onClick={toggleExpanded}
        aria-expanded={isExpanded}
        aria-label={
          isExpanded
            ? `Collapse ${suite.filePath} test suite`
            : `Expand ${suite.filePath} test suite`
        }
      >
        {isExpanded ? (
          <ChevronDown className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
        ) : (
          <ChevronRight className="h-4 w-4 flex-shrink-0 text-muted-foreground" />
        )}
        <FileText className="h-5 w-5 flex-shrink-0 text-muted-foreground" />
        <span className="flex-1 text-left text-sm font-medium truncate">
          {suite.filePath}
        </span>
        <FrameworkBadge framework={suite.framework} />
        <span className="text-xs text-muted-foreground flex-shrink-0">
          {testCount} {testCount === 1 ? "test" : "tests"}
        </span>
      </button>

      {isExpanded && (
        <div className="border-t bg-muted/20 px-4 py-2">
          <div className="space-y-1 pl-6">
            {suite.tests.length === 0 ? (
              <div className="py-2 text-sm text-muted-foreground">No tests in this suite.</div>
            ) : (
              suite.tests.map((test) => (
                <TestItem key={`${test.line}-${test.name}`} test={test} />
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
};
