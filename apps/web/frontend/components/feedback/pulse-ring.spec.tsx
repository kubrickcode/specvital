import { render } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { PulseRing } from "./pulse-ring";

let mockShouldReduceMotion = false;

vi.mock("motion/react", () => ({
  m: {
    span: ({ children, className, ...props }: React.PropsWithChildren<{ className?: string }>) => (
      <span className={className} data-testid="animated-span" {...props}>
        {children}
      </span>
    ),
  },
  useReducedMotion: () => mockShouldReduceMotion,
}));

vi.mock("@/lib/motion", () => ({
  useReducedMotion: () => mockShouldReduceMotion,
}));

describe("PulseRing", () => {
  it("renders with default sm size classes", () => {
    const { container } = render(<PulseRing />);

    const wrapper = container.firstChild as HTMLElement;
    expect(wrapper.className).toContain("size-4");
  });

  it("renders with md size classes when specified", () => {
    const { container } = render(<PulseRing size="md" />);

    const wrapper = container.firstChild as HTMLElement;
    expect(wrapper.className).toContain("size-5");
  });

  it("renders with xs size classes when specified", () => {
    const { container } = render(<PulseRing size="xs" />);

    const wrapper = container.firstChild as HTMLElement;
    expect(wrapper.className).toContain("size-2");
  });

  it("applies custom className", () => {
    const { container } = render(<PulseRing className="text-red-500" />);

    const wrapper = container.firstChild as HTMLElement;
    expect(wrapper.className).toContain("text-red-500");
  });

  it("renders animated ring when motion is not reduced", () => {
    mockShouldReduceMotion = false;
    const { queryByTestId } = render(<PulseRing />);

    expect(queryByTestId("animated-span")).toBeInTheDocument();
  });

  it("hides animated ring when reduced motion is preferred", () => {
    mockShouldReduceMotion = true;
    const { queryByTestId } = render(<PulseRing />);

    expect(queryByTestId("animated-span")).not.toBeInTheDocument();
    mockShouldReduceMotion = false;
  });

  it("always renders the inner dot element", () => {
    const { container } = render(<PulseRing />);

    const dot = container.querySelector("span > span:last-child") as HTMLElement;
    expect(dot).toBeInTheDocument();
    expect(dot.className).toContain("rounded-full");
    expect(dot.className).toContain("bg-current");
  });
});
