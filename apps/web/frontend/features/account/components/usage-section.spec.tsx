import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { UsageStatusResponse } from "../api";

import { UsageSection } from "./usage-section";

vi.mock("./usage-progress", () => ({
  UsageProgress: (props: {
    label: string;
    limit: number | null;
    percentage: number | null;
    reserved: number;
    unit: string;
    used: number;
  }) => (
    <div data-testid={`usage-progress-${props.label}`}>
      <span>{props.label}</span>
      <span>
        {props.used}/{props.limit ?? "unlimited"}
      </span>
      <span>{props.percentage != null ? `${props.percentage}%` : "unlimited"}</span>
      {props.reserved > 0 && <span>{props.reserved} reserved</span>}
    </div>
  ),
}));

type UsageSectionProps = {
  error?: Error | null;
  isLoading: boolean;
  usage?: UsageStatusResponse;
};

const messages = {
  account: {
    usage: {
      analysis: "Analysis Runs",
      resetIn: "Resets in {days} days ({date})",
      specview: "AI Spec Docs",
      title: "Usage This Period",
      unavailable: "Usage information unavailable",
    },
  },
  pricing: {
    features: {
      analysis: {
        unit: "runs",
      },
      specview: {
        unit: "tests",
      },
    },
  },
};

const createMockUsage = (overrides?: Partial<UsageStatusResponse>): UsageStatusResponse => ({
  analysis: {
    limit: 50,
    percentage: 20,
    reserved: 0,
    used: 10,
  },
  plan: {
    analysisMonthlyLimit: 50,
    retentionDays: 30,
    specviewMonthlyLimit: 100,
    tier: "free",
  },
  resetAt: "2026-03-01T00:00:00Z",
  specview: {
    limit: 100,
    percentage: 45,
    reserved: 0,
    used: 45,
  },
  ...overrides,
});

const renderUsageSection = (props: Partial<UsageSectionProps> = {}) => {
  const defaultProps: UsageSectionProps = {
    isLoading: false,
    usage: createMockUsage(),
    ...props,
  };

  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <UsageSection {...defaultProps} />
    </NextIntlClientProvider>
  );
};

describe("UsageSection", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("loading state", () => {
    it("renders title but no usage details", () => {
      renderUsageSection({ isLoading: true, usage: undefined });

      expect(screen.getByText("Usage This Period")).toBeInTheDocument();
      expect(screen.queryByText("AI Spec Docs")).not.toBeInTheDocument();
      expect(screen.queryByText("Analysis Runs")).not.toBeInTheDocument();
    });
  });

  describe("error state", () => {
    it("renders unavailable message when error is present", () => {
      renderUsageSection({ error: new Error("network failure"), usage: undefined });

      expect(screen.getByText("Usage information unavailable")).toBeInTheDocument();
    });

    it("renders unavailable message when usage is missing", () => {
      renderUsageSection({ usage: undefined });

      expect(screen.getByText("Usage information unavailable")).toBeInTheDocument();
    });
  });

  describe("usage display", () => {
    it("renders specview usage progress", () => {
      renderUsageSection({
        usage: createMockUsage({
          specview: { limit: 100, percentage: 45, reserved: 0, used: 45 },
        }),
      });

      const specviewProgress = screen.getByTestId("usage-progress-AI Spec Docs");
      expect(specviewProgress).toBeInTheDocument();
      expect(screen.getByText("45/100")).toBeInTheDocument();
    });

    it("renders analysis usage progress", () => {
      renderUsageSection({
        usage: createMockUsage({
          analysis: { limit: 50, percentage: 20, reserved: 0, used: 10 },
        }),
      });

      const analysisProgress = screen.getByTestId("usage-progress-Analysis Runs");
      expect(analysisProgress).toBeInTheDocument();
      expect(screen.getByText("10/50")).toBeInTheDocument();
    });

    it("passes reserved count to usage progress", () => {
      renderUsageSection({
        usage: createMockUsage({
          specview: { limit: 100, percentage: 47, reserved: 2, used: 45 },
        }),
      });

      expect(screen.getByText("2 reserved")).toBeInTheDocument();
    });
  });

  describe("reset date", () => {
    it("renders reset info with days and date", () => {
      renderUsageSection({
        usage: createMockUsage({ resetAt: "2026-03-01T00:00:00Z" }),
      });

      // getResetInfo calculates days from now and formats date
      // The exact text depends on the current date, so check for the pattern
      expect(screen.getByText(/Resets in/)).toBeInTheDocument();
    });
  });

  describe("unlimited usage", () => {
    it("passes null limit for unlimited metrics", () => {
      renderUsageSection({
        usage: createMockUsage({
          analysis: { limit: null, percentage: null, reserved: 0, used: 10 },
          specview: { limit: null, percentage: null, reserved: 0, used: 45 },
        }),
      });

      expect(screen.getByText("45/unlimited")).toBeInTheDocument();
      expect(screen.getByText("10/unlimited")).toBeInTheDocument();
    });
  });
});
