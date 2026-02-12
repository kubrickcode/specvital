"use client";

import { motion } from "motion/react";

import type { AnalysisResult } from "@/lib/api";
import { createStaggerContainer, fadeInUp, useReducedMotion } from "@/lib/motion";

import { AnalysisHeader } from "./analysis-header";
import { InlineStats } from "./inline-stats";
import { SpecPanel } from "./spec-panel";
import { TabNavigation } from "./tab-navigation";
import { TestsPanel } from "./tests-panel";
import { UpdateBanner } from "./update-banner";
import { useCommitSelect } from "../hooks/use-commit-select";
import { usePrimaryTab } from "../hooks/use-primary-tab";

type AnalysisContentProps = {
  result: AnalysisResult;
};

const pageStaggerContainer = createStaggerContainer(0.08, 0);

export const AnalysisContent = ({ result }: AnalysisContentProps) => {
  const { setTab, tab } = usePrimaryTab();
  const { selectCommit } = useCommitSelect();
  const shouldReduceMotion = useReducedMotion();

  const containerVariants = shouldReduceMotion ? {} : pageStaggerContainer;
  const itemVariants = shouldReduceMotion ? {} : fadeInUp;

  const availableFrameworks = result.summary.frameworks.map((f) => f.framework);

  return (
    <motion.main
      animate="visible"
      className="container mx-auto px-4 py-6"
      initial={shouldReduceMotion ? false : "hidden"}
      variants={containerVariants}
    >
      <UpdateBanner owner={result.owner} repo={result.repo} />

      <div className="space-y-6">
        {/* Header Section */}
        <motion.div variants={itemVariants}>
          <AnalysisHeader
            analyzedAt={result.analyzedAt}
            branchName={result.branchName}
            commitSha={result.commitSha}
            committedAt={result.committedAt}
            onCommitSelect={selectCommit}
            owner={result.owner}
            parserVersion={result.parserVersion}
            repo={result.repo}
          />
        </motion.div>

        {/* Stats Section */}
        <motion.div variants={itemVariants}>
          <InlineStats summary={result.summary} />
        </motion.div>

        {/* Content Section */}
        <motion.section
          className="rounded-lg border border-border/60 bg-card/30"
          variants={itemVariants}
        >
          <div className="px-5 py-4 border-b border-border/40">
            <TabNavigation activeTab={tab} onTabChange={setTab} />
          </div>

          {tab === "tests" ? (
            <TestsPanel
              availableFrameworks={availableFrameworks}
              data={result}
              suites={result.suites}
              totalCount={result.summary.total}
            />
          ) : (
            <SpecPanel
              analysisCommitSha={result.commitSha}
              analysisId={result.id}
              availableFrameworks={availableFrameworks}
              owner={result.owner}
              repo={result.repo}
              totalTests={result.summary.total}
            />
          )}
        </motion.section>
      </div>
    </motion.main>
  );
};
