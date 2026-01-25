import { render, screen } from "@testing-library/react";
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
    status: {
      focused: "Focused",
      newCommits: "New commits",
      skippedHigh: "High skip",
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

      expect(screen.getByRole("button", { name: /remove bookmark/i })).toBeInTheDocument();
      expect(screen.getByLabelText("Bookmarked")).toBeInTheDocument();
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
});
