import type { TestSuite } from "@/lib/api";

import type { FileTreeNode } from "../types";

const sortNodes = (nodes: FileTreeNode[]): FileTreeNode[] =>
  [...nodes].sort((a, b) => {
    if (a.type !== b.type) {
      return a.type === "directory" ? -1 : 1;
    }
    return a.name.localeCompare(b.name);
  });

const normalizePath = (filePath: string): string[] => filePath.split("/").filter(Boolean);

const getDirectoryPath = (filePath: string): string[] => {
  const parts = normalizePath(filePath);
  return parts.slice(0, -1);
};

const getFileName = (filePath: string): string => {
  const parts = normalizePath(filePath);
  return parts[parts.length - 1] ?? filePath;
};

const createFileNode = (suite: TestSuite): FileTreeNode => ({
  children: [],
  name: getFileName(suite.filePath),
  path: suite.filePath,
  suites: [suite],
  testCount: suite.tests.length,
  type: "file",
});

const ensureDirectoryPath = (root: FileTreeNode, pathParts: string[]): FileTreeNode => {
  let current = root;
  let currentPath = "";

  for (const part of pathParts) {
    currentPath = currentPath ? `${currentPath}/${part}` : part;
    let child = current.children.find((c) => c.type === "directory" && c.name === part);

    if (!child) {
      child = {
        children: [],
        name: part,
        path: currentPath,
        suites: [],
        testCount: 0,
        type: "directory",
      };
      current.children.push(child);
    }

    current = child;
  }

  return current;
};

const updateTestCounts = (node: FileTreeNode): number => {
  if (node.type === "file") {
    return node.testCount;
  }

  let count = 0;
  for (const child of node.children) {
    count += updateTestCounts(child);
  }
  node.testCount = count;
  return count;
};

const sortTreeRecursively = (node: FileTreeNode): void => {
  node.children = sortNodes(node.children);
  for (const child of node.children) {
    sortTreeRecursively(child);
  }
};

export const buildFileTree = (suites: TestSuite[]): FileTreeNode[] => {
  if (suites.length === 0) {
    return [];
  }

  const root: FileTreeNode = {
    children: [],
    name: "",
    path: "",
    suites: [],
    testCount: 0,
    type: "directory",
  };

  for (const suite of suites) {
    const dirParts = getDirectoryPath(suite.filePath);
    const parentDir = ensureDirectoryPath(root, dirParts);

    const existingFile = parentDir.children.find(
      (c) => c.type === "file" && c.path === suite.filePath
    );

    if (existingFile) {
      existingFile.suites.push(suite);
      existingFile.testCount += suite.tests.length;
    } else {
      parentDir.children.push(createFileNode(suite));
    }
  }

  updateTestCounts(root);
  sortTreeRecursively(root);

  return root.children;
};
