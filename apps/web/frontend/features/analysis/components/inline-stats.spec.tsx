import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it } from "vitest";

import type { Summary } from "@/lib/api/types";

import { InlineStats } from "./inline-stats";

const messages = {
  stats: {
    active: "Active",
    focused: "Focused",
    skipped: "Skipped",
    todo: "Todo",
    total: "Total",
    xfail: "Xfail",
  },
};

const baseSummary: Summary = {
  active: 80,
  focused: 0,
  frameworks: [
    { active: 80, focused: 0, framework: "vitest", skipped: 10, todo: 0, total: 90, xfail: 0 },
  ],
  skipped: 10,
  todo: 0,
  total: 90,
  xfail: 0,
};

const renderInlineStats = (summary: Summary) => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <InlineStats summary={summary} />
    </NextIntlClientProvider>
  );
};

describe("InlineStats", () => {
  it("renders total, active, and skipped counts", () => {
    renderInlineStats(baseSummary);

    expect(screen.getByText("90")).toBeInTheDocument();
    expect(screen.getByText("Total")).toBeInTheDocument();
    expect(screen.getByText("80")).toBeInTheDocument();
    expect(screen.getByText("Active")).toBeInTheDocument();
    expect(screen.getByText("10")).toBeInTheDocument();
    expect(screen.getByText("Skipped")).toBeInTheDocument();
  });

  it("formats large numbers with locale separators", () => {
    const summary: Summary = {
      ...baseSummary,
      active: 1234,
      skipped: 56,
      total: 1290,
    };

    renderInlineStats(summary);

    expect(screen.getByText("1,290")).toBeInTheDocument();
    expect(screen.getByText("1,234")).toBeInTheDocument();
  });

  it("hides focused count when value is 0", () => {
    renderInlineStats(baseSummary);

    expect(screen.queryByText("Focused")).not.toBeInTheDocument();
  });

  it("shows focused count when value is greater than 0", () => {
    const summary: Summary = { ...baseSummary, focused: 3, total: 93 };

    renderInlineStats(summary);

    expect(screen.getByText("Focused")).toBeInTheDocument();
    expect(screen.getByText("3")).toBeInTheDocument();
  });

  it("hides xfail count when value is 0", () => {
    renderInlineStats(baseSummary);

    expect(screen.queryByText("Xfail")).not.toBeInTheDocument();
  });

  it("shows xfail count when value is greater than 0", () => {
    const summary: Summary = { ...baseSummary, xfail: 5, total: 95 };

    renderInlineStats(summary);

    expect(screen.getByText("Xfail")).toBeInTheDocument();
    expect(screen.getByText("5")).toBeInTheDocument();
  });

  it("displays all 5 status types when all have non-zero counts", () => {
    const summary: Summary = {
      active: 100,
      focused: 2,
      frameworks: [
        {
          active: 100,
          focused: 2,
          framework: "jest",
          skipped: 10,
          todo: 3,
          total: 120,
          xfail: 5,
        },
      ],
      skipped: 10,
      todo: 3,
      total: 120,
      xfail: 5,
    };

    renderInlineStats(summary);

    expect(screen.getByText("Active")).toBeInTheDocument();
    expect(screen.getByText("Focused")).toBeInTheDocument();
    expect(screen.getByText("Skipped")).toBeInTheDocument();
    expect(screen.getByText("Xfail")).toBeInTheDocument();
    expect(screen.getByText("Todo")).toBeInTheDocument();
  });

  it("hides framework breakdown when only one framework exists", () => {
    renderInlineStats(baseSummary);

    expect(screen.queryByText("vitest")).not.toBeInTheDocument();
  });

  it("shows framework breakdown with percentages when multiple frameworks exist", () => {
    const summary: Summary = {
      ...baseSummary,
      frameworks: [
        { active: 60, focused: 0, framework: "vitest", skipped: 0, todo: 0, total: 60, xfail: 0 },
        { active: 30, focused: 0, framework: "jest", skipped: 0, todo: 0, total: 30, xfail: 0 },
      ],
      total: 90,
    };

    renderInlineStats(summary);

    expect(screen.getByText("vitest")).toBeInTheDocument();
    expect(screen.getByText("jest")).toBeInTheDocument();
    expect(screen.getByText("(67%)")).toBeInTheDocument();
    expect(screen.getByText("(33%)")).toBeInTheDocument();
  });

  it("renders framework names as badge-like elements with count and percentage", () => {
    const summary: Summary = {
      ...baseSummary,
      frameworks: [
        { active: 80, focused: 0, framework: "jest", skipped: 10, todo: 0, total: 90, xfail: 0 },
        { active: 10, focused: 0, framework: "vitest", skipped: 0, todo: 0, total: 10, xfail: 0 },
      ],
      total: 100,
    };

    renderInlineStats(summary);

    // Each framework renders in a styled container span with name, count, and percentage
    const jestLabel = screen.getByText("jest");
    const jestContainer = jestLabel.parentElement!;
    expect(jestContainer).toHaveClass("bg-muted/50");
    expect(jestContainer).toHaveTextContent("jest");
    expect(jestContainer).toHaveTextContent("90");
    expect(jestContainer).toHaveTextContent("(90%)");

    const vitestLabel = screen.getByText("vitest");
    const vitestContainer = vitestLabel.parentElement!;
    expect(vitestContainer).toHaveClass("bg-muted/50");
    expect(vitestContainer).toHaveTextContent("vitest");
    expect(vitestContainer).toHaveTextContent("10");
    expect(vitestContainer).toHaveTextContent("(10%)");
  });
});
