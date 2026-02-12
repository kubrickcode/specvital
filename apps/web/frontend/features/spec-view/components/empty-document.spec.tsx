import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { EmptyDocument } from "./empty-document";

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => {
    const messages: Record<string, string> = {
      description: "Generate an AI-powered specification document from your test cases.",
      generateButton: "Generate Document",
      generating: "Generating...",
      title: "No Specification Document",
    };
    return messages[key] || key;
  },
}));

vi.mock("./quota-indicator", () => ({
  QuotaIndicator: ({
    limit,
    percentage,
    used,
  }: {
    limit: number | null;
    percentage: number | null;
    used: number;
  }) => (
    <div data-testid="quota-indicator">
      {percentage !== null ? `${used}/${limit} (${percentage}%)` : `${used} (unlimited)`}
    </div>
  ),
}));

describe("EmptyDocument", () => {
  let mockOnGenerate: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockOnGenerate = vi.fn();
  });

  afterEach(() => {
    cleanup();
  });

  it("shows generate button when no document exists", () => {
    render(<EmptyDocument onGenerate={mockOnGenerate} />);

    expect(screen.getByText("No Specification Document")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Generate Document/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Generate Document/i })).toBeEnabled();
  });

  it("calls onGenerate when generate button is clicked", () => {
    render(<EmptyDocument onGenerate={mockOnGenerate} />);

    fireEvent.click(screen.getByRole("button", { name: /Generate Document/i }));

    expect(mockOnGenerate).toHaveBeenCalledOnce();
  });

  it("disables generate button when isLoading is true", () => {
    render(<EmptyDocument isLoading onGenerate={mockOnGenerate} />);

    const button = screen.getByRole("button", { name: /Generating.../i });
    expect(button).toBeDisabled();
  });

  it("shows loading text when isLoading is true", () => {
    render(<EmptyDocument isLoading onGenerate={mockOnGenerate} />);

    expect(screen.getByText("Generating...")).toBeInTheDocument();
  });

  it("renders quota indicator when quota is provided", () => {
    const quota = { limit: 100, percentage: 50, reserved: 0, used: 50 };

    render(<EmptyDocument onGenerate={mockOnGenerate} quota={quota} />);

    expect(screen.getByTestId("quota-indicator")).toBeInTheDocument();
    expect(screen.getByTestId("quota-indicator")).toHaveTextContent("50/100 (50%)");
  });

  it("does not render quota indicator when quota is null", () => {
    render(<EmptyDocument onGenerate={mockOnGenerate} quota={null} />);

    expect(screen.queryByTestId("quota-indicator")).not.toBeInTheDocument();
  });

  it("does not render quota indicator when quota is undefined", () => {
    render(<EmptyDocument onGenerate={mockOnGenerate} />);

    expect(screen.queryByTestId("quota-indicator")).not.toBeInTheDocument();
  });

  it("disables generate button when quota is exceeded (100%)", () => {
    const quota = { limit: 100, percentage: 100, reserved: 0, used: 100 };

    render(<EmptyDocument onGenerate={mockOnGenerate} quota={quota} />);

    expect(screen.getByRole("button", { name: /Generate Document/i })).toBeDisabled();
  });

  it("disables generate button when quota is over limit (120%)", () => {
    const quota = { limit: 100, percentage: 120, reserved: 0, used: 120 };

    render(<EmptyDocument onGenerate={mockOnGenerate} quota={quota} />);

    expect(screen.getByRole("button", { name: /Generate Document/i })).toBeDisabled();
  });

  it("keeps generate button enabled when quota is below limit", () => {
    const quota = { limit: 100, percentage: 80, reserved: 0, used: 80 };

    render(<EmptyDocument onGenerate={mockOnGenerate} quota={quota} />);

    expect(screen.getByRole("button", { name: /Generate Document/i })).toBeEnabled();
  });

  it("keeps generate button enabled for unlimited quota (percentage null)", () => {
    const quota = { limit: null, percentage: null, reserved: 0, used: 50 };

    render(<EmptyDocument onGenerate={mockOnGenerate} quota={quota} />);

    expect(screen.getByRole("button", { name: /Generate Document/i })).toBeEnabled();
  });
});
