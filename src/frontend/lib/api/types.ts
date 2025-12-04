// Synced with backend: src/backend/analyzer/types.go

// RFC 7807 Problem Details for HTTP APIs
export type ProblemDetail = {
  detail: string;
  instance?: string;
  status: number;
  title: string;
  type?: string;
};

export type TestStatus = "active" | "skipped" | "todo";

export type Framework = "jest" | "vitest" | "playwright" | "go";

export type TestCase = {
  filePath: string;
  framework: Framework;
  line: number;
  name: string;
  status: TestStatus;
};

export type TestSuite = {
  filePath: string;
  framework: Framework;
  tests: TestCase[];
};

export type FrameworkSummary = {
  active: number;
  framework: Framework;
  skipped: number;
  todo: number;
  total: number;
};

export type Summary = {
  active: number;
  frameworks: FrameworkSummary[];
  skipped: number;
  todo: number;
  total: number;
};

export type AnalysisResult = {
  analyzedAt: string;
  commitSha: string;
  owner: string;
  repo: string;
  suites: TestSuite[];
  summary: Summary;
};
