import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { RepositoryCard as RepositoryCardType } from "@/lib/api/types";

import { RepositoryCard, type RepositoryCardProps } from "./repository-card";

const mockUseAuth = vi.fn();

vi.mock("@/features/auth/hooks/use-auth", () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock("@/features/auth/hooks/use-login-modal", () => ({
  useLoginModal: () => ({
    open: vi.fn(),
  }),
}));

vi.mock("@/i18n/navigation", () => ({
  Link: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
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
    status: {
      analyzing: "Analyzing...",
      focused: "Focused",
      newCommits: "New commits",
      skippedHigh: "High skip",
      unknown: "Status unknown",
      upToDate: "Up to date",
    },
  },
};

const createMockRepo = (overrides?: Partial<RepositoryCardType>): RepositoryCardType => ({
  fullName: "facebook/react",
  id: "1",
  isAnalyzedByMe: false,
  isBookmarked: false,
  latestAnalysis: undefined,
  name: "react",
  owner: "facebook",
  updateStatus: "up-to-date",
  ...overrides,
});

type AnalyzedCardTestProps = {
  onBookmarkToggle?: (owner: string, repo: string, isBookmarked: boolean) => void;
  onReanalyze?: (owner: string, repo: string) => void;
  repo?: RepositoryCardType;
  variant?: "dashboard" | "explore";
};

const renderRepositoryCard = (props: AnalyzedCardTestProps = {}) => {
  const cardProps: RepositoryCardProps = {
    onBookmarkToggle: props.onBookmarkToggle,
    onReanalyze: props.onReanalyze,
    repo: props.repo ?? createMockRepo(),
    variant: props.variant,
  };

  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <RepositoryCard {...cardProps} />
    </NextIntlClientProvider>
  );
};

describe("RepositoryCard", () => {
  beforeEach(() => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
    });
  });

  describe("dashboard variant (default)", () => {
    it("shows star button for bookmark toggle", () => {
      renderRepositoryCard({ variant: "dashboard" });

      expect(screen.getByRole("button", { name: /add bookmark/i })).toBeInTheDocument();
    });

    it("shows filled star when bookmarked", () => {
      renderRepositoryCard({
        repo: createMockRepo({ isBookmarked: true }),
        variant: "dashboard",
      });

      const button = screen.getByRole("button", { name: /remove bookmark/i });

      expect(button).toBeInTheDocument();
      expect(button).toHaveAttribute("aria-pressed", "true");
    });
  });

  describe("explore variant", () => {
    it("does not show action button", () => {
      renderRepositoryCard({ variant: "explore" });

      expect(screen.queryByRole("button")).not.toBeInTheDocument();
    });

    it("does not show bookmarked star icon", () => {
      renderRepositoryCard({
        repo: createMockRepo({ isBookmarked: true }),
        variant: "explore",
      });

      expect(screen.queryByLabelText("Bookmarked")).not.toBeInTheDocument();
    });
  });

  describe("unauthenticated user", () => {
    beforeEach(() => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: false,
        isLoading: false,
      });
    });

    it("shows login prompt for dashboard variant", () => {
      renderRepositoryCard({ variant: "dashboard" });

      expect(screen.getByRole("button", { name: /sign in to bookmark/i })).toBeInTheDocument();
    });
  });

  describe("summary stats section", () => {
    it("renders test count and status badge when analysis exists", () => {
      renderRepositoryCard({
        repo: createMockRepo({
          latestAnalysis: {
            analyzedAt: "2024-06-01T12:00:00Z",
            change: 5,
            commitSha: "abc123",
            testCount: 42,
          },
          updateStatus: "up-to-date",
        }),
      });

      expect(screen.getByText("42")).toBeInTheDocument();
      expect(screen.getByText("tests")).toBeInTheDocument();
      expect(screen.getByText("Up to date")).toBeInTheDocument();
    });

    it("shows no analysis message when latestAnalysis is absent", () => {
      renderRepositoryCard({
        repo: createMockRepo({ latestAnalysis: undefined }),
      });

      expect(screen.getByText("No analysis yet")).toBeInTheDocument();
    });
  });

  describe("reanalysis", () => {
    it("shows reanalyze button for repos with new-commits status", () => {
      renderRepositoryCard({
        repo: createMockRepo({
          latestAnalysis: {
            analyzedAt: "2024-06-01T12:00:00Z",
            change: 0,
            commitSha: "abc123",
            testCount: 10,
          },
          updateStatus: "new-commits",
        }),
      });

      expect(screen.getByRole("button", { name: /reanalyze/i })).toBeInTheDocument();
      expect(screen.getByText("Update")).toBeInTheDocument();
    });

    it("calls onReanalyze when reanalyze button is clicked", () => {
      const onReanalyze = vi.fn();
      renderRepositoryCard({
        onReanalyze,
        repo: createMockRepo({
          latestAnalysis: {
            analyzedAt: "2024-06-01T12:00:00Z",
            change: 0,
            commitSha: "abc123",
            testCount: 10,
          },
          updateStatus: "new-commits",
        }),
      });

      fireEvent.click(screen.getByRole("button", { name: /reanalyze/i }));

      expect(onReanalyze).toHaveBeenCalledWith("facebook", "react");
    });

    it("does not show reanalyze button for up-to-date repos", () => {
      renderRepositoryCard({
        repo: createMockRepo({
          latestAnalysis: {
            analyzedAt: "2024-06-01T12:00:00Z",
            change: 0,
            commitSha: "abc123",
            testCount: 10,
          },
          updateStatus: "up-to-date",
        }),
      });

      expect(screen.queryByRole("button", { name: /reanalyze/i })).not.toBeInTheDocument();
    });

    it("hides reanalyze button and shows analyzing badge when status becomes analyzing", () => {
      renderRepositoryCard({
        repo: createMockRepo({
          latestAnalysis: {
            analyzedAt: "2024-06-01T12:00:00Z",
            change: 0,
            commitSha: "abc123",
            testCount: 10,
          },
          updateStatus: "analyzing" as "up-to-date",
        }),
      });

      expect(screen.queryByRole("button", { name: /reanalyze/i })).not.toBeInTheDocument();
      expect(screen.getByText("Analyzing...")).toBeInTheDocument();
    });
  });
});
