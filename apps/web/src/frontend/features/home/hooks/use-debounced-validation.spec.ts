import { renderHook, act } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";

import { useDebouncedValidation } from "./use-debounced-validation";

describe("useDebouncedValidation", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("returns idle when value is empty", () => {
    const validate = vi.fn().mockReturnValue(true);
    const { result } = renderHook(() => useDebouncedValidation("", validate));

    expect(result.current).toBe("idle");
    expect(validate).not.toHaveBeenCalled();
  });

  it("returns idle when value length is below minLength", () => {
    const validate = vi.fn().mockReturnValue(true);
    const { result } = renderHook(() => useDebouncedValidation("abc", validate, { minLength: 5 }));

    expect(result.current).toBe("idle");
    expect(validate).not.toHaveBeenCalled();
  });

  it("returns valid after debounce delay when validation passes", () => {
    const validate = vi.fn().mockReturnValue(true);
    const { result } = renderHook(() =>
      useDebouncedValidation("https://github.com/test/repo", validate, { delay: 500 })
    );

    expect(result.current).toBe("idle");

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(result.current).toBe("valid");
    expect(validate).toHaveBeenCalledWith("https://github.com/test/repo");
  });

  it("returns invalid after debounce delay when validation fails", () => {
    const validate = vi.fn().mockReturnValue(false);
    const { result } = renderHook(() =>
      useDebouncedValidation("invalid-url", validate, { delay: 500, minLength: 1 })
    );

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(result.current).toBe("invalid");
  });

  it("debounces validation on rapid value changes", () => {
    const validate = vi.fn().mockReturnValue(true);
    const { rerender, result } = renderHook(
      ({ value }) => useDebouncedValidation(value, validate, { delay: 500, minLength: 1 }),
      { initialProps: { value: "a" } }
    );

    act(() => {
      vi.advanceTimersByTime(200);
    });
    rerender({ value: "ab" });

    act(() => {
      vi.advanceTimersByTime(200);
    });
    rerender({ value: "abc" });

    expect(validate).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(result.current).toBe("valid");
    expect(validate).toHaveBeenCalledTimes(1);
    expect(validate).toHaveBeenCalledWith("abc");
  });

  it("resets to idle when value becomes shorter than minLength", () => {
    const validate = vi.fn().mockReturnValue(true);
    const { rerender, result } = renderHook(
      ({ value }) => useDebouncedValidation(value, validate, { delay: 500, minLength: 5 }),
      { initialProps: { value: "https://github.com" } }
    );

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(result.current).toBe("valid");

    rerender({ value: "http" });
    expect(result.current).toBe("idle");
  });
});
