const GITHUB_IDENTIFIER_MAX_LENGTH = 39;
const GITHUB_IDENTIFIER_REGEX = /^[a-zA-Z0-9]([a-zA-Z0-9._-]*[a-zA-Z0-9])?$/;

export const isValidGitHubIdentifier = (value: string): boolean =>
  GITHUB_IDENTIFIER_REGEX.test(value) && value.length <= GITHUB_IDENTIFIER_MAX_LENGTH;

export const validateRepositoryIdentifiers = (owner: string, repo: string): void => {
  if (!isValidGitHubIdentifier(owner) || !isValidGitHubIdentifier(repo)) {
    throw new Error("Invalid repository identifier");
  }
};
