import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

export const RepositorySkeleton = () => {
  return (
    <Card aria-label="Loading repository" className="relative h-full p-4" role="status">
      <div className="flex items-start justify-between gap-2 mb-3">
        <Skeleton className="h-5 w-32" />
        <Skeleton className="size-8 rounded-md" />
      </div>

      <div className="space-y-3">
        <div className="flex items-baseline gap-2">
          <Skeleton className="h-8 w-12" />
          <Skeleton className="h-4 w-10" />
          <Skeleton className="h-5 w-10 rounded-full" />
        </div>

        <div className="flex items-center gap-2">
          <Skeleton className="h-5 w-8 rounded-full" />
        </div>

        <div className="flex items-center justify-between pt-2 border-t">
          <Skeleton className="h-3 w-16" />
          <Skeleton className="h-7 w-16 rounded-md" />
        </div>
      </div>
    </Card>
  );
};
