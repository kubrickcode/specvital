import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { EmptyStateVariant } from "./empty-state-variant";

const messages = {
  dashboard: {
    emptyStateVariant: {
      error: {
        description: "Something went wrong. Please try again.",
        title: "Error",
      },
      noBookmarks: {
        description: "You haven't bookmarked any repositories yet.",
        title: "No Bookmarks",
      },
      noRepos: {
        description: "Get started by analyzing a repository.",
        title: "No Repositories",
      },
      noSearchResults: {
        description: "No repositories match your search.",
        descriptionWithQuery: 'No results for "{query}".',
        title: "No Results",
      },
    },
  },
};

type RenderOptions = {
  action?: React.ReactNode;
  searchQuery?: string;
  variant: "no-repos" | "no-search-results" | "no-bookmarks" | "error";
};

const renderEmptyState = (options: RenderOptions) => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <EmptyStateVariant {...options} />
    </NextIntlClientProvider>
  );
};

describe("EmptyStateVariant", () => {
  describe("no-repos variant", () => {
    it("displays the title and description", () => {
      renderEmptyState({ variant: "no-repos" });

      expect(screen.getByText("No Repositories")).toBeInTheDocument();
      expect(screen.getByText("Get started by analyzing a repository.")).toBeInTheDocument();
    });

    it("renders action button when provided", () => {
      renderEmptyState({
        action: <button type="button">Analyze a repository</button>,
        variant: "no-repos",
      });

      expect(screen.getByRole("button", { name: "Analyze a repository" })).toBeInTheDocument();
    });

    it("triggers callback when action button is clicked", () => {
      const onClick = vi.fn();
      renderEmptyState({
        action: (
          <button onClick={onClick} type="button">
            Analyze a repository
          </button>
        ),
        variant: "no-repos",
      });

      fireEvent.click(screen.getByRole("button", { name: "Analyze a repository" }));

      expect(onClick).toHaveBeenCalledOnce();
    });
  });

  describe("no-search-results variant", () => {
    it("shows query-specific description when searchQuery is provided", () => {
      renderEmptyState({ searchQuery: "nonexistent", variant: "no-search-results" });

      expect(screen.getByText('No results for "nonexistent".')).toBeInTheDocument();
    });

    it("shows generic description when searchQuery is not provided", () => {
      renderEmptyState({ variant: "no-search-results" });

      expect(screen.getByText("No repositories match your search.")).toBeInTheDocument();
    });
  });

  describe("no-bookmarks variant", () => {
    it("displays the no bookmarks title and description", () => {
      renderEmptyState({ variant: "no-bookmarks" });

      expect(screen.getByText("No Bookmarks")).toBeInTheDocument();
      expect(screen.getByText("You haven't bookmarked any repositories yet.")).toBeInTheDocument();
    });
  });

  describe("error variant", () => {
    it("displays the error title and description", () => {
      renderEmptyState({ variant: "error" });

      expect(screen.getByText("Error")).toBeInTheDocument();
      expect(screen.getByText("Something went wrong. Please try again.")).toBeInTheDocument();
    });
  });

  it("does not render action area when action is not provided", () => {
    const { container } = renderEmptyState({ variant: "no-repos" });

    expect(container.querySelectorAll("button")).toHaveLength(0);
  });
});
