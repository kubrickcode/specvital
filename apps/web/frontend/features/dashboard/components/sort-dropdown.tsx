"use client";

import { ArrowUpDown } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import type { SortOption } from "../types";

const SORT_OPTIONS: SortOption[] = ["name", "recent", "tests"];

const isSortOption = (value: string): value is SortOption =>
  SORT_OPTIONS.includes(value as SortOption);

type SortDropdownProps = {
  isMobile?: boolean;
  onSortChange: (sort: SortOption) => void;
  sortBy: SortOption;
};

export const SortDropdown = ({ isMobile = false, onSortChange, sortBy }: SortDropdownProps) => {
  const t = useTranslations("dashboard");

  const sortLabels: Record<SortOption, string> = {
    name: t("sort.name"),
    recent: t("sort.recent"),
    tests: t("sort.tests"),
  };

  const handleSortChange = (value: string) => {
    if (isSortOption(value)) {
      onSortChange(value);
    }
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          aria-label={t("sort.label")}
          className={isMobile ? "h-11 flex-1" : "h-9"}
          variant="outline"
        >
          <ArrowUpDown aria-hidden="true" />
          <span>
            {t("sort.label")}: {sortLabels[sortBy]}
          </span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-40">
        <DropdownMenuRadioGroup onValueChange={handleSortChange} value={sortBy}>
          <DropdownMenuRadioItem value="recent">{sortLabels.recent}</DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="name">{sortLabels.name}</DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="tests">{sortLabels.tests}</DropdownMenuRadioItem>
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
