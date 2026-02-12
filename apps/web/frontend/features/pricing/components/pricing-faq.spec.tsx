import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi, beforeEach } from "vitest";

import messages from "@/i18n/messages/en.json";

import { PricingFaq } from "./pricing-faq";

vi.mock("motion/react", () => ({
  AnimatePresence: ({ children }: React.PropsWithChildren) => <>{children}</>,
  motion: {
    div: ({ children, ...props }: React.PropsWithChildren) => <div {...props}>{children}</div>,
  },
}));

const renderWithProvider = (ui: React.ReactElement) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      {ui}
    </NextIntlClientProvider>
  );

const getButtonByText = (text: string) => {
  const span = screen.getByText(text);
  return span.closest("button")!;
};

describe("PricingFaq", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders FAQ heading", () => {
    renderWithProvider(<PricingFaq />);

    expect(screen.getByText("Frequently Asked Questions")).toBeInTheDocument();
  });

  it("renders all 4 FAQ items as collapsed buttons", () => {
    renderWithProvider(<PricingFaq />);

    const buttons = screen.getAllByRole("button");
    expect(buttons).toHaveLength(4);

    expect(screen.getByText("What is an AI Spec Document?")).toBeInTheDocument();
    expect(screen.getByText("What is Analysis?")).toBeInTheDocument();
    expect(screen.getByText("When will payments go live?")).toBeInTheDocument();
    expect(screen.getByText("What is data retention?")).toBeInTheDocument();
  });

  it("all items start collapsed with aria-expanded=false", () => {
    renderWithProvider(<PricingFaq />);

    const buttons = screen.getAllByRole("button");
    buttons.forEach((button) => {
      expect(button).toHaveAttribute("aria-expanded", "false");
    });
  });

  it("expands item and shows answer text on click", () => {
    renderWithProvider(<PricingFaq />);

    const button = getButtonByText("What is an AI Spec Document?");
    fireEvent.click(button);

    expect(button).toHaveAttribute("aria-expanded", "true");
    expect(screen.getByText(/automatically organizes your test cases/)).toBeInTheDocument();
  });

  it("collapses expanded item on second click", () => {
    renderWithProvider(<PricingFaq />);

    const button = getButtonByText("What is an AI Spec Document?");
    fireEvent.click(button);

    expect(button).toHaveAttribute("aria-expanded", "true");

    fireEvent.click(button);

    expect(button).toHaveAttribute("aria-expanded", "false");
  });

  it("closes previously open item when another item is clicked", () => {
    renderWithProvider(<PricingFaq />);

    const firstButton = getButtonByText("What is an AI Spec Document?");
    const secondButton = getButtonByText("What is Analysis?");

    fireEvent.click(firstButton);
    expect(firstButton).toHaveAttribute("aria-expanded", "true");

    fireEvent.click(secondButton);
    expect(firstButton).toHaveAttribute("aria-expanded", "false");
    expect(secondButton).toHaveAttribute("aria-expanded", "true");
    expect(screen.getByText(/parses your repository's test files/)).toBeInTheDocument();
  });
});
