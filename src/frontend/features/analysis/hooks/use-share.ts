"use client";

import { useTranslations } from "next-intl";
import { useCallback, useRef, useState } from "react";
import { toast } from "sonner";

import { shouldUseNativeShare } from "@/lib/utils/device";

const SUCCESS_DISPLAY_DURATION_MS = 2000;

type ShareState = "idle" | "success";

type UseShareReturn = {
  copyToClipboard: () => Promise<void>;
  isNativeShareSupported: boolean;
  shareNative: () => Promise<void>;
  state: ShareState;
};

export const useShare = (): UseShareReturn => {
  const t = useTranslations("share");
  const [state, setState] = useState<ShareState>("idle");
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const isNativeShareSupported = shouldUseNativeShare();

  const clearSuccessTimeout = () => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
  };

  const copyToClipboard = useCallback(async () => {
    const url = window.location.href;

    if (!navigator.clipboard) {
      toast.error(t("copyFailed"), {
        description: t("copyFailedDescription"),
      });
      return;
    }

    clearSuccessTimeout();

    try {
      await navigator.clipboard.writeText(url);
      setState("success");
      toast.success(t("linkCopied"), {
        description: t("linkCopiedDescription"),
      });

      timeoutRef.current = setTimeout(() => {
        setState("idle");
        timeoutRef.current = null;
      }, SUCCESS_DISPLAY_DURATION_MS);
    } catch (error) {
      console.error("Failed to copy to clipboard:", error);
      setState("idle");
      toast.error(t("copyFailed"), {
        description: t("copyFailedDescription"),
      });
    }
  }, [t]);

  const shareNative = useCallback(async () => {
    if (!isNativeShareSupported) {
      return copyToClipboard();
    }

    try {
      await navigator.share({
        title: document.title,
        url: window.location.href,
      });
    } catch (error) {
      if (error instanceof Error && error.name === "AbortError") {
        return;
      }
      console.error("Native share failed:", error);
      return copyToClipboard();
    }
  }, [copyToClipboard, isNativeShareSupported]);

  return {
    copyToClipboard,
    isNativeShareSupported,
    shareNative,
    state,
  };
};
