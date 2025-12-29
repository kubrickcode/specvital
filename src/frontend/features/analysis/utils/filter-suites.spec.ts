import { describe, expect, it } from "vitest";

import type { TestSuite } from "@/lib/api";

import { filterSuites } from "./filter-suites";

const createTestSuite = (overrides: Partial<TestSuite> = {}): TestSuite => ({
  filePath: "src/auth/login.spec.ts",
  framework: "vitest",
  suiteName: "LoginService",
  tests: [
    {
      filePath: "src/auth/login.spec.ts",
      framework: "vitest",
      line: 10,
      name: "should authenticate user with valid credentials",
      status: "active",
    },
    {
      filePath: "src/auth/login.spec.ts",
      framework: "vitest",
      line: 20,
      name: "should reject invalid password",
      status: "active",
    },
  ],
  ...overrides,
});

describe("filterSuites", () => {
  describe("empty query", () => {
    it("should return all suites when query is empty", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, { query: "" });
      expect(result).toEqual(suites);
    });

    it("should return all suites when query is whitespace only", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, { query: "   " });
      expect(result).toEqual(suites);
    });
  });

  describe("test name matching", () => {
    it("should filter by test name", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, { query: "authenticate" });

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("should authenticate user with valid credentials");
    });

    it("should be case insensitive", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, { query: "AUTHENTICATE" });

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
    });

    it("should match multiple words (AND logic)", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, { query: "authenticate user" });

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
    });

    it("should return empty when no match found", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, { query: "nonexistent" });

      expect(result).toHaveLength(0);
    });
  });

  describe("file path matching", () => {
    it("should match by file path", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, { query: "auth/login" });

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(2);
    });

    it("should match partial file path", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, { query: "login.spec" });

      expect(result).toHaveLength(1);
    });
  });

  describe("suite name matching", () => {
    it("should match by suite name", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, { query: "LoginService" });

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(2);
    });
  });

  describe("multiple suites", () => {
    it("should filter across multiple suites", () => {
      const suites = [
        createTestSuite(),
        createTestSuite({
          filePath: "src/users/profile.spec.ts",
          suiteName: "ProfileService",
          tests: [
            {
              filePath: "src/users/profile.spec.ts",
              framework: "vitest",
              line: 5,
              name: "should update user profile",
              status: "active",
            },
          ],
        }),
      ];

      const result = filterSuites(suites, { query: "profile" });

      expect(result).toHaveLength(1);
      expect(result[0]?.filePath).toBe("src/users/profile.spec.ts");
    });

    it("should include tests from test content and file path matches separately", () => {
      const suites = [
        createTestSuite(),
        createTestSuite({
          filePath: "src/auth/session.spec.ts",
          suiteName: "SessionService",
          tests: [
            {
              filePath: "src/auth/session.spec.ts",
              framework: "vitest",
              line: 5,
              name: "should create session",
              status: "active",
            },
          ],
        }),
      ];

      const result = filterSuites(suites, { query: "auth" });

      expect(result).toHaveLength(2);
    });
  });

  describe("edge cases", () => {
    it("should handle empty suites array", () => {
      const result = filterSuites([], { query: "test" });
      expect(result).toHaveLength(0);
    });

    it("should handle suite with no tests", () => {
      const suites = [createTestSuite({ tests: [] })];
      const result = filterSuites(suites, { query: "login" });

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(0);
    });

    it("should handle suite with empty suiteName", () => {
      const suites = [createTestSuite({ suiteName: "" })];
      const result = filterSuites(suites, { query: "authenticate" });

      expect(result).toHaveLength(1);
    });
  });
});
