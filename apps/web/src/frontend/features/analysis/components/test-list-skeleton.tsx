import { Skeleton } from "@/components/ui/skeleton";

const SKELETON_ITEMS = 6;

const TestSuiteAccordionSkeleton = () => {
  return (
    <div className="rounded-lg border bg-card shadow-sm">
      <div className="flex items-center gap-3 px-4 py-3">
        {/* Chevron icon */}
        <Skeleton className="h-4 w-4" />

        {/* File icon */}
        <Skeleton className="h-5 w-5" />

        {/* File path */}
        <Skeleton className="h-4 flex-1 max-w-[60%]" />

        {/* Framework badge */}
        <Skeleton className="h-5 w-16 rounded-full" />

        {/* Test count */}
        <Skeleton className="h-3 w-12" />
      </div>
    </div>
  );
};

export const TestListSkeleton = () => {
  return (
    <div aria-label="Loading test suites" aria-live="polite" className="space-y-3" role="status">
      {Array.from({ length: SKELETON_ITEMS }).map((_, index) => (
        <TestSuiteAccordionSkeleton key={index} />
      ))}
    </div>
  );
};
