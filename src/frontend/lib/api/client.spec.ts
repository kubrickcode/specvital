import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { apiFetch, parseJsonResponse } from "./client";

describe("apiFetch", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("returns response on success", async () => {
    const mockResponse = { ok: true } as Response;
    vi.mocked(fetch).mockResolvedValueOnce(mockResponse);

    const result = await apiFetch("/api/test");
    expect(result).toBe(mockResponse);
  });

  it("throws timeout error on AbortError", async () => {
    const abortError = new Error("Aborted");
    abortError.name = "AbortError";
    vi.mocked(fetch).mockRejectedValueOnce(abortError);

    await expect(apiFetch("/api/test")).rejects.toThrow("Request timed out");
  });

  it("throws network error on fetch failure", async () => {
    vi.mocked(fetch).mockRejectedValueOnce(new Error("Network failure"));

    await expect(apiFetch("/api/test")).rejects.toThrow("Network error: Network failure");
  });

  it("uses POST method when specified", async () => {
    vi.mocked(fetch).mockResolvedValueOnce({ ok: true } as Response);

    await apiFetch("/api/test", { method: "POST" });

    expect(fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({ method: "POST" })
    );
  });
});

describe("parseJsonResponse", () => {
  it("returns parsed JSON on success", async () => {
    const data = { foo: "bar" };
    const response = {
      json: () => Promise.resolve(data),
      ok: true,
    } as Response;

    const result = await parseJsonResponse(response);
    expect(result).toEqual(data);
  });

  it("throws error with detail on failure", async () => {
    const response = {
      json: () => Promise.resolve({ detail: "Resource not found" }),
      ok: false,
      statusText: "Not Found",
    } as Response;

    await expect(parseJsonResponse(response)).rejects.toThrow("Resource not found");
  });

  it("throws error with statusText when no detail", async () => {
    const response = {
      json: () => Promise.reject(new Error("Parse error")),
      ok: false,
      statusText: "Internal Server Error",
    } as Response;

    await expect(parseJsonResponse(response)).rejects.toThrow("Internal Server Error");
  });
});
