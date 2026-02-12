"use client";

import { Skeleton } from "@/components/ui/skeleton";

const SKELETON_ITEMS_COUNT = 3;

export const SearchSkeleton = () => {
  return (
    <div className="py-2 px-1 motion-safe:animate-pulse">
      {Array.from({ length: SKELETON_ITEMS_COUNT }).map((_, index) => (
        <div className="flex items-center gap-2 px-2 py-1.5" key={index}>
          <Skeleton className="size-4 shrink-0 rounded" />
          <Skeleton className="h-4 flex-1 rounded" />
        </div>
      ))}
    </div>
  );
};
