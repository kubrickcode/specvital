import { AlertCircle } from "lucide-react";

import { cn } from "@/lib/utils";

type ErrorFallbackProps = {
  action?: React.ReactNode;
  className?: string;
  description?: string;
  fullScreen?: boolean;
  title: string;
};

export const ErrorFallback = ({
  action,
  className,
  description,
  fullScreen = true,
  title,
}: ErrorFallbackProps) => {
  return (
    <main
      className={cn(
        "flex flex-col items-center justify-center p-8",
        fullScreen && "min-h-screen",
        className
      )}
    >
      <div className="flex flex-col items-center gap-4 text-center">
        <AlertCircle className="h-12 w-12 text-destructive" />
        <h1 className="text-xl font-semibold">{title}</h1>
        {description && <p className="text-muted-foreground max-w-md">{description}</p>}
        {action}
      </div>
    </main>
  );
};
