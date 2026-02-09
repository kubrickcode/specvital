"use client";

import type { GitHubRepository, RepositoryCard } from "@/lib/api/types";

type UseUnanalyzedReposParams = {
  analyzedRepositories: RepositoryCard[];
  githubRepositories: GitHubRepository[];
};

type UseUnanalyzedReposReturn = {
  count: number;
  data: GitHubRepository[];
};

export const useUnanalyzedRepos = ({
  analyzedRepositories,
  githubRepositories,
}: UseUnanalyzedReposParams): UseUnanalyzedReposReturn => {
  const analyzedFullNames = new Set(analyzedRepositories.map((r) => r.fullName.toLowerCase()));
  const unanalyzedRepos = githubRepositories.filter(
    (repo) => !analyzedFullNames.has(repo.fullName.toLowerCase())
  );

  return {
    count: unanalyzedRepos.length,
    data: unanalyzedRepos,
  };
};
