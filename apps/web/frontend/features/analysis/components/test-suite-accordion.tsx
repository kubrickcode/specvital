"use client";

import { ChevronDown, ChevronRight, FileText } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { FrameworkBadge } from "@/components/ui/framework-badge";
import type { TestSuite } from "@/lib/api";
import { collapseTransition, expandCollapse, useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

import { StatusMiniBar } from "./status-mini-bar";
import { TestItem } from "./test-item";
import { calculateStatusCounts } from "../utils/calculate-status-counts";

type TestSuiteAccordionProps = {
  suite: TestSuite;
};

export const TestSuiteAccordion = ({ suite }: TestSuiteAccordionProps) => {
  const [isExpanded, setIsExpanded] = useState(true);
  const shouldReduceMotion = useReducedMotion();
  const t = useTranslations("analyze");

  const toggleExpanded = () => {
    setIsExpanded((prev) => !prev);
  };

  const testCount = suite.tests.length;
  const statusCounts = calculateStatusCounts(suite.tests);

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
          "w-full px-4 py-3 rounded-t-lg",
          "transition-all duration-200 ease-in-out",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
          isExpanded ? "bg-accent/40 hover:bg-accent/60" : "hover:bg-muted/70"
        )}
        onClick={toggleExpanded}
      >
        {/* Desktop: single row */}
        <div className="hidden sm:flex items-center gap-3">
          {isExpanded ? (
            <ChevronDown className="h-4 w-4 flex-shrink-0 text-muted-foreground transition-transform duration-200" />
          ) : (
            <ChevronRight className="h-4 w-4 flex-shrink-0 text-muted-foreground transition-transform duration-200" />
          )}
          <FileText className="h-5 w-5 flex-shrink-0 text-muted-foreground" />
          <span
            className={cn(
              "min-w-0 flex-1 text-left text-sm font-medium",
              !isExpanded && "truncate"
            )}
            title={suite.filePath + (suite.suiteName ? ` › ${suite.suiteName}` : "")}
          >
            {suite.filePath}
            {suite.suiteName && (
              <span className="ml-2 text-muted-foreground font-normal">› {suite.suiteName}</span>
            )}
          </span>
          <FrameworkBadge framework={suite.framework} />
          <StatusMiniBar counts={statusCounts} />
          <span className="text-xs text-muted-foreground flex-shrink-0">
            {t("testCount", { count: testCount })}
          </span>
        </div>
        {/* Mobile: icon area + content area */}
        <div className="sm:hidden flex items-start gap-3">
          {/* Icon area - top aligned */}
          <div className="flex items-center gap-2 pt-0.5">
            {isExpanded ? (
              <ChevronDown className="h-4 w-4 flex-shrink-0 text-muted-foreground transition-transform duration-200" />
            ) : (
              <ChevronRight className="h-4 w-4 flex-shrink-0 text-muted-foreground transition-transform duration-200" />
            )}
            <FileText className="h-5 w-5 flex-shrink-0 text-muted-foreground" />
          </div>
          {/* Content area */}
          <div className="min-w-0 flex-1">
            <span
              className={cn("block text-left text-sm font-medium", !isExpanded && "truncate")}
              title={suite.filePath + (suite.suiteName ? ` › ${suite.suiteName}` : "")}
            >
              {suite.filePath}
              {suite.suiteName && (
                <span className="ml-2 text-muted-foreground font-normal">› {suite.suiteName}</span>
              )}
            </span>
            <div className="mt-1.5 flex items-center gap-2">
              <FrameworkBadge framework={suite.framework} />
              <StatusMiniBar counts={statusCounts} />
              <span className="text-xs text-muted-foreground">
                {t("testCount", { count: testCount })}
              </span>
            </div>
          </div>
        </div>
      </button>

      <AnimatePresence initial={false}>
        {isExpanded && (
          <motion.div
            animate={shouldReduceMotion ? {} : "expanded"}
            exit={shouldReduceMotion ? {} : "collapsed"}
            initial={shouldReduceMotion ? false : "collapsed"}
            transition={shouldReduceMotion ? undefined : collapseTransition}
            variants={shouldReduceMotion ? undefined : expandCollapse}
          >
            <div className="border-t border-accent/20 bg-accent/10 px-4 py-2">
              <div className="space-y-1 pl-6">
                {suite.tests.length === 0 ? (
                  <div className="py-2 text-sm text-muted-foreground">{t("noTestsInSuite")}</div>
                ) : (
                  suite.tests.map((test) => (
                    <TestItem key={`${test.line}-${test.name}`} test={test} />
                  ))
                )}
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
