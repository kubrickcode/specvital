import type { QueryClient } from "@tanstack/react-query";

export const AUTH_QUERY_KEY = "auth" as const;

export const authKeys = {
  all: [AUTH_QUERY_KEY] as const,
  user: () => [...authKeys.all, "user"] as const,
};

// Type guard to safely check if error has status property
const hasStatus = (error: unknown): error is Error & { status: number } => {
  return (
    error instanceof Error &&
    "status" in error &&
    typeof (error as Error & { status: unknown }).status === "number"
  );
};

export const isUnauthorizedError = (error: unknown): boolean => {
  return hasStatus(error) && error.status === 401;
};

export const handleUnauthorizedError = (queryClient: QueryClient): void => {
  if (typeof document !== "undefined") {
    document.cookie = "has_session=; path=/; max-age=0";
  }
  queryClient.setQueryData(authKeys.user(), null);
};
