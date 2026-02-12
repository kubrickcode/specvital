import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { AccountContent } from "./account-content";

const mockUseSubscription = vi.fn();
const mockUseUsage = vi.fn();

vi.mock("../hooks", () => ({
  useSubscription: (...args: unknown[]) => mockUseSubscription(...args),
  useUsage: (...args: unknown[]) => mockUseUsage(...args),
}));

vi.mock("@/features/auth", () => ({
  RequireAuth: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

vi.mock("./plan-section", () => ({
  PlanSection: (props: { error?: Error | null; isLoading: boolean }) => (
    <div data-testid="plan-section" data-loading={props.isLoading} data-error={!!props.error}>
      PlanSection
    </div>
  ),
}));

vi.mock("./usage-section", () => ({
  UsageSection: (props: { error?: Error | null; isLoading: boolean }) => (
    <div data-testid="usage-section" data-loading={props.isLoading} data-error={!!props.error}>
      UsageSection
    </div>
  ),
}));

const messages = {
  account: {
    plan: { title: "Current Plan" },
    usage: { title: "Usage This Period" },
  },
};

const renderAccountContent = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <AccountContent />
    </NextIntlClientProvider>
  );
};

describe("AccountContent", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders both PlanSection and UsageSection", () => {
    mockUseSubscription.mockReturnValue({
      data: { plan: { tier: "free" } },
      error: null,
      isLoading: false,
    });
    mockUseUsage.mockReturnValue({
      data: { analysis: {}, specview: {}, resetAt: "2026-03-01T00:00:00Z" },
      error: null,
      isLoading: false,
    });

    renderAccountContent();

    expect(screen.getByTestId("plan-section")).toBeInTheDocument();
    expect(screen.getByTestId("usage-section")).toBeInTheDocument();
  });

  it("passes isLoading=true when subscription is loading", () => {
    mockUseSubscription.mockReturnValue({
      data: undefined,
      error: null,
      isLoading: true,
    });
    mockUseUsage.mockReturnValue({
      data: undefined,
      error: null,
      isLoading: false,
    });

    renderAccountContent();

    expect(screen.getByTestId("plan-section")).toHaveAttribute("data-loading", "true");
    expect(screen.getByTestId("usage-section")).toHaveAttribute("data-loading", "true");
  });

  it("passes isLoading=true when usage is loading", () => {
    mockUseSubscription.mockReturnValue({
      data: undefined,
      error: null,
      isLoading: false,
    });
    mockUseUsage.mockReturnValue({
      data: undefined,
      error: null,
      isLoading: true,
    });

    renderAccountContent();

    expect(screen.getByTestId("plan-section")).toHaveAttribute("data-loading", "true");
    expect(screen.getByTestId("usage-section")).toHaveAttribute("data-loading", "true");
  });

  it("passes error from subscription to PlanSection", () => {
    mockUseSubscription.mockReturnValue({
      data: undefined,
      error: new Error("subscription failed"),
      isLoading: false,
    });
    mockUseUsage.mockReturnValue({
      data: undefined,
      error: null,
      isLoading: false,
    });

    renderAccountContent();

    expect(screen.getByTestId("plan-section")).toHaveAttribute("data-error", "true");
    expect(screen.getByTestId("usage-section")).toHaveAttribute("data-error", "false");
  });

  it("passes error from usage to UsageSection", () => {
    mockUseSubscription.mockReturnValue({
      data: undefined,
      error: null,
      isLoading: false,
    });
    mockUseUsage.mockReturnValue({
      data: undefined,
      error: new Error("usage failed"),
      isLoading: false,
    });

    renderAccountContent();

    expect(screen.getByTestId("plan-section")).toHaveAttribute("data-error", "false");
    expect(screen.getByTestId("usage-section")).toHaveAttribute("data-error", "true");
  });
});
