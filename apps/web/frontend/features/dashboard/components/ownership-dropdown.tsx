"use client";

import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import {
  OWNERSHIP_FILTER_ICONS,
  OWNERSHIP_FILTER_OPTIONS,
  type OwnershipFilter,
  useOwnershipFilter,
} from "../hooks/use-ownership-filter";

const isOwnershipFilter = (value: string): value is OwnershipFilter =>
  OWNERSHIP_FILTER_OPTIONS.includes(value as OwnershipFilter);

export const OwnershipDropdown = () => {
  const t = useTranslations("dashboard.filter.ownership");
  const { ownershipFilter, setOwnershipFilter } = useOwnershipFilter();

  const ownershipLabels: Record<OwnershipFilter, string> = {
    all: t("all"),
    mine: t("mine"),
    organization: t("organization"),
    others: t("others"),
  };

  const handleOwnershipChange = (value: string) => {
    if (isOwnershipFilter(value)) {
      setOwnershipFilter(value);
    }
  };

  const CurrentIcon = OWNERSHIP_FILTER_ICONS[ownershipFilter];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          aria-label={t("label")}
          className="h-11 flex-1 sm:h-9 sm:flex-none"
          variant="outline"
        >
          <CurrentIcon aria-hidden="true" className="size-4" />
          <span>{ownershipLabels[ownershipFilter]}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" className="w-44">
        <DropdownMenuRadioGroup onValueChange={handleOwnershipChange} value={ownershipFilter}>
          {OWNERSHIP_FILTER_OPTIONS.map((option) => {
            const Icon = OWNERSHIP_FILTER_ICONS[option];
            return (
              <DropdownMenuRadioItem key={option} value={option}>
                <Icon aria-hidden="true" className="mr-2 size-4" />
                {ownershipLabels[option]}
              </DropdownMenuRadioItem>
            );
          })}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
