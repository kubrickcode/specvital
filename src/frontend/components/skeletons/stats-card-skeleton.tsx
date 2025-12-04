import { cn } from "@/lib/utils";

const DEFAULT_STATS_COUNT = 4;

type StatsCardSkeletonProps = {
  count?: number;
};

export const StatsCardSkeleton = ({ count = DEFAULT_STATS_COUNT }: StatsCardSkeletonProps) => {
  return (
    <div
      className={cn("rounded-lg border bg-card p-6 shadow-xs")}
      role="status"
      aria-label="Loading statistics"
    >
      <div className="mb-4 h-7 w-32 animate-pulse rounded bg-muted" />
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        {Array.from({ length: count }, (_, index) => (
          <div key={index} className="flex flex-col gap-2">
            <div className="h-9 w-16 animate-pulse rounded bg-muted" />
            <div className="h-5 w-12 animate-pulse rounded bg-muted" />
          </div>
        ))}
      </div>
      <span className="sr-only">Loading...</span>
    </div>
  );
};
