import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { fetchAnalysis } from "./client";

describe("fetchAnalysis", () => {
  const mockAnalysisResult = {
    analyzedAt: "2024-01-01T00:00:00Z",
    commitSha: "abc123",
    owner: "test",
    repo: "repo",
    suites: [],
    summary: {
      active: 10,
      focused: 0,
      frameworks: [],
      skipped: 2,
      todo: 1,
      total: 13,
      xfail: 0,
    },
  };

  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("returns completed response on success", async () => {
    const response = { data: mockAnalysisResult, status: "completed" };
    vi.mocked(fetch).mockResolvedValueOnce({
      json: () => Promise.resolve(response),
      ok: true,
    } as Response);

    const result = await fetchAnalysis("owner", "repo");
    expect(result).toEqual(response);
  });

  it("returns queued response with 202 status", async () => {
    const response = { status: "queued" };
    vi.mocked(fetch).mockResolvedValueOnce({
      json: () => Promise.resolve(response),
      ok: false,
      status: 202,
    } as Response);

    const result = await fetchAnalysis("owner", "repo");
    expect(result).toEqual(response);
  });

  it("returns analyzing response", async () => {
    const response = { status: "analyzing" };
    vi.mocked(fetch).mockResolvedValueOnce({
      json: () => Promise.resolve(response),
      ok: true,
    } as Response);

    const result = await fetchAnalysis("owner", "repo");
    expect(result).toEqual(response);
  });

  it("returns failed response", async () => {
    const response = { error: "Analysis failed", status: "failed" };
    vi.mocked(fetch).mockResolvedValueOnce({
      json: () => Promise.resolve(response),
      ok: true,
    } as Response);

    const result = await fetchAnalysis("owner", "repo");
    expect(result).toEqual(response);
  });

  it("throws Error on HTTP error with ProblemDetail", async () => {
    const problemDetail = {
      detail: "Repository not found",
      status: 404,
      title: "Not Found",
    };

    vi.mocked(fetch).mockResolvedValueOnce({
      json: () => Promise.resolve(problemDetail),
      ok: false,
      status: 404,
      statusText: "Not Found",
    } as Response);

    await expect(fetchAnalysis("owner", "notfound")).rejects.toThrow("Repository not found");
  });

  it("throws timeout error on AbortError", async () => {
    const abortError = new Error("Aborted");
    abortError.name = "AbortError";

    vi.mocked(fetch).mockRejectedValueOnce(abortError);

    await expect(fetchAnalysis("owner", "repo", 100)).rejects.toThrow("Request timed out");
  });

  it("throws network error on fetch failure", async () => {
    vi.mocked(fetch).mockRejectedValueOnce(new Error("Network failure"));

    await expect(fetchAnalysis("owner", "repo")).rejects.toThrow(
      "Failed to fetch analysis: Network failure"
    );
  });

  it("throws error on JSON parse failure", async () => {
    vi.mocked(fetch).mockResolvedValueOnce({
      json: () => Promise.reject(new Error("Invalid JSON")),
      ok: true,
    } as Response);

    await expect(fetchAnalysis("owner", "repo")).rejects.toThrow(
      "Failed to parse response as JSON: Invalid JSON"
    );
  });

  it("handles non-JSON error response with statusText", async () => {
    vi.mocked(fetch).mockResolvedValueOnce({
      json: () => Promise.reject(new Error("Not JSON")),
      ok: false,
      status: 500,
      statusText: "Internal Server Error",
    } as unknown as Response);

    await expect(fetchAnalysis("owner", "repo")).rejects.toThrow(
      "API request failed: Internal Server Error"
    );
  });
});
