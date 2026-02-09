"use client";

import { AlertCircle, ExternalLink, Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import type { OrganizationAccessStatus } from "@/lib/api/types";

import { useGitHubAppInstallUrl } from "../hooks";

type OrgConnectionBannerProps = {
  accessStatus: OrganizationAccessStatus;
  orgLogin: string;
};

export const OrgConnectionBanner = ({ accessStatus, orgLogin }: OrgConnectionBannerProps) => {
  const t = useTranslations("dashboard.githubApp");
  const { install, isLoading } = useGitHubAppInstallUrl();

  if (accessStatus === "accessible") {
    return null;
  }

  const isPending = accessStatus === "pending";

  return (
    <Alert className="mb-4">
      <AlertCircle className="size-4" />
      <AlertTitle>{isPending ? t("pendingTitle") : t("restrictedTitle")}</AlertTitle>
      <AlertDescription className="mt-2">
        <p className="mb-3">{isPending ? t("pendingDescription") : t("restrictedDescription")}</p>
        <Button disabled={isLoading} onClick={install} size="sm" variant="outline">
          {isLoading ? (
            <Loader2 className="mr-2 size-4 animate-spin" />
          ) : (
            <ExternalLink className="mr-2 size-4" />
          )}
          {t("connectButton", { org: orgLogin })}
        </Button>
      </AlertDescription>
    </Alert>
  );
};
