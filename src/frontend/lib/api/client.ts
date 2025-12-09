import { z } from "zod";
import type { AnalysisResponse, ProblemDetail } from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000";
const DEFAULT_TIMEOUT_MS = 30000;

const getApiUrl = (path: string): string => {
  if (typeof window === "undefined") {
    return `${API_URL}${path}`;
  }
  return path;
};

export const ERROR_TYPES = {
  BAD_REQUEST: "BAD_REQUEST",
  FORBIDDEN: "FORBIDDEN",
  NETWORK: "NETWORK",
  NOT_FOUND: "NOT_FOUND",
  PARSE_ERROR: "PARSE_ERROR",
  RATE_LIMIT: "RATE_LIMIT",
  SERVER_ERROR: "SERVER_ERROR",
  TIMEOUT: "TIMEOUT",
  VALIDATION_ERROR: "VALIDATION_ERROR",
} as const;

export type ErrorType = (typeof ERROR_TYPES)[keyof typeof ERROR_TYPES];

const testStatusSchema = z.enum(["active", "skipped", "todo"]);
const frameworkSchema = z.string();

const testCaseSchema = z.object({
  filePath: z.string(),
  framework: frameworkSchema,
  line: z.number(),
  name: z.string(),
  status: testStatusSchema,
});

const testSuiteSchema = z.object({
  filePath: z.string(),
  framework: frameworkSchema,
  tests: z.array(testCaseSchema),
});

const frameworkSummarySchema = z.object({
  active: z.number(),
  framework: frameworkSchema,
  skipped: z.number(),
  todo: z.number(),
  total: z.number(),
});

const summarySchema = z.object({
  active: z.number(),
  frameworks: z.array(frameworkSummarySchema),
  skipped: z.number(),
  todo: z.number(),
  total: z.number(),
});

const analysisResultSchema = z.object({
  analyzedAt: z.string(),
  commitSha: z.string(),
  owner: z.string(),
  repo: z.string(),
  suites: z.array(testSuiteSchema),
  summary: summarySchema,
});

const analysisResponseSchema = z.discriminatedUnion("status", [
  z.object({ status: z.literal("completed"), data: analysisResultSchema }),
  z.object({ status: z.literal("analyzing") }),
  z.object({ status: z.literal("queued") }),
  z.object({ status: z.literal("failed"), error: z.string() }),
]);

const rateLimitInfoSchema = z.object({
  limit: z.number(),
  remaining: z.number(),
  resetAt: z.number(),
});

const problemDetailSchema = z.object({
  detail: z.string(),
  instance: z.string().optional(),
  rateLimit: rateLimitInfoSchema.optional(),
  status: z.number(),
  title: z.string(),
  type: z.string().optional(),
});

export class ApiError extends Error {
  public readonly errorType: ErrorType;
  public readonly problemDetail?: ProblemDetail;

  constructor(
    message: string,
    public readonly status?: number,
    public readonly responseBody?: unknown,
    errorTypeOverride?: ErrorType
  ) {
    super(message);
    this.name = "ApiError";
    this.errorType = errorTypeOverride ?? this.categorizeError();

    if (this.isProblemDetail(responseBody)) {
      this.problemDetail = responseBody;
    }
  }

  static timeout(message = "Request timed out"): ApiError {
    return new ApiError(message, undefined, undefined, ERROR_TYPES.TIMEOUT);
  }

  static parseError(message: string, data?: unknown): ApiError {
    return new ApiError(message, undefined, data, ERROR_TYPES.PARSE_ERROR);
  }

  private categorizeError(): ErrorType {
    if (!this.status) {
      return ERROR_TYPES.NETWORK;
    }

    switch (this.status) {
      case 400:
        return ERROR_TYPES.BAD_REQUEST;
      case 403:
        return ERROR_TYPES.FORBIDDEN;
      case 404:
        return ERROR_TYPES.NOT_FOUND;
      case 422:
        return ERROR_TYPES.VALIDATION_ERROR;
      case 429:
        return ERROR_TYPES.RATE_LIMIT;
      default:
        if (this.status >= 500) {
          return ERROR_TYPES.SERVER_ERROR;
        }
        return ERROR_TYPES.BAD_REQUEST;
    }
  }

  private isProblemDetail(body: unknown): body is ProblemDetail {
    return problemDetailSchema.safeParse(body).success;
  }

  public getUserMessage(): string {
    if (this.problemDetail) {
      return this.problemDetail.detail;
    }

    switch (this.errorType) {
      case ERROR_TYPES.BAD_REQUEST:
        return "Invalid request. Please check the repository owner and name.";
      case ERROR_TYPES.FORBIDDEN:
        return "Access denied. This repository may be private or require authentication.";
      case ERROR_TYPES.NOT_FOUND:
        return "Repository not found. Please verify the owner and repository name.";
      case ERROR_TYPES.NETWORK:
        return "Network error. Please check your connection and try again.";
      case ERROR_TYPES.PARSE_ERROR:
        return "Failed to parse server response. Please try again.";
      case ERROR_TYPES.RATE_LIMIT:
        return "Rate limit exceeded. Please try again later.";
      case ERROR_TYPES.SERVER_ERROR:
        return "Server error. Please try again later.";
      case ERROR_TYPES.TIMEOUT:
        return "Request timed out. Please try again.";
      case ERROR_TYPES.VALIDATION_ERROR:
        return "Validation error. The provided data is invalid.";
      default:
        return this.message;
    }
  }
}

export const getErrorMessage = (error: unknown): string => {
  if (error instanceof ApiError) {
    return error.getUserMessage();
  }

  if (error instanceof Error) {
    return error.message;
  }

  return "An unexpected error occurred. Please try again.";
};

export const fetchAnalysis = async (
  owner: string,
  repo: string,
  timeoutMs = DEFAULT_TIMEOUT_MS
): Promise<AnalysisResponse> => {
  const url = getApiUrl(`/api/analyze/${owner}/${repo}`);

  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs);

  let response: Response;

  try {
    response = await fetch(url, {
      cache: "no-store",
      headers: { Accept: "application/json" },
      signal: controller.signal,
    });
  } catch (error) {
    if (error instanceof Error && error.name === "AbortError") {
      throw ApiError.timeout();
    }
    throw new ApiError(
      `Failed to fetch analysis: ${error instanceof Error ? error.message : "Network error"}`
    );
  } finally {
    clearTimeout(timeoutId);
  }

  // 202 Accepted is valid for queued/analyzing status
  if (!response.ok && response.status !== 202) {
    let responseBody: unknown;
    try {
      responseBody = await response.json();
    } catch {
      responseBody = await response.text();
    }
    throw new ApiError(`API request failed: ${response.statusText}`, response.status, responseBody);
  }

  let data: unknown;
  try {
    data = await response.json();
  } catch (error) {
    throw new ApiError(
      `Failed to parse response as JSON: ${error instanceof Error ? error.message : "Parse error"}`
    );
  }

  const parseResult = analysisResponseSchema.safeParse(data);
  if (!parseResult.success) {
    throw ApiError.parseError(`Invalid response format: ${parseResult.error.message}`, data);
  }

  return parseResult.data;
};
