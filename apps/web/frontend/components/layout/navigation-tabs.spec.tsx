import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { NavigationTabs } from "./navigation-tabs";

const mockUseAuth = vi.fn();
const mockUsePathname = vi.fn();

vi.mock("@/features/auth", () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock("@/i18n/navigation", () => ({
  Link: ({
    children,
    href,
    ...props
  }: {
    children: React.ReactNode;
    href: string;
    [key: string]: unknown;
  }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
  usePathname: () => mockUsePathname(),
}));

const messages = {
  navigation: {
    ariaLabel: "Main navigation",
    dashboard: "Dashboard",
    docs: "Docs",
    explore: "Explore",
    pricing: "Pricing",
  },
};

const renderNavigationTabs = () =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <NavigationTabs />
    </NextIntlClientProvider>
  );

describe("NavigationTabs", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUsePathname.mockReturnValue("/");
  });

  describe("Desktop Viewport", () => {
    it("renders main navigation with accessible label", () => {
      mockUseAuth.mockReturnValue({ isAuthenticated: false });

      renderNavigationTabs();

      const nav = screen.getByRole("navigation", { name: "Main navigation" });
      expect(nav).toBeInTheDocument();
    });

    it("displays Explore, Docs, and Pricing links for unauthenticated users", () => {
      mockUseAuth.mockReturnValue({ isAuthenticated: false });

      renderNavigationTabs();

      expect(screen.getByText("Explore")).toBeInTheDocument();
      expect(screen.getByText("Docs")).toBeInTheDocument();
      expect(screen.getByText("Pricing")).toBeInTheDocument();
      expect(screen.queryByText("Dashboard")).not.toBeInTheDocument();
    });

    it("displays Dashboard link along with others for authenticated users", () => {
      mockUseAuth.mockReturnValue({ isAuthenticated: true });

      renderNavigationTabs();

      expect(screen.getByText("Dashboard")).toBeInTheDocument();
      expect(screen.getByText("Explore")).toBeInTheDocument();
      expect(screen.getByText("Docs")).toBeInTheDocument();
      expect(screen.getByText("Pricing")).toBeInTheDocument();
    });

    it("marks the active route with aria-current=page", () => {
      mockUseAuth.mockReturnValue({ isAuthenticated: false });
      mockUsePathname.mockReturnValue("/explore");

      renderNavigationTabs();

      const exploreLink = screen.getByText("Explore");
      expect(exploreLink).toHaveAttribute("aria-current", "page");

      const pricingLink = screen.getByText("Pricing");
      expect(pricingLink).not.toHaveAttribute("aria-current");
    });
  });
});
