import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import type { StatusCounts } from "../utils/calculate-status-counts";

import { StatusMiniBar } from "./status-mini-bar";

const renderStatusMiniBar = (counts: StatusCounts) => {
  return render(<StatusMiniBar counts={counts} />);
};

describe("StatusMiniBar", () => {
  it("includes all 5 status types in aria-label", () => {
    const counts: StatusCounts = { active: 10, focused: 2, skipped: 3, todo: 1, xfail: 4 };

    renderStatusMiniBar(counts);

    const bar = screen.getByRole("img");
    expect(bar).toHaveAttribute(
      "aria-label",
      "10 active, 2 focused, 3 skipped, 4 xfail, 1 todo out of 20 tests"
    );
  });

  it("renders nothing when total count is 0", () => {
    const counts: StatusCounts = { active: 0, focused: 0, skipped: 0, todo: 0, xfail: 0 };

    const { container } = renderStatusMiniBar(counts);

    expect(container.firstChild).toBeNull();
  });

  it("does not render focused segment when focused count is 0", () => {
    const counts: StatusCounts = { active: 10, focused: 0, skipped: 5, todo: 0, xfail: 0 };

    renderStatusMiniBar(counts);

    const bar = screen.getByRole("img");
    const segments = bar.querySelectorAll("div");
    const classNames = Array.from(segments).map((s) => s.className);

    expect(classNames).toContain("bg-status-active");
    expect(classNames).toContain("bg-status-skipped");
    expect(classNames).not.toContain("bg-status-focused");
  });

  it("does not render xfail segment when xfail count is 0", () => {
    const counts: StatusCounts = { active: 10, focused: 0, skipped: 5, todo: 0, xfail: 0 };

    renderStatusMiniBar(counts);

    const bar = screen.getByRole("img");
    const segments = bar.querySelectorAll("div");
    const classNames = Array.from(segments).map((s) => s.className);

    expect(classNames).not.toContain("bg-status-xfail");
  });

  it("renders focused segment when focused count is greater than 0", () => {
    const counts: StatusCounts = { active: 10, focused: 3, skipped: 0, todo: 0, xfail: 0 };

    renderStatusMiniBar(counts);

    const bar = screen.getByRole("img");
    const segments = bar.querySelectorAll("div");
    const classNames = Array.from(segments).map((s) => s.className);

    expect(classNames).toContain("bg-status-focused");
  });

  it("renders xfail segment when xfail count is greater than 0", () => {
    const counts: StatusCounts = { active: 10, focused: 0, skipped: 0, todo: 0, xfail: 2 };

    renderStatusMiniBar(counts);

    const bar = screen.getByRole("img");
    const segments = bar.querySelectorAll("div");
    const classNames = Array.from(segments).map((s) => s.className);

    expect(classNames).toContain("bg-status-xfail");
  });

  it("only renders segments for statuses with non-zero counts", () => {
    const counts: StatusCounts = { active: 5, focused: 0, skipped: 0, todo: 0, xfail: 0 };

    renderStatusMiniBar(counts);

    const bar = screen.getByRole("img");
    const segments = bar.querySelectorAll("div");

    expect(segments).toHaveLength(1);
    expect(segments[0]).toHaveClass("bg-status-active");
  });

  it("sets correct width percentage for each segment", () => {
    const counts: StatusCounts = { active: 50, focused: 0, skipped: 50, todo: 0, xfail: 0 };

    renderStatusMiniBar(counts);

    const bar = screen.getByRole("img");
    const segments = bar.querySelectorAll("div");

    expect(segments[0]).toHaveStyle({ width: "50%" });
    expect(segments[1]).toHaveStyle({ width: "50%" });
  });
});
