import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { OwnershipDropdown } from "./ownership-dropdown";

const mockSetOwnershipFilter = vi.fn();
const mockOwnershipFilter = vi.fn();

vi.mock("../hooks/use-ownership-filter", async (importOriginal) => {
  const actual = await importOriginal<typeof import("../hooks/use-ownership-filter")>();
  return {
    ...actual,
    useOwnershipFilter: () => ({
      ownershipFilter: mockOwnershipFilter(),
      setOwnershipFilter: mockSetOwnershipFilter,
    }),
  };
});

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
    filter: {
      ownership: {
        all: "All",
        label: "Ownership",
        mine: "Mine",
        organization: "Organization",
        others: "Others",
      },
    },
  },
};

const renderOwnershipDropdown = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <OwnershipDropdown />
    </NextIntlClientProvider>
  );
};

const openDropdown = (trigger: HTMLElement) => {
  fireEvent.pointerDown(trigger, { button: 0, ctrlKey: false, pointerType: "mouse" });
};

describe("OwnershipDropdown", () => {
  beforeEach(() => {
    mockOwnershipFilter.mockReturnValue("all");
    mockSetOwnershipFilter.mockClear();
  });

  it("displays all ownership options when opened", () => {
    renderOwnershipDropdown();

    openDropdown(screen.getByRole("button", { name: /ownership/i }));

    expect(screen.getByRole("menuitemradio", { name: /All/i })).toBeInTheDocument();
    expect(screen.getByRole("menuitemradio", { name: /Mine/i })).toBeInTheDocument();
    expect(screen.getByRole("menuitemradio", { name: /Organization/i })).toBeInTheDocument();
    expect(screen.getByRole("menuitemradio", { name: /Others/i })).toBeInTheDocument();
  });

  it("calls setOwnershipFilter with 'mine' when Mine is selected", () => {
    renderOwnershipDropdown();

    openDropdown(screen.getByRole("button", { name: /ownership/i }));
    fireEvent.click(screen.getByRole("menuitemradio", { name: /Mine/i }));

    expect(mockSetOwnershipFilter).toHaveBeenCalledWith("mine");
  });

  it("calls setOwnershipFilter with 'organization' when Organization is selected", () => {
    renderOwnershipDropdown();

    openDropdown(screen.getByRole("button", { name: /ownership/i }));
    fireEvent.click(screen.getByRole("menuitemradio", { name: /Organization/i }));

    expect(mockSetOwnershipFilter).toHaveBeenCalledWith("organization");
  });

  it("marks the current filter as checked", () => {
    mockOwnershipFilter.mockReturnValue("mine");
    renderOwnershipDropdown();

    openDropdown(screen.getByRole("button", { name: /ownership/i }));

    expect(screen.getByRole("menuitemradio", { name: /Mine/i })).toHaveAttribute(
      "aria-checked",
      "true"
    );
    expect(screen.getByRole("menuitemradio", { name: /All/i })).toHaveAttribute(
      "aria-checked",
      "false"
    );
  });

  it("shows current filter label in the trigger button", () => {
    mockOwnershipFilter.mockReturnValue("organization");
    renderOwnershipDropdown();

    expect(screen.getByRole("button", { name: /ownership/i })).toHaveTextContent("Organization");
  });
});
