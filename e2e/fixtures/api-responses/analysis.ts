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

export interface VersionInfo {
  version: number;
  createdAt: string;
  modelId?: string;
}

export interface VersionHistoryResponse {
  data: VersionInfo[];
  language: SpecLanguage;
  latestVersion: number;
}

// Repository-based spec document types (cross-analysis version access)
export interface RepoSpecDocument {
  id: string;
  analysisId: string;
  language: SpecLanguage;
  version: number;
  executiveSummary?: string;
  modelId?: string;
  createdAt: string;
  commitSha: string;
  availableLanguages?: AvailableLanguageInfo[];
  domains: SpecDomain[];
}

export interface RepoSpecDocumentCompleted {
  status: "completed";
  data: RepoSpecDocument;
}

export interface RepoSpecDocumentEmpty {
  status: "empty";
  message?: string;
}

export type RepoSpecDocumentResponse = RepoSpecDocumentCompleted | RepoSpecDocumentEmpty;

export interface RepoVersionInfo {
  id: string;
  analysisId: string;
  version: number;
  language?: SpecLanguage;
  modelId?: string;
  createdAt: string;
  commitSha: string;
}

export interface RepoVersionHistoryResponse {
  data: RepoVersionInfo[];
  language: SpecLanguage;
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

// Analysis with focused tests (for conditional display testing)
export const mockAnalysisWithFocused: AnalysisCompletedResponse = {
  status: "completed",
  data: {
    id: "focused-analysis-id",
    analyzedAt: now,
    branchName: "feature/test",
    commitSha: "focused123",
    committedAt: commitTime,
    isInMyHistory: true,
    owner: "test-owner",
    repo: "focused-test-repo",
    parserVersion: "v1.0.0",
    suites: [
      {
        filePath: "src/__tests__/focused.test.ts",
        framework: "jest",
        suiteName: "Focused Test Suite",
        tests: [
          {
            filePath: "src/__tests__/focused.test.ts",
            framework: "jest",
            line: 10,
            name: "should run only this test",
            status: "focused",
            modifier: "only",
          },
          {
            filePath: "src/__tests__/focused.test.ts",
            framework: "jest",
            line: 20,
            name: "another focused test",
            status: "focused",
            modifier: "only",
          },
          {
            filePath: "src/__tests__/focused.test.ts",
            framework: "jest",
            line: 30,
            name: "regular active test",
            status: "active",
          },
          {
            filePath: "src/__tests__/focused.test.ts",
            framework: "jest",
            line: 40,
            name: "skipped test",
            status: "skipped",
            modifier: "skip",
          },
        ],
      },
    ],
    summary: {
      active: 1,
      focused: 2,
      skipped: 1,
      todo: 0,
      total: 4,
      xfail: 0,
      frameworks: [
        {
          framework: "jest",
          active: 1,
          focused: 2,
          skipped: 1,
          todo: 0,
          total: 4,
          xfail: 0,
        },
      ],
    },
  },
};

// Analysis with xfail tests (for conditional display testing)
export const mockAnalysisWithXfail: AnalysisCompletedResponse = {
  status: "completed",
  data: {
    id: "xfail-analysis-id",
    analyzedAt: now,
    branchName: "main",
    commitSha: "xfail456",
    committedAt: commitTime,
    isInMyHistory: true,
    owner: "test-owner",
    repo: "xfail-test-repo",
    parserVersion: "v1.0.0",
    suites: [
      {
        filePath: "tests/test_xfail.py",
        framework: "pytest",
        suiteName: "Xfail Test Suite",
        tests: [
          {
            filePath: "tests/test_xfail.py",
            framework: "pytest",
            line: 10,
            name: "test_expected_failure",
            status: "xfail",
            modifier: "xfail",
          },
          {
            filePath: "tests/test_xfail.py",
            framework: "pytest",
            line: 20,
            name: "test_another_xfail",
            status: "xfail",
            modifier: "xfail",
          },
          {
            filePath: "tests/test_xfail.py",
            framework: "pytest",
            line: 30,
            name: "test_working",
            status: "active",
          },
          {
            filePath: "tests/test_xfail.py",
            framework: "pytest",
            line: 40,
            name: "test_skipped",
            status: "skipped",
            modifier: "skip",
          },
        ],
      },
    ],
    summary: {
      active: 1,
      focused: 0,
      skipped: 1,
      todo: 0,
      total: 4,
      xfail: 2,
      frameworks: [
        {
          framework: "pytest",
          active: 1,
          focused: 0,
          skipped: 1,
          todo: 0,
          total: 4,
          xfail: 2,
        },
      ],
    },
  },
};

// Analysis with all 5 statuses (comprehensive testing)
export const mockAnalysisWithAllStatuses: AnalysisCompletedResponse = {
  status: "completed",
  data: {
    id: "all-status-analysis-id",
    analyzedAt: now,
    branchName: "develop",
    commitSha: "all789",
    committedAt: commitTime,
    isInMyHistory: true,
    owner: "test-owner",
    repo: "all-status-repo",
    parserVersion: "v1.0.0",
    suites: [
      {
        filePath: "src/__tests__/comprehensive.test.ts",
        framework: "vitest",
        suiteName: "Comprehensive Status Suite",
        tests: [
          {
            filePath: "src/__tests__/comprehensive.test.ts",
            framework: "vitest",
            line: 10,
            name: "active test 1",
            status: "active",
          },
          {
            filePath: "src/__tests__/comprehensive.test.ts",
            framework: "vitest",
            line: 20,
            name: "active test 2",
            status: "active",
          },
          {
            filePath: "src/__tests__/comprehensive.test.ts",
            framework: "vitest",
            line: 30,
            name: "focused test",
            status: "focused",
            modifier: "only",
          },
          {
            filePath: "src/__tests__/comprehensive.test.ts",
            framework: "vitest",
            line: 40,
            name: "skipped test",
            status: "skipped",
            modifier: "skip",
          },
          {
            filePath: "src/__tests__/comprehensive.test.ts",
            framework: "vitest",
            line: 50,
            name: "xfail test",
            status: "xfail",
            modifier: "xfail",
          },
          {
            filePath: "src/__tests__/comprehensive.test.ts",
            framework: "vitest",
            line: 60,
            name: "todo test",
            status: "todo",
            modifier: "todo",
          },
        ],
      },
    ],
    summary: {
      active: 2,
      focused: 1,
      skipped: 1,
      todo: 1,
      total: 6,
      xfail: 1,
      frameworks: [
        {
          framework: "vitest",
          active: 2,
          focused: 1,
          skipped: 1,
          todo: 1,
          total: 6,
          xfail: 1,
        },
      ],
    },
  },
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

// Version history mock data
const versionHistoryCreatedAt = new Date().toISOString();

export const mockVersionHistoryMultiple: VersionHistoryResponse = {
  data: [
    {
      version: 2,
      createdAt: versionHistoryCreatedAt,
      modelId: "gemini-2.0-flash",
    },
    {
      version: 1,
      createdAt: new Date(Date.now() - 86400000 * 2).toISOString(), // 2 days ago
      modelId: "gemini-2.0-flash",
    },
  ],
  language: "English",
  latestVersion: 2,
};

export const mockVersionHistorySingle: VersionHistoryResponse = {
  data: [
    {
      version: 1,
      createdAt: versionHistoryCreatedAt,
      modelId: "gemini-2.0-flash",
    },
  ],
  language: "English",
  latestVersion: 1,
};

// Spec document for version 1 (older version)
export const mockSpecDocumentVersion1: SpecDocumentResponse = {
  status: "completed",
  data: {
    id: "doc-123-v1",
    analysisId: sampleAnalysisId,
    language: "English",
    version: 1,
    createdAt: new Date(Date.now() - 86400000 * 2).toISOString(),
    executiveSummary: "Test Repository Specification v1 - Previous version of the document.",
    modelId: "gemini-2.0-flash",
    availableLanguages: [
      {
        language: "English",
        latestVersion: 2,
        createdAt: versionHistoryCreatedAt,
      },
      {
        language: "Korean",
        latestVersion: 1,
        createdAt: new Date(Date.now() - 86400000).toISOString(),
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
            description: "Validates user credentials",
            behaviors: [
              {
                id: "spec-1",
                description: "should authenticate user",
                status: "active",
              },
            ],
          },
        ],
      },
    ],
  },
};

// Repository-based spec document (completed) - includes commitSha
export const mockRepoSpecDocumentCompleted: RepoSpecDocumentResponse = {
  status: "completed",
  data: {
    id: "doc-123",
    analysisId: sampleAnalysisId,
    language: "English",
    version: 2,
    createdAt: now,
    commitSha: "abc123def",
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
        createdAt: new Date(Date.now() - 86400000).toISOString(),
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

// Repository-based spec document (empty) - no document exists
export const mockRepoSpecDocumentEmpty: RepoSpecDocumentResponse = {
  status: "empty",
  message: "No spec document found for this repository",
};

// Repository-based spec document for version 1
export const mockRepoSpecDocumentVersion1: RepoSpecDocumentResponse = {
  status: "completed",
  data: {
    id: "doc-123-v1",
    analysisId: sampleAnalysisId,
    language: "English",
    version: 1,
    createdAt: new Date(Date.now() - 86400000 * 2).toISOString(),
    commitSha: "abc123def",
    executiveSummary: "Test Repository Specification v1 - Previous version of the document.",
    modelId: "gemini-2.0-flash",
    availableLanguages: [
      {
        language: "English",
        latestVersion: 2,
        createdAt: versionHistoryCreatedAt,
      },
      {
        language: "Korean",
        latestVersion: 1,
        createdAt: new Date(Date.now() - 86400000).toISOString(),
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
            description: "Validates user credentials",
            behaviors: [
              {
                id: "spec-1",
                description: "should authenticate user",
                status: "active",
              },
            ],
          },
        ],
      },
    ],
  },
};

// Repository version history mock data
export const mockRepoVersionHistoryMultiple: RepoVersionHistoryResponse = {
  data: [
    {
      id: "doc-123",
      analysisId: sampleAnalysisId,
      version: 2,
      language: "English",
      createdAt: versionHistoryCreatedAt,
      commitSha: "abc123def",
      modelId: "gemini-2.0-flash",
    },
    {
      id: "doc-123-v1",
      analysisId: sampleAnalysisId,
      version: 1,
      language: "English",
      createdAt: new Date(Date.now() - 86400000 * 2).toISOString(),
      commitSha: "older123",
      modelId: "gemini-2.0-flash",
    },
  ],
  language: "English",
};
