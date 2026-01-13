"use client";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

type DocumentSkeletonProps = {
  domainCount?: number;
};

const BehaviorSkeleton = () => (
  <div className="flex items-start gap-3 px-3 py-2.5">
    <Skeleton className="h-4 w-4 flex-shrink-0 rounded-full mt-0.5" />
    <div className="flex-1 space-y-2">
      <Skeleton className="h-4 w-full" />
      <div className="flex items-center gap-2">
        <Skeleton className="h-3 w-32" />
        <Skeleton className="h-3 w-12" />
        <Skeleton className="h-4 w-16 rounded-full" />
      </div>
    </div>
  </div>
);

const FeatureSkeleton = ({ behaviorCount = 3 }: { behaviorCount?: number }) => (
  <div className="border-l-2 border-muted-foreground/20 ml-2">
    <div className="flex items-center gap-2 px-3 py-2">
      <Skeleton className="h-4 w-4" />
      <Skeleton className="h-5 w-48" />
      <Skeleton className="h-5 w-8 rounded-full ml-auto" />
    </div>
    <div className="pl-4 space-y-0.5">
      {Array.from({ length: behaviorCount }).map((_, i) => (
        <BehaviorSkeleton key={i} />
      ))}
    </div>
  </div>
);

const DomainSkeleton = ({ featureCount = 2 }: { featureCount?: number }) => (
  <Card className="overflow-hidden">
    <CardHeader className="pb-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Skeleton className="h-5 w-5" />
          <Skeleton className="h-5 w-5" />
          <Skeleton className="h-6 w-40" />
        </div>
        <div className="flex items-center gap-2">
          <Skeleton className="h-5 w-12 rounded-full" />
          <Skeleton className="h-5 w-24 rounded-full" />
        </div>
      </div>
      <Skeleton className="h-4 w-3/4 mt-2 ml-12" />
    </CardHeader>
    <CardContent className="pt-0 space-y-3">
      {Array.from({ length: featureCount }).map((_, i) => (
        <FeatureSkeleton behaviorCount={2 + i} key={i} />
      ))}
    </CardContent>
  </Card>
);

const ExecutiveSummarySkeleton = () => (
  <Card>
    <CardHeader className="pb-3">
      <div className="flex items-center gap-2">
        <Skeleton className="h-5 w-5" />
        <Skeleton className="h-6 w-32" />
      </div>
    </CardHeader>
    <CardContent>
      <div className="space-y-2">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-3/4" />
      </div>
    </CardContent>
  </Card>
);

const TocSidebarSkeleton = () => (
  <aside className="hidden lg:block w-64 flex-shrink-0">
    <div className="sticky top-4 space-y-4">
      <Skeleton className="h-6 w-32" />
      <div className="space-y-2">
        {Array.from({ length: 4 }).map((_, i) => (
          <div className="space-y-1" key={i}>
            <Skeleton className="h-5 w-full" />
            <div className="pl-4 space-y-1">
              <Skeleton className="h-4 w-3/4" />
              <Skeleton className="h-4 w-2/3" />
            </div>
          </div>
        ))}
      </div>
    </div>
  </aside>
);

export const DocumentSkeleton = ({ domainCount = 3 }: DocumentSkeletonProps) => {
  return (
    <div
      aria-busy="true"
      aria-label="Loading specification document"
      className="flex gap-6"
      role="status"
    >
      <TocSidebarSkeleton />

      <div className="flex-1 space-y-6 min-w-0">
        <ExecutiveSummarySkeleton />

        {/* Search & Filter placeholder */}
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
          <Skeleton className="h-10 w-full sm:w-80" />
          <Skeleton className="h-10 flex-1" />
        </div>

        {/* Domain cards */}
        <div className="space-y-4">
          {Array.from({ length: domainCount }).map((_, i) => (
            <DomainSkeleton featureCount={2 + (i % 2)} key={i} />
          ))}
        </div>
      </div>
    </div>
  );
};
