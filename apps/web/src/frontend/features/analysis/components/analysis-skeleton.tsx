import { Skeleton } from "@/components/ui/skeleton";

import type { AnalysisStatus } from "../types";
import { AnalysisWaitingCard } from "./analysis-waiting-card";
import { InlineStatsSkeleton } from "./inline-stats-skeleton";
import { TestListSkeleton } from "./test-list-skeleton";

type AnalysisSkeletonProps = {
  owner?: string;
  repo?: string;
  startedAt?: string | null;
  status?: AnalysisStatus;
};

export const AnalysisSkeleton = ({
  owner,
  repo,
  startedAt,
  status = "loading",
}: AnalysisSkeletonProps) => {
  // Render AnalysisWaitingCard for queued/analyzing states (backend status known)
  if ((status === "queued" || status === "analyzing") && owner && repo) {
    return (
      <AnalysisWaitingCard
        owner={owner}
        repo={repo}
        startedAt={startedAt ?? null}
        status={status}
      />
    );
  }

  // Pure skeleton for loading state (waiting for backend response)
  return (
    <main aria-busy="true" className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        {/* Header skeleton */}
        <header className="space-y-2">
          <div className="flex items-center justify-between">
            <Skeleton className="h-8 w-48" />
            <Skeleton className="h-9 w-32" />
          </div>
          <div className="flex items-center gap-4">
            <Skeleton className="h-4 w-20" />
            <Skeleton className="h-4 w-36" />
          </div>
        </header>

        {/* Loading indicator - simple skeleton bar instead of status text */}
        <Skeleton className="h-16 w-full rounded-lg" />

        <InlineStatsSkeleton />

        {/* Test suites section */}
        <section className="space-y-4">
          <Skeleton className="h-7 w-32" />
          <TestListSkeleton />
        </section>
      </div>
    </main>
  );
};
