"use client";

import { Search } from "lucide-react";
import { useTranslations } from "next-intl";

import { Input } from "@/components/ui/input";

import type { SortOption } from "../types";
import { BookmarkToggle } from "./bookmark-toggle";
import { MobileFilterDrawer } from "./mobile-filter-drawer";
import { OwnershipDropdown } from "./ownership-dropdown";
import { PaginationStatus } from "./pagination-status";
import { SortDropdown } from "./sort-dropdown";

type FilterBarProps = {
  hasNextPage: boolean;
  isLoading: boolean;
  onSearchChange: (query: string) => void;
  onSortChange: (sort: SortOption) => void;
  searchQuery: string;
  sortBy: SortOption;
  totalLoaded: number;
};

export const FilterBar = ({
  hasNextPage,
  isLoading,
  onSearchChange,
  onSortChange,
  searchQuery,
  sortBy,
  totalLoaded,
}: FilterBarProps) => {
  const t = useTranslations("dashboard");

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

        <div className="flex w-full gap-2 sm:hidden">
          <MobileFilterDrawer />
          <SortDropdown isMobile onSortChange={onSortChange} sortBy={sortBy} />
        </div>

        <div className="hidden gap-2 sm:flex">
          <OwnershipDropdown />
          <BookmarkToggle />
          <SortDropdown onSortChange={onSortChange} sortBy={sortBy} />
        </div>
      </div>

      <PaginationStatus hasNextPage={hasNextPage} isLoading={isLoading} totalLoaded={totalLoaded} />
    </div>
  );
};
