import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import type { SortOption } from "../types";
import { FilterBar } from "./filter-bar";

vi.mock("./bookmark-toggle", () => ({
  BookmarkToggle: () => <div data-testid="bookmark-toggle" />,
}));

vi.mock("./mobile-filter-drawer", () => ({
  MobileFilterDrawer: () => <div data-testid="mobile-filter-drawer" />,
}));

vi.mock("./ownership-dropdown", () => ({
  OwnershipDropdown: () => <div data-testid="ownership-dropdown" />,
}));

vi.mock("./sort-dropdown", () => ({
  SortDropdown: () => <div data-testid="sort-dropdown" />,
}));

vi.mock("./pagination-status", () => ({
  PaginationStatus: () => <div data-testid="pagination-status" />,
}));

Object.defineProperty(window, "matchMedia", {
  value: vi.fn().mockImplementation((query: string) => ({
    addEventListener: vi.fn(),
    addListener: vi.fn(),
    dispatchEvent: vi.fn(),
    matches: false,
    media: query,
    onchange: null,
    removeEventListener: vi.fn(),
    removeListener: vi.fn(),
  })),
  writable: true,
});

const messages = {
  dashboard: {
    searchPlaceholder: "Search repositories...",
  },
};

type RenderOptions = {
  hasNextPage?: boolean;
  isLoading?: boolean;
  onSearchChange?: (query: string) => void;
  onSortChange?: (sort: SortOption) => void;
  searchQuery?: string;
  sortBy?: SortOption;
  totalLoaded?: number;
};

const renderFilterBar = (options: RenderOptions = {}) => {
  const props = {
    hasNextPage: options.hasNextPage ?? false,
    isLoading: options.isLoading ?? false,
    onSearchChange: options.onSearchChange ?? vi.fn(),
    onSortChange: options.onSortChange ?? vi.fn(),
    searchQuery: options.searchQuery ?? "",
    sortBy: options.sortBy ?? ("recent" as SortOption),
    totalLoaded: options.totalLoaded ?? 0,
  };

  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <FilterBar {...props} />
    </NextIntlClientProvider>
  );
};

describe("FilterBar", () => {
  it("renders the search input with placeholder", () => {
    renderFilterBar();

    expect(screen.getByRole("searchbox", { name: /search repositories/i })).toBeInTheDocument();
  });

  it("displays the current search query value", () => {
    renderFilterBar({ searchQuery: "react" });

    expect(screen.getByRole("searchbox", { name: /search repositories/i })).toHaveValue("react");
  });

  it("calls onSearchChange when input value changes", () => {
    const onSearchChange = vi.fn();
    renderFilterBar({ onSearchChange });

    fireEvent.change(screen.getByRole("searchbox", { name: /search repositories/i }), {
      target: { value: "vue" },
    });

    expect(onSearchChange).toHaveBeenCalledWith("vue");
  });

  it("calls onSearchChange with empty string when input is cleared", () => {
    const onSearchChange = vi.fn();
    renderFilterBar({ onSearchChange, searchQuery: "react" });

    fireEvent.change(screen.getByRole("searchbox", { name: /search repositories/i }), {
      target: { value: "" },
    });

    expect(onSearchChange).toHaveBeenCalledWith("");
  });

  it("renders child filter components", () => {
    renderFilterBar();

    expect(screen.getByTestId("pagination-status")).toBeInTheDocument();
  });
});
