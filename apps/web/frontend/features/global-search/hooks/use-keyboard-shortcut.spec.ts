import { fireEvent } from "@testing-library/react";
import { renderHook } from "@testing-library/react";
import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";

import { useKeyboardShortcut } from "./use-keyboard-shortcut";

const mockToggle = vi.fn();

vi.mock("./use-global-search-store", () => ({
  globalSearchStore: {
    toggle: () => mockToggle(),
  },
}));

describe("useKeyboardShortcut", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("toggles search dialog when Ctrl+K is pressed on non-Mac", () => {
    Object.defineProperty(navigator, "userAgent", {
      configurable: true,
      value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    });

    renderHook(() => useKeyboardShortcut());

    fireEvent.keyDown(document, { key: "k", ctrlKey: true });

    expect(mockToggle).toHaveBeenCalledTimes(1);
  });

  it("toggles search dialog when Meta+K (Cmd+K) is pressed on Mac", () => {
    Object.defineProperty(navigator, "userAgent", {
      configurable: true,
      value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
    });

    renderHook(() => useKeyboardShortcut());

    fireEvent.keyDown(document, { key: "k", metaKey: true });

    expect(mockToggle).toHaveBeenCalledTimes(1);
  });

  it("does not toggle when only K key is pressed without modifier", () => {
    renderHook(() => useKeyboardShortcut());

    fireEvent.keyDown(document, { key: "k" });

    expect(mockToggle).not.toHaveBeenCalled();
  });

  it("does not toggle when user is typing in an input field", () => {
    renderHook(() => useKeyboardShortcut());

    const input = document.createElement("input");
    document.body.appendChild(input);

    fireEvent.keyDown(input, { key: "k", ctrlKey: true });

    expect(mockToggle).not.toHaveBeenCalled();

    document.body.removeChild(input);
  });
});
