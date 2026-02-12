import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ExploreContent } from "./explore-content";

const mockUseAuth = vi.fn();
const mockReanalyze = vi.fn();

vi.mock("@/features/auth", () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock("@/features/dashboard", () => ({
  EmptyStateVariant: ({ variant, searchQuery }: { searchQuery?: string; variant: string }) => (
    <div data-testid={`empty-state-${variant}`}>{searchQuery && <span>{searchQuery}</span>}</div>
  ),
  InfiniteScrollLoader: () => <div data-testid="infinite-scroll-loader" />,
  PaginationStatus: () => <div data-testid="pagination-status" />,
  RepositoryList: ({
    isLoading,
    repositories,
  }: {
    isLoading?: boolean;
    repositories: Array<{ name: string; owner: string }>;
  }) => (
    <div data-testid="repository-list">
      {isLoading ? (
        <span>Loading skeleton</span>
      ) : (
        repositories.map((r) => (
          <div data-testid="repo-card" key={`${r.owner}/${r.name}`}>
            {r.owner}/{r.name}
          </div>
        ))
      )}
    </div>
  ),
  useReanalyze: () => ({ reanalyze: mockReanalyze }),
}));

const mockUseExploreRepositories = vi.fn();

vi.mock("../hooks", () => ({
  useExploreRepositories: (...args: unknown[]) => mockUseExploreRepositories(...args),
}));

vi.mock("./login-required-state", () => ({
  LoginRequiredState: ({ titleKey }: { titleKey: string }) => (
    <div data-testid="login-required-state">{titleKey}</div>
  ),
}));

vi.mock("./my-repos-tab", () => ({
  MyReposTab: () => <div data-testid="my-repos-tab">My Repos Content</div>,
}));

vi.mock("./org-repos-tab", () => ({
  OrgReposTab: () => <div data-testid="org-repos-tab">Org Repos Content</div>,
}));

vi.mock("./search-sort-controls", () => ({
  SearchSortControls: () => <div data-testid="search-sort-controls" />,
}));

const messages = {
  explore: {
    community: {
      visibilityDisclosure: "Showing public repositories. Visibility determined at analysis time.",
    },
    loginRequired: {
      myReposDescription:
        "Connect your GitHub account to see and analyze your personal repositories.",
      myReposTitle: "Sign in to view your repositories",
      organizationsDescription:
        "Connect your GitHub account to access your organization repositories.",
      organizationsTitle: "Sign in to view organizations",
      signIn: "Sign in to continue",
    },
    searchPlaceholder: "Search repositories...",
    sort: {
      label: "Sort",
      name: "Name",
      recent: "Recent",
      tests: "Tests",
    },
    tabs: {
      community: "Community",
      myRepos: "My Repos",
      organizations: "Organizations",
    },
  },
};

const renderExploreContent = () =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <ExploreContent />
    </NextIntlClientProvider>
  );

describe("ExploreContent", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseAuth.mockReturnValue({ isAuthenticated: false });
  });

  describe("Community Tab > Empty State", () => {
    it("displays empty state when no repositories exist", () => {
      mockUseExploreRepositories.mockReturnValue({
        data: [],
        fetchNextPage: vi.fn(),
        hasNextPage: false,
        isError: false,
        isFetchingNextPage: false,
        isLoading: false,
        refetch: vi.fn(),
      });

      renderExploreContent();

      expect(screen.getByTestId("empty-state-no-repos")).toBeInTheDocument();
    });
  });

  describe("Community Tab > Repository Card List", () => {
    it("displays public repository card list", () => {
      mockUseExploreRepositories.mockReturnValue({
        data: [
          { name: "react", owner: "facebook" },
          { name: "vue", owner: "vuejs" },
        ],
        fetchNextPage: vi.fn(),
        hasNextPage: false,
        isError: false,
        isFetchingNextPage: false,
        isLoading: false,
        refetch: vi.fn(),
      });

      renderExploreContent();

      const cards = screen.getAllByTestId("repo-card");
      expect(cards).toHaveLength(2);
      expect(screen.getByText("facebook/react")).toBeInTheDocument();
      expect(screen.getByText("vuejs/vue")).toBeInTheDocument();
    });
  });
});
