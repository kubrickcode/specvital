import { describe, expect, it, vi } from "vitest";

import type { RepoSpecDocument, SpecDocument } from "../types";
import {
  copySpecToClipboard,
  downloadSpecMarkdown,
  exportSpecToMarkdown,
  generateSpecFilename,
} from "./export-spec-markdown";

const createMockDocument = (overrides: Partial<SpecDocument> = {}): SpecDocument => ({
  analysisId: "analysis-1",
  createdAt: "2024-01-20T10:30:00Z",
  domains: [
    {
      classificationConfidence: 0.92,
      description: "User authentication and authorization",
      features: [
        {
          behaviors: [
            {
              convertedDescription: "Authenticates user with valid credentials",
              id: "behavior-1",
              originalName: "should authenticate user with valid credentials",
              sortOrder: 0,
              sourceInfo: {
                filePath: "src/auth/login.spec.ts",
                framework: "jest",
                lineNumber: 45,
                status: "active",
              },
            },
            {
              convertedDescription: "Handles expired session gracefully",
              id: "behavior-2",
              originalName: "should handle expired session",
              sortOrder: 1,
              sourceInfo: {
                filePath: "src/auth/session.spec.ts",
                framework: "jest",
                lineNumber: 23,
                status: "skipped",
              },
            },
          ],
          description: "User login with various methods",
          id: "feature-1",
          name: "Login Flow",
          sortOrder: 0,
        },
      ],
      id: "domain-1",
      name: "Authentication",
      sortOrder: 0,
    },
  ],
  executiveSummary: "This is the executive summary.",
  id: "doc-1",
  language: "English",
  modelId: "gemini-2.0-flash",
  version: 1,
  ...overrides,
});

describe("exportSpecToMarkdown", () => {
  it("generates markdown with all sections", () => {
    const document = createMockDocument();
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("# acme/app - Specification Document");
    expect(result).toContain("## Executive Summary");
    expect(result).toContain("This is the executive summary.");
    expect(result).toContain("## 📊 Overview");
    expect(result).toContain("| 1 | 1 | 2 |");
    expect(result).toContain("## Authentication");
    expect(result).toContain("### Login Flow");
    expect(result).toContain("#### Behaviors");
    expect(result).toContain("✅ Authenticates user with valid credentials");
    expect(result).toContain("⏭️ ~~Handles expired session gracefully~~");
  });

  it("includes metadata in header", () => {
    const document = createMockDocument();
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("> Language: English");
    expect(result).toContain("> Version: 1");
    expect(result).toContain("> Generated: 2024-01-20T10:30:00Z");
    expect(result).toContain("> Model: gemini-2.0-flash");
  });

  it("includes source info by default", () => {
    const document = createMockDocument();
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("- Source: `src/auth/login.spec.ts:45` (jest)");
    expect(result).toContain("- Source: `src/auth/session.spec.ts:23` (jest)");
  });

  it("excludes source info when option is false", () => {
    const document = createMockDocument();
    const result = exportSpecToMarkdown(
      document,
      { owner: "acme", repo: "app" },
      { includeSourceInfo: false }
    );

    expect(result).not.toContain("- Source:");
  });

  it("excludes timestamp when option is false", () => {
    const document = createMockDocument();
    const result = exportSpecToMarkdown(
      document,
      { owner: "acme", repo: "app" },
      { includeTimestamp: false }
    );

    expect(result).not.toContain("> Language:");
    expect(result).not.toContain("> Version:");
    expect(result).not.toContain("> Generated:");
  });

  it("handles document without executive summary", () => {
    const document = createMockDocument({ executiveSummary: undefined });
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).not.toContain("## Executive Summary");
  });

  it("handles domain without confidence score", () => {
    const document = createMockDocument({
      domains: [
        {
          features: [],
          id: "domain-1",
          name: "Core",
          sortOrder: 0,
        },
      ],
    });
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("## Core");
    expect(result).not.toContain("Classification Confidence");
  });

  it("renders all test status icons correctly", () => {
    const document = createMockDocument({
      domains: [
        {
          features: [
            {
              behaviors: [
                {
                  convertedDescription: "Active behavior",
                  id: "b1",
                  originalName: "active test",
                  sortOrder: 0,
                  sourceInfo: {
                    filePath: "test.ts",
                    framework: "jest",
                    lineNumber: 1,
                    status: "active",
                  },
                },
                {
                  convertedDescription: "Focused behavior",
                  id: "b2",
                  originalName: "focused test",
                  sortOrder: 1,
                  sourceInfo: {
                    filePath: "test.ts",
                    framework: "jest",
                    lineNumber: 2,
                    status: "focused",
                  },
                },
                {
                  convertedDescription: "Skipped behavior",
                  id: "b3",
                  originalName: "skipped test",
                  sortOrder: 2,
                  sourceInfo: {
                    filePath: "test.ts",
                    framework: "jest",
                    lineNumber: 3,
                    status: "skipped",
                  },
                },
                {
                  convertedDescription: "Todo behavior",
                  id: "b4",
                  originalName: "todo test",
                  sortOrder: 3,
                  sourceInfo: {
                    filePath: "test.ts",
                    framework: "jest",
                    lineNumber: 4,
                    status: "todo",
                  },
                },
                {
                  convertedDescription: "Xfail behavior",
                  id: "b5",
                  originalName: "xfail test",
                  sortOrder: 4,
                  sourceInfo: {
                    filePath: "test.ts",
                    framework: "jest",
                    lineNumber: 5,
                    status: "xfail",
                  },
                },
              ],
              id: "feature-1",
              name: "All Statuses",
              sortOrder: 0,
            },
          ],
          id: "domain-1",
          name: "Test Statuses",
          sortOrder: 0,
        },
      ],
    });
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("✅ Active behavior");
    expect(result).toContain("🎯 Focused behavior");
    expect(result).toContain("⏭️ ~~Skipped behavior~~");
    expect(result).toContain("📝 Todo behavior");
    expect(result).toContain("❌ Xfail behavior");
  });

  it("includes footer with link", () => {
    const document = createMockDocument();
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("*Exported from [SpecVital](https://specvital.com/acme/app)*");
  });

  it("calculates stats correctly for multiple domains", () => {
    const document = createMockDocument({
      domains: [
        {
          features: [
            {
              behaviors: [
                { convertedDescription: "b1", id: "b1", originalName: "t1", sortOrder: 0 },
                { convertedDescription: "b2", id: "b2", originalName: "t2", sortOrder: 1 },
              ],
              id: "f1",
              name: "Feature 1",
              sortOrder: 0,
            },
            {
              behaviors: [
                { convertedDescription: "b3", id: "b3", originalName: "t3", sortOrder: 0 },
              ],
              id: "f2",
              name: "Feature 2",
              sortOrder: 1,
            },
          ],
          id: "d1",
          name: "Domain 1",
          sortOrder: 0,
        },
        {
          features: [
            {
              behaviors: [
                { convertedDescription: "b4", id: "b4", originalName: "t4", sortOrder: 0 },
                { convertedDescription: "b5", id: "b5", originalName: "t5", sortOrder: 1 },
              ],
              id: "f3",
              name: "Feature 3",
              sortOrder: 0,
            },
          ],
          id: "d2",
          name: "Domain 2",
          sortOrder: 1,
        },
      ],
    });
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("| 2 | 3 | 5 |");
  });

  it("includes commitSha for RepoSpecDocument", () => {
    const document: RepoSpecDocument = {
      ...createMockDocument(),
      commitSha: "abc123def456",
    };
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("> Commit: abc123def456");
  });

  it("handles feature with empty behaviors array", () => {
    const document = createMockDocument({
      domains: [
        {
          features: [
            {
              behaviors: [],
              id: "f1",
              name: "Empty Feature",
              sortOrder: 0,
            },
          ],
          id: "d1",
          name: "Empty Domain",
          sortOrder: 0,
        },
      ],
    });
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("### Empty Feature");
    expect(result).toContain("#### Behaviors");
    expect(result).toContain("| 1 | 1 | 0 |");
  });

  it("handles behavior without sourceInfo", () => {
    const document = createMockDocument({
      domains: [
        {
          features: [
            {
              behaviors: [
                {
                  convertedDescription: "No source behavior",
                  id: "b1",
                  originalName: "test",
                  sortOrder: 0,
                },
              ],
              id: "f1",
              name: "Feature",
              sortOrder: 0,
            },
          ],
          id: "d1",
          name: "Domain",
          sortOrder: 0,
        },
      ],
    });
    const result = exportSpecToMarkdown(document, { owner: "acme", repo: "app" });

    expect(result).toContain("No source behavior");
    expect(result).not.toContain("- Source:");
  });
});

describe("generateSpecFilename", () => {
  it("generates filename with language suffix", () => {
    expect(generateSpecFilename("acme", "app", "English")).toBe("acme-app-spec-English.md");
    expect(generateSpecFilename("acme", "app", "Korean")).toBe("acme-app-spec-Korean.md");
  });
});

describe("downloadSpecMarkdown", () => {
  it("creates blob and triggers download", () => {
    const mockUrl = "blob:mock-url";
    const createObjectURLMock = vi.fn(() => mockUrl);
    const revokeObjectURLMock = vi.fn();
    const clickMock = vi.fn();
    const createElementMock = vi.fn(() => ({
      click: clickMock,
      download: "",
      href: "",
    }));

    vi.stubGlobal("URL", {
      createObjectURL: createObjectURLMock,
      revokeObjectURL: revokeObjectURLMock,
    });
    vi.spyOn(document, "createElement").mockImplementation(
      createElementMock as unknown as typeof document.createElement
    );

    downloadSpecMarkdown("# Test content", "test.md");

    expect(createObjectURLMock).toHaveBeenCalledWith(expect.any(Blob));
    expect(createElementMock).toHaveBeenCalledWith("a");
    expect(clickMock).toHaveBeenCalled();
    expect(revokeObjectURLMock).toHaveBeenCalledWith(mockUrl);

    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });
});

describe("copySpecToClipboard", () => {
  it("returns true on successful copy", async () => {
    const writeTextMock = vi.fn().mockResolvedValue(undefined);
    vi.stubGlobal("navigator", {
      clipboard: { writeText: writeTextMock },
    });

    const result = await copySpecToClipboard("test content");

    expect(result).toBe(true);
    expect(writeTextMock).toHaveBeenCalledWith("test content");

    vi.unstubAllGlobals();
  });

  it("returns false on copy failure", async () => {
    const writeTextMock = vi.fn().mockRejectedValue(new Error("Copy failed"));
    vi.stubGlobal("navigator", {
      clipboard: { writeText: writeTextMock },
    });

    const result = await copySpecToClipboard("test content");

    expect(result).toBe(false);

    vi.unstubAllGlobals();
  });
});
