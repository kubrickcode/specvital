"use client";

import { Building2, ChevronRight, FolderGit2, RefreshCw } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";
import { cn } from "@/lib/utils";

type DiscoveryCardProps = {
  count: number;
  isLoading?: boolean;
  isRefreshing?: boolean;
  onClick: () => void;
  onRefresh: () => void;
  type: "organization" | "personal";
};

export const DiscoveryCard = ({
  count,
  isLoading,
  isRefreshing,
  onClick,
  onRefresh,
  type,
}: DiscoveryCardProps) => {
  const t = useTranslations("dashboard.discovery");
  const isPersonal = type === "personal";
  const Icon = isPersonal ? FolderGit2 : Building2;
  const isDisabled = isLoading || (isPersonal && count === 0);

  const handleRefreshClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    onRefresh();
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      onClick();
    }
  };

  const ariaLabel = isPersonal ? t("myRepos") : t("orgRepos");

  return (
    <Card
      aria-label={ariaLabel}
      className={cn(
        "relative p-4 transition-all duration-200",
        !isDisabled && "cursor-pointer hover:shadow-md hover:border-primary/20",
        !isDisabled &&
          "focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 focus-visible:outline-none",
        isDisabled && "opacity-60"
      )}
      onClick={isDisabled ? undefined : onClick}
      onKeyDown={isDisabled ? undefined : handleKeyDown}
      role={isDisabled ? undefined : "button"}
      tabIndex={isDisabled ? undefined : 0}
    >
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-3 min-w-0">
          <div
            className={cn(
              "flex items-center justify-center size-10 rounded-lg shrink-0",
              isPersonal ? "bg-blue-100 text-blue-600" : "bg-purple-100 text-purple-600"
            )}
          >
            <Icon className="size-5" />
          </div>
          <div className="min-w-0">
            <h3 className="font-medium text-sm">{isPersonal ? t("myRepos") : t("orgRepos")}</h3>
            <p className="text-xs text-muted-foreground">
              {isLoading ? t("loading") : count > 0 ? t("reposFound", { count }) : t("noRepos")}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-1 shrink-0">
          <ResponsiveTooltip content={t("refresh")}>
            <Button
              aria-label={t("refresh")}
              className="size-8"
              disabled={isRefreshing}
              onClick={handleRefreshClick}
              size="icon"
              variant="ghost"
            >
              <RefreshCw className={cn("size-4", isRefreshing && "animate-spin")} />
            </Button>
          </ResponsiveTooltip>

          {!isDisabled && <ChevronRight className="size-4 text-muted-foreground" />}
        </div>
      </div>
    </Card>
  );
};
