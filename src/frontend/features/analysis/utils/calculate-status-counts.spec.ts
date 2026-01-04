import { describe, expect, it } from "vitest";

import { calculateStatusCounts } from "./calculate-status-counts";

describe("calculateStatusCounts", () => {
  it("should count active tests correctly", () => {
    const tests = [{ status: "active" }, { status: "active" }, { status: "skipped" }];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 2, skipped: 1, todo: 0 });
  });

  it("should treat focused as active", () => {
    const tests = [{ status: "focused" }, { status: "active" }];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 2, skipped: 0, todo: 0 });
  });

  it("should treat xfail as skipped", () => {
    const tests = [{ status: "xfail" }, { status: "skipped" }];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 0, skipped: 2, todo: 0 });
  });

  it("should count todo tests correctly", () => {
    const tests = [{ status: "todo" }, { status: "todo" }, { status: "active" }];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 1, skipped: 0, todo: 2 });
  });

  it("should return zero counts for empty array", () => {
    const result = calculateStatusCounts([]);

    expect(result).toEqual({ active: 0, skipped: 0, todo: 0 });
  });

  it("should handle all status types together", () => {
    const tests = [
      { status: "active" },
      { status: "focused" },
      { status: "skipped" },
      { status: "xfail" },
      { status: "todo" },
    ];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 2, skipped: 2, todo: 1 });
  });
});
