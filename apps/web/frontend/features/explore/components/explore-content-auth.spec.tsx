import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import React from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ExploreContent } from "./explore-content";

const mockUseAuth = vi.fn();
const mockReanalyze = vi.fn();

vi.mock("@/features/auth", () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock("@/features/dashboard", () => ({
  EmptyStateVariant: ({ variant }: { variant: string }) => (
    <div data-testid={`empty-state-${variant}`} />
  ),
  InfiniteScrollLoader: () => <div data-testid="infinite-scroll-loader" />,
  PaginationStatus: () => <div data-testid="pagination-status" />,
  RepositoryList: ({ repositories }: { repositories: unknown[] }) => (
    <div data-testid="repository-list">{repositories.length} repos</div>
  ),
  useReanalyze: () => ({ reanalyze: mockReanalyze }),
}));

// Radix Tabs does not respond to pointer/click events in jsdom.
// Mock the Tabs components with a simple implementation to test the auth gating behavior.
vi.mock("@/components/ui/tabs", () => {
  const TabsContext = React.createContext({
    activeTab: "",
    setActiveTab: (_v: string) => {},
  });

  return {
    Tabs: ({
      children,
      defaultValue,
      onValueChange,
      value,
    }: {
      children: React.ReactNode;
      defaultValue?: string;
      onValueChange?: (value: string) => void;
      value?: string;
    }) => {
      const [internalValue, setInternalValue] = React.useState(defaultValue ?? "");
      const activeTab = value ?? internalValue;
      const setActiveTab = (v: string) => {
        setInternalValue(v);
        onValueChange?.(v);
      };
      return (
        <TabsContext.Provider value={{ activeTab, setActiveTab }}>
          <div data-testid="tabs">{children}</div>
        </TabsContext.Provider>
      );
    },
    TabsContent: ({ children, value }: { children: React.ReactNode; value: string }) => {
      const { activeTab } = React.useContext(TabsContext);
      if (activeTab !== value) return null;
      return <div data-testid={`tab-content-${value}`}>{children}</div>;
    },
    TabsList: ({ children }: { children: React.ReactNode }) => (
      <div data-testid="tabs-list" role="tablist">
        {children}
      </div>
    ),
    TabsTrigger: ({ children, value }: { children: React.ReactNode; value: string }) => {
      const { activeTab, setActiveTab } = React.useContext(TabsContext);
      return (
        <button
          aria-selected={activeTab === value}
          data-state={activeTab === value ? "active" : "inactive"}
          onClick={() => setActiveTab(value)}
          role="tab"
          type="button"
        >
          {children}
        </button>
      );
    },
  };
});

const mockUseExploreRepositories = vi.fn();

vi.mock("../hooks", () => ({
  useExploreRepositories: (...args: unknown[]) => mockUseExploreRepositories(...args),
}));

vi.mock("./login-required-state", () => ({
  LoginRequiredState: ({
    descriptionKey,
    titleKey,
  }: {
    descriptionKey: string;
    titleKey: string;
  }) => (
    <div data-testid="login-required-state">
      <span data-testid="login-title-key">{titleKey}</span>
      <span data-testid="login-description-key">{descriptionKey}</span>
    </div>
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

const defaultExploreReturn = {
  data: [{ name: "react", owner: "facebook" }],
  fetchNextPage: vi.fn(),
  hasNextPage: false,
  isError: false,
  isFetchingNextPage: false,
  isLoading: false,
  refetch: vi.fn(),
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
    mockUseExploreRepositories.mockReturnValue(defaultExploreReturn);
  });

  describe("Authenticated User > My Repos Tab", () => {
    it("renders MyReposTab when authenticated and My Repos tab is selected", () => {
      mockUseAuth.mockReturnValue({ isAuthenticated: true });

      renderExploreContent();

      fireEvent.click(screen.getByRole("tab", { name: /my repos/i }));

      expect(screen.getByTestId("my-repos-tab")).toBeInTheDocument();
      expect(screen.getByText("My Repos Content")).toBeInTheDocument();
    });
  });

  describe("Authenticated User > Organizations Tab", () => {
    it("renders OrgReposTab when authenticated and Organizations tab is selected", () => {
      mockUseAuth.mockReturnValue({ isAuthenticated: true });

      renderExploreContent();

      fireEvent.click(screen.getByRole("tab", { name: /organizations/i }));

      expect(screen.getByTestId("org-repos-tab")).toBeInTheDocument();
      expect(screen.getByText("Org Repos Content")).toBeInTheDocument();
    });
  });

  describe("Unauthenticated User > Login Required State", () => {
    it("shows login required state on My Repos tab when not authenticated", () => {
      mockUseAuth.mockReturnValue({ isAuthenticated: false });

      renderExploreContent();

      fireEvent.click(screen.getByRole("tab", { name: /my repos/i }));

      expect(screen.getByTestId("login-required-state")).toBeInTheDocument();
      expect(screen.getByTestId("login-title-key")).toHaveTextContent("myReposTitle");
    });

    it("shows login required state on Organizations tab when not authenticated", () => {
      mockUseAuth.mockReturnValue({ isAuthenticated: false });

      renderExploreContent();

      fireEvent.click(screen.getByRole("tab", { name: /organizations/i }));

      expect(screen.getByTestId("login-title-key")).toHaveTextContent("organizationsTitle");
    });
  });
});
