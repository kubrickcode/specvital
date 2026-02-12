import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { ShimmerBar } from "./shimmer-bar";

let mockShouldReduceMotion = false;

vi.mock("motion/react", () => ({
  m: {
    div: ({
      children,
      className,
      style,
      ...props
    }: React.PropsWithChildren<{ className?: string; style?: React.CSSProperties }>) => (
      <div className={className} data-testid="animated-div" style={style} {...props}>
        {children}
      </div>
    ),
  },
  useReducedMotion: () => mockShouldReduceMotion,
}));

vi.mock("@/lib/motion", () => ({
  useReducedMotion: () => mockShouldReduceMotion,
}));

describe("ShimmerBar", () => {
  it("renders with progressbar role", () => {
    render(<ShimmerBar />);

    expect(screen.getByRole("progressbar")).toBeInTheDocument();
  });

  it("renders with default sm height class", () => {
    render(<ShimmerBar />);

    const progressbar = screen.getByRole("progressbar");
    expect(progressbar.className).toContain("h-1");
  });

  it("renders with xs height class when specified", () => {
    render(<ShimmerBar height="xs" />);

    const progressbar = screen.getByRole("progressbar");
    expect(progressbar.className).toContain("h-0.5");
  });

  it("renders with md height class when specified", () => {
    render(<ShimmerBar height="md" />);

    const progressbar = screen.getByRole("progressbar");
    expect(progressbar.className).toContain("h-1.5");
  });

  it("sets correct aria attributes", () => {
    render(<ShimmerBar />);

    const progressbar = screen.getByRole("progressbar");
    expect(progressbar).toHaveAttribute("aria-valuemin", "0");
    expect(progressbar).toHaveAttribute("aria-valuemax", "100");
    expect(progressbar).toHaveAttribute("aria-valuetext", "Loading");
  });

  it("renders animated shimmer when motion is not reduced", () => {
    mockShouldReduceMotion = false;
    const { queryByTestId } = render(<ShimmerBar />);

    expect(queryByTestId("animated-div")).toBeInTheDocument();
  });

  it("renders static fallback when reduced motion is preferred", () => {
    mockShouldReduceMotion = true;
    const { queryByTestId } = render(<ShimmerBar />);

    expect(queryByTestId("animated-div")).not.toBeInTheDocument();

    const progressbar = screen.getByRole("progressbar");
    const staticDiv = progressbar.querySelector("div");
    expect(staticDiv).toBeInTheDocument();
    expect(staticDiv!.className).toContain("opacity-30");

    mockShouldReduceMotion = false;
  });

  it("applies custom className", () => {
    render(<ShimmerBar className="mt-4" />);

    const progressbar = screen.getByRole("progressbar");
    expect(progressbar.className).toContain("mt-4");
  });

  it("applies custom color style", () => {
    mockShouldReduceMotion = true;
    render(<ShimmerBar color="var(--ai-primary)" />);

    const progressbar = screen.getByRole("progressbar");
    const staticDiv = progressbar.querySelector("div");
    expect(staticDiv).toHaveStyle({ background: "var(--ai-primary)" });

    mockShouldReduceMotion = false;
  });
});
