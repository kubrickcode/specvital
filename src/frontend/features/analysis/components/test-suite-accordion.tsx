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
    <div
      className={cn(
        "rounded-lg border bg-card transition-shadow duration-200",
        isExpanded ? "shadow-md" : "shadow-sm"
      )}
    >
      <button
        aria-expanded={isExpanded}
        aria-label={
          isExpanded
            ? `Collapse ${suite.filePath} test suite`
            : `Expand ${suite.filePath} test suite`
        }
        className={cn(
          "flex w-full items-center gap-3 px-4 py-3 rounded-t-lg",
          "transition-all duration-200 ease-in-out",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
          isExpanded ? "bg-accent/40 hover:bg-accent/60" : "hover:bg-muted/70"
        )}
        onClick={toggleExpanded}
      >
        {isExpanded ? (
          <ChevronDown className="h-4 w-4 flex-shrink-0 text-muted-foreground transition-transform duration-200" />
        ) : (
          <ChevronRight className="h-4 w-4 flex-shrink-0 text-muted-foreground transition-transform duration-200" />
        )}
        <FileText className="h-5 w-5 flex-shrink-0 text-muted-foreground" />
        <span className="flex-1 text-left text-sm font-medium truncate">
          {suite.filePath}
          {suite.suiteName && (
            <span className="ml-2 text-muted-foreground font-normal">› {suite.suiteName}</span>
          )}
        </span>
        <FrameworkBadge framework={suite.framework} />
        <span className="text-xs text-muted-foreground flex-shrink-0">
          {testCount} {testCount === 1 ? "test" : "tests"}
        </span>
      </button>

      {isExpanded && (
        <div className="border-t border-accent/20 bg-accent/10 px-4 py-2">
          <div className="space-y-1 pl-6">
            {suite.tests.length === 0 ? (
              <div className="py-2 text-sm text-muted-foreground">No tests in this suite.</div>
            ) : (
              suite.tests.map((test) => <TestItem key={`${test.line}-${test.name}`} test={test} />)
            )}
          </div>
        </div>
      )}
    </div>
  );
};
