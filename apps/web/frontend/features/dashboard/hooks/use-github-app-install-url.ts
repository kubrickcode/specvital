"use client";

import { useMutation } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useState } from "react";
import { toast } from "sonner";

import { fetchGitHubAppInstallUrl } from "../api";

type UseGitHubAppInstallUrlReturn = {
  install: () => Promise<void>;
  isLoading: boolean;
};

export const useGitHubAppInstallUrl = (): UseGitHubAppInstallUrlReturn => {
  const t = useTranslations("dashboard.githubApp");
  const [isRedirecting, setIsRedirecting] = useState(false);

  const mutation = useMutation({
    mutationFn: fetchGitHubAppInstallUrl,
  });

  const install = async (): Promise<void> => {
    setIsRedirecting(true);
    try {
      const result = await mutation.mutateAsync();
      window.location.href = result.installUrl;
    } catch (error) {
      setIsRedirecting(false);
      toast.error(t("installError"), {
        description: error instanceof Error ? error.message : undefined,
      });
    }
  };

  return {
    install,
    isLoading: mutation.isPending || isRedirecting,
  };
};
