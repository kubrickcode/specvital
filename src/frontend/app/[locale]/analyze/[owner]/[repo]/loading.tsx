import { StatsCardSkeleton, TestListSkeleton } from "@/components/skeletons";

export const Loading = () => {
  return (
    <main className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        <header className="space-y-2">
          <div className="h-8 w-64 animate-pulse rounded bg-muted" />
          <div className="h-5 w-48 animate-pulse rounded bg-muted" />
        </header>

        <div className="space-y-6">
          <StatsCardSkeleton />

          <div className="space-y-4">
            <div className="h-7 w-48 animate-pulse rounded bg-muted" />
            <TestListSkeleton />
          </div>
        </div>
      </div>
    </main>
  );
};

export default Loading;
