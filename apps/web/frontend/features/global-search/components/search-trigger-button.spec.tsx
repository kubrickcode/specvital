import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi, beforeEach } from "vitest";

import messages from "@/i18n/messages/en.json";

import { SearchTriggerButton } from "./search-trigger-button";

const mockOpen = vi.fn();

vi.mock("../hooks", () => ({
  useGlobalSearchStore: () => ({
    open: mockOpen,
  }),
}));

vi.mock("@/components/ui/responsive-tooltip", () => ({
  ResponsiveTooltip: ({
    children,
  }: {
    children: React.ReactNode;
    content: React.ReactNode;
    side?: string;
    sideOffset?: number;
  }) => <>{children}</>,
}));

const renderWithProvider = (ui: React.ReactElement) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      {ui}
    </NextIntlClientProvider>
  );

describe("SearchTriggerButton", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Trigger button click opens dialog", () => {
    it("calls open when desktop button is clicked", () => {
      renderWithProvider(<SearchTriggerButton />);

      const buttons = screen.getAllByRole("button");
      fireEvent.click(buttons[0]);

      expect(mockOpen).toHaveBeenCalledTimes(1);
    });

    it("calls open when mobile icon button is clicked", () => {
      renderWithProvider(<SearchTriggerButton />);

      const mobileButton = screen.getByLabelText("Search");
      fireEvent.click(mobileButton);

      expect(mockOpen).toHaveBeenCalledTimes(1);
    });
  });

  describe("Desktop shortcut badge", () => {
    it("renders keyboard shortcut badge with K suffix", () => {
      renderWithProvider(<SearchTriggerButton />);

      const kbd = document.querySelector("kbd");
      expect(kbd).toBeInTheDocument();
      expect(kbd?.textContent).toMatch(/K$/);
    });

    it("displays Ctrl prefix by default on non-Mac environments", () => {
      renderWithProvider(<SearchTriggerButton />);

      const kbd = document.querySelector("kbd");
      expect(kbd?.textContent).toContain("Ctrl");
    });
  });
});
