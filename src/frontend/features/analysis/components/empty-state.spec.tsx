import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

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

  it("displays supported frameworks list", () => {
    renderEmptyState();

    expect(screen.getByText("Supported frameworks:")).toBeInTheDocument();
    expect(screen.getByText("Jest")).toBeInTheDocument();
    expect(screen.getByText("Vitest")).toBeInTheDocument();
    expect(screen.getByText("Playwright")).toBeInTheDocument();
    expect(screen.getByText("Go")).toBeInTheDocument();
  });

  it("shows framework file patterns", () => {
    renderEmptyState();

    expect(screen.getByText("*.test.ts, *.test.tsx, *.test.js, __tests__/*")).toBeInTheDocument();
    expect(screen.getByText("*_test.go")).toBeInTheDocument();
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
