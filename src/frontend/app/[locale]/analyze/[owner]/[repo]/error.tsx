"use client";

import { AlertCircle, Ban, Clock, FileQuestion, Lock, ServerCrash } from "lucide-react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { RateLimitWarning } from "@/components/rate-limit-warning";
import { ApiError, ERROR_TYPES, getErrorMessage } from "@/lib/api";

const ERROR_TITLE_KEYS = {
  ACCESS_DENIED: "accessDenied",
  DEFAULT: "error",
  INVALID_REQUEST: "invalidRequest",
  NETWORK_ERROR: "networkError",
  NOT_FOUND: "notFound",
  PARSE_ERROR: "parseError",
  RATE_LIMIT_EXCEEDED: "rateLimitExceeded",
  SERVER_ERROR: "serverError",
  SOMETHING_WENT_WRONG: "somethingWentWrong",
  VALIDATION_ERROR: "validationError",
} as const;

type ErrorProps = {
  error: Error & { digest?: string };
  reset: () => void;
};

// Next.js App Router requires default export for error boundaries
const Error = ({ error, reset }: ErrorProps) => {
  const t = useTranslations("errors");
  const tCommon = useTranslations("common");

  const errorMessage = getErrorMessage(error);
  const errorTitleKey = getErrorTitleKey(error);
  const icon = getErrorIcon(error);
  const showRetry = shouldShowRetry(error);

  const rateLimit = error instanceof ApiError ? error.problemDetail?.rateLimit : undefined;

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <div className="w-full max-w-md space-y-6 text-center">
        <div className="flex justify-center">{icon}</div>

        <div className="space-y-2">
          <h1 className="text-2xl font-bold">{t(errorTitleKey)}</h1>
          <p className="text-muted-foreground">{errorMessage}</p>

          {error instanceof ApiError && error.status && (
            <p className="text-xs text-muted-foreground">
              {t("errorCode", { code: error.status })}
            </p>
          )}
        </div>

        {rateLimit && <RateLimitWarning rateLimit={rateLimit} />}

        <div className="flex flex-col gap-3 sm:flex-row sm:justify-center">
          {showRetry && (
            <Button onClick={reset} variant="default">
              {tCommon("tryAgain")}
            </Button>
          )}

          <Button asChild variant="outline">
            <Link href="/">{tCommon("goHome")}</Link>
          </Button>
        </div>

        {process.env.NODE_ENV === "development" && error.digest && (
          <details className="mt-4 text-left">
            <summary className="cursor-pointer text-sm text-muted-foreground">
              {t("technicalDetails")}
            </summary>
            <pre className="mt-2 overflow-auto rounded-md bg-muted p-3 text-xs">
              {error.digest}
            </pre>
          </details>
        )}
      </div>
    </main>
  );
};

export default Error;

// Helper functions

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

const getErrorTitleKey = (error: Error) => {
  if (!(error instanceof ApiError)) {
    return ERROR_TITLE_KEYS.SOMETHING_WENT_WRONG;
  }

  switch (error.errorType) {
    case ERROR_TYPES.BAD_REQUEST:
      return ERROR_TITLE_KEYS.INVALID_REQUEST;
    case ERROR_TYPES.FORBIDDEN:
      return ERROR_TITLE_KEYS.ACCESS_DENIED;
    case ERROR_TYPES.NOT_FOUND:
      return ERROR_TITLE_KEYS.NOT_FOUND;
    case ERROR_TYPES.RATE_LIMIT:
      return ERROR_TITLE_KEYS.RATE_LIMIT_EXCEEDED;
    case ERROR_TYPES.NETWORK:
      return ERROR_TITLE_KEYS.NETWORK_ERROR;
    case ERROR_TYPES.SERVER_ERROR:
      return ERROR_TITLE_KEYS.SERVER_ERROR;
    case ERROR_TYPES.VALIDATION_ERROR:
      return ERROR_TITLE_KEYS.VALIDATION_ERROR;
    case ERROR_TYPES.PARSE_ERROR:
      return ERROR_TITLE_KEYS.PARSE_ERROR;
    default:
      return ERROR_TITLE_KEYS.DEFAULT;
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
