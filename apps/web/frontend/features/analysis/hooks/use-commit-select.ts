"use client";

import { parseAsString, useQueryState } from "nuqs";

const commitParser = parseAsString.withDefault("");

const isValidCommitSha = (sha: string): boolean => /^[a-f0-9]{7,40}$/.test(sha);

export const useCommitSelect = () => {
  const [rawCommitSha, setCommitSha] = useQueryState("commit", commitParser);

  // Only return valid commit SHAs, treat invalid ones as null
  const commitSha = rawCommitSha && isValidCommitSha(rawCommitSha) ? rawCommitSha : null;

  const selectCommit = (sha: string | null) => {
    setCommitSha(sha);
  };

  return {
    commitSha,
    selectCommit,
  } as const;
};
