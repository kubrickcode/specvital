const DEFAULT_SUITE_COUNT = 5;

type TestListSkeletonProps = {
  count?: number;
};

export const TestListSkeleton = ({ count = DEFAULT_SUITE_COUNT }: TestListSkeletonProps) => {
  return (
    <div className="space-y-4" role="status" aria-label="Loading test suites">
      {Array.from({ length: count }, (_, index) => (
        <div key={index} className="rounded-lg border bg-card shadow-sm">
          <div className="flex items-center gap-3 px-4 py-3">
            <div className="h-4 w-4 flex-shrink-0 animate-pulse rounded bg-muted" />
            <div className="h-5 w-5 flex-shrink-0 animate-pulse rounded bg-muted" />
            <div className="h-5 flex-1 animate-pulse rounded bg-muted" />
            <div className="h-6 w-16 flex-shrink-0 animate-pulse rounded-full bg-muted" />
            <div className="h-4 w-12 flex-shrink-0 animate-pulse rounded bg-muted" />
          </div>
        </div>
      ))}
      <span className="sr-only">Loading...</span>
    </div>
  );
};
