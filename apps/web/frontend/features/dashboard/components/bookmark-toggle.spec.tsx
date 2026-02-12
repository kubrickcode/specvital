import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { BookmarkToggle } from "./bookmark-toggle";

const mockSetBookmarkOnly = vi.fn();
const mockBookmarkOnly = vi.fn();

vi.mock("../hooks/use-bookmark-filter", () => ({
  useBookmarkFilter: () => ({
    bookmarkOnly: mockBookmarkOnly(),
    setBookmarkOnly: mockSetBookmarkOnly,
  }),
}));

const messages = {
  dashboard: {
    filter: {
      bookmarked: "Bookmarked",
      bookmarkedLabel: "Show bookmarked only",
    },
  },
};

const renderBookmarkToggle = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <BookmarkToggle />
    </NextIntlClientProvider>
  );
};

describe("BookmarkToggle", () => {
  beforeEach(() => {
    mockBookmarkOnly.mockReturnValue(false);
    mockSetBookmarkOnly.mockClear();
  });

  it("renders with aria-pressed false when not active", () => {
    renderBookmarkToggle();

    const toggle = screen.getByRole("button", { name: /show bookmarked only/i });

    expect(toggle).toHaveAttribute("aria-pressed", "false");
  });

  it("renders with aria-pressed true when bookmark filter is active", () => {
    mockBookmarkOnly.mockReturnValue(true);
    renderBookmarkToggle();

    const toggle = screen.getByRole("button", { name: /show bookmarked only/i });

    expect(toggle).toHaveAttribute("aria-pressed", "true");
  });

  it("calls setBookmarkOnly with true when toggled on", () => {
    renderBookmarkToggle();

    fireEvent.click(screen.getByRole("button", { name: /show bookmarked only/i }));

    expect(mockSetBookmarkOnly).toHaveBeenCalledWith(true);
  });

  it("calls setBookmarkOnly with null when toggled off", () => {
    mockBookmarkOnly.mockReturnValue(true);
    renderBookmarkToggle();

    fireEvent.click(screen.getByRole("button", { name: /show bookmarked only/i }));

    expect(mockSetBookmarkOnly).toHaveBeenCalledWith(null);
  });

  it("displays the bookmarked label text", () => {
    renderBookmarkToggle();

    expect(screen.getByText("Bookmarked")).toBeInTheDocument();
  });
});
