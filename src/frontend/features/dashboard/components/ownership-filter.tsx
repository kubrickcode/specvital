"use client";

import { Filter } from "lucide-react";
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
  OWNERSHIP_FILTER_VALUES,
  type OwnershipFilter as OwnershipFilterType,
} from "@/lib/api/types";

const isOwnershipOption = (value: string): value is OwnershipFilterType =>
  (OWNERSHIP_FILTER_VALUES as readonly string[]).includes(value);

type OwnershipFilterProps = {
  onChange: (value: OwnershipFilterType) => void;
  value: OwnershipFilterType;
};

export const OwnershipFilter = ({ onChange, value }: OwnershipFilterProps) => {
  const t = useTranslations("dashboard.ownership");

  const ownershipLabels: Record<OwnershipFilterType, string> = {
    all: t("all"),
    mine: t("mine"),
    organization: t("organization"),
  };

  const handleValueChange = (newValue: string) => {
    if (isOwnershipOption(newValue)) {
      onChange(newValue);
    }
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button className="sm:w-auto w-full" variant="outline">
          <Filter aria-hidden="true" className="size-4" />
          <span>
            {t("label")}: {ownershipLabels[value]}
          </span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-44">
        <DropdownMenuRadioGroup onValueChange={handleValueChange} value={value}>
          <DropdownMenuRadioItem value="all">{ownershipLabels.all}</DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="mine">{ownershipLabels.mine}</DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="organization">
            {ownershipLabels.organization}
          </DropdownMenuRadioItem>
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
