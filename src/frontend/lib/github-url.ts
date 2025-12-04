import { z } from "zod";

// GitHub naming rules:
// - Owner: 1-39 chars, alphanumeric or hyphen, cannot start/end with hyphen
// - Repo: 1-100 chars, alphanumeric, hyphen, underscore, dot, cannot be ".." or end with ".git"
const OWNER_PATTERN = /^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?$/;
const REPO_PATTERN = /^(?!\.\.)[a-zA-Z0-9._-]{1,100}(?<!\.git)$/;
const GITHUB_URL_PATTERN = /^https?:\/\/github\.com\/([^/]+)\/([^/]+?)\/?$/;

export type ParsedGitHubUrl = {
  owner: string;
  repo: string;
};

export type ParseGitHubUrlResult =
  | { success: true; data: ParsedGitHubUrl }
  | { success: false; error: string };

const gitHubUrlSchema = z
  .string()
  .trim()
  .min(1, "URL is required")
  .superRefine((url, ctx) => {
    const match = url.match(GITHUB_URL_PATTERN);

    if (!match) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Invalid GitHub repository URL (e.g., https://github.com/owner/repo)",
      });
      return z.NEVER;
    }

    const [, owner, repo] = match;

    if (!OWNER_PATTERN.test(owner)) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Invalid GitHub username format",
      });
      return z.NEVER;
    }

    if (!REPO_PATTERN.test(repo)) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Invalid repository name format",
      });
      return z.NEVER;
    }

    if (owner.includes("..") || repo.includes("..")) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Invalid path in URL",
      });
      return z.NEVER;
    }
  })
  .transform((url) => {
    const match = url.match(GITHUB_URL_PATTERN)!;
    return { owner: match[1], repo: match[2] };
  });

export const parseGitHubUrl = (url: string): ParseGitHubUrlResult => {
  const result = gitHubUrlSchema.safeParse(url);

  if (!result.success) {
    return {
      success: false,
      error: result.error.errors[0]?.message ?? "Invalid URL",
    };
  }

  return {
    success: true,
    data: result.data,
  };
};

export const isValidGitHubUrl = (url: string): boolean => {
  return gitHubUrlSchema.safeParse(url).success;
};
