import { Skeleton } from "@/components/ui/skeleton";

export const StatsCardSkeleton = () => {
  return (
    <div
      aria-label="Loading statistics"
      aria-live="polite"
      className="rounded-lg border bg-card p-6 shadow-xs"
      role="status"
    >
      {/* Header */}
      <Skeleton className="h-5 w-24 mb-4" />

      {/* Stats Grid */}
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        {/* Total - emphasized */}
        <div className="flex flex-col">
          <Skeleton className="h-10 w-16 mb-1" />
          <Skeleton className="h-4 w-12" />
        </div>

        {/* Active */}
        <div className="flex flex-col">
          <Skeleton className="h-9 w-14 mb-1" />
          <Skeleton className="h-4 w-12" />
        </div>

        {/* Skipped */}
        <div className="flex flex-col">
          <Skeleton className="h-9 w-14 mb-1" />
          <Skeleton className="h-4 w-14" />
        </div>

        {/* Todo */}
        <div className="flex flex-col">
          <Skeleton className="h-9 w-14 mb-1" />
          <Skeleton className="h-4 w-10" />
        </div>
      </div>

      {/* Framework Breakdown */}
      <div className="mt-6 pt-6 border-t">
        <Skeleton className="h-4 w-32 mb-4" />
        <div className="space-y-4">
          {[0, 1].map((index) => (
            <div className="space-y-2" key={index}>
              {/* Header row */}
              <div className="flex items-center justify-between">
                <span className="flex items-center gap-2">
                  <Skeleton className="h-3 w-3 rounded-full" />
                  <Skeleton className="h-4 w-20" />
                  <Skeleton className="h-4 w-10" />
                </span>
                <Skeleton className="h-4 w-16" />
              </div>

              {/* Progress bar */}
              <div className="h-2 w-full rounded-full bg-muted overflow-hidden">
                <Skeleton
                  className="h-full rounded-full"
                  style={{ width: index === 0 ? "65%" : "35%" }}
                />
              </div>

              {/* Stats row */}
              <Skeleton className="h-4 w-48" />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
