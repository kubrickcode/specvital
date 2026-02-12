import { fireEvent, render, screen, within } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi, beforeEach } from "vitest";

import messages from "@/i18n/messages/en.json";

import { GlobalSearchDialog } from "./global-search-dialog";

const mockClose = vi.fn();
const mockSetQuery = vi.fn();
const mockNavigateToRepository = vi.fn();
const mockLogin = vi.fn();
const mockAddItem = vi.fn();

let mockIsOpen = true;
let mockHasResults = false;
let mockGroupedResults = { bookmarks: [], community: [], repositories: [] };

vi.mock("../hooks", () => ({
  useGlobalSearchStore: () => ({
    close: mockClose,
    isOpen: mockIsOpen,
  }),
  useDebouncedSearch: () => ({
    groupedResults: mockGroupedResults,
    hasError: false,
    hasResults: mockHasResults,
    isAuthenticated: true,
    isLoading: false,
    navigateToRepository: mockNavigateToRepository,
    setQuery: mockSetQuery,
  }),
  useRecentItems: () => ({
    addItem: mockAddItem,
    recentItems: [],
  }),
  useStaticActions: () => ({
    analyzeDialogOpen: false,
    commandItems: [
      {
        id: "action-toggle-theme-dark",
        label: "Switch to Dark Mode",
        onSelect: vi.fn(),
      },
    ],
    navigationItems: [
      {
        id: "nav-explore",
        label: "Go to Explore",
        onSelect: vi.fn(),
      },
      {
        id: "nav-pricing",
        label: "Go to Pricing",
        onSelect: vi.fn(),
      },
      {
        id: "nav-home",
        label: "Go to Home",
        onSelect: vi.fn(),
      },
    ],
    setAnalyzeDialogOpen: vi.fn(),
  }),
}));

vi.mock("@/features/auth", () => ({
  useAuth: () => ({
    isAuthenticated: true,
    login: mockLogin,
  }),
}));

vi.mock("@/features/home", () => ({
  AnalyzeDialog: () => null,
}));

const renderWithProvider = (ui: React.ReactElement) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      {ui}
    </NextIntlClientProvider>
  );

describe("GlobalSearchDialog", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockIsOpen = true;
    mockHasResults = false;
    mockGroupedResults = { bookmarks: [], community: [], repositories: [] };
  });

  it("renders dialog with search input when open", () => {
    renderWithProvider(<GlobalSearchDialog />);

    expect(screen.getByPlaceholderText("Search...")).toBeInTheDocument();
  });

  it("displays navigation section with Explore, Pricing, and Home items", () => {
    renderWithProvider(<GlobalSearchDialog />);

    expect(screen.getByText("Navigation")).toBeInTheDocument();
    expect(screen.getByText("Go to Explore")).toBeInTheDocument();
    expect(screen.getByText("Go to Pricing")).toBeInTheDocument();
    expect(screen.getByText("Go to Home")).toBeInTheDocument();
  });

  it("displays commands section with theme toggle", () => {
    renderWithProvider(<GlobalSearchDialog />);

    expect(screen.getByText("Commands")).toBeInTheDocument();
    expect(screen.getByText("Switch to Dark Mode")).toBeInTheDocument();
  });

  it("renders keyboard hints footer", () => {
    renderWithProvider(<GlobalSearchDialog />);

    expect(screen.getByText("Navigate")).toBeInTheDocument();
    expect(screen.getByText("Select")).toBeInTheDocument();
    expect(screen.getByText("Close")).toBeInTheDocument();
  });

  it("renders mobile close button with correct aria-label", () => {
    renderWithProvider(<GlobalSearchDialog />);

    expect(screen.getByLabelText("Close")).toBeInTheDocument();
  });

  describe("Escape closes dialog", () => {
    it("calls close when Escape key is pressed", () => {
      renderWithProvider(<GlobalSearchDialog />);

      fireEvent.keyDown(document, { key: "Escape" });

      expect(mockClose).toHaveBeenCalled();
    });
  });

  describe("No results display", () => {
    it("shows no-results message when query yields empty results", () => {
      renderWithProvider(<GlobalSearchDialog />);

      const input = screen.getByPlaceholderText("Search...");
      fireEvent.change(input, { target: { value: "nonexistent" } });

      expect(screen.getByText(/No results for/)).toBeInTheDocument();
    });
  });

  describe("Search reset restores static actions", () => {
    it("restores navigation and command items after clearing search input", () => {
      renderWithProvider(<GlobalSearchDialog />);

      const input = screen.getByPlaceholderText("Search...");

      fireEvent.change(input, { target: { value: "test" } });
      expect(screen.queryByText("Navigation")).not.toBeInTheDocument();

      fireEvent.change(input, { target: { value: "" } });
      expect(screen.getByText("Navigation")).toBeInTheDocument();
      expect(screen.getByText("Commands")).toBeInTheDocument();
    });
  });

  describe("Keyboard arrow navigation", () => {
    it("changes data-selected attribute when arrow keys are pressed", () => {
      renderWithProvider(<GlobalSearchDialog />);

      const commandItems = screen.getAllByRole("option");
      expect(commandItems.length).toBeGreaterThan(1);

      const firstItem = commandItems[0];
      expect(firstItem).toHaveAttribute("data-selected", "true");

      const listEl = screen.getByRole("listbox");
      fireEvent.keyDown(listEl, { key: "ArrowDown" });

      const updatedItems = screen.getAllByRole("option");
      const selectedAfterDown = updatedItems.find(
        (item) => item.getAttribute("data-selected") === "true"
      );
      expect(selectedAfterDown).toBeDefined();
      expect(selectedAfterDown).not.toBe(firstItem);
    });
  });

  describe("Fullscreen dialog on mobile", () => {
    it("renders dialog content with fullscreen mobile classes", () => {
      renderWithProvider(<GlobalSearchDialog />);

      const dialogContent = document.querySelector("[data-slot='dialog-content']");
      expect(dialogContent).toBeInTheDocument();
      expect(dialogContent).toHaveClass("h-[100dvh]");
      expect(dialogContent).toHaveClass("w-screen");
    });
  });
});
