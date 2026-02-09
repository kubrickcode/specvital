"use client";

import { AlertCircle, Building2, ChevronRight, Link2, RefreshCw } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import type { GitHubOrganization } from "@/lib/api/types";
import { cn } from "@/lib/utils";

import { OrgInstallDialog } from "./org-install-dialog";

type OrganizationPickerProps = {
  isLoading?: boolean;
  isOpen: boolean;
  isRefreshing?: boolean;
  onOpenChange: (open: boolean) => void;
  onRefreshOrg: (org: string) => void;
  onSelectOrg: (org: GitHubOrganization) => void;
  organizations: GitHubOrganization[];
  unanalyzedCounts: Record<string, number>;
};

export const OrganizationPicker = ({
  isLoading,
  isOpen,
  isRefreshing,
  onOpenChange,
  onRefreshOrg,
  onSelectOrg,
  organizations,
  unanalyzedCounts,
}: OrganizationPickerProps) => {
  const t = useTranslations("dashboard.discovery");
  const tApp = useTranslations("dashboard.githubApp");
  const [installDialogOrg, setInstallDialogOrg] = useState<GitHubOrganization | null>(null);

  const handleConnectClick = (e: React.MouseEvent, org: GitHubOrganization) => {
    e.stopPropagation();
    setInstallDialogOrg(org);
  };

  const getAccessStatusBadge = (org: GitHubOrganization) => {
    if (org.accessStatus === "accessible") {
      return null;
    }

    if (org.accessStatus === "pending") {
      return (
        <Badge className="bg-yellow-100 text-yellow-800 hover:bg-yellow-100" variant="secondary">
          <AlertCircle className="mr-1 size-3" />
          {tApp("statusPending")}
        </Badge>
      );
    }

    return (
      <Badge className="bg-orange-100 text-orange-800 hover:bg-orange-100" variant="secondary">
        <Link2 className="mr-1 size-3" />
        {tApp("statusRestricted")}
      </Badge>
    );
  };

  return (
    <>
      <Sheet onOpenChange={onOpenChange} open={isOpen}>
        <SheetContent className="w-full sm:max-w-lg overflow-y-auto" side="right">
          <SheetHeader className="mb-4">
            <SheetTitle>{t("selectOrganization")}</SheetTitle>
            <SheetDescription>
              {t("selectOrganizationDescription", { count: organizations.length })}
            </SheetDescription>
          </SheetHeader>

          <ul className="space-y-2">
            {isLoading ? (
              <li className="text-center py-8 text-sm text-muted-foreground">{t("loading")}</li>
            ) : organizations.length === 0 ? (
              <li className="text-center py-8 text-sm text-muted-foreground">
                {t("noOrganizations")}
              </li>
            ) : (
              organizations.map((org) => {
                const count = unanalyzedCounts[org.login] ?? 0;
                const isRestricted = org.accessStatus !== "accessible";
                const isDisabled = count === 0 && !isRestricted;

                const handleKeyDown = (e: React.KeyboardEvent) => {
                  if (e.key === "Enter" || e.key === " ") {
                    e.preventDefault();
                    if (isRestricted) {
                      setInstallDialogOrg(org);
                    } else {
                      onSelectOrg(org);
                    }
                  }
                };

                const handleCardClick = () => {
                  if (isRestricted) {
                    setInstallDialogOrg(org);
                  } else {
                    onSelectOrg(org);
                  }
                };

                return (
                  <li key={org.id}>
                    <Card
                      aria-label={org.login}
                      className={cn(
                        "relative p-4 transition-all duration-200",
                        !isDisabled && "cursor-pointer hover:shadow-md hover:border-primary/20",
                        !isDisabled &&
                          "focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 focus-visible:outline-none",
                        isDisabled && "opacity-60",
                        isRestricted && "border-orange-200"
                      )}
                      onClick={isDisabled ? undefined : handleCardClick}
                      onKeyDown={isDisabled ? undefined : handleKeyDown}
                      role={isDisabled ? undefined : "button"}
                      tabIndex={isDisabled ? undefined : 0}
                    >
                      <div className="flex items-center justify-between gap-3">
                        <div className="flex items-center gap-3 min-w-0">
                          {org.avatarUrl ? (
                            <img
                              alt={org.login}
                              className="size-10 rounded-lg shrink-0 object-cover"
                              src={org.avatarUrl}
                            />
                          ) : (
                            <div className="flex items-center justify-center size-10 rounded-lg shrink-0 bg-purple-100 text-purple-600">
                              <Building2 className="size-5" />
                            </div>
                          )}
                          <div className="min-w-0">
                            <div className="flex items-center gap-2">
                              <h3 className="font-medium text-sm">{org.login}</h3>
                              {getAccessStatusBadge(org)}
                            </div>
                            <p className="text-xs text-muted-foreground">
                              {isRestricted
                                ? tApp("connectToAccess")
                                : count > 0
                                  ? t("reposFound", { count })
                                  : t("noRepos")}
                            </p>
                          </div>
                        </div>

                        <div className="flex items-center gap-1 shrink-0">
                          {isRestricted ? (
                            <Button
                              onClick={(e) => handleConnectClick(e, org)}
                              size="sm"
                              variant="outline"
                            >
                              {tApp("connect")}
                            </Button>
                          ) : (
                            <>
                              <ResponsiveTooltip content={t("refresh")}>
                                <Button
                                  aria-label={t("refresh")}
                                  className="size-8"
                                  disabled={isRefreshing}
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    onRefreshOrg(org.login);
                                  }}
                                  size="icon"
                                  variant="ghost"
                                >
                                  <RefreshCw
                                    className={cn("size-4", isRefreshing && "animate-spin")}
                                  />
                                </Button>
                              </ResponsiveTooltip>

                              {!isDisabled && (
                                <ChevronRight className="size-4 text-muted-foreground" />
                              )}
                            </>
                          )}
                        </div>
                      </div>
                    </Card>
                  </li>
                );
              })
            )}
          </ul>
        </SheetContent>
      </Sheet>

      <OrgInstallDialog
        isOpen={!!installDialogOrg}
        onOpenChange={(open) => !open && setInstallDialogOrg(null)}
        organization={installDialogOrg}
      />
    </>
  );
};
