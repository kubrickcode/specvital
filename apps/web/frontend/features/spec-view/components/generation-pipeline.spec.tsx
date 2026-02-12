import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import messages from "@/i18n/messages/en.json";

import { GenerationPipeline } from "./generation-pipeline";

vi.mock("motion/react", () => ({
  m: {
    div: ({ children, ...props }: React.PropsWithChildren) => <div {...props}>{children}</div>,
    span: ({ children, ...props }: React.PropsWithChildren) => <span {...props}>{children}</span>,
  },
  useReducedMotion: () => false,
}));

vi.mock("@/lib/motion", () => ({
  useReducedMotion: () => false,
}));

vi.mock("@/components/feedback/pulse-ring", () => ({
  PulseRing: () => <span data-testid="pulse-ring" />,
}));

vi.mock("@/components/feedback/shimmer-bar", () => ({
  ShimmerBar: () => <div data-testid="shimmer-bar" role="progressbar" />,
}));

const renderWithProvider = (ui: React.ReactElement) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      {ui}
    </NextIntlClientProvider>
  );

describe("GenerationPipeline", () => {
  it("renders three pipeline steps", () => {
    renderWithProvider(<GenerationPipeline status={null} />);

    expect(screen.getByText("Queued")).toBeInTheDocument();
    expect(screen.getByText("Generating")).toBeInTheDocument();
    expect(screen.getByText("Complete")).toBeInTheDocument();
  });

  it("renders list with accessible label", () => {
    renderWithProvider(<GenerationPipeline status={null} />);

    expect(
      screen.getByRole("list", { name: "Spec generation progress steps" })
    ).toBeInTheDocument();
  });

  describe("pending status", () => {
    it("marks step 1 as active and shows its number", () => {
      renderWithProvider(<GenerationPipeline status="pending" />);

      const listItems = screen.getAllByRole("listitem");
      expect(listItems[0]).toHaveAttribute("aria-current", "step");
      expect(listItems[0]).toHaveTextContent("1");
    });

    it("keeps steps 2 and 3 as upcoming", () => {
      renderWithProvider(<GenerationPipeline status="pending" />);

      const listItems = screen.getAllByRole("listitem");
      expect(listItems[1]).not.toHaveAttribute("aria-current");
      expect(listItems[2]).not.toHaveAttribute("aria-current");
    });
  });

  describe("running status", () => {
    it("marks step 1 as completed and step 2 as active", () => {
      renderWithProvider(<GenerationPipeline status="running" />);

      const listItems = screen.getAllByRole("listitem");
      // Step 1 completed: no aria-current, has check icon
      expect(listItems[0]).not.toHaveAttribute("aria-current");
      // Step 2 active
      expect(listItems[1]).toHaveAttribute("aria-current", "step");
      expect(listItems[1]).toHaveTextContent("2");
    });

    it("renders shimmer bar for the active connector", () => {
      renderWithProvider(<GenerationPipeline status="running" />);

      expect(screen.getByTestId("shimmer-bar")).toBeInTheDocument();
    });
  });

  describe("completed status", () => {
    it("marks all steps as completed with no active step", () => {
      renderWithProvider(<GenerationPipeline status="completed" />);

      const listItems = screen.getAllByRole("listitem");
      listItems.forEach((item) => {
        expect(item).not.toHaveAttribute("aria-current");
      });
    });
  });

  describe("failed status", () => {
    it("marks step 2 as failed with no active step", () => {
      renderWithProvider(<GenerationPipeline status="failed" />);

      const listItems = screen.getAllByRole("listitem");
      // Step 1: completed
      expect(listItems[0]).not.toHaveAttribute("aria-current");
      // Step 2: failed (no aria-current since it's not "active")
      expect(listItems[1]).not.toHaveAttribute("aria-current");
      // Step 3: upcoming
      expect(listItems[2]).not.toHaveAttribute("aria-current");
    });
  });

  it("accepts className prop", () => {
    renderWithProvider(<GenerationPipeline className="custom-class" status={null} />);

    const list = screen.getByRole("list");
    expect(list.className).toContain("custom-class");
  });
});
