"use client";

import { Building2, ChevronRight, FolderGit2, RefreshCw } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
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
  const isDisabled = isLoading || count === 0;

  const handleRefreshClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    onRefresh();
  };

  return (
    <Card
      className={cn(
        "relative p-4 transition-all duration-200",
        !isDisabled && "cursor-pointer hover:shadow-md hover:border-primary/20",
        isDisabled && "opacity-60"
      )}
      onClick={isDisabled ? undefined : onClick}
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
          <Tooltip>
            <TooltipTrigger asChild>
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
            </TooltipTrigger>
            <TooltipContent>{t("refresh")}</TooltipContent>
          </Tooltip>

          {!isDisabled && <ChevronRight className="size-4 text-muted-foreground" />}
        </div>
      </div>
    </Card>
  );
};
