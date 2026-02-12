import type { FuseResult, FuseResultMatch } from "fuse.js";

import type { RepositorySearchItem } from "./fuse-config";

type RepoFuseResult = FuseResult<RepositorySearchItem>;

const CATEGORY_WEIGHTS: Record<RepositorySearchItem["category"], number> = {
  bookmarks: 2,
  community: 3,
  repositories: 1,
};

export type ScoredResult = {
  categoryRank: number;
  combinedScore: number;
  item: RepositorySearchItem;
  matches?: readonly FuseResultMatch[];
};

export const scoreAndSortResults = (results: RepoFuseResult[]): ScoredResult[] => {
  const scored = results.map((result) => {
    const fuseScore = result.score ?? 0;
    const categoryWeight = CATEGORY_WEIGHTS[result.item.category];

    return {
      categoryRank: categoryWeight,
      combinedScore: fuseScore + categoryWeight * 0.1,
      item: result.item,
      matches: result.matches,
    };
  });

  return scored.sort((a, b) => {
    if (a.categoryRank !== b.categoryRank) {
      return a.categoryRank - b.categoryRank;
    }
    return a.combinedScore - b.combinedScore;
  });
};

export const groupResultsByCategory = (
  results: ScoredResult[]
): Map<RepositorySearchItem["category"], ScoredResult[]> => {
  const grouped = new Map<RepositorySearchItem["category"], ScoredResult[]>();

  for (const result of results) {
    const existing = grouped.get(result.item.category) ?? [];
    existing.push(result);
    grouped.set(result.item.category, existing);
  }

  return grouped;
};
