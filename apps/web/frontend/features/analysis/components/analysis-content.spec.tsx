import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi, beforeEach } from "vitest";

import type { AnalysisResult } from "@/lib/api/types";

import { AnalysisContent } from "./analysis-content";

let mockTab = "tests";
const mockSetTab = vi.fn((newTab: string) => {
  mockTab = newTab;
});

vi.mock("../hooks/use-primary-tab", () => ({
  usePrimaryTab: () => ({
    isPending: false,
    setTab: mockSetTab,
    tab: mockTab,
  }),
}));

vi.mock("../hooks/use-commit-select", () => ({
  useCommitSelect: () => ({
    commitSha: null,
    selectCommit: vi.fn(),
  }),
}));

vi.mock("@/lib/motion", () => ({
  createStaggerContainer: () => ({}),
  fadeInUp: {},
  useReducedMotion: () => true,
}));

vi.mock("./analysis-header", () => ({
  AnalysisHeader: () => <div data-testid="analysis-header">Header</div>,
}));

vi.mock("./inline-stats", () => ({
  InlineStats: () => <div data-testid="inline-stats">Stats</div>,
}));

vi.mock("./update-banner", () => ({
  UpdateBanner: () => null,
}));

vi.mock("./tests-panel", () => ({
  TestsPanel: ({ availableFrameworks }: { availableFrameworks: string[] }) => (
    <div data-testid="tests-panel">
      Tests Panel Content
      <div data-testid="available-frameworks">{availableFrameworks.join(",")}</div>
    </div>
  ),
}));

vi.mock("./spec-panel", () => ({
  SpecPanel: () => <div data-testid="spec-panel">Spec Panel Content</div>,
}));

const messages = {
  analyze: {
    tabs: {
      spec: "AI Spec",
      tests: "Tests",
    },
  },
};

const mockResult: AnalysisResult = {
  analyzedAt: "2024-06-01T12:00:00Z",
  branchName: "main",
  commitSha: "abc123def456",
  committedAt: "2024-05-31T10:00:00Z",
  id: "test-id-001",
  owner: "facebook",
  parserVersion: "v1.5.1",
  repo: "react",
  suites: [],
  summary: {
    active: 100,
    focused: 0,
    frameworks: [
      { active: 100, focused: 0, framework: "jest", skipped: 5, todo: 0, total: 105, xfail: 0 },
    ],
    skipped: 5,
    todo: 0,
    total: 105,
    xfail: 0,
  },
};

const renderAnalysisContent = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <AnalysisContent result={mockResult} />
    </NextIntlClientProvider>
  );
};

describe("AnalysisContent", () => {
  beforeEach(() => {
    mockTab = "tests";
    mockSetTab.mockClear();
  });

  it("selects Tests tab by default", () => {
    renderAnalysisContent();

    const testsTab = screen.getByRole("tab", { name: /Tests/i });
    expect(testsTab).toHaveAttribute("aria-selected", "true");

    const specTab = screen.getByRole("tab", { name: /AI Spec/i });
    expect(specTab).toHaveAttribute("aria-selected", "false");
  });

  it("shows tests panel when Tests tab is active", () => {
    renderAnalysisContent();

    expect(screen.getByTestId("tests-panel")).toBeInTheDocument();
    expect(screen.queryByTestId("spec-panel")).not.toBeInTheDocument();
  });

  it("calls setTab with 'spec' when AI Spec tab is clicked", () => {
    renderAnalysisContent();

    fireEvent.click(screen.getByRole("tab", { name: /AI Spec/i }));

    expect(mockSetTab).toHaveBeenCalledWith("spec");
  });

  it("shows spec panel when spec tab is active", () => {
    mockTab = "spec";

    renderAnalysisContent();

    expect(screen.getByTestId("spec-panel")).toBeInTheDocument();
    expect(screen.queryByTestId("tests-panel")).not.toBeInTheDocument();
  });

  it("renders tab navigation with tablist role", () => {
    renderAnalysisContent();

    expect(screen.getByRole("tablist")).toBeInTheDocument();
    expect(screen.getAllByRole("tab")).toHaveLength(2);
  });

  it("passes available frameworks to Tests panel for list/tree toggle", () => {
    renderAnalysisContent();

    // TestsPanel receives availableFrameworks which enables the DataViewToggle (list/tree)
    expect(screen.getByTestId("available-frameworks")).toHaveTextContent("jest");
  });
});
