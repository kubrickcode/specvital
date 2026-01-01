"use client";

import { List, Star, User } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import type { ViewFilter } from "../hooks/use-view-filter";
import { useViewFilter } from "../hooks/use-view-filter";

const VIEW_FILTER_OPTIONS: ViewFilter[] = ["all", "mine", "starred"];

const VIEW_FILTER_ICONS: Record<ViewFilter, React.ComponentType<{ className?: string }>> = {
  all: List,
  mine: User,
  starred: Star,
};

const isViewFilter = (value: string): value is ViewFilter =>
  VIEW_FILTER_OPTIONS.includes(value as ViewFilter);

export const ViewFilterDropdown = () => {
  const t = useTranslations("dashboard.view");
  const { setViewFilter, viewFilter } = useViewFilter();

  const viewLabels: Record<ViewFilter, string> = {
    all: t("all"),
    mine: t("mine"),
    starred: t("starred"),
  };

  const handleViewChange = (value: string) => {
    if (isViewFilter(value)) {
      setViewFilter(value);
    }
  };

  const CurrentIcon = VIEW_FILTER_ICONS[viewFilter];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          aria-label={t("label")}
          className="h-11 flex-1 sm:h-9 sm:flex-none"
          variant="outline"
        >
          <CurrentIcon aria-hidden="true" className="size-4" />
          <span>
            {t("label")}: {viewLabels[viewFilter]}
          </span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-44">
        <DropdownMenuRadioGroup onValueChange={handleViewChange} value={viewFilter}>
          {VIEW_FILTER_OPTIONS.map((option) => {
            const Icon = VIEW_FILTER_ICONS[option];
            return (
              <DropdownMenuRadioItem key={option} value={option}>
                <Icon aria-hidden="true" className="mr-2 size-4" />
                {viewLabels[option]}
              </DropdownMenuRadioItem>
            );
          })}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
