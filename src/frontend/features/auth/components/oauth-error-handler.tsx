"use client";

import { useSearchParams } from "next/navigation";
import { useEffect } from "react";
import { toast } from "sonner";

const ERROR_MESSAGES: Record<string, string> = {
  auth_failed: "Authentication failed",
  invalid_code: "Invalid authorization code",
  invalid_state: "Session expired. Please try again.",
  missing_oauth_params: "Missing OAuth parameters",
  network_error: "Network error during authentication",
};

export function OAuthErrorHandler() {
  const searchParams = useSearchParams();

  useEffect(() => {
    if (!searchParams) return;

    const error = searchParams.get("error");
    if (error) {
      const message = ERROR_MESSAGES[error] || `Authentication error: ${error}`;
      toast.error(message);

      // Clean up URL (remove error param)
      const url = new URL(window.location.href);
      url.searchParams.delete("error");
      window.history.replaceState({}, "", url.toString());
    }
  }, [searchParams]);

  return null;
}
