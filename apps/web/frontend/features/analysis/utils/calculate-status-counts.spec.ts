import { describe, expect, it } from "vitest";

import { calculateStatusCounts } from "./calculate-status-counts";

describe("calculateStatusCounts", () => {
  it("should count active tests correctly", () => {
    const tests = [{ status: "active" }, { status: "active" }, { status: "skipped" }];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 2, focused: 0, skipped: 1, todo: 0, xfail: 0 });
  });

  it("should count focused tests separately from active", () => {
    const tests = [{ status: "focused" }, { status: "active" }];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 1, focused: 1, skipped: 0, todo: 0, xfail: 0 });
  });

  it("should count xfail tests separately from skipped", () => {
    const tests = [{ status: "xfail" }, { status: "skipped" }];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 0, focused: 0, skipped: 1, todo: 0, xfail: 1 });
  });

  it("should count todo tests correctly", () => {
    const tests = [{ status: "todo" }, { status: "todo" }, { status: "active" }];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 1, focused: 0, skipped: 0, todo: 2, xfail: 0 });
  });

  it("should return zero counts for empty array", () => {
    const result = calculateStatusCounts([]);

    expect(result).toEqual({ active: 0, focused: 0, skipped: 0, todo: 0, xfail: 0 });
  });

  it("should count all 5 status types individually", () => {
    const tests = [
      { status: "active" },
      { status: "focused" },
      { status: "skipped" },
      { status: "xfail" },
      { status: "todo" },
    ];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 1, focused: 1, skipped: 1, todo: 1, xfail: 1 });
  });

  it("should count multiple tests of each status type", () => {
    const tests = [
      { status: "active" },
      { status: "active" },
      { status: "focused" },
      { status: "focused" },
      { status: "focused" },
      { status: "skipped" },
      { status: "xfail" },
      { status: "xfail" },
      { status: "todo" },
    ];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 2, focused: 3, skipped: 1, todo: 1, xfail: 2 });
  });

  it("should ignore unknown status values", () => {
    const tests = [{ status: "active" }, { status: "unknown" }, { status: "invalid" }];

    const result = calculateStatusCounts(tests);

    expect(result).toEqual({ active: 1, focused: 0, skipped: 0, todo: 0, xfail: 0 });
  });
});
