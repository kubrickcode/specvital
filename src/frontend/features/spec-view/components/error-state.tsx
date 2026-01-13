"use client";

import { AlertCircle, RefreshCw, WifiOff } from "lucide-react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

import type { ErrorType } from "../utils/error";

type ErrorStateProps = {
  className?: string;
  errorMessage?: string;
  errorType?: ErrorType;
  onRetry?: () => void;
};

const ERROR_CONFIG: Record<
  ErrorType,
  {
    description: string;
    icon: typeof AlertCircle;
    retryLabel: string;
    title: string;
  }
> = {
  generation: {
    description: "An error occurred while generating the specification document. Please try again.",
    icon: AlertCircle,
    retryLabel: "Retry Generation",
    title: "Generation Failed",
  },
  network: {
    description:
      "Unable to connect to the server. Please check your internet connection and try again.",
    icon: WifiOff,
    retryLabel: "Retry",
    title: "Connection Error",
  },
  unknown: {
    description: "An unexpected error occurred. Please try again or contact support.",
    icon: AlertCircle,
    retryLabel: "Try Again",
    title: "Something Went Wrong",
  },
};

export const ErrorState = ({
  className,
  errorMessage,
  errorType = "unknown",
  onRetry,
}: ErrorStateProps) => {
  const config = ERROR_CONFIG[errorType];
  const Icon = config.icon;

  return (
    <Alert className={cn("", className)} role="alert" variant="destructive">
      <Icon aria-hidden="true" className="h-4 w-4" />
      <AlertTitle>{config.title}</AlertTitle>
      <AlertDescription className="flex flex-col gap-4">
        <p>{errorMessage || config.description}</p>
        {onRetry && (
          <Button className="w-fit" onClick={onRetry} size="sm" type="button" variant="outline">
            <RefreshCw aria-hidden="true" className="h-4 w-4 mr-2" />
            {config.retryLabel}
          </Button>
        )}
      </AlertDescription>
    </Alert>
  );
};
