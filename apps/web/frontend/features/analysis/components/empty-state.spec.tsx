import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { HIGHLIGHTED_FRAMEWORKS, TOTAL_FRAMEWORK_COUNT } from "@/lib/constants/frameworks";

import { EmptyState } from "./empty-state";

const mockPush = vi.fn();

vi.mock("@/i18n/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

const messages = {
  emptyState: {
    analyzeAnother: "Analyze Another Repository",
    description: "We couldn't find any test files in this repository.",
    supportedFrameworks: "Supported frameworks:",
    title: "No tests found",
  },
};

const renderEmptyState = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <EmptyState />
    </NextIntlClientProvider>
  );
};

describe("EmptyState", () => {
  it("renders empty state title and description", () => {
    renderEmptyState();

    expect(screen.getByText("No tests found")).toBeInTheDocument();
    expect(
      screen.getByText("We couldn't find any test files in this repository.")
    ).toBeInTheDocument();
  });

  it("displays highlighted frameworks as badges", () => {
    renderEmptyState();

    expect(screen.getByText("Supported frameworks:")).toBeInTheDocument();
    for (const framework of HIGHLIGHTED_FRAMEWORKS) {
      expect(screen.getAllByText(framework).length).toBeGreaterThanOrEqual(1);
    }
  });

  it("shows remaining framework count badge", () => {
    renderEmptyState();

    const remainingCount = TOTAL_FRAMEWORK_COUNT - HIGHLIGHTED_FRAMEWORKS.length;
    expect(screen.getByText(`+${remainingCount} more`)).toBeInTheDocument();
  });

  it("displays all framework names in detail section", () => {
    renderEmptyState();

    expect(screen.getByText(/Jest, Vitest, Mocha, Playwright, Cypress/)).toBeInTheDocument();
    expect(screen.getByText(/pytest, unittest/)).toBeInTheDocument();
    expect(screen.getAllByText(/Go Testing/).length).toBeGreaterThanOrEqual(1);
  });

  it("renders analyze another repository button", () => {
    renderEmptyState();

    const button = screen.getByRole("button", { name: "Analyze Another Repository" });
    expect(button).toBeInTheDocument();
  });

  it("has accessible icon with aria-hidden", () => {
    renderEmptyState();

    const svg = document.querySelector("svg");
    expect(svg).toBeInTheDocument();
    expect(svg).toHaveAttribute("aria-hidden", "true");
  });
});
