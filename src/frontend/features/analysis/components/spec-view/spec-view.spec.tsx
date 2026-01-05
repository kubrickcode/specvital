import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { LanguageSelector } from "./language-selector";
import { SpecItem } from "./spec-item";
import { SpecViewSkeleton } from "./spec-view-skeleton";

const messages = {
  analyze: {
    specView: {
      cancel: "Cancel",
      convert: "Convert",
      errorTitle: "Conversion Failed",
      languageLabel: "Language",
      retry: "Try Again",
      summary: "{total} tests ({cached} cached, {converted} converted)",
    },
  },
};

describe("SpecItem", () => {
  it("renders converted name and original name on hover", () => {
    const item = {
      convertedName: "User authentication check",
      isFromCache: true,
      line: 42,
      originalName: "should authenticate user when credentials valid",
      status: "active" as const,
    };

    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <SpecItem item={item} />
      </NextIntlClientProvider>
    );

    expect(screen.getByText("User authentication check")).toBeInTheDocument();
    expect(screen.getByText("should authenticate user when credentials valid")).toBeInTheDocument();
    expect(screen.getByText("L:42")).toBeInTheDocument();
  });

  it("displays correct status icon for different statuses", () => {
    const statuses = ["active", "focused", "skipped", "todo", "xfail"] as const;

    for (const status of statuses) {
      const item = {
        convertedName: `Test for ${status}`,
        isFromCache: false,
        line: 1,
        originalName: "original",
        status,
      };

      const { unmount } = render(
        <NextIntlClientProvider locale="en" messages={messages}>
          <SpecItem item={item} />
        </NextIntlClientProvider>
      );

      const statusLabels: Record<typeof status, string> = {
        active: "Active test",
        focused: "Focused test",
        skipped: "Skipped test",
        todo: "Todo test",
        xfail: "Expected failure",
      };

      expect(screen.getByLabelText(statusLabels[status])).toBeInTheDocument();
      unmount();
    }
  });
});

describe("LanguageSelector", () => {
  it("renders with current language", () => {
    const onChange = vi.fn();

    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <LanguageSelector onChange={onChange} value="en" />
      </NextIntlClientProvider>
    );

    expect(screen.getByRole("button")).toHaveTextContent("English");
  });
});

describe("SpecViewSkeleton", () => {
  it("renders loading skeleton", () => {
    render(<SpecViewSkeleton />);

    expect(screen.getByRole("status")).toHaveAttribute("aria-label", "Loading spec view");
  });
});
