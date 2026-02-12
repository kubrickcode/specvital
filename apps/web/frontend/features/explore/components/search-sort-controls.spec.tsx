import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi, beforeEach } from "vitest";

import messages from "@/i18n/messages/en.json";

import { SearchSortControls } from "./search-sort-controls";

vi.mock("@/features/dashboard", () => ({
  PaginationStatus: ({
    hasNextPage,
    isLoading,
    totalLoaded,
  }: {
    hasNextPage: boolean;
    isLoading: boolean;
    totalLoaded: number;
  }) => (
    <div data-testid="pagination-status">
      {isLoading ? "Loading..." : `${totalLoaded} loaded${hasNextPage ? " (more)" : ""}`}
    </div>
  ),
}));

const defaultProps = {
  hasNextPage: true,
  isLoading: false,
  onSearchChange: vi.fn(),
  onSortChange: vi.fn(),
  searchQuery: "",
  sortBy: "recent" as const,
  totalLoaded: 10,
};

const renderWithProvider = (ui: React.ReactElement) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      {ui}
    </NextIntlClientProvider>
  );

describe("SearchSortControls", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders search input with placeholder", () => {
    renderWithProvider(<SearchSortControls {...defaultProps} />);

    expect(screen.getByPlaceholderText("Search repositories...")).toBeInTheDocument();
  });

  it("renders sort button showing current sort option as Recent", () => {
    renderWithProvider(<SearchSortControls {...defaultProps} />);

    expect(screen.getByText(/Sort.*Recent/)).toBeInTheDocument();
  });

  it("shows Name in sort button when sortBy is name", () => {
    renderWithProvider(<SearchSortControls {...defaultProps} sortBy="name" />);

    expect(screen.getByText(/Sort.*Name/)).toBeInTheDocument();
  });

  it("shows Tests in sort button when sortBy is tests", () => {
    renderWithProvider(<SearchSortControls {...defaultProps} sortBy="tests" />);

    expect(screen.getByText(/Sort.*Tests/)).toBeInTheDocument();
  });

  it("calls onSearchChange when typing in search input", () => {
    const onSearchChange = vi.fn();
    renderWithProvider(<SearchSortControls {...defaultProps} onSearchChange={onSearchChange} />);

    const searchInput = screen.getByPlaceholderText("Search repositories...");
    fireEvent.change(searchInput, { target: { value: "react" } });

    expect(onSearchChange).toHaveBeenCalledWith("react");
  });

  it("renders pagination status component", () => {
    renderWithProvider(<SearchSortControls {...defaultProps} />);

    expect(screen.getByTestId("pagination-status")).toBeInTheDocument();
    expect(screen.getByText("10 loaded (more)")).toBeInTheDocument();
  });

  it("reflects current search query value", () => {
    renderWithProvider(<SearchSortControls {...defaultProps} searchQuery="vitest" />);

    expect(screen.getByDisplayValue("vitest")).toBeInTheDocument();
  });

  it("shows loading state in pagination when isLoading", () => {
    renderWithProvider(<SearchSortControls {...defaultProps} isLoading />);

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });
});
