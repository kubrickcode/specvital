import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import type { TestStatus } from "@/lib/api/types";

import { FilterEmptyState } from "./filter-empty-state";
import type { FilterInfo } from "./filter-empty-state";

const messages = {
  analyze: {
    filter: {
      noResults: "No matching tests found",
      noResultsDescription: "Try adjusting your filters",
      resetFilters: "Reset Filters",
      searchLabel: "Search",
      statusActive: "Active",
      statusFocused: "Focused",
      statusSkipped: "Skipped",
      statusTodo: "Todo",
      statusXfail: "Expected Fail",
    },
  },
};

const defaultFilterInfo: FilterInfo = {
  frameworks: [],
  query: "",
  statuses: [],
};

const renderFilterEmptyState = (filterInfo: FilterInfo = defaultFilterInfo, onReset = vi.fn()) => {
  return {
    onReset,
    ...render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <FilterEmptyState filterInfo={filterInfo} onReset={onReset} />
      </NextIntlClientProvider>
    ),
  };
};

describe("FilterEmptyState", () => {
  it("displays the empty search results message", () => {
    renderFilterEmptyState({ frameworks: [], query: "nonexistent", statuses: [] });

    expect(screen.getByText("No matching tests found")).toBeInTheDocument();
    expect(screen.getByText("Try adjusting your filters")).toBeInTheDocument();
  });

  it("displays search query as a filter badge", () => {
    renderFilterEmptyState({ frameworks: [], query: "my search term", statuses: [] });

    const badge = screen.getByText(/my search term/);
    expect(badge).toBeInTheDocument();
    expect(badge.textContent).toContain("Search");
    expect(badge.textContent).toContain("my search term");
  });

  it("displays framework names as filter badges", () => {
    renderFilterEmptyState({
      frameworks: ["jest", "vitest"],
      query: "",
      statuses: [],
    });

    expect(screen.getByText("jest")).toBeInTheDocument();
    expect(screen.getByText("vitest")).toBeInTheDocument();
  });

  it("displays status names as filter badges", () => {
    const statuses: TestStatus[] = ["active", "skipped"];
    renderFilterEmptyState({ frameworks: [], query: "", statuses });

    expect(screen.getByText("Active")).toBeInTheDocument();
    expect(screen.getByText("Skipped")).toBeInTheDocument();
  });

  it("displays all filter types as badges simultaneously", () => {
    const statuses: TestStatus[] = ["focused"];
    renderFilterEmptyState({
      frameworks: ["pytest"],
      query: "login",
      statuses,
    });

    expect(screen.getByText(/login/)).toBeInTheDocument();
    expect(screen.getByText("pytest")).toBeInTheDocument();
    expect(screen.getByText("Focused")).toBeInTheDocument();
  });

  it("calls onReset when 'Reset Filters' button is clicked", () => {
    const onReset = vi.fn();
    renderFilterEmptyState(
      {
        frameworks: ["jest"],
        query: "test",
        statuses: [],
      },
      onReset
    );

    fireEvent.click(screen.getByRole("button", { name: /Reset Filters/i }));

    expect(onReset).toHaveBeenCalledOnce();
  });

  it("does not show filter badges or reset button when no active filters exist", () => {
    renderFilterEmptyState({ frameworks: [], query: "", statuses: [] });

    expect(screen.queryByRole("button", { name: /Reset Filters/i })).not.toBeInTheDocument();
    expect(screen.queryByText("Search")).not.toBeInTheDocument();
  });

  it("trims whitespace-only query and treats it as no search filter", () => {
    renderFilterEmptyState({ frameworks: [], query: "   ", statuses: [] });

    expect(screen.queryByRole("button", { name: /Reset Filters/i })).not.toBeInTheDocument();
  });
});
