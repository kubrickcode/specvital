import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import { ApiError, ERROR_TYPES, fetchAnalysis, getErrorMessage } from "./client";

describe("ApiError", () => {
  describe("categorizeError", () => {
    it("categorizes 400 as BAD_REQUEST", () => {
      const error = new ApiError("msg", 400);
      expect(error.errorType).toBe(ERROR_TYPES.BAD_REQUEST);
    });

    it("categorizes 403 as FORBIDDEN", () => {
      const error = new ApiError("msg", 403);
      expect(error.errorType).toBe(ERROR_TYPES.FORBIDDEN);
    });

    it("categorizes 404 as NOT_FOUND", () => {
      const error = new ApiError("msg", 404);
      expect(error.errorType).toBe(ERROR_TYPES.NOT_FOUND);
    });

    it("categorizes 422 as VALIDATION_ERROR", () => {
      const error = new ApiError("msg", 422);
      expect(error.errorType).toBe(ERROR_TYPES.VALIDATION_ERROR);
    });

    it("categorizes 429 as RATE_LIMIT", () => {
      const error = new ApiError("msg", 429);
      expect(error.errorType).toBe(ERROR_TYPES.RATE_LIMIT);
    });

    it("categorizes 5xx as SERVER_ERROR", () => {
      const error500 = new ApiError("msg", 500);
      const error502 = new ApiError("msg", 502);
      const error503 = new ApiError("msg", 503);

      expect(error500.errorType).toBe(ERROR_TYPES.SERVER_ERROR);
      expect(error502.errorType).toBe(ERROR_TYPES.SERVER_ERROR);
      expect(error503.errorType).toBe(ERROR_TYPES.SERVER_ERROR);
    });

    it("categorizes unknown 4xx as BAD_REQUEST", () => {
      const error = new ApiError("msg", 418);
      expect(error.errorType).toBe(ERROR_TYPES.BAD_REQUEST);
    });

    it("categorizes undefined status as NETWORK", () => {
      const error = new ApiError("msg");
      expect(error.errorType).toBe(ERROR_TYPES.NETWORK);
    });
  });

  describe("factory methods", () => {
    it("timeout() creates TIMEOUT error", () => {
      const error = ApiError.timeout();
      expect(error.errorType).toBe(ERROR_TYPES.TIMEOUT);
      expect(error.message).toBe("Request timed out");
    });

    it("timeout() accepts custom message", () => {
      const error = ApiError.timeout("Custom timeout");
      expect(error.message).toBe("Custom timeout");
    });

    it("parseError() creates PARSE_ERROR error", () => {
      const data = { invalid: "data" };
      const error = ApiError.parseError("Parse failed", data);
      expect(error.errorType).toBe(ERROR_TYPES.PARSE_ERROR);
      expect(error.responseBody).toBe(data);
    });
  });

  describe("ProblemDetail extraction", () => {
    it("extracts valid ProblemDetail from responseBody", () => {
      const problemDetail = {
        status: 404,
        title: "Not Found",
        detail: "Repository not found",
      };
      const error = new ApiError("msg", 404, problemDetail);
      expect(error.problemDetail).toEqual(problemDetail);
    });

    it("ignores invalid ProblemDetail", () => {
      const invalid = { foo: "bar" };
      const error = new ApiError("msg", 400, invalid);
      expect(error.problemDetail).toBeUndefined();
    });

    it("handles ProblemDetail with optional fields", () => {
      const problemDetail = {
        status: 403,
        title: "Forbidden",
        detail: "Access denied",
        instance: "/api/analyze/owner/repo",
        type: "about:blank",
      };
      const error = new ApiError("msg", 403, problemDetail);
      expect(error.problemDetail).toEqual(problemDetail);
    });
  });

  describe("getUserMessage", () => {
    it("returns ProblemDetail.detail when available", () => {
      const problemDetail = {
        status: 404,
        title: "Not Found",
        detail: "Custom error detail from server",
      };
      const error = new ApiError("msg", 404, problemDetail);
      expect(error.getUserMessage()).toBe("Custom error detail from server");
    });

    it("returns fallback message for BAD_REQUEST", () => {
      const error = new ApiError("msg", 400);
      expect(error.getUserMessage()).toBe(
        "Invalid request. Please check the repository owner and name.",
      );
    });

    it("returns fallback message for FORBIDDEN", () => {
      const error = new ApiError("msg", 403);
      expect(error.getUserMessage()).toBe(
        "Access denied. This repository may be private or require authentication.",
      );
    });

    it("returns fallback message for NOT_FOUND", () => {
      const error = new ApiError("msg", 404);
      expect(error.getUserMessage()).toBe(
        "Repository not found. Please verify the owner and repository name.",
      );
    });

    it("returns fallback message for RATE_LIMIT", () => {
      const error = new ApiError("msg", 429);
      expect(error.getUserMessage()).toBe("Rate limit exceeded. Please try again later.");
    });

    it("returns fallback message for NETWORK", () => {
      const error = new ApiError("msg");
      expect(error.getUserMessage()).toBe(
        "Network error. Please check your connection and try again.",
      );
    });

    it("returns fallback message for SERVER_ERROR", () => {
      const error = new ApiError("msg", 500);
      expect(error.getUserMessage()).toBe("Server error. Please try again later.");
    });

    it("returns fallback message for TIMEOUT", () => {
      const error = ApiError.timeout();
      expect(error.getUserMessage()).toBe("Request timed out. Please try again.");
    });

    it("returns fallback message for PARSE_ERROR", () => {
      const error = ApiError.parseError("msg");
      expect(error.getUserMessage()).toBe("Failed to parse server response. Please try again.");
    });

    it("returns fallback message for VALIDATION_ERROR", () => {
      const error = new ApiError("msg", 422);
      expect(error.getUserMessage()).toBe("Validation error. The provided data is invalid.");
    });
  });
});

describe("getErrorMessage", () => {
  it("returns getUserMessage for ApiError", () => {
    const error = new ApiError("msg", 404);
    expect(getErrorMessage(error)).toBe(
      "Repository not found. Please verify the owner and repository name.",
    );
  });

  it("returns message for generic Error", () => {
    const error = new Error("Generic error");
    expect(getErrorMessage(error)).toBe("Generic error");
  });

  it("returns default message for unknown error", () => {
    expect(getErrorMessage("string error")).toBe("An unexpected error occurred. Please try again.");
    expect(getErrorMessage(null)).toBe("An unexpected error occurred. Please try again.");
    expect(getErrorMessage(undefined)).toBe("An unexpected error occurred. Please try again.");
  });
});

describe("fetchAnalysis", () => {
  const mockAnalysisResult = {
    analyzedAt: "2024-01-01T00:00:00Z",
    commitSha: "abc123",
    owner: "test",
    repo: "repo",
    suites: [],
    summary: {
      active: 10,
      skipped: 2,
      todo: 1,
      total: 13,
      frameworks: [],
    },
  };

  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("returns parsed result on success", async () => {
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockAnalysisResult),
    } as Response);

    const result = await fetchAnalysis("owner", "repo");
    expect(result).toEqual(mockAnalysisResult);
  });

  it("throws ApiError on HTTP error with ProblemDetail", async () => {
    const problemDetail = {
      status: 404,
      title: "Not Found",
      detail: "Repository not found",
    };

    vi.mocked(fetch).mockResolvedValueOnce({
      ok: false,
      status: 404,
      statusText: "Not Found",
      json: () => Promise.resolve(problemDetail),
    } as Response);

    try {
      await fetchAnalysis("owner", "notfound");
      expect.fail("Should have thrown");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiError);
      if (error instanceof ApiError) {
        expect(error.status).toBe(404);
        expect(error.problemDetail).toEqual(problemDetail);
      }
    }
  });

  it("throws TIMEOUT error on AbortError", async () => {
    const abortError = new Error("Aborted");
    abortError.name = "AbortError";

    vi.mocked(fetch).mockRejectedValueOnce(abortError);

    try {
      await fetchAnalysis("owner", "repo", 100);
      expect.fail("Should have thrown");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiError);
      if (error instanceof ApiError) {
        expect(error.errorType).toBe(ERROR_TYPES.TIMEOUT);
      }
    }
  });

  it("throws NETWORK error on fetch failure", async () => {
    vi.mocked(fetch).mockRejectedValueOnce(new Error("Network failure"));

    try {
      await fetchAnalysis("owner", "repo");
      expect.fail("Should have thrown");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiError);
      if (error instanceof ApiError) {
        expect(error.errorType).toBe(ERROR_TYPES.NETWORK);
      }
    }
  });

  it("throws PARSE_ERROR on invalid response format", async () => {
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ invalid: "data" }),
    } as Response);

    try {
      await fetchAnalysis("owner", "repo");
      expect.fail("Should have thrown");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiError);
      if (error instanceof ApiError) {
        expect(error.errorType).toBe(ERROR_TYPES.PARSE_ERROR);
      }
    }
  });

  it("throws error on JSON parse failure", async () => {
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.reject(new Error("Invalid JSON")),
    } as Response);

    await expect(fetchAnalysis("owner", "repo")).rejects.toThrow(ApiError);
  });

  it("handles non-JSON error response", async () => {
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: false,
      status: 500,
      statusText: "Internal Server Error",
      json: () => Promise.reject(new Error("Not JSON")),
      text: () => Promise.resolve("Internal server error"),
    } as unknown as Response);

    try {
      await fetchAnalysis("owner", "repo");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiError);
      if (error instanceof ApiError) {
        expect(error.status).toBe(500);
        expect(error.responseBody).toBe("Internal server error");
      }
    }
  });
});
