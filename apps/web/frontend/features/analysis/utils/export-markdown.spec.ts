import { describe, expect, it } from "vitest";

import type { components } from "@/lib/api/generated-types";

import { exportToMarkdown } from "./export-markdown";

type AnalysisResult = components["schemas"]["AnalysisResult"];

const createMockAnalysisResult = (overrides: Partial<AnalysisResult> = {}): AnalysisResult => ({
  analyzedAt: "2024-01-15T10:30:00Z",
  branchName: "main",
  commitSha: "abc123def456",
  committedAt: "2024-01-14T09:00:00Z",
  id: "550e8400-e29b-41d4-a716-446655440000",
  owner: "facebook",
  repo: "react",
  suites: [
    {
      filePath: "src/__tests__/App.test.tsx",
      framework: "vitest",
      suiteName: "App Component",
      tests: [
        {
          filePath: "src/__tests__/App.test.tsx",
          framework: "vitest",
          line: 10,
          name: "should render correctly",
          status: "active",
        },
        {
          filePath: "src/__tests__/App.test.tsx",
          framework: "vitest",
          line: 20,
          name: "should handle click",
          status: "skipped",
        },
        {
          filePath: "src/__tests__/App.test.tsx",
          framework: "vitest",
          line: 30,
          name: "should validate input",
          status: "todo",
        },
      ],
    },
    {
      filePath: "src/__tests__/utils.test.ts",
      framework: "jest",
      suiteName: "Utility Functions",
      tests: [
        {
          filePath: "src/__tests__/utils.test.ts",
          framework: "jest",
          line: 5,
          name: "formatDate works correctly",
          status: "active",
        },
      ],
    },
  ],
  summary: {
    active: 2,
    focused: 0,
    frameworks: [
      {
        active: 1,
        focused: 0,
        framework: "vitest",
        skipped: 1,
        todo: 1,
        total: 3,
        xfail: 0,
      },
      {
        active: 1,
        focused: 0,
        framework: "jest",
        skipped: 0,
        todo: 0,
        total: 1,
        xfail: 0,
      },
    ],
    skipped: 1,
    todo: 1,
    total: 4,
    xfail: 0,
  },
  ...overrides,
});

describe("exportToMarkdown", () => {
  it("should generate valid markdown header", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("# facebook/react - Test Specification");
    expect(markdown).toContain("> Branch: main");
    expect(markdown).toContain("> Commit: abc123def456");
  });

  it("should include summary table", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("## 📊 Summary");
    expect(markdown).toContain("| Total | 4 |");
    expect(markdown).toContain("| Active | 2 |");
    expect(markdown).toContain("| Skipped | 1 |");
    expect(markdown).toContain("| Todo | 1 |");
  });

  it("should include framework breakdown table", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("### By Framework");
    expect(markdown).toContain("| vitest | 3 | 1 | 1 | 1 |");
    expect(markdown).toContain("| jest | 1 | 1 | 0 | 0 |");
  });

  it("should group suites by framework", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("## vitest");
    expect(markdown).toContain("## jest");
  });

  it("should format test cases with correct checkboxes", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data, { includeLineNumbers: false });

    expect(markdown).toContain("- [x] should render correctly");
    expect(markdown).toContain("- [ ] ~~should handle click~~");
    expect(markdown).toContain("- [ ] should validate input *(todo)*");
    expect(markdown).toContain("- [x] formatDate works correctly");
  });

  it("should include line numbers by default", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("(L10)");
    expect(markdown).toContain("(L20)");
    expect(markdown).toContain("(L30)");
  });

  it("should exclude line numbers when option is false", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data, { includeLineNumbers: false });

    expect(markdown).not.toContain("(L10)");
    expect(markdown).not.toContain("(L20)");
  });

  it("should include timestamps by default", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("> Generated:");
    expect(markdown).toContain("> Analyzed: 2024-01-15T10:30:00Z");
  });

  it("should exclude timestamps when option is false", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data, { includeTimestamp: false });

    expect(markdown).not.toContain("> Generated:");
    expect(markdown).not.toContain("> Analyzed:");
  });

  it("should handle focused tests", () => {
    const data = createMockAnalysisResult({
      suites: [
        {
          filePath: "test.ts",
          framework: "vitest",
          suiteName: "Focus Test",
          tests: [
            {
              filePath: "test.ts",
              framework: "vitest",
              line: 1,
              name: "focused test",
              status: "focused",
            },
          ],
        },
      ],
      summary: {
        active: 0,
        focused: 1,
        frameworks: [
          {
            active: 0,
            focused: 1,
            framework: "vitest",
            skipped: 0,
            todo: 0,
            total: 1,
            xfail: 0,
          },
        ],
        skipped: 0,
        todo: 0,
        total: 1,
        xfail: 0,
      },
    });
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("- [!] focused test *(focused)*");
    expect(markdown).toContain("| Focused | 1 |");
  });

  it("should handle xfail tests", () => {
    const data = createMockAnalysisResult({
      suites: [
        {
          filePath: "test.ts",
          framework: "pytest",
          suiteName: "XFail Test",
          tests: [
            {
              filePath: "test.ts",
              framework: "pytest",
              line: 1,
              name: "expected failure",
              status: "xfail",
            },
          ],
        },
      ],
      summary: {
        active: 0,
        focused: 0,
        frameworks: [
          {
            active: 0,
            focused: 0,
            framework: "pytest",
            skipped: 0,
            todo: 0,
            total: 1,
            xfail: 1,
          },
        ],
        skipped: 0,
        todo: 0,
        total: 1,
        xfail: 1,
      },
    });
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("- [?] expected failure *(expected fail)*");
    expect(markdown).toContain("| Expected Fail | 1 |");
  });

  it("should include file path and framework info for each suite", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("📁 `src/__tests__/App.test.tsx`");
    expect(markdown).toContain("🔧 vitest");
    expect(markdown).toContain("📁 `src/__tests__/utils.test.ts`");
    expect(markdown).toContain("🔧 jest");
  });

  it("should include footer with SpecVital link", () => {
    const data = createMockAnalysisResult();
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain(
      "*Exported from [SpecVital](https://specvital.com/analyze/facebook/react)*"
    );
  });

  it("should handle missing branch name", () => {
    const data = createMockAnalysisResult({ branchName: undefined });
    const markdown = exportToMarkdown(data);

    expect(markdown).not.toContain("> Branch:");
  });

  it("should handle empty suites", () => {
    const data = createMockAnalysisResult({
      suites: [],
      summary: {
        active: 0,
        focused: 0,
        frameworks: [],
        skipped: 0,
        todo: 0,
        total: 0,
        xfail: 0,
      },
    });
    const markdown = exportToMarkdown(data);

    expect(markdown).toContain("| Total | 0 |");
    expect(markdown).toContain("## 📋 Test Suites");
  });
});
