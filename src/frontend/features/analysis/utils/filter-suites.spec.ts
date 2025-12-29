import { describe, expect, it } from "vitest";

import type { TestCase, TestSuite } from "@/lib/api";

import type { FilterState } from "../types";
import { filterSuites } from "./filter-suites";

const createTestCase = (overrides: Partial<TestCase> = {}): TestCase => ({
  filePath: "src/auth/login.spec.ts",
  framework: "vitest",
  line: 10,
  name: "should authenticate user with valid credentials",
  status: "active",
  ...overrides,
});

const createTestSuite = (overrides: Partial<TestSuite> = {}): TestSuite => ({
  filePath: "src/auth/login.spec.ts",
  framework: "vitest",
  suiteName: "LoginService",
  tests: [
    createTestCase(),
    createTestCase({
      line: 20,
      name: "should reject invalid password",
    }),
  ],
  ...overrides,
});

const createFilter = (overrides: Partial<FilterState> = {}): FilterState => ({
  frameworks: [],
  query: "",
  statuses: [],
  ...overrides,
});

describe("filterSuites", () => {
  describe("empty filters", () => {
    it("should return all suites when all filters are empty", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, createFilter());
      expect(result).toEqual(suites);
    });

    it("should return all suites when query is whitespace only", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, createFilter({ query: "   " }));
      expect(result).toEqual(suites);
    });
  });

  describe("query filtering", () => {
    it("should filter by test name", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, createFilter({ query: "authenticate" }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("should authenticate user with valid credentials");
    });

    it("should be case insensitive", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, createFilter({ query: "AUTHENTICATE" }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
    });

    it("should match multiple words (AND logic)", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, createFilter({ query: "authenticate user" }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
    });

    it("should return empty when no match found", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, createFilter({ query: "nonexistent" }));

      expect(result).toHaveLength(0);
    });

    it("should match by file path", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, createFilter({ query: "auth/login" }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(2);
    });

    it("should match by suite name", () => {
      const suites = [createTestSuite()];
      const result = filterSuites(suites, createFilter({ query: "LoginService" }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(2);
    });
  });

  describe("status filtering", () => {
    it("should filter by single status", () => {
      const suites = [
        createTestSuite({
          tests: [
            createTestCase({ name: "active test", status: "active" }),
            createTestCase({ name: "skipped test", status: "skipped" }),
            createTestCase({ name: "todo test", status: "todo" }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ statuses: ["skipped"] }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("skipped test");
    });

    it("should filter by multiple statuses (OR logic within statuses)", () => {
      const suites = [
        createTestSuite({
          tests: [
            createTestCase({ name: "active test", status: "active" }),
            createTestCase({ name: "skipped test", status: "skipped" }),
            createTestCase({ name: "todo test", status: "todo" }),
            createTestCase({ name: "focused test", status: "focused" }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ statuses: ["skipped", "todo"] }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(2);
      expect(result[0]?.tests.map((t) => t.name)).toEqual(["skipped test", "todo test"]);
    });

    it("should return empty when no tests match the status", () => {
      const suites = [
        createTestSuite({
          tests: [createTestCase({ status: "active" }), createTestCase({ status: "active" })],
        }),
      ];

      const result = filterSuites(suites, createFilter({ statuses: ["xfail"] }));

      expect(result).toHaveLength(0);
    });

    it("should return all statuses when statuses array is empty", () => {
      const suites = [
        createTestSuite({
          tests: [createTestCase({ status: "active" }), createTestCase({ status: "skipped" })],
        }),
      ];

      const result = filterSuites(suites, createFilter({ statuses: [] }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(2);
    });

    it("should filter across multiple suites", () => {
      const suites = [
        createTestSuite({
          filePath: "suite1.spec.ts",
          tests: [
            createTestCase({ name: "active 1", status: "active" }),
            createTestCase({ name: "skipped 1", status: "skipped" }),
          ],
        }),
        createTestSuite({
          filePath: "suite2.spec.ts",
          tests: [
            createTestCase({ name: "todo 2", status: "todo" }),
            createTestCase({ name: "skipped 2", status: "skipped" }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ statuses: ["skipped"] }));

      expect(result).toHaveLength(2);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("skipped 1");
      expect(result[1]?.tests).toHaveLength(1);
      expect(result[1]?.tests[0]?.name).toBe("skipped 2");
    });
  });

  describe("framework filtering", () => {
    it("should filter by single framework", () => {
      const suites = [
        createTestSuite({
          tests: [
            createTestCase({ framework: "vitest", name: "vitest test" }),
            createTestCase({ framework: "jest", name: "jest test" }),
            createTestCase({ framework: "pytest", name: "pytest test" }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ frameworks: ["jest"] }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("jest test");
    });

    it("should filter by multiple frameworks (OR logic within frameworks)", () => {
      const suites = [
        createTestSuite({
          tests: [
            createTestCase({ framework: "vitest", name: "vitest test" }),
            createTestCase({ framework: "jest", name: "jest test" }),
            createTestCase({ framework: "pytest", name: "pytest test" }),
            createTestCase({ framework: "mocha", name: "mocha test" }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ frameworks: ["vitest", "jest"] }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(2);
      expect(result[0]?.tests.map((t) => t.name)).toEqual(["vitest test", "jest test"]);
    });

    it("should return empty when no tests match the framework", () => {
      const suites = [
        createTestSuite({
          tests: [createTestCase({ framework: "vitest" }), createTestCase({ framework: "vitest" })],
        }),
      ];

      const result = filterSuites(suites, createFilter({ frameworks: ["rspec"] }));

      expect(result).toHaveLength(0);
    });

    it("should return all frameworks when frameworks array is empty", () => {
      const suites = [
        createTestSuite({
          tests: [createTestCase({ framework: "vitest" }), createTestCase({ framework: "jest" })],
        }),
      ];

      const result = filterSuites(suites, createFilter({ frameworks: [] }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(2);
    });

    it("should filter across multiple suites", () => {
      const suites = [
        createTestSuite({
          filePath: "suite1.spec.ts",
          tests: [
            createTestCase({ framework: "vitest", name: "vitest 1" }),
            createTestCase({ framework: "jest", name: "jest 1" }),
          ],
        }),
        createTestSuite({
          filePath: "suite2.spec.ts",
          tests: [
            createTestCase({ framework: "pytest", name: "pytest 2" }),
            createTestCase({ framework: "vitest", name: "vitest 2" }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ frameworks: ["vitest"] }));

      expect(result).toHaveLength(2);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("vitest 1");
      expect(result[1]?.tests).toHaveLength(1);
      expect(result[1]?.tests[0]?.name).toBe("vitest 2");
    });
  });

  describe("combined filters", () => {
    it("should apply AND logic between query and status", () => {
      const suites = [
        createTestSuite({
          filePath: "src/user/user.spec.ts",
          tests: [
            createTestCase({
              filePath: "src/user/user.spec.ts",
              name: "skipped login test",
              status: "skipped",
            }),
            createTestCase({
              filePath: "src/user/user.spec.ts",
              name: "active login test",
              status: "active",
            }),
            createTestCase({
              filePath: "src/user/user.spec.ts",
              name: "skipped auth test",
              status: "skipped",
            }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ query: "login", statuses: ["skipped"] }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("skipped login test");
    });

    it("should apply AND logic between query and framework", () => {
      const suites = [
        createTestSuite({
          filePath: "src/user/user.spec.ts",
          tests: [
            createTestCase({
              filePath: "src/user/user.spec.ts",
              framework: "vitest",
              name: "vitest login test",
            }),
            createTestCase({
              filePath: "src/user/user.spec.ts",
              framework: "jest",
              name: "jest login test",
            }),
            createTestCase({
              filePath: "src/user/user.spec.ts",
              framework: "vitest",
              name: "vitest auth test",
            }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ frameworks: ["vitest"], query: "login" }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("vitest login test");
    });

    it("should apply AND logic between status and framework", () => {
      const suites = [
        createTestSuite({
          tests: [
            createTestCase({ framework: "vitest", name: "skipped vitest", status: "skipped" }),
            createTestCase({ framework: "vitest", name: "active vitest", status: "active" }),
            createTestCase({ framework: "jest", name: "skipped jest", status: "skipped" }),
            createTestCase({ framework: "jest", name: "active jest", status: "active" }),
          ],
        }),
      ];

      const result = filterSuites(
        suites,
        createFilter({ frameworks: ["vitest"], statuses: ["skipped"] })
      );

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("skipped vitest");
    });

    it("should apply AND logic between all three filters", () => {
      const suites = [
        createTestSuite({
          filePath: "src/user/user.spec.ts",
          tests: [
            createTestCase({
              filePath: "src/user/user.spec.ts",
              framework: "vitest",
              name: "skipped vitest login test",
              status: "skipped",
            }),
            createTestCase({
              filePath: "src/user/user.spec.ts",
              framework: "vitest",
              name: "skipped vitest auth test",
              status: "skipped",
            }),
            createTestCase({
              filePath: "src/user/user.spec.ts",
              framework: "vitest",
              name: "active vitest login test",
              status: "active",
            }),
            createTestCase({
              filePath: "src/user/user.spec.ts",
              framework: "jest",
              name: "skipped jest login test",
              status: "skipped",
            }),
          ],
        }),
      ];

      const result = filterSuites(
        suites,
        createFilter({
          frameworks: ["vitest"],
          query: "login",
          statuses: ["skipped"],
        })
      );

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(1);
      expect(result[0]?.tests[0]?.name).toBe("skipped vitest login test");
    });

    it("should return empty when combined filters match nothing", () => {
      const suites = [
        createTestSuite({
          tests: [
            createTestCase({ framework: "vitest", name: "test 1", status: "active" }),
            createTestCase({ framework: "jest", name: "test 2", status: "skipped" }),
          ],
        }),
      ];

      const result = filterSuites(
        suites,
        createFilter({ frameworks: ["vitest"], statuses: ["skipped"] })
      );

      expect(result).toHaveLength(0);
    });

    it("should not match file path when status or framework filters are active", () => {
      const suites = [
        createTestSuite({
          filePath: "src/auth/login.spec.ts",
          tests: [createTestCase({ status: "active" })],
        }),
      ];

      const result = filterSuites(suites, createFilter({ query: "auth", statuses: ["skipped"] }));

      expect(result).toHaveLength(0);
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
            createTestCase({
              filePath: "src/users/profile.spec.ts",
              line: 5,
              name: "should update user profile",
            }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ query: "profile" }));

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
            createTestCase({
              filePath: "src/auth/session.spec.ts",
              line: 5,
              name: "should create session",
            }),
          ],
        }),
      ];

      const result = filterSuites(suites, createFilter({ query: "auth" }));

      expect(result).toHaveLength(2);
    });
  });

  describe("edge cases", () => {
    it("should handle empty suites array", () => {
      const result = filterSuites([], createFilter({ query: "test" }));
      expect(result).toHaveLength(0);
    });

    it("should handle suite with no tests", () => {
      const suites = [createTestSuite({ tests: [] })];
      const result = filterSuites(suites, createFilter({ query: "login" }));

      expect(result).toHaveLength(1);
      expect(result[0]?.tests).toHaveLength(0);
    });

    it("should handle suite with empty suiteName", () => {
      const suites = [createTestSuite({ suiteName: "" })];
      const result = filterSuites(suites, createFilter({ query: "authenticate" }));

      expect(result).toHaveLength(1);
    });

    it("should exclude suite with empty tests when status filter is active", () => {
      const suites = [createTestSuite({ tests: [] })];
      const result = filterSuites(suites, createFilter({ statuses: ["skipped"] }));

      expect(result).toHaveLength(0);
    });

    it("should exclude suite with empty tests when framework filter is active", () => {
      const suites = [createTestSuite({ tests: [] })];
      const result = filterSuites(suites, createFilter({ frameworks: ["vitest"] }));

      expect(result).toHaveLength(0);
    });
  });
});
