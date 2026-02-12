import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { components } from "@/lib/api/generated-types";

import { PlanSection } from "./plan-section";

type PlanInfo = components["schemas"]["PlanInfo"];
type PlanSectionProps = {
  error?: Error | null;
  isLoading: boolean;
  plan?: PlanInfo;
};

vi.mock("@/i18n/navigation", () => ({
  Link: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}));

const messages = {
  account: {
    plan: {
      analysisLimit: "Analysis",
      contactUs: "Contact Us",
      days: "{days} days",
      month: "month",
      retention: "Data Retention",
      specviewLimit: "AI Spec Docs",
      tier: {
        enterprise: "Enterprise",
        free: "Free",
        pro: "Pro",
        pro_plus: "Pro+",
      },
      title: "Current Plan",
      unavailable: "Plan information unavailable",
      unlimited: "Unlimited",
      upgrade: "Upgrade Plan",
    },
  },
};

const createMockPlan = (overrides?: Partial<PlanInfo>): PlanInfo => ({
  analysisMonthlyLimit: 50,
  retentionDays: 30,
  specviewMonthlyLimit: 100,
  tier: "free",
  ...overrides,
});

const renderPlanSection = (props: Partial<PlanSectionProps> = {}) => {
  const defaultProps: PlanSectionProps = {
    isLoading: false,
    plan: createMockPlan(),
    ...props,
  };

  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <PlanSection {...defaultProps} />
    </NextIntlClientProvider>
  );
};

describe("PlanSection", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("loading state", () => {
    it("renders title but no plan details", () => {
      renderPlanSection({ isLoading: true, plan: undefined });

      expect(screen.getByText("Current Plan")).toBeInTheDocument();
      expect(screen.queryByText("Free")).not.toBeInTheDocument();
      expect(screen.queryByText("AI Spec Docs")).not.toBeInTheDocument();
    });
  });

  describe("error state", () => {
    it("renders unavailable message when error is present", () => {
      renderPlanSection({ error: new Error("fetch failed"), plan: undefined });

      expect(screen.getByText("Plan information unavailable")).toBeInTheDocument();
    });

    it("renders unavailable message when plan is missing", () => {
      renderPlanSection({ plan: undefined });

      expect(screen.getByText("Plan information unavailable")).toBeInTheDocument();
    });
  });

  describe("plan display", () => {
    it("renders tier badge with plan name", () => {
      renderPlanSection({ plan: createMockPlan({ tier: "free" }) });

      expect(screen.getByText("Free")).toBeInTheDocument();
    });

    it("renders specview monthly limit", () => {
      renderPlanSection({ plan: createMockPlan({ specviewMonthlyLimit: 100 }) });

      expect(screen.getByText("AI Spec Docs")).toBeInTheDocument();
      expect(screen.getByText("100/month")).toBeInTheDocument();
    });

    it("renders analysis monthly limit", () => {
      renderPlanSection({ plan: createMockPlan({ analysisMonthlyLimit: 50 }) });

      expect(screen.getByText("Analysis")).toBeInTheDocument();
      expect(screen.getByText("50/month")).toBeInTheDocument();
    });

    it("renders retention days", () => {
      renderPlanSection({ plan: createMockPlan({ retentionDays: 30 }) });

      expect(screen.getByText("Data Retention")).toBeInTheDocument();
      expect(screen.getByText("30 days")).toBeInTheDocument();
    });

    it("renders unlimited retention when retentionDays is null", () => {
      renderPlanSection({ plan: createMockPlan({ retentionDays: null }) });

      expect(screen.getByText("Unlimited")).toBeInTheDocument();
    });

    it("renders unlimited retention when retentionDays is undefined", () => {
      renderPlanSection({ plan: createMockPlan({ retentionDays: undefined }) });

      expect(screen.getByText("Unlimited")).toBeInTheDocument();
    });
  });

  describe("tier-specific footer", () => {
    it("shows upgrade link for free tier", () => {
      renderPlanSection({ plan: createMockPlan({ tier: "free" }) });

      const upgradeLink = screen.getByRole("link", { name: /upgrade plan/i });
      expect(upgradeLink).toBeInTheDocument();
      expect(upgradeLink).toHaveAttribute("href", "/pricing");
    });

    it("shows upgrade link for pro tier", () => {
      renderPlanSection({ plan: createMockPlan({ tier: "pro" }) });

      expect(screen.getByRole("link", { name: /upgrade plan/i })).toBeInTheDocument();
    });

    it("shows upgrade link for pro_plus tier", () => {
      renderPlanSection({ plan: createMockPlan({ tier: "pro_plus" }) });

      expect(screen.getByRole("link", { name: /upgrade plan/i })).toBeInTheDocument();
    });

    it("shows contact us mailto link for enterprise tier", () => {
      renderPlanSection({ plan: createMockPlan({ tier: "enterprise" }) });

      const contactLink = screen.getByRole("link", { name: /contact us/i });
      expect(contactLink).toBeInTheDocument();
      expect(contactLink).toHaveAttribute("href", "mailto:support@specvital.com");
      expect(screen.queryByRole("link", { name: /upgrade plan/i })).not.toBeInTheDocument();
    });
  });

  describe("tier badge variants", () => {
    it("renders pro tier badge", () => {
      renderPlanSection({ plan: createMockPlan({ tier: "pro" }) });

      expect(screen.getByText("Pro")).toBeInTheDocument();
    });

    it("renders pro_plus tier badge", () => {
      renderPlanSection({ plan: createMockPlan({ tier: "pro_plus" }) });

      expect(screen.getByText("Pro+")).toBeInTheDocument();
    });

    it("renders enterprise tier badge", () => {
      renderPlanSection({ plan: createMockPlan({ tier: "enterprise" }) });

      expect(screen.getByText("Enterprise")).toBeInTheDocument();
    });
  });
});
