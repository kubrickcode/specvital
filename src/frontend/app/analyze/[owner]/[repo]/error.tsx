"use client";

import { AlertCircle, Ban, Clock, FileQuestion, Lock, ServerCrash } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ApiError, ERROR_TYPES, getErrorMessage } from "@/lib/api";

type ErrorProps = {
  error: Error & { digest?: string };
  reset: () => void;
};

const getErrorIcon = (error: Error) => {
  if (!(error instanceof ApiError)) {
    return <AlertCircle className="h-12 w-12 text-destructive" />;
  }

  switch (error.errorType) {
    case ERROR_TYPES.BAD_REQUEST:
    case ERROR_TYPES.VALIDATION_ERROR:
      return <AlertCircle className="h-12 w-12 text-amber-600" />;
    case ERROR_TYPES.FORBIDDEN:
      return <Lock className="h-12 w-12 text-amber-600" />;
    case ERROR_TYPES.NOT_FOUND:
      return <FileQuestion className="h-12 w-12 text-amber-600" />;
    case ERROR_TYPES.RATE_LIMIT:
      return <Clock className="h-12 w-12 text-amber-600" />;
    case ERROR_TYPES.NETWORK:
      return <Ban className="h-12 w-12 text-destructive" />;
    case ERROR_TYPES.SERVER_ERROR:
      return <ServerCrash className="h-12 w-12 text-destructive" />;
    default:
      return <AlertCircle className="h-12 w-12 text-destructive" />;
  }
};

const getErrorTitle = (error: Error) => {
  if (!(error instanceof ApiError)) {
    return "Something went wrong";
  }

  switch (error.errorType) {
    case ERROR_TYPES.BAD_REQUEST:
      return "Invalid Request";
    case ERROR_TYPES.FORBIDDEN:
      return "Access Denied";
    case ERROR_TYPES.NOT_FOUND:
      return "Repository Not Found";
    case ERROR_TYPES.RATE_LIMIT:
      return "Rate Limit Exceeded";
    case ERROR_TYPES.NETWORK:
      return "Network Error";
    case ERROR_TYPES.SERVER_ERROR:
      return "Server Error";
    case ERROR_TYPES.VALIDATION_ERROR:
      return "Validation Error";
    case ERROR_TYPES.PARSE_ERROR:
      return "Parse Error";
    default:
      return "Error";
  }
};

const shouldShowRetry = (error: Error) => {
  if (!(error instanceof ApiError)) {
    return true;
  }

  // Don't show retry for permanent errors
  switch (error.errorType) {
    case ERROR_TYPES.BAD_REQUEST:
    case ERROR_TYPES.FORBIDDEN:
    case ERROR_TYPES.NOT_FOUND:
    case ERROR_TYPES.VALIDATION_ERROR:
      return false;
    default:
      return true;
  }
};

// Next.js App Router requires default export for error boundaries
export default function Error({ error, reset }: ErrorProps) {
  const errorMessage = getErrorMessage(error);
  const errorTitle = getErrorTitle(error);
  const icon = getErrorIcon(error);
  const showRetry = shouldShowRetry(error);

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <div className="w-full max-w-md space-y-6 text-center">
        <div className="flex justify-center">{icon}</div>

        <div className="space-y-2">
          <h1 className="text-2xl font-bold">{errorTitle}</h1>
          <p className="text-muted-foreground">{errorMessage}</p>

          {error instanceof ApiError && error.status && (
            <p className="text-xs text-muted-foreground">Error code: {error.status}</p>
          )}
        </div>

        <div className="flex flex-col gap-3 sm:flex-row sm:justify-center">
          {showRetry && (
            <Button onClick={reset} variant="default">
              Try again
            </Button>
          )}

          <Button asChild variant="outline">
            <Link href="/">Go home</Link>
          </Button>
        </div>

        {process.env.NODE_ENV === "development" && error.digest && (
          <details className="mt-4 text-left">
            <summary className="cursor-pointer text-sm text-muted-foreground">
              Technical details
            </summary>
            <pre className="mt-2 overflow-auto rounded-md bg-muted p-3 text-xs">
              {error.digest}
            </pre>
          </details>
        )}
      </div>
    </main>
  );
}
