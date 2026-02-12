import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import { QuotaIndicator } from "./quota-indicator";

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => {
    const messages: Record<string, string> = {
      almostFull: "Almost at limit",
      processing: "processing",
      thisMonth: "this month",
      unlimited: "(No limit)",
      viewAccount: "View Account",
      warning: "Running low",
    };
    return messages[key] || key;
  },
}));

vi.mock("@/i18n/navigation", () => ({
  Link: ({
    children,
    href,
    ...props
  }: {
    children: React.ReactNode;
    className?: string;
    href: string;
  }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}));

describe("QuotaIndicator", () => {
  afterEach(() => {
    cleanup();
  });

  describe("normal level (percentage < 70)", () => {
    it("displays usage count and limit", () => {
      render(<QuotaIndicator limit={100} percentage={30} reserved={0} used={30} />);

      expect(screen.getByText(/30/)).toBeInTheDocument();
      expect(screen.getByText(/100/)).toBeInTheDocument();
      expect(screen.getByText(/this month/)).toBeInTheDocument();
    });

    it("does not show warning or danger indicators", () => {
      render(<QuotaIndicator limit={100} percentage={30} reserved={0} used={30} />);

      expect(screen.queryByText("Running low")).not.toBeInTheDocument();
      expect(screen.queryByText("Almost at limit")).not.toBeInTheDocument();
      expect(screen.queryByText("View Account")).not.toBeInTheDocument();
    });
  });

  describe("warning level (70 <= percentage < 90)", () => {
    it("shows running low indicator", () => {
      render(<QuotaIndicator limit={100} percentage={75} reserved={0} used={75} />);

      expect(screen.getByText("Running low")).toBeInTheDocument();
    });

    it("does not show view account link", () => {
      render(<QuotaIndicator limit={100} percentage={75} reserved={0} used={75} />);

      expect(screen.queryByText("View Account")).not.toBeInTheDocument();
    });
  });

  describe("danger level (90 <= percentage < 100)", () => {
    it("shows almost at limit indicator", () => {
      render(<QuotaIndicator limit={100} percentage={95} reserved={0} used={95} />);

      expect(screen.getByText("Almost at limit")).toBeInTheDocument();
    });

    it("does not show view account link when not yet exceeded", () => {
      render(<QuotaIndicator limit={100} percentage={95} reserved={0} used={95} />);

      expect(screen.queryByText("View Account")).not.toBeInTheDocument();
    });
  });

  describe("exceeded (percentage >= 100)", () => {
    it("shows view account link when quota is exactly at limit", () => {
      render(<QuotaIndicator limit={100} percentage={100} reserved={0} used={100} />);

      const link = screen.getByText("View Account");
      expect(link).toBeInTheDocument();
      expect(link.closest("a")).toHaveAttribute("href", "/account");
    });

    it("shows view account link when quota exceeds limit", () => {
      render(<QuotaIndicator limit={100} percentage={120} reserved={0} used={120} />);

      const link = screen.getByText("View Account");
      expect(link).toBeInTheDocument();
      expect(link.closest("a")).toHaveAttribute("href", "/account");
    });

    it("does not show almost at limit text when exceeded", () => {
      render(<QuotaIndicator limit={100} percentage={100} reserved={0} used={100} />);

      expect(screen.queryByText("Almost at limit")).not.toBeInTheDocument();
    });
  });

  describe("unlimited plan (percentage null)", () => {
    it("displays used count with no limit indicator", () => {
      render(<QuotaIndicator limit={null} percentage={null} reserved={0} used={42} />);

      expect(screen.getByText(/42/)).toBeInTheDocument();
      expect(screen.getByText(/\(No limit\)/)).toBeInTheDocument();
    });

    it("does not show warning or view account", () => {
      render(<QuotaIndicator limit={null} percentage={null} reserved={0} used={500} />);

      expect(screen.queryByText("Running low")).not.toBeInTheDocument();
      expect(screen.queryByText("View Account")).not.toBeInTheDocument();
    });
  });

  describe("reserved (processing) display", () => {
    it("shows reserved count for limited plans", () => {
      render(<QuotaIndicator limit={100} percentage={50} reserved={5} used={50} />);

      expect(screen.getByText(/processing/)).toBeInTheDocument();
      expect(screen.getByText(/· 5 processing/)).toBeInTheDocument();
    });

    it("shows reserved count for unlimited plans", () => {
      render(<QuotaIndicator limit={null} percentage={null} reserved={3} used={42} />);

      expect(screen.getByText(/processing/)).toBeInTheDocument();
      expect(screen.getByText(/· 3 processing/)).toBeInTheDocument();
    });

    it("does not show processing text when reserved is zero", () => {
      render(<QuotaIndicator limit={100} percentage={50} reserved={0} used={50} />);

      expect(screen.queryByText("processing")).not.toBeInTheDocument();
    });
  });
});
