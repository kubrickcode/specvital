import type { FileTreeNode, FlatTreeItem } from "../types";

const flattenNodes = (
  nodes: FileTreeNode[],
  expandedPaths: Set<string>,
  depth: number,
  result: FlatTreeItem[]
): void => {
  for (const node of nodes) {
    const hasChildren = node.children.length > 0;
    const isExpanded = expandedPaths.has(node.path);

    result.push({
      depth,
      hasChildren,
      isExpanded,
      node,
    });

    if (hasChildren && isExpanded) {
      flattenNodes(node.children, expandedPaths, depth + 1, result);
    }
  }
};

export const flattenTree = (nodes: FileTreeNode[], expandedPaths: Set<string>): FlatTreeItem[] => {
  const result: FlatTreeItem[] = [];
  flattenNodes(nodes, expandedPaths, 0, result);
  return result;
};
