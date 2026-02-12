import { Loader2 } from "lucide-react";

import { cn } from "@/lib/utils";

type LoadingFallbackProps = {
  className?: string;
  fullScreen?: boolean;
  message?: string;
};

export const LoadingFallback = ({
  className,
  fullScreen = true,
  message,
}: LoadingFallbackProps) => {
  return (
    <main
      className={cn(
        "flex flex-col items-center justify-center p-8",
        fullScreen && "flex-1",
        className
      )}
    >
      <div className="flex flex-col items-center gap-4">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        {message && (
          <p aria-live="polite" className="text-muted-foreground">
            {message}
          </p>
        )}
      </div>
    </main>
  );
};
