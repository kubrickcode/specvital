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
import { PaginationStatus, type SortOption } from "@/features/dashboard";

const SORT_OPTIONS: SortOption[] = ["name", "recent", "tests"];

const isSortOption = (value: string): value is SortOption =>
  SORT_OPTIONS.includes(value as SortOption);

type SearchSortControlsProps = {
  hasNextPage: boolean;
  isLoading: boolean;
  onSearchChange: (query: string) => void;
  onSortChange: (sort: SortOption) => void;
  searchQuery: string;
  sortBy: SortOption;
  totalLoaded: number;
};

export const SearchSortControls = ({
  hasNextPage,
  isLoading,
  onSearchChange,
  onSortChange,
  searchQuery,
  sortBy,
  totalLoaded,
}: SearchSortControlsProps) => {
  const t = useTranslations("explore");

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
    <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:flex-1">
        <div className="relative flex-1 sm:max-w-sm">
          <Search
            aria-hidden="true"
            className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
          />
          <Input
            aria-label={t("searchPlaceholder")}
            className="h-11 pl-10 sm:h-9 sm:pl-9"
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder={t("searchPlaceholder")}
            type="search"
            value={searchQuery}
          />
        </div>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button className="h-11 flex-1 sm:h-9 sm:flex-none" variant="outline">
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
      </div>

      <PaginationStatus hasNextPage={hasNextPage} isLoading={isLoading} totalLoaded={totalLoaded} />
    </div>
  );
};
