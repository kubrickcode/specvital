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
        fullScreen && "min-h-screen",
        className
      )}
    >
      <div className="flex flex-col items-center gap-4">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        {message && (
          <p className="text-muted-foreground" aria-live="polite">
            {message}
          </p>
        )}
      </div>
    </main>
  );
};
