import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { RepositoryCard } from "@/lib/api/types";

vi.mock("@/i18n/navigation", () => ({
  Link: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}));

vi.mock("@/features/auth/hooks/use-auth", () => ({
  useAuth: () => ({
    isAuthenticated: true,
    isLoading: false,
  }),
}));

vi.mock("@/features/auth/hooks/use-login-modal", () => ({
  useLoginModal: () => ({
    open: vi.fn(),
  }),
}));

const mockReanalyze = vi.fn();
const mockAddBookmark = vi.fn();
const mockRemoveBookmark = vi.fn();
const mockFetchNextPage = vi.fn();
const mockRefetch = vi.fn();
const mockRepositories = vi.fn<() => RepositoryCard[]>();
const mockIsLoading = vi.fn<() => boolean>();
const mockHasNextPage = vi.fn<() => boolean>();
const mockIsError = vi.fn<() => boolean>();
const mockBookmarkOnly = vi.fn<() => boolean>();
const mockSetBookmarkOnly = vi.fn();
const mockOwnershipFilter = vi.fn();
const mockStatsData = vi.fn();
const mockStatsIsLoading = vi.fn<() => boolean>();

vi.mock("../hooks", () => ({
  useAddBookmark: () => ({ addBookmark: mockAddBookmark }),
  usePaginatedRepositories: () => ({
    data: mockRepositories(),
    fetchNextPage: mockFetchNextPage,
    hasNextPage: mockHasNextPage(),
    isError: mockIsError(),
    isFetchingNextPage: false,
    isLoading: mockIsLoading(),
    refetch: mockRefetch,
  }),
  useReanalyze: () => ({ isPending: false, reanalyze: mockReanalyze }),
  useRemoveBookmark: () => ({ removeBookmark: mockRemoveBookmark }),
  useRepositoryStats: () => ({
    data: mockStatsData(),
    isLoading: mockStatsIsLoading(),
  }),
}));

vi.mock("../hooks/use-bookmark-filter", () => ({
  useBookmarkFilter: () => ({
    bookmarkOnly: mockBookmarkOnly(),
    setBookmarkOnly: mockSetBookmarkOnly,
  }),
}));

vi.mock("../hooks/use-ownership-filter", () => ({
  OWNERSHIP_FILTER_ICONS: {
    all: () => null,
    mine: () => null,
    organization: () => null,
    others: () => null,
  },
  OWNERSHIP_FILTER_OPTIONS: ["all", "mine", "organization", "others"],
  useOwnershipFilter: () => ({
    ownershipFilter: mockOwnershipFilter(),
    setOwnershipFilter: vi.fn(),
  }),
}));

vi.mock("./active-tasks-section", () => ({
  ActiveTasksSection: () => <div data-testid="active-tasks-section" />,
}));

vi.mock("@/features/home", () => ({
  AnalyzeDialog: ({ variant }: { variant?: string }) => (
    <button data-testid="analyze-dialog" data-variant={variant} type="button">
      Analyze a repository
    </button>
  ),
}));

vi.mock("./summary-section", () => ({
  SummarySection: () => <div data-testid="summary-section">Summary Section</div>,
}));

type FilterBarMockProps = {
  onSearchChange: (query: string) => void;
  searchQuery: string;
};

vi.mock("./filter-bar", () => ({
  FilterBar: ({ onSearchChange, searchQuery }: FilterBarMockProps) => (
    <div data-testid="filter-bar">
      <input
        aria-label="Search repositories..."
        onChange={(e) => onSearchChange(e.target.value)}
        type="search"
        value={searchQuery}
      />
    </div>
  ),
}));

vi.mock("./infinite-scroll-loader", () => ({
  InfiniteScrollLoader: () => <div data-testid="infinite-scroll-loader" />,
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

import { DashboardContent } from "./dashboard-content";

const messages = {
  dashboard: {
    card: {
      addBookmark: "Add bookmark",
      analyzeLatest: "Analyze latest commits",
      bookmarked: "Bookmarked",
      loginToBookmark: "Sign in to bookmark",
      noAnalysis: "No analysis yet",
      reanalyze: "Reanalyze repository",
      removeBookmark: "Remove bookmark",
      tests: "tests",
      update: "Update",
    },
    delta: {
      decreased: "Tests decreased by {count}",
      increased: "Tests increased by {count}",
      unchanged: "Tests unchanged",
    },
    emptyStateVariant: {
      error: {
        description: "Something went wrong. Please try again.",
        title: "Error",
      },
      noBookmarks: {
        description: "Bookmark repositories to see them here",
        title: "No bookmarked repositories",
      },
      noRepos: {
        description: "Analyze a repository to get started.",
        title: "No repositories yet",
      },
      noSearchResults: {
        description: "Try adjusting your search query or filters",
        descriptionWithQuery: 'No results found for "{query}".',
        title: "No matching repositories",
      },
    },
    exploreCta: "Browse other repositories",
    status: {
      analyzing: "Analyzing...",
      newCommits: "New commits",
      unknown: "Status unknown",
      upToDate: "Up to date",
    },
  },
};

const createMockRepo = (overrides?: Partial<RepositoryCard>): RepositoryCard => ({
  fullName: "facebook/react",
  id: "1",
  isAnalyzedByMe: false,
  isBookmarked: false,
  latestAnalysis: {
    analyzedAt: "2024-06-01T12:00:00Z",
    change: 0,
    commitSha: "abc123",
    testCount: 42,
  },
  name: "react",
  owner: "facebook",
  updateStatus: "up-to-date",
  ...overrides,
});

const createQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });

const renderDashboardContent = () => {
  const queryClient = createQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
        <DashboardContent />
      </NextIntlClientProvider>
    </QueryClientProvider>
  );
};

describe("DashboardContent", () => {
  beforeEach(() => {
    mockRepositories.mockReturnValue([]);
    mockIsLoading.mockReturnValue(false);
    mockHasNextPage.mockReturnValue(false);
    mockIsError.mockReturnValue(false);
    mockBookmarkOnly.mockReturnValue(false);
    mockOwnershipFilter.mockReturnValue("all");
    mockStatsData.mockReturnValue({ totalRepositories: 5, totalTests: 150 });
    mockStatsIsLoading.mockReturnValue(false);
    mockReanalyze.mockClear();
    mockAddBookmark.mockClear();
    mockRemoveBookmark.mockClear();
  });

  describe("dashboard page structure", () => {
    it("renders summary stats section", () => {
      renderDashboardContent();

      expect(screen.getByTestId("summary-section")).toBeInTheDocument();
    });

    it("renders repository list when repositories exist", () => {
      mockRepositories.mockReturnValue([
        createMockRepo({ fullName: "facebook/react", id: "1", name: "react", owner: "facebook" }),
        createMockRepo({ fullName: "vuejs/core", id: "2", name: "core", owner: "vuejs" }),
      ]);
      renderDashboardContent();

      expect(screen.getByText("facebook/react")).toBeInTheDocument();
      expect(screen.getByText("vuejs/core")).toBeInTheDocument();
    });

    it("renders empty state when no repositories exist", () => {
      mockRepositories.mockReturnValue([]);
      renderDashboardContent();

      expect(screen.getByText("No repositories yet")).toBeInTheDocument();
      expect(screen.getByTestId("analyze-dialog")).toBeInTheDocument();
    });
  });

  describe("bookmark filter behavior", () => {
    it("shows only bookmarked repos when bookmark filter is active", () => {
      mockBookmarkOnly.mockReturnValue(true);
      mockRepositories.mockReturnValue([
        createMockRepo({
          fullName: "facebook/react",
          id: "1",
          isBookmarked: true,
          name: "react",
          owner: "facebook",
        }),
        createMockRepo({
          fullName: "vuejs/core",
          id: "2",
          isBookmarked: false,
          name: "core",
          owner: "vuejs",
        }),
      ]);
      renderDashboardContent();

      expect(screen.getByText("facebook/react")).toBeInTheDocument();
      expect(screen.queryByText("vuejs/core")).not.toBeInTheDocument();
    });

    it("shows no-bookmarks empty state when bookmark filter yields zero results", () => {
      mockBookmarkOnly.mockReturnValue(true);
      mockRepositories.mockReturnValue([createMockRepo({ id: "1", isBookmarked: false })]);
      renderDashboardContent();

      expect(screen.getByText("No bookmarked repositories")).toBeInTheDocument();
    });
  });

  describe("search empty state", () => {
    it("shows no-search-results empty state when search matches nothing", () => {
      mockRepositories.mockReturnValue([
        createMockRepo({ fullName: "facebook/react", id: "1", name: "react", owner: "facebook" }),
      ]);
      renderDashboardContent();

      const searchInput = screen.getByRole("searchbox", { name: /search repositories/i });
      fireEvent.change(searchInput, { target: { value: "nonexistent-query" } });

      expect(screen.getByText("No matching repositories")).toBeInTheDocument();
    });
  });

  describe("empty state action", () => {
    it("renders analyze dialog button in no-repos empty state", () => {
      mockRepositories.mockReturnValue([]);
      renderDashboardContent();

      const analyzeButton = screen.getByTestId("analyze-dialog");

      expect(analyzeButton).toBeInTheDocument();
      expect(analyzeButton).toHaveAttribute("data-variant", "empty-state");
    });
  });

  describe("reanalysis flow via repository card", () => {
    it("triggers reanalyze when card reanalyze button is clicked", () => {
      mockRepositories.mockReturnValue([
        createMockRepo({
          id: "1",
          latestAnalysis: {
            analyzedAt: "2024-06-01T12:00:00Z",
            change: 0,
            commitSha: "abc123",
            testCount: 10,
          },
          updateStatus: "new-commits",
        }),
      ]);
      renderDashboardContent();

      fireEvent.click(screen.getByRole("button", { name: /reanalyze/i }));

      expect(mockReanalyze).toHaveBeenCalledWith("facebook", "react");
    });
  });
});
