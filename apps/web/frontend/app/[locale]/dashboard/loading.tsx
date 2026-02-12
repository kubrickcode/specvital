import { Skeleton } from "@/components/ui/skeleton";
import { RepositorySkeleton } from "@/features/dashboard";

const DashboardLoading = () => {
  return (
    <main className="container mx-auto px-4 py-8">
      <Skeleton className="mb-6 h-8 w-40" />

      <div className="space-y-8">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <Skeleton className="h-10 w-full sm:max-w-xs" />
          <Skeleton className="h-10 w-32" />
        </div>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {Array.from({ length: 8 }).map((_, index) => (
            <RepositorySkeleton key={index} />
          ))}
        </div>
      </div>
    </main>
  );
};

export default DashboardLoading;
