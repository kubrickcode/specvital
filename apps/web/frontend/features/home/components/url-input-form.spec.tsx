import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi, beforeEach } from "vitest";

import messages from "@/i18n/messages/en.json";

import { UrlInputForm } from "./url-input-form";

const mockPush = vi.fn();

vi.mock("@/i18n/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

vi.mock("@/lib/hooks", () => ({
  useMediaQuery: () => true,
}));

const renderWithProvider = (ui: React.ReactElement) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      {ui}
    </NextIntlClientProvider>
  );

describe("UrlInputForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders heading label, text input, and submit button", () => {
    renderWithProvider(<UrlInputForm />);

    expect(screen.getByLabelText("GitHub Repository URL")).toBeInTheDocument();
    expect(screen.getByRole("textbox")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Analyze" })).toBeInTheDocument();
  });

  it("renders placeholder text on desktop", () => {
    renderWithProvider(<UrlInputForm />);

    expect(screen.getByPlaceholderText("https://github.com/facebook/react")).toBeInTheDocument();
  });

  it("displays error message when submitting an invalid URL", () => {
    renderWithProvider(<UrlInputForm />);

    const input = screen.getByRole("textbox");
    const button = screen.getByRole("button", { name: "Analyze" });

    fireEvent.change(input, { target: { value: "not-a-valid-url" } });
    fireEvent.click(button);

    expect(screen.getByRole("alert")).toBeInTheDocument();
  });

  it("clears error when user modifies input after submission failure", () => {
    renderWithProvider(<UrlInputForm />);

    const input = screen.getByRole("textbox");
    const button = screen.getByRole("button", { name: "Analyze" });

    fireEvent.change(input, { target: { value: "bad" } });
    fireEvent.click(button);

    expect(screen.getByRole("alert")).toBeInTheDocument();

    fireEvent.change(input, { target: { value: "bad2" } });

    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });

  it("navigates to analyze page on valid URL submit", () => {
    renderWithProvider(<UrlInputForm />);

    const input = screen.getByRole("textbox");
    const button = screen.getByRole("button", { name: "Analyze" });

    fireEvent.change(input, {
      target: { value: "https://github.com/facebook/react" },
    });
    fireEvent.click(button);

    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });

  it("calls onSuccess callback when submit succeeds", () => {
    const onSuccess = vi.fn();
    renderWithProvider(<UrlInputForm onSuccess={onSuccess} />);

    const input = screen.getByRole("textbox");
    const button = screen.getByRole("button", { name: "Analyze" });

    fireEvent.change(input, {
      target: { value: "https://github.com/facebook/react" },
    });
    fireEvent.click(button);

    expect(onSuccess).toHaveBeenCalledOnce();
  });

  it("accepts shorthand owner/repo format", () => {
    const onSuccess = vi.fn();
    renderWithProvider(<UrlInputForm onSuccess={onSuccess} />);

    const input = screen.getByRole("textbox");
    const button = screen.getByRole("button", { name: "Analyze" });

    fireEvent.change(input, { target: { value: "facebook/react" } });
    fireEvent.click(button);

    expect(onSuccess).toHaveBeenCalledOnce();
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });

  it("renders help button for supported formats", () => {
    renderWithProvider(<UrlInputForm />);

    expect(screen.getByLabelText("View supported formats")).toBeInTheDocument();
  });
});
