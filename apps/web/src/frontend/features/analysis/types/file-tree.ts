import type { TestSuite } from "@/lib/api";

export type FileTreeNodeType = "directory" | "file";

export type FileTreeNode = {
  children: FileTreeNode[];
  name: string;
  path: string;
  suites: TestSuite[];
  testCount: number;
  type: FileTreeNodeType;
};

export type FlatTreeItem = {
  depth: number;
  hasChildren: boolean;
  isExpanded: boolean;
  node: FileTreeNode;
};
