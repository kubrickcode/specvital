import { Loader2 } from "lucide-react";

import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

import { StatsCardSkeleton } from "./stats-card-skeleton";
import { TestListSkeleton } from "./test-list-skeleton";

type AnalysisStatus = "loading" | "queued" | "analyzing";

type AnalysisSkeletonProps = {
  description?: string;
  status?: AnalysisStatus;
  title?: string;
};

const STATUS_CONFIG: Record<
  AnalysisStatus,
  { bg: string; border: string; defaultDescription: string; defaultTitle: string; icon: string }
> = {
  analyzing: {
    bg: "bg-chart-1/10",
    border: "border-l-chart-1",
    defaultDescription: "Scanning test files...",
    defaultTitle: "Analyzing",
    icon: "text-chart-1",
  },
  loading: {
    bg: "bg-accent/30",
    border: "border-l-muted-foreground",
    defaultDescription: "Preparing analysis...",
    defaultTitle: "Loading",
    icon: "text-muted-foreground",
  },
  queued: {
    bg: "bg-chart-2/10",
    border: "border-l-chart-2",
    defaultDescription: "Analysis will start shortly",
    defaultTitle: "Queued",
    icon: "text-chart-2",
  },
};

export const AnalysisSkeleton = ({
  description,
  status = "loading",
  title,
}: AnalysisSkeletonProps) => {
  const config = STATUS_CONFIG[status];
  const displayTitle = title ?? config.defaultTitle;
  const displayDescription = description ?? config.defaultDescription;

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

        {/* Status banner */}
        <div
          aria-live="polite"
          className={cn("rounded-lg border-l-4 px-4 py-3", config.border, config.bg)}
          role="status"
        >
          <div className="flex items-center gap-3">
            <Loader2 className={cn("h-5 w-5 animate-spin shrink-0", config.icon)} />
            <div className="min-w-0">
              <p className="font-medium text-foreground">{displayTitle}</p>
              <p className="text-sm text-muted-foreground">{displayDescription}</p>
            </div>
          </div>
        </div>

        <StatsCardSkeleton />

        {/* Test suites section */}
        <section className="space-y-4">
          <Skeleton className="h-7 w-32" />
          <TestListSkeleton />
        </section>
      </div>
    </main>
  );
};
