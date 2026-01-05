import { Skeleton } from "@/components/ui/skeleton";

const SKELETON_ITEMS = 6;

const SpecFileGroupSkeleton = () => {
  return (
    <div className="rounded-lg border bg-card shadow-sm">
      <div className="flex items-center gap-3 px-4 py-3">
        <Skeleton className="h-4 w-4" />
        <Skeleton className="h-5 w-5" />
        <Skeleton className="h-4 flex-1 max-w-[60%]" />
        <Skeleton className="h-5 w-16 rounded-full" />
        <Skeleton className="h-3 w-12" />
      </div>
    </div>
  );
};

export const SpecViewSkeleton = () => {
  return (
    <div aria-label="Loading spec view" aria-live="polite" className="space-y-3" role="status">
      {Array.from({ length: SKELETON_ITEMS }).map((_, index) => (
        <SpecFileGroupSkeleton key={index} />
      ))}
    </div>
  );
};
