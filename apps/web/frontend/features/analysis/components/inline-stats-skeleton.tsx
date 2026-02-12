import { Skeleton } from "@/components/ui/skeleton";

export const InlineStatsSkeleton = () => {
  return (
    <div
      aria-label="Loading statistics"
      aria-live="polite"
      className="rounded-lg border border-border/60 bg-card/50 p-5"
      role="status"
    >
      <div className="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
        {/* Left: Total + Status */}
        <div className="flex flex-wrap items-center gap-x-6 gap-y-3">
          {/* Total */}
          <div className="flex items-baseline gap-2">
            <Skeleton className="h-10 w-20" />
            <Skeleton className="h-4 w-8" />
          </div>

          {/* Separator */}
          <div className="h-8 w-px bg-border/60 hidden sm:block" />

          {/* Status breakdown */}
          <div className="flex items-center gap-5">
            {[0, 1, 2].map((i) => (
              <span className="inline-flex items-center gap-2" key={i}>
                <Skeleton className="w-2.5 h-2.5 rounded-full" />
                <Skeleton className="h-4 w-10" />
                <Skeleton className="h-6 w-8" />
              </span>
            ))}
          </div>
        </div>

        {/* Right: Frameworks */}
        <div className="flex flex-wrap items-center gap-3">
          {[0, 1, 2].map((i) => (
            <Skeleton className="h-8 w-28 rounded-md" key={i} />
          ))}
        </div>
      </div>
    </div>
  );
};
