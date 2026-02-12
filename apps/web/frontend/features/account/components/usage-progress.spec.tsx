import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { UsageProgress } from "./usage-progress";

type UsageProgressProps = {
  label: string;
  limit: number | null;
  percentage: number | null;
  reserved: number;
  unit: string;
  used: number;
};

const messages = {
  specView: {
    quota: {
      processing: "processing",
    },
  },
};

const createDefaultProps = (overrides?: Partial<UsageProgressProps>): UsageProgressProps => ({
  label: "AI Spec Docs",
  limit: 100,
  percentage: 45,
  reserved: 0,
  unit: "tests",
  used: 45,
  ...overrides,
});

const renderUsageProgress = (overrides?: Partial<UsageProgressProps>) => {
  const props = createDefaultProps(overrides);

  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <UsageProgress {...props} />
    </NextIntlClientProvider>
  );
};

describe("UsageProgress", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("label and usage text", () => {
    it("renders label", () => {
      renderUsageProgress({ label: "Analysis Runs" });

      expect(screen.getByText("Analysis Runs")).toBeInTheDocument();
    });

    it("renders used / limit with unit", () => {
      renderUsageProgress({ limit: 200, unit: "runs", used: 75 });

      expect(screen.getByText("75 / 200 runs")).toBeInTheDocument();
    });
  });

  describe("percentage display", () => {
    it("shows rounded percentage for limited plans", () => {
      renderUsageProgress({ percentage: 45.7 });

      expect(screen.getByText("46%")).toBeInTheDocument();
    });

    it("does not show percentage for unlimited plans", () => {
      renderUsageProgress({ limit: null, percentage: null });

      expect(screen.queryByText(/%/)).not.toBeInTheDocument();
    });
  });

  describe("unlimited plan", () => {
    it("renders infinity symbol in usage text", () => {
      renderUsageProgress({ limit: null, percentage: null, used: 30 });

      const usageTexts = screen.getAllByText((content) => content.includes("30 /"));
      expect(usageTexts.length).toBeGreaterThan(0);
    });
  });

  describe("high usage warning", () => {
    it("applies destructive styling at 90% usage", () => {
      renderUsageProgress({ percentage: 90 });

      const percentageEl = screen.getByText("90%");
      expect(percentageEl.className).toContain("text-destructive");
    });

    it("does not apply destructive styling below 90%", () => {
      renderUsageProgress({ percentage: 89 });

      const percentageEl = screen.getByText("89%");
      expect(percentageEl.className).not.toContain("text-destructive");
    });
  });

  describe("reserved items", () => {
    it("shows processing text when reserved > 0", () => {
      renderUsageProgress({ reserved: 3 });

      expect(screen.getByText("3 processing")).toBeInTheDocument();
    });

    it("does not show processing text when reserved is 0", () => {
      renderUsageProgress({ reserved: 0 });

      expect(screen.queryByText(/processing/)).not.toBeInTheDocument();
    });
  });
});
