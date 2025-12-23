"use client";

import { ArrowUpDown, Search } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";

import type { SortOption } from "../types";

const SORT_OPTIONS: SortOption[] = ["name", "recent", "tests"];

const isSortOption = (value: string): value is SortOption =>
  SORT_OPTIONS.includes(value as SortOption);

type DashboardHeaderProps = {
  onSearchChange: (query: string) => void;
  onSortChange: (sort: SortOption) => void;
  searchQuery: string;
  sortBy: SortOption;
};

export const DashboardHeader = ({
  onSearchChange,
  onSortChange,
  searchQuery,
  sortBy,
}: DashboardHeaderProps) => {
  const t = useTranslations("dashboard");

  const sortLabels: Record<SortOption, string> = {
    name: t("sort.name", { defaultValue: "Name" }),
    recent: t("sort.recent", { defaultValue: "Recent" }),
    tests: t("sort.tests", { defaultValue: "Tests" }),
  };

  const searchPlaceholder = t("searchPlaceholder", { defaultValue: "Search repositories..." });

  const handleSortChange = (value: string) => {
    if (isSortOption(value)) {
      onSortChange(value);
    }
  };

  return (
    <div className="flex flex-col sm:flex-row gap-3 mb-6">
      <div className="relative flex-1">
        <Search
          aria-hidden="true"
          className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground"
        />
        <Input
          aria-label={searchPlaceholder}
          className="pl-9"
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder={searchPlaceholder}
          type="search"
          value={searchQuery}
        />
      </div>

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button className="sm:w-auto w-full" variant="outline">
            <ArrowUpDown aria-hidden="true" />
            <span>
              {t("sort.label", { defaultValue: "Sort" })}: {sortLabels[sortBy]}
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
    </div>
  );
};
