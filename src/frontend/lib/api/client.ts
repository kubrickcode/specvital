import { z } from "zod";
import type { AnalysisResult } from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3000";

// Zod schema for runtime validation matching Go backend types
const testStatusSchema = z.enum(["active", "skipped", "todo"]);

const frameworkSchema = z.enum(["jest", "vitest", "playwright", "go"]);

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

class ApiError extends Error {
  constructor(
    message: string,
    public status?: number,
    public responseBody?: unknown,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

export const fetchAnalysis = async (
  owner: string,
  repo: string,
): Promise<AnalysisResult> => {
  const url = `${API_URL}/api/analyze/${owner}/${repo}`;

  let response: Response;

  try {
    response = await fetch(url, {
      cache: "no-store",
      headers: {
        Accept: "application/json",
      },
    });
  } catch (error) {
    throw new ApiError(
      `Failed to fetch analysis: ${error instanceof Error ? error.message : "Network error"}`,
    );
  }

  if (!response.ok) {
    let responseBody: unknown;

    try {
      responseBody = await response.json();
    } catch {
      responseBody = await response.text();
    }

    throw new ApiError(
      `API request failed: ${response.statusText}`,
      response.status,
      responseBody,
    );
  }

  let data: unknown;

  try {
    data = await response.json();
  } catch (error) {
    throw new ApiError(
      `Failed to parse response as JSON: ${error instanceof Error ? error.message : "Parse error"}`,
    );
  }

  const parseResult = analysisResultSchema.safeParse(data);

  if (!parseResult.success) {
    throw new ApiError(
      `Invalid response format: ${parseResult.error.message}`,
      undefined,
      data,
    );
  }

  return parseResult.data;
};
