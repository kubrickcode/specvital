import { describe, expect, it } from "vitest";

import type { TestCase, TestSuite } from "@/lib/api";

import { buildFileTree } from "./build-file-tree";

const createTestCase = (overrides: Partial<TestCase> = {}): TestCase => ({
  filePath: "src/auth/login.spec.ts",
  framework: "vitest",
  line: 10,
  name: "should authenticate user",
  status: "active",
  ...overrides,
});

const createTestSuite = (overrides: Partial<TestSuite> = {}): TestSuite => ({
  filePath: "src/auth/login.spec.ts",
  framework: "vitest",
  suiteName: "LoginService",
  tests: [createTestCase()],
  ...overrides,
});

describe("buildFileTree", () => {
  describe("empty input", () => {
    it("should return empty array for empty suites", () => {
      const result = buildFileTree([]);
      expect(result).toEqual([]);
    });
  });

  describe("single file", () => {
    it("should create directory nodes for file path", () => {
      const suites = [createTestSuite({ filePath: "src/auth/login.spec.ts" })];
      const result = buildFileTree(suites);

      expect(result).toHaveLength(1);
      expect(result[0]?.name).toBe("src");
      expect(result[0]?.type).toBe("directory");
      expect(result[0]?.path).toBe("src");
    });

    it("should create nested directory structure", () => {
      const suites = [createTestSuite({ filePath: "src/auth/login.spec.ts" })];
      const result = buildFileTree(suites);

      expect(result[0]?.children).toHaveLength(1);
      expect(result[0]?.children[0]?.name).toBe("auth");
      expect(result[0]?.children[0]?.type).toBe("directory");
      expect(result[0]?.children[0]?.path).toBe("src/auth");
    });

    it("should create file node at leaf level", () => {
      const suites = [
        createTestSuite({
          filePath: "src/auth/login.spec.ts",
          tests: [createTestCase(), createTestCase({ line: 20 })],
        }),
      ];
      const result = buildFileTree(suites);

      const fileNode = result[0]?.children[0]?.children[0];
      expect(fileNode?.name).toBe("login.spec.ts");
      expect(fileNode?.type).toBe("file");
      expect(fileNode?.path).toBe("src/auth/login.spec.ts");
      expect(fileNode?.testCount).toBe(2);
      expect(fileNode?.suites).toHaveLength(1);
    });
  });

  describe("test counts", () => {
    it("should aggregate test counts in parent directories", () => {
      const suites = [
        createTestSuite({
          filePath: "src/auth/login.spec.ts",
          tests: [createTestCase(), createTestCase({ line: 20 })],
        }),
        createTestSuite({
          filePath: "src/auth/logout.spec.ts",
          tests: [createTestCase({ filePath: "src/auth/logout.spec.ts" })],
        }),
      ];
      const result = buildFileTree(suites);

      expect(result[0]?.testCount).toBe(3);
      expect(result[0]?.children[0]?.testCount).toBe(3);
    });

    it("should count tests across different directories", () => {
      const suites = [
        createTestSuite({
          filePath: "src/auth/login.spec.ts",
          tests: [createTestCase(), createTestCase({ line: 20 })],
        }),
        createTestSuite({
          filePath: "src/user/profile.spec.ts",
          tests: [createTestCase({ filePath: "src/user/profile.spec.ts" })],
        }),
      ];
      const result = buildFileTree(suites);

      expect(result[0]?.testCount).toBe(3);
    });
  });

  describe("sorting", () => {
    it("should sort directories before files", () => {
      const suites = [
        createTestSuite({ filePath: "src/a.spec.ts" }),
        createTestSuite({ filePath: "src/utils/helper.spec.ts" }),
        createTestSuite({ filePath: "src/b.spec.ts" }),
      ];
      const result = buildFileTree(suites);

      const srcChildren = result[0]?.children;
      expect(srcChildren).toHaveLength(3);
      expect(srcChildren?.[0]?.type).toBe("directory");
      expect(srcChildren?.[0]?.name).toBe("utils");
      expect(srcChildren?.[1]?.type).toBe("file");
      expect(srcChildren?.[1]?.name).toBe("a.spec.ts");
      expect(srcChildren?.[2]?.type).toBe("file");
      expect(srcChildren?.[2]?.name).toBe("b.spec.ts");
    });

    it("should sort directories alphabetically", () => {
      const suites = [
        createTestSuite({ filePath: "src/z-dir/test.spec.ts" }),
        createTestSuite({ filePath: "src/a-dir/test.spec.ts" }),
        createTestSuite({ filePath: "src/m-dir/test.spec.ts" }),
      ];
      const result = buildFileTree(suites);

      const srcChildren = result[0]?.children;
      expect(srcChildren?.[0]?.name).toBe("a-dir");
      expect(srcChildren?.[1]?.name).toBe("m-dir");
      expect(srcChildren?.[2]?.name).toBe("z-dir");
    });

    it("should sort files alphabetically", () => {
      const suites = [
        createTestSuite({ filePath: "src/z.spec.ts" }),
        createTestSuite({ filePath: "src/a.spec.ts" }),
        createTestSuite({ filePath: "src/m.spec.ts" }),
      ];
      const result = buildFileTree(suites);

      const srcChildren = result[0]?.children;
      expect(srcChildren?.[0]?.name).toBe("a.spec.ts");
      expect(srcChildren?.[1]?.name).toBe("m.spec.ts");
      expect(srcChildren?.[2]?.name).toBe("z.spec.ts");
    });
  });

  describe("multiple files", () => {
    it("should group files under shared directories", () => {
      const suites = [
        createTestSuite({ filePath: "src/auth/login.spec.ts" }),
        createTestSuite({ filePath: "src/auth/logout.spec.ts" }),
      ];
      const result = buildFileTree(suites);

      expect(result).toHaveLength(1);
      expect(result[0]?.children).toHaveLength(1);
      expect(result[0]?.children[0]?.children).toHaveLength(2);
    });

    it("should create separate branches for different root directories", () => {
      const suites = [
        createTestSuite({ filePath: "src/auth/login.spec.ts" }),
        createTestSuite({ filePath: "tests/e2e/home.spec.ts" }),
      ];
      const result = buildFileTree(suites);

      expect(result).toHaveLength(2);
      expect(result[0]?.name).toBe("src");
      expect(result[1]?.name).toBe("tests");
    });
  });

  describe("suite merging", () => {
    it("should merge suites with the same file path", () => {
      const suites = [
        createTestSuite({
          filePath: "src/auth/login.spec.ts",
          suiteName: "LoginService",
          tests: [createTestCase({ line: 10 }), createTestCase({ line: 20 })],
        }),
        createTestSuite({
          filePath: "src/auth/login.spec.ts",
          suiteName: "AuthService",
          tests: [createTestCase({ line: 30 })],
        }),
      ];
      const result = buildFileTree(suites);

      const fileNode = result[0]?.children[0]?.children[0];
      expect(fileNode?.suites).toHaveLength(2);
      expect(fileNode?.testCount).toBe(3);
    });

    it("should not merge suites with different file paths", () => {
      const suites = [
        createTestSuite({ filePath: "src/auth/login.spec.ts", suiteName: "Login" }),
        createTestSuite({ filePath: "src/auth/logout.spec.ts", suiteName: "Logout" }),
      ];
      const result = buildFileTree(suites);

      const authDir = result[0]?.children[0];
      expect(authDir?.children).toHaveLength(2);
      expect(authDir?.children[0]?.suites).toHaveLength(1);
      expect(authDir?.children[1]?.suites).toHaveLength(1);
    });
  });

  describe("edge cases", () => {
    it("should handle root-level files", () => {
      const suites = [createTestSuite({ filePath: "test.spec.ts" })];
      const result = buildFileTree(suites);

      expect(result).toHaveLength(1);
      expect(result[0]?.type).toBe("file");
      expect(result[0]?.name).toBe("test.spec.ts");
    });

    it("should handle deeply nested paths", () => {
      const suites = [
        createTestSuite({ filePath: "src/features/auth/components/__tests__/login.spec.ts" }),
      ];
      const result = buildFileTree(suites);

      let node = result[0];
      const expectedPath = ["src", "features", "auth", "components", "__tests__"];
      for (const part of expectedPath) {
        expect(node?.name).toBe(part);
        expect(node?.type).toBe("directory");
        node = node?.children[0];
      }
      expect(node?.name).toBe("login.spec.ts");
      expect(node?.type).toBe("file");
    });

    it("should handle files with same name in different directories", () => {
      const suites = [
        createTestSuite({ filePath: "src/auth/index.spec.ts" }),
        createTestSuite({ filePath: "src/user/index.spec.ts" }),
      ];
      const result = buildFileTree(suites);

      expect(result[0]?.children).toHaveLength(2);
      expect(result[0]?.children[0]?.name).toBe("auth");
      expect(result[0]?.children[1]?.name).toBe("user");
      expect(result[0]?.children[0]?.children[0]?.name).toBe("index.spec.ts");
      expect(result[0]?.children[1]?.children[0]?.name).toBe("index.spec.ts");
    });

    it("should handle suite with zero tests", () => {
      const suites = [createTestSuite({ filePath: "src/empty.spec.ts", tests: [] })];
      const result = buildFileTree(suites);

      expect(result[0]?.children[0]?.testCount).toBe(0);
    });
  });

  describe("malformed paths", () => {
    it("should handle path with leading slash", () => {
      const suites = [createTestSuite({ filePath: "/src/test.spec.ts" })];
      const result = buildFileTree(suites);

      expect(result).toHaveLength(1);
      expect(result[0]?.name).toBe("src");
      expect(result[0]?.children[0]?.name).toBe("test.spec.ts");
    });

    it("should handle path with trailing slash", () => {
      const suites = [createTestSuite({ filePath: "src/test.spec.ts/" })];
      const result = buildFileTree(suites);

      expect(result).toHaveLength(1);
      expect(result[0]?.name).toBe("src");
      expect(result[0]?.children[0]?.name).toBe("test.spec.ts");
    });

    it("should handle path with consecutive slashes", () => {
      const suites = [createTestSuite({ filePath: "src//auth//test.spec.ts" })];
      const result = buildFileTree(suites);

      expect(result).toHaveLength(1);
      expect(result[0]?.name).toBe("src");
      expect(result[0]?.children[0]?.name).toBe("auth");
      expect(result[0]?.children[0]?.children[0]?.name).toBe("test.spec.ts");
    });

    it("should handle empty file path", () => {
      const suites = [createTestSuite({ filePath: "" })];
      const result = buildFileTree(suites);

      expect(result).toHaveLength(1);
      expect(result[0]?.type).toBe("file");
      expect(result[0]?.name).toBe("");
    });
  });
});
