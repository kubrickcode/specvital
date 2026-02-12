import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it } from "vitest";

import messages from "@/i18n/messages/en.json";

import { SearchKeyboardHints } from "./search-keyboard-hints";

const renderWithProvider = (ui: React.ReactElement) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      {ui}
    </NextIntlClientProvider>
  );

describe("SearchKeyboardHints", () => {
  describe("Desktop: keyboard hints visible", () => {
    it("renders all keyboard hint labels", () => {
      renderWithProvider(<SearchKeyboardHints />);

      expect(screen.getByText("Navigate")).toBeInTheDocument();
      expect(screen.getByText("Select")).toBeInTheDocument();
      expect(screen.getByText("Close")).toBeInTheDocument();
    });

    it("renders arrow key and enter key indicators", () => {
      renderWithProvider(<SearchKeyboardHints />);

      expect(screen.getByText("\u2191")).toBeInTheDocument();
      expect(screen.getByText("\u2193")).toBeInTheDocument();
      expect(screen.getByText("\u21B5")).toBeInTheDocument();
    });

    it("renders Esc key indicator", () => {
      renderWithProvider(<SearchKeyboardHints />);

      expect(screen.getByText("Esc")).toBeInTheDocument();
    });
  });

  describe("Mobile: keyboard hints hidden via CSS", () => {
    it("has hidden class that hides on mobile and shows on md breakpoint", () => {
      renderWithProvider(<SearchKeyboardHints />);

      const hintsContainer = screen.getByText("Navigate").closest("div.hidden.md\\:flex");
      expect(hintsContainer).toBeInTheDocument();
      expect(hintsContainer).toHaveClass("hidden");
      expect(hintsContainer).toHaveClass("md:flex");
    });
  });
});
