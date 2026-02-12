"use client";

import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useEffect } from "react";
import { toast } from "sonner";

const ERROR_KEYS: Record<string, string> = {
  auth_failed: "authFailed",
  invalid_code: "invalidCode",
  invalid_state: "invalidState",
  missing_oauth_params: "missingOAuthParams",
  network_error: "networkError",
};

export function OAuthErrorHandler() {
  const t = useTranslations("auth.toast");
  const searchParams = useSearchParams();

  useEffect(() => {
    if (!searchParams) return;

    const error = searchParams.get("error");
    if (error) {
      const messageKey = ERROR_KEYS[error];
      const message = messageKey ? t(messageKey) : t("unknownError", { error });
      toast.error(message);

      // Clean up URL (remove error param)
      const url = new URL(window.location.href);
      url.searchParams.delete("error");
      window.history.replaceState({}, "", url.toString());
    }
  }, [searchParams, t]);

  return null;
}
