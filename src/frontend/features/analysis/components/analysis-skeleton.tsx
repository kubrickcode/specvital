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
  { border: string; bg: string; icon: string; defaultTitle: string; defaultDescription: string }
> = {
  loading: {
    border: "border-l-muted-foreground",
    bg: "bg-accent/30",
    icon: "text-muted-foreground",
    defaultTitle: "Loading",
    defaultDescription: "Preparing analysis...",
  },
  queued: {
    border: "border-l-chart-2",
    bg: "bg-chart-2/10",
    icon: "text-chart-2",
    defaultTitle: "Queued",
    defaultDescription: "Analysis will start shortly",
  },
  analyzing: {
    border: "border-l-chart-1",
    bg: "bg-chart-1/10",
    icon: "text-chart-1",
    defaultTitle: "Analyzing",
    defaultDescription: "Scanning test files...",
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
    <main className="container mx-auto px-4 py-8" aria-busy="true">
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
          className={cn(
            "rounded-lg border-l-4 px-4 py-3",
            config.border,
            config.bg
          )}
          role="status"
          aria-live="polite"
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
