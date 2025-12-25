"use client";

import { useMemo } from "react";

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
  const unanalyzedRepos = useMemo(() => {
    const analyzedFullNames = new Set(analyzedRepositories.map((r) => r.fullName.toLowerCase()));

    return githubRepositories.filter((repo) => !analyzedFullNames.has(repo.fullName.toLowerCase()));
  }, [analyzedRepositories, githubRepositories]);

  return {
    count: unanalyzedRepos.length,
    data: unanalyzedRepos,
  };
};
