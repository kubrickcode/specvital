import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { TrustBadgesWithDialog } from "./trust-badges-with-dialog";

vi.mock("@/lib/motion", () => ({
  staggerContainer: {},
  staggerItem: {},
  useReducedMotion: () => true,
}));

vi.mock("motion/react", () => ({
  motion: {
    button: ({ children, ...props }: { children: React.ReactNode; [key: string]: unknown }) => (
      <button {...props}>{children}</button>
    ),
    div: ({ children, ...props }: { children: React.ReactNode; [key: string]: unknown }) => (
      <div {...props}>{children}</div>
    ),
  },
}));

const messages = {
  home: {
    frameworks: {
      cpp: "C++",
      csharp: "C#",
      description: "List of supported testing frameworks and their file patterns",
      go: "Go",
      java: "Java",
      javascript: "JavaScript / TypeScript",
      kotlin: "Kotlin",
      php: "PHP",
      python: "Python",
      ruby: "Ruby",
      rust: "Rust",
      swift: "Swift",
      title: "Supported Frameworks",
    },
    trustBadges: {
      accurate: "AST-powered",
      free: "Free to start",
      instant: "Instant analysis",
      multiFramework: "20+ frameworks",
    },
  },
};

const renderComponent = () =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <TrustBadgesWithDialog />
    </NextIntlClientProvider>
  );

describe("TrustBadgesWithDialog", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("opens supported frameworks dialog when 20+ frameworks badge is clicked", () => {
    renderComponent();

    const frameworkButton = screen.getByRole("button", { name: /20\+ frameworks/i });
    expect(frameworkButton).toBeInTheDocument();

    fireEvent.click(frameworkButton);

    expect(screen.getByText("Supported Frameworks")).toBeInTheDocument();
    expect(screen.getByText("JavaScript / TypeScript")).toBeInTheDocument();
  });

  it("displays framework categories inside the dialog", () => {
    renderComponent();

    fireEvent.click(screen.getByRole("button", { name: /20\+ frameworks/i }));

    expect(screen.getByText("Python")).toBeInTheDocument();
    expect(screen.getByText("Go")).toBeInTheDocument();
    expect(screen.getByText("Java")).toBeInTheDocument();
  });

  it("renders individual framework names as badges in the dialog", () => {
    renderComponent();

    fireEvent.click(screen.getByRole("button", { name: /20\+ frameworks/i }));

    expect(screen.getByText("Jest")).toBeInTheDocument();
    expect(screen.getByText("Vitest")).toBeInTheDocument();
    expect(screen.getByText("pytest")).toBeInTheDocument();
    expect(screen.getByText("JUnit 5")).toBeInTheDocument();
    expect(screen.getByText("Go Testing")).toBeInTheDocument();
  });
});
