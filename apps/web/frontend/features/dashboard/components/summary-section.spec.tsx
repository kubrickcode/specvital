import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

const mockStatsData = vi.fn();
const mockStatsIsLoading = vi.fn<() => boolean>();

vi.mock("../hooks", () => ({
  useRepositoryStats: () => ({
    data: mockStatsData(),
    isLoading: mockStatsIsLoading(),
  }),
}));

vi.mock("motion/react", () => ({
  motion: {
    div: ({
      children,
      className,
      ...rest
    }: {
      children: React.ReactNode;
      className?: string;
      [key: string]: unknown;
    }) => (
      <div className={className} data-testid="motion-div" {...filterDomProps(rest)}>
        {children}
      </div>
    ),
    span: ({
      children,
      className,
      ...rest
    }: {
      children: React.ReactNode;
      className?: string;
      [key: string]: unknown;
    }) => (
      <span className={className} {...filterDomProps(rest)}>
        {children}
      </span>
    ),
  },
}));

vi.mock("@/lib/motion", () => ({
  easeOutTransition: {},
  staggerContainer: {},
  staggerItem: {},
  useReducedMotion: () => true,
}));

const filterDomProps = (props: Record<string, unknown>) => {
  const domSafe: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(props)) {
    if (
      key === "className" ||
      key === "style" ||
      key.startsWith("data-") ||
      key.startsWith("aria-") ||
      key === "role" ||
      key === "id"
    ) {
      domSafe[key] = value;
    }
  }
  return domSafe;
};

import { SummarySection } from "./summary-section";

const messages = {
  dashboard: {
    summary: {
      activeRepos: "Active Repositories",
      loading: "Loading statistics...",
      totalTests: "Total Tests",
    },
  },
};

const renderSummarySection = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <SummarySection />
    </NextIntlClientProvider>
  );
};

describe("SummarySection", () => {
  beforeEach(() => {
    mockStatsData.mockReturnValue(undefined);
    mockStatsIsLoading.mockReturnValue(false);
  });

  it("renders loading skeletons while data is loading", () => {
    mockStatsIsLoading.mockReturnValue(true);
    renderSummarySection();

    expect(screen.getByRole("status")).toBeInTheDocument();
    expect(screen.getByText("Loading statistics...")).toBeInTheDocument();
  });

  it("renders nothing when data is undefined and not loading", () => {
    mockStatsData.mockReturnValue(undefined);
    mockStatsIsLoading.mockReturnValue(false);
    const { container } = renderSummarySection();

    expect(container.firstChild).toBeNull();
  });

  it("renders total tests and active repositories counts", () => {
    mockStatsData.mockReturnValue({ totalRepositories: 12, totalTests: 543 });
    renderSummarySection();

    expect(screen.getByText("Total Tests")).toBeInTheDocument();
    expect(screen.getByText("Active Repositories")).toBeInTheDocument();
  });
});
