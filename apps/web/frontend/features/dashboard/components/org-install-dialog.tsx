"use client";

import { Building2, ExternalLink, Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import type { GitHubOrganization } from "@/lib/api/types";

import { useGitHubAppInstallUrl } from "../hooks";

type OrgInstallDialogProps = {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  organization: GitHubOrganization | null;
};

export const OrgInstallDialog = ({ isOpen, onOpenChange, organization }: OrgInstallDialogProps) => {
  const t = useTranslations("dashboard.githubApp");
  const { install, isLoading } = useGitHubAppInstallUrl();

  if (!organization) {
    return null;
  }

  return (
    <Dialog onOpenChange={onOpenChange} open={isOpen}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            {organization.avatarUrl ? (
              <img
                alt={organization.login}
                className="size-6 rounded"
                src={organization.avatarUrl}
              />
            ) : (
              <Building2 className="size-6" />
            )}
            {t("installDialogTitle", { org: organization.login })}
          </DialogTitle>
          <DialogDescription>{t("installDialogDescription")}</DialogDescription>
        </DialogHeader>

        <div className="space-y-3 text-sm text-muted-foreground">
          <p>{t("installDialogBenefits")}</p>
          <ul className="list-disc list-inside space-y-1 ml-2">
            <li>{t("installDialogBenefit1")}</li>
            <li>{t("installDialogBenefit2")}</li>
            <li>{t("installDialogBenefit3")}</li>
          </ul>
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          <Button onClick={() => onOpenChange(false)} variant="outline">
            {t("cancel")}
          </Button>
          <Button disabled={isLoading} onClick={install}>
            {isLoading ? (
              <Loader2 className="mr-2 size-4 animate-spin" />
            ) : (
              <ExternalLink className="mr-2 size-4" />
            )}
            {t("installButton")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
