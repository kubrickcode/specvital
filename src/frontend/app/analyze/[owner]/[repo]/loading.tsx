export const Loading = () => {
  return (
    <main className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        <header className="space-y-2">
          <div className="h-8 w-64 animate-pulse rounded bg-muted" />
          <div className="h-5 w-48 animate-pulse rounded bg-muted" />
        </header>

        <div className="rounded-lg border bg-card p-6">
          <div className="space-y-4">
            <div className="h-4 w-full animate-pulse rounded bg-muted" />
            <div className="h-4 w-3/4 animate-pulse rounded bg-muted" />
            <div className="h-4 w-1/2 animate-pulse rounded bg-muted" />
          </div>
        </div>
      </div>
    </main>
  );
};

export default Loading;
