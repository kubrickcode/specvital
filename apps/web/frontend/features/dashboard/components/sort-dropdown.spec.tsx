import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import type { SortOption } from "../types";
import { SortDropdown } from "./sort-dropdown";

// Radix DropdownMenu checks `event instanceof PointerEvent` which jsdom lacks
class MockPointerEvent extends Event {
  button: number;
  ctrlKey: boolean;
  pointerType: string;
  constructor(type: string, props: PointerEventInit & EventInit = {}) {
    super(type, props);
    this.button = props.button ?? 0;
    this.ctrlKey = props.ctrlKey ?? false;
    this.pointerType = props.pointerType ?? "mouse";
  }
}
window.PointerEvent = MockPointerEvent as never;

Object.defineProperty(window, "matchMedia", {
  value: vi.fn().mockImplementation((query: string) => ({
    addEventListener: vi.fn(),
    addListener: vi.fn(),
    dispatchEvent: vi.fn(),
    matches: false,
    media: query,
    onchange: null,
    removeEventListener: vi.fn(),
    removeListener: vi.fn(),
  })),
  writable: true,
});

const messages = {
  dashboard: {
    sort: {
      label: "Sort by",
      name: "Name",
      recent: "Recent",
      tests: "Tests count",
    },
  },
};

type RenderOptions = {
  onSortChange?: (sort: SortOption) => void;
  sortBy?: SortOption;
};

const renderSortDropdown = (options: RenderOptions = {}) => {
  const props = {
    onSortChange: options.onSortChange ?? vi.fn(),
    sortBy: options.sortBy ?? "recent",
  };

  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <SortDropdown {...props} />
    </NextIntlClientProvider>
  );
};

const openDropdown = (trigger: HTMLElement) => {
  fireEvent.pointerDown(trigger, { button: 0, ctrlKey: false, pointerType: "mouse" });
};

describe("SortDropdown", () => {
  it("displays all sort options when opened", () => {
    renderSortDropdown();

    openDropdown(screen.getByRole("button", { name: /sort by/i }));

    expect(screen.getByRole("menuitemradio", { name: "Recent" })).toBeInTheDocument();
    expect(screen.getByRole("menuitemradio", { name: "Name" })).toBeInTheDocument();
    expect(screen.getByRole("menuitemradio", { name: "Tests count" })).toBeInTheDocument();
  });

  it("marks the current sort option as checked", () => {
    renderSortDropdown({ sortBy: "recent" });

    openDropdown(screen.getByRole("button", { name: /sort by/i }));

    expect(screen.getByRole("menuitemradio", { name: "Recent" })).toHaveAttribute(
      "aria-checked",
      "true"
    );
    expect(screen.getByRole("menuitemradio", { name: "Name" })).toHaveAttribute(
      "aria-checked",
      "false"
    );
  });

  it("marks Name as checked when sortBy is name", () => {
    renderSortDropdown({ sortBy: "name" });

    openDropdown(screen.getByRole("button", { name: /sort by/i }));

    expect(screen.getByRole("menuitemradio", { name: "Name" })).toHaveAttribute(
      "aria-checked",
      "true"
    );
    expect(screen.getByRole("menuitemradio", { name: "Recent" })).toHaveAttribute(
      "aria-checked",
      "false"
    );
  });

  it("marks Tests count as checked when sortBy is tests", () => {
    renderSortDropdown({ sortBy: "tests" });

    openDropdown(screen.getByRole("button", { name: /sort by/i }));

    expect(screen.getByRole("menuitemradio", { name: "Tests count" })).toHaveAttribute(
      "aria-checked",
      "true"
    );
  });

  it("calls onSortChange when selecting a different option", () => {
    const onSortChange = vi.fn();
    renderSortDropdown({ onSortChange, sortBy: "recent" });

    openDropdown(screen.getByRole("button", { name: /sort by/i }));
    fireEvent.click(screen.getByRole("menuitemradio", { name: "Name" }));

    expect(onSortChange).toHaveBeenCalledWith("name");
  });

  it("shows current sort label in the trigger button", () => {
    renderSortDropdown({ sortBy: "name" });

    expect(screen.getByRole("button", { name: /sort by/i })).toHaveTextContent("Sort by: Name");
  });
});
