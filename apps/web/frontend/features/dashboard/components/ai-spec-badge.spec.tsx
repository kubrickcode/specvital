import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { AiSpecSummary } from "@/lib/api/types";

import { AiSpecBadge } from "./ai-spec-badge";

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
    aiSpec: {
      badge: "AI Spec",
      languageCount: "{count} languages",
      lastGenerated: "Generated {time}",
      viewDetails: "View details",
    },
  },
};

const createSummary = (overrides?: Partial<AiSpecSummary>): AiSpecSummary => ({
  hasSpec: true,
  languageCount: 3,
  latestGeneratedAt: "2025-01-15T10:00:00Z",
  ...overrides,
});

type RenderOptions = {
  owner?: string;
  repo?: string;
  summary?: AiSpecSummary;
};

const renderAiSpecBadge = (options: RenderOptions = {}) => {
  const props = {
    owner: options.owner ?? "facebook",
    repo: options.repo ?? "react",
    summary: options.summary ?? createSummary(),
  };

  return render(
    <NextIntlClientProvider
      locale="en"
      messages={messages}
      now={new Date("2025-01-15T12:00:00Z")}
      timeZone="UTC"
    >
      <AiSpecBadge {...props} />
    </NextIntlClientProvider>
  );
};

describe("AiSpecBadge", () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2025-01-15T12:00:00Z"));
  });

  it("renders the badge when hasSpec is true", () => {
    renderAiSpecBadge();

    expect(screen.getByLabelText("AI Spec")).toBeInTheDocument();
    expect(screen.getByText("AI Spec")).toBeInTheDocument();
  });

  it("renders nothing when hasSpec is false", () => {
    const { container } = renderAiSpecBadge({
      summary: createSummary({ hasSpec: false }),
    });

    expect(container.querySelector("[aria-label='AI Spec']")).not.toBeInTheDocument();
  });

  it("shows popover with language count on badge click", () => {
    renderAiSpecBadge({
      summary: createSummary({ languageCount: 5 }),
    });

    fireEvent.click(screen.getByLabelText("AI Spec"));

    expect(screen.getByText("5 languages")).toBeInTheDocument();
  });

  it("shows view details link pointing to spec tab", () => {
    renderAiSpecBadge({ owner: "vercel", repo: "next.js" });

    fireEvent.click(screen.getByLabelText("AI Spec"));

    const link = screen.getByRole("link", { name: "View details" });

    expect(link).toHaveAttribute("href", "/analyze/vercel/next.js?tab=spec");
  });

  it("does not show language count when languageCount is undefined", () => {
    renderAiSpecBadge({
      summary: createSummary({ languageCount: undefined }),
    });

    fireEvent.click(screen.getByLabelText("AI Spec"));

    expect(screen.queryByText(/languages/)).not.toBeInTheDocument();
  });

  it("does not show language count when languageCount is 0", () => {
    renderAiSpecBadge({
      summary: createSummary({ languageCount: 0 }),
    });

    fireEvent.click(screen.getByLabelText("AI Spec"));

    expect(screen.queryByText(/languages/)).not.toBeInTheDocument();
  });
});
