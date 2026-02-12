import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import type { TestSuite } from "@/lib/api";

import { TestList } from "./test-list";

vi.mock("./empty-state", () => ({
  EmptyState: () => <div data-testid="empty-state">No tests found</div>,
}));

vi.mock("./test-suite-accordion", () => ({
  TestSuiteAccordion: ({ suite }: { suite: TestSuite }) => (
    <div data-testid={`suite-${suite.suiteName}`}>
      <span>{suite.suiteName}</span>
      <span>{suite.tests.length} tests</span>
    </div>
  ),
}));

// Mock useWindowVirtualizer to avoid jsdom limitations (scrollTo, measureElement loops)
const mockGetVirtualItems = vi.fn();
const mockGetTotalSize = vi.fn();
const mockMeasureElement = vi.fn();

vi.mock("@tanstack/react-virtual", () => ({
  useWindowVirtualizer: ({ count }: { count: number; estimateSize: () => number }) => {
    // Simulate rendering the first few items (overscan) like the real virtualizer
    const itemCount = Math.min(count, 10);
    const items = Array.from({ length: itemCount }, (_, i) => ({
      index: i,
      key: i,
      size: 72,
      start: i * 72,
    }));

    mockGetVirtualItems.mockReturnValue(items);
    mockGetTotalSize.mockReturnValue(count * 72);

    return {
      getVirtualItems: () => items,
      getTotalSize: () => count * 72,
      measureElement: mockMeasureElement,
    };
  },
}));

const createSuites = (count: number): TestSuite[] =>
  Array.from({ length: count }, (_, i) => ({
    filePath: `src/test-${i}.test.ts`,
    framework: "vitest" as const,
    suiteName: `Suite ${i}`,
    tests: [
      {
        filePath: `src/test-${i}.test.ts`,
        framework: "vitest" as const,
        name: `test ${i}`,
        status: "active" as const,
        suiteName: `Suite ${i}`,
      },
    ],
  }));

describe("TestList", () => {
  it("renders virtual scroll container for large test lists (1000 suites)", () => {
    const suites = createSuites(1000);
    const { container } = render(<TestList suites={suites} />);

    // Virtual scroll container should have relative positioning and calculated total height
    const scrollContainer = container.firstElementChild as HTMLElement;
    expect(scrollContainer).toBeInTheDocument();
    expect(scrollContainer.style.position).toBe("relative");

    // Total height should reflect all 1000 items (each 72px)
    const totalHeight = parseInt(scrollContainer.style.height, 10);
    expect(totalHeight).toBe(1000 * 72);

    // Only a subset of items should be rendered (virtual window)
    const renderedItems = scrollContainer.querySelectorAll("[data-index]");
    expect(renderedItems.length).toBeLessThan(1000);
    expect(renderedItems.length).toBeGreaterThan(0);

    // Verify actual suite content renders for the visible items
    expect(screen.getByTestId("suite-Suite 0")).toBeInTheDocument();
  });

  it("computes accurate total size for large datasets", () => {
    const suites = createSuites(1000);
    const { container } = render(<TestList suites={suites} />);

    const scrollContainer = container.firstElementChild as HTMLElement;
    // 1000 items x 72px per item = 72000px total
    expect(scrollContainer.style.height).toBe("72000px");
  });
});
