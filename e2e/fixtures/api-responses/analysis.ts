/**
 * Mock data for Analysis page APIs
 * Endpoints:
 * - /api/analyze/{owner}/{repo}
 * - /api/spec-view/{analysisId}
 * - /api/spec-view/generate
 * - /api/spec-view/status/{analysisId}
 */

export interface AnalysisResult {
  id: string;
  analyzedAt: string;
  branchName?: string;
  commitSha: string;
  committedAt?: string;
  isInMyHistory?: boolean;
  owner: string;
  parserVersion?: string;
  repo: string;
  suites: TestSuite[];
  summary: Summary;
}

export interface TestSuite {
  filePath: string;
  framework: string;
  suiteName: string;
  tests: TestCase[];
}

export interface TestCase {
  filePath: string;
  framework: string;
  line: number;
  modifier?: string;
  name: string;
  status: "active" | "focused" | "skipped" | "todo" | "xfail";
}

export interface Summary {
  active: number;
  focused: number;
  skipped: number;
  todo: number;
  total: number;
  xfail: number;
  frameworks: FrameworkSummary[];
}

export interface FrameworkSummary {
  framework: string;
  active: number;
  focused: number;
  skipped: number;
  todo: number;
  total: number;
  xfail: number;
}

export interface AnalysisCompletedResponse {
  status: "completed";
  data: AnalysisResult;
}

export interface SpecDocumentResponse {
  status: "completed" | "generating";
  data?: SpecDocument;
  generationStatus?: {
    status: "pending" | "running" | "completed" | "failed" | "not_found";
  };
}

export interface SpecDocument {
  id: string;
  analysisId: string;
  language: SpecLanguage;
  version: number;
  executiveSummary?: string;
  modelId?: string;
  availableLanguages?: AvailableLanguageInfo[];
  domains: SpecDomain[];
  createdAt: string;
}

export interface AvailableLanguageInfo {
  language: SpecLanguage;
  latestVersion: number;
  createdAt: string;
}

export interface SpecDomain {
  id: string;
  name: string;
  description?: string;
  features: SpecFeature[];
}

export interface SpecFeature {
  id: string;
  name: string;
  description?: string;
  behaviors: SpecBehavior[];
}

export interface SpecBehavior {
  id: string;
  description: string;
  status?: "active" | "skipped" | "todo" | "focused" | "xfail";
}

export interface RequestSpecGenerationResponse {
  status: "pending" | "running" | "completed" | "failed" | "not_found";
  analysisId: string;
}

// Sample analysis data
const sampleAnalysisId = "550e8400-e29b-41d4-a716-446655440000";
const now = new Date().toISOString();
const commitTime = new Date(Date.now() - 86400000).toISOString();

export const mockAnalysisCompleted: AnalysisCompletedResponse = {
  status: "completed",
  data: {
    id: sampleAnalysisId,
    analyzedAt: now,
    branchName: "main",
    commitSha: "abc123def",
    committedAt: commitTime,
    isInMyHistory: true,
    owner: "test-owner",
    repo: "test-repo",
    parserVersion: "v1.0.0",
    suites: [
      {
        filePath: "src/__tests__/example.test.ts",
        framework: "jest",
        suiteName: "Example Test Suite",
        tests: [
          {
            filePath: "src/__tests__/example.test.ts",
            framework: "jest",
            line: 10,
            name: "should work correctly",
            status: "active",
          },
          {
            filePath: "src/__tests__/example.test.ts",
            framework: "jest",
            line: 20,
            name: "should handle errors",
            status: "active",
          },
          {
            filePath: "src/__tests__/example.test.ts",
            framework: "jest",
            line: 30,
            name: "should be skipped",
            status: "skipped",
            modifier: "skip",
          },
        ],
      },
      {
        filePath: "src/__tests__/another.test.ts",
        framework: "jest",
        suiteName: "Another Test Suite",
        tests: [
          {
            filePath: "src/__tests__/another.test.ts",
            framework: "jest",
            line: 5,
            name: "another test case",
            status: "active",
          },
        ],
      },
    ],
    summary: {
      active: 3,
      focused: 0,
      skipped: 1,
      todo: 0,
      total: 4,
      xfail: 0,
      frameworks: [
        {
          framework: "jest",
          active: 3,
          focused: 0,
          skipped: 1,
          todo: 0,
          total: 4,
          xfail: 0,
        },
      ],
    },
  },
};

// Large dataset for virtual scroll testing (100 suites)
const generateLargeSuites = (): TestSuite[] => {
  return Array.from({ length: 100 }, (_, suiteIndex) => ({
    filePath: `src/__tests__/suite-${suiteIndex + 1}.test.ts`,
    framework: "jest",
    suiteName: `Test Suite ${suiteIndex + 1}`,
    tests: Array.from({ length: 10 }, (_, testIndex) => ({
      filePath: `src/__tests__/suite-${suiteIndex + 1}.test.ts`,
      framework: "jest",
      line: (testIndex + 1) * 10,
      name: `test case ${testIndex + 1} in suite ${suiteIndex + 1}`,
      status: testIndex % 5 === 0 ? "skipped" : "active" as const,
    })),
  }));
};

export const mockAnalysisLarge: AnalysisCompletedResponse = {
  status: "completed",
  data: {
    id: "large-analysis-id",
    analyzedAt: now,
    branchName: "main",
    commitSha: "def456abc",
    committedAt: commitTime,
    isInMyHistory: true,
    owner: "test-owner",
    repo: "large-test-repo",
    parserVersion: "v1.0.0",
    suites: generateLargeSuites(),
    summary: {
      active: 800,
      focused: 0,
      skipped: 200,
      todo: 0,
      total: 1000,
      xfail: 0,
      frameworks: [
        {
          framework: "jest",
          active: 800,
          focused: 0,
          skipped: 200,
          todo: 0,
          total: 1000,
          xfail: 0,
        },
      ],
    },
  },
};

// Spec document not found (for generate flow)
export const mockSpecDocumentNotFound: SpecDocumentResponse = {
  status: "generating",
  generationStatus: {
    status: "not_found",
  },
};

// Spec document completed with content (for language switch, regeneration, TOC tests)
export const mockSpecDocumentCompleted: SpecDocumentResponse = {
  status: "completed",
  data: {
    id: "doc-123",
    analysisId: sampleAnalysisId,
    language: "English",
    version: 2,
    createdAt: now,
    executiveSummary: "Test Repository Specification - This document describes the test specifications for the test repository.",
    modelId: "gemini-2.0-flash",
    availableLanguages: [
      {
        language: "English",
        latestVersion: 2,
        createdAt: now,
      },
      {
        language: "Korean",
        latestVersion: 1,
        createdAt: new Date(Date.now() - 86400000).toISOString(), // 1 day ago
      },
    ],
    domains: [
      {
        id: "domain-1",
        name: "User Authentication",
        description: "Handles user login and registration flows",
        features: [
          {
            id: "feature-1",
            name: "Login Flow",
            description: "Validates user credentials and issues session tokens",
            behaviors: [
              {
                id: "spec-1",
                description: "should authenticate with valid email and password",
                status: "active",
              },
              {
                id: "spec-2",
                description: "should reject invalid credentials",
                status: "active",
              },
            ],
          },
          {
            id: "feature-2",
            name: "Registration Flow",
            description: "Creates new user accounts with email verification",
            behaviors: [
              {
                id: "spec-3",
                description: "should create user with valid email",
                status: "active",
              },
            ],
          },
        ],
      },
      {
        id: "domain-2",
        name: "Payment Processing",
        description: "Handles payment transactions and billing",
        features: [
          {
            id: "feature-3",
            name: "Checkout Flow",
            description: "Processes payment and creates orders",
            behaviors: [
              {
                id: "spec-4",
                description: "should process payment with valid card",
                status: "active",
              },
            ],
          },
        ],
      },
    ],
  },
};

// Spec generation status - not found (document not yet generated)
export const mockSpecStatusNotFound: RequestSpecGenerationResponse = {
  status: "not_found",
  analysisId: sampleAnalysisId,
};

// Spec generation accepted
export const mockSpecGenerationAccepted: RequestSpecGenerationResponse = {
  status: "pending",
  analysisId: sampleAnalysisId,
};

// Spec generation running
export const mockSpecGenerationRunning: RequestSpecGenerationResponse = {
  status: "running",
  analysisId: sampleAnalysisId,
};

// Spec generation completed
export const mockSpecGenerationCompleted: RequestSpecGenerationResponse = {
  status: "completed",
  analysisId: sampleAnalysisId,
};

// Spec languages (24 languages)
export const SPEC_LANGUAGES = [
  "Arabic",
  "Chinese",
  "Czech",
  "Danish",
  "Dutch",
  "English",
  "Finnish",
  "French",
  "German",
  "Greek",
  "Hindi",
  "Indonesian",
  "Italian",
  "Japanese",
  "Korean",
  "Polish",
  "Portuguese",
  "Russian",
  "Spanish",
  "Swedish",
  "Thai",
  "Turkish",
  "Ukrainian",
  "Vietnamese",
] as const;

export type SpecLanguage = (typeof SPEC_LANGUAGES)[number];
