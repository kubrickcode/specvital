import { describe, expect, it } from "vitest";

import type { FileTreeNode } from "../types";
import { flattenTree } from "./flatten-tree";

const createFileNode = (name: string, path: string): FileTreeNode => ({
  children: [],
  name,
  path,
  suites: [],
  testCount: 1,
  type: "file",
});

const createDirectoryNode = (
  name: string,
  path: string,
  children: FileTreeNode[] = []
): FileTreeNode => ({
  children,
  name,
  path,
  suites: [],
  testCount: children.reduce((sum, c) => sum + c.testCount, 0),
  type: "directory",
});

describe("flattenTree", () => {
  describe("empty input", () => {
    it("should return empty array for empty nodes", () => {
      const result = flattenTree([], new Set());
      expect(result).toEqual([]);
    });
  });

  describe("flat structure", () => {
    it("should flatten single file node", () => {
      const nodes = [createFileNode("test.spec.ts", "test.spec.ts")];
      const result = flattenTree(nodes, new Set());

      expect(result).toHaveLength(1);
      expect(result[0]?.node.name).toBe("test.spec.ts");
      expect(result[0]?.depth).toBe(0);
      expect(result[0]?.hasChildren).toBe(false);
      expect(result[0]?.isExpanded).toBe(false);
    });

    it("should flatten multiple file nodes at root", () => {
      const nodes = [
        createFileNode("a.spec.ts", "a.spec.ts"),
        createFileNode("b.spec.ts", "b.spec.ts"),
      ];
      const result = flattenTree(nodes, new Set());

      expect(result).toHaveLength(2);
      expect(result[0]?.node.name).toBe("a.spec.ts");
      expect(result[1]?.node.name).toBe("b.spec.ts");
    });
  });

  describe("collapsed directories", () => {
    it("should not include children of collapsed directory", () => {
      const nodes = [
        createDirectoryNode("src", "src", [createFileNode("test.spec.ts", "src/test.spec.ts")]),
      ];
      const result = flattenTree(nodes, new Set());

      expect(result).toHaveLength(1);
      expect(result[0]?.node.name).toBe("src");
      expect(result[0]?.hasChildren).toBe(true);
      expect(result[0]?.isExpanded).toBe(false);
    });
  });

  describe("expanded directories", () => {
    it("should include children of expanded directory", () => {
      const nodes = [
        createDirectoryNode("src", "src", [createFileNode("test.spec.ts", "src/test.spec.ts")]),
      ];
      const result = flattenTree(nodes, new Set(["src"]));

      expect(result).toHaveLength(2);
      expect(result[0]?.node.name).toBe("src");
      expect(result[0]?.isExpanded).toBe(true);
      expect(result[1]?.node.name).toBe("test.spec.ts");
      expect(result[1]?.depth).toBe(1);
    });

    it("should handle deeply nested expanded paths", () => {
      const nodes = [
        createDirectoryNode("src", "src", [
          createDirectoryNode("auth", "src/auth", [
            createFileNode("login.spec.ts", "src/auth/login.spec.ts"),
          ]),
        ]),
      ];
      const result = flattenTree(nodes, new Set(["src", "src/auth"]));

      expect(result).toHaveLength(3);
      expect(result[0]?.node.name).toBe("src");
      expect(result[0]?.depth).toBe(0);
      expect(result[1]?.node.name).toBe("auth");
      expect(result[1]?.depth).toBe(1);
      expect(result[2]?.node.name).toBe("login.spec.ts");
      expect(result[2]?.depth).toBe(2);
    });

    it("should stop at collapsed nested directory", () => {
      const nodes = [
        createDirectoryNode("src", "src", [
          createDirectoryNode("auth", "src/auth", [
            createFileNode("login.spec.ts", "src/auth/login.spec.ts"),
          ]),
        ]),
      ];
      const result = flattenTree(nodes, new Set(["src"]));

      expect(result).toHaveLength(2);
      expect(result[0]?.node.name).toBe("src");
      expect(result[1]?.node.name).toBe("auth");
      expect(result[1]?.hasChildren).toBe(true);
      expect(result[1]?.isExpanded).toBe(false);
    });
  });

  describe("mixed structure", () => {
    it("should handle mixed directories and files", () => {
      const nodes = [
        createDirectoryNode("src", "src", [
          createDirectoryNode("auth", "src/auth", [
            createFileNode("login.spec.ts", "src/auth/login.spec.ts"),
          ]),
          createFileNode("index.spec.ts", "src/index.spec.ts"),
        ]),
      ];
      const result = flattenTree(nodes, new Set(["src"]));

      expect(result).toHaveLength(3);
      expect(result[0]?.node.name).toBe("src");
      expect(result[1]?.node.name).toBe("auth");
      expect(result[2]?.node.name).toBe("index.spec.ts");
    });

    it("should handle multiple root nodes with different expansion states", () => {
      const nodes = [
        createDirectoryNode("src", "src", [createFileNode("a.spec.ts", "src/a.spec.ts")]),
        createDirectoryNode("tests", "tests", [createFileNode("b.spec.ts", "tests/b.spec.ts")]),
      ];
      const result = flattenTree(nodes, new Set(["src"]));

      expect(result).toHaveLength(3);
      expect(result[0]?.node.name).toBe("src");
      expect(result[0]?.isExpanded).toBe(true);
      expect(result[1]?.node.name).toBe("a.spec.ts");
      expect(result[2]?.node.name).toBe("tests");
      expect(result[2]?.isExpanded).toBe(false);
    });
  });

  describe("depth calculation", () => {
    it("should increment depth for each level", () => {
      const nodes = [
        createDirectoryNode("a", "a", [
          createDirectoryNode("b", "a/b", [
            createDirectoryNode("c", "a/b/c", [
              createFileNode("file.spec.ts", "a/b/c/file.spec.ts"),
            ]),
          ]),
        ]),
      ];
      const result = flattenTree(nodes, new Set(["a", "a/b", "a/b/c"]));

      expect(result).toHaveLength(4);
      expect(result[0]?.depth).toBe(0);
      expect(result[1]?.depth).toBe(1);
      expect(result[2]?.depth).toBe(2);
      expect(result[3]?.depth).toBe(3);
    });
  });

  describe("edge cases", () => {
    it("should handle directory with no children", () => {
      const nodes = [createDirectoryNode("empty", "empty", [])];
      const result = flattenTree(nodes, new Set(["empty"]));

      expect(result).toHaveLength(1);
      expect(result[0]?.hasChildren).toBe(false);
    });

    it("should handle path not in expanded set", () => {
      const nodes = [
        createDirectoryNode("src", "src", [createFileNode("test.spec.ts", "src/test.spec.ts")]),
      ];
      const result = flattenTree(nodes, new Set(["other-path"]));

      expect(result).toHaveLength(1);
      expect(result[0]?.isExpanded).toBe(false);
    });
  });
});
