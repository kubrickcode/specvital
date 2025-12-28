import { z } from "zod";

const OWNER_PATTERN = /^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?$/;
const REPO_PATTERN = /^(?!\.\.)[a-zA-Z0-9._-]{1,100}(?<!\.git)$/;

const SHORTHAND_PATTERN = /^([^/]+)\/([^/]+)$/;
const GITHUB_URL_PATTERN = /^https?:\/\/github\.com\/([^/]+)\/([^/]+?)(?:\.git)?(?:\/.*)?$/;
const GITHUB_DOMAIN_PATTERN = /^github\.com\/([^/]+)\/([^/]+?)(?:\.git)?(?:\/.*)?$/;

export type ParsedGitHubUrl = {
  owner: string;
  repo: string;
};

export type ParseGitHubUrlResult =
  | { data: ParsedGitHubUrl; success: true }
  | { error: string; success: false };

type NormalizeResult = { owner: string; repo: string; success: true } | { success: false };

const extractOwnerRepo = (input: string): NormalizeResult => {
  const trimmed = input.trim();
  if (!trimmed) return { success: false };

  let match = trimmed.match(GITHUB_URL_PATTERN);
  if (match) {
    return { owner: match[1], repo: match[2], success: true };
  }

  match = trimmed.match(GITHUB_DOMAIN_PATTERN);
  if (match) {
    return { owner: match[1], repo: match[2], success: true };
  }

  match = trimmed.match(SHORTHAND_PATTERN);
  if (match) {
    return { owner: match[1], repo: match[2], success: true };
  }

  return { success: false };
};

export const normalizeGitHubInput = (input: string): string | null => {
  const result = extractOwnerRepo(input);
  if (!result.success) return null;
  return `${result.owner}/${result.repo}`;
};

const gitHubInputSchema = z
  .string()
  .trim()
  .min(1, "URL is required")
  .superRefine((input, ctx) => {
    const result = extractOwnerRepo(input);

    if (!result.success) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Invalid GitHub repository (e.g., owner/repo or https://github.com/owner/repo)",
      });
      return z.NEVER;
    }

    const { owner, repo } = result;

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
  .transform((input) => {
    const result = extractOwnerRepo(input);
    if (!result.success) throw new Error("Unexpected parse failure");
    return { owner: result.owner, repo: result.repo };
  });

export const parseGitHubUrl = (input: string): ParseGitHubUrlResult => {
  const result = gitHubInputSchema.safeParse(input);

  if (!result.success) {
    return {
      error: result.error.issues[0]?.message ?? "Invalid URL",
      success: false,
    };
  }

  return {
    data: result.data,
    success: true,
  };
};

export const isValidGitHubUrl = (input: string): boolean => {
  return gitHubInputSchema.safeParse(input).success;
};

export type InputFeedback =
  | { type: "empty" }
  | { normalized: string; type: "shorthand" }
  | { normalized: string; type: "deeplink" }
  | { type: "url" }
  | { type: "invalid" };

const DEEPLINK_SUFFIX_PATTERN =
  /\/(?:tree|blob|issues|pull|pulls|actions|releases|commits|branches|tags)\b/;

export const getInputFeedback = (input: string): InputFeedback => {
  const trimmed = input.trim();
  if (!trimmed) return { type: "empty" };

  if (GITHUB_URL_PATTERN.test(trimmed) || GITHUB_DOMAIN_PATTERN.test(trimmed)) {
    if (DEEPLINK_SUFFIX_PATTERN.test(trimmed)) {
      const result = parseGitHubUrl(trimmed);
      if (result.success) {
        return {
          normalized: `${result.data.owner}/${result.data.repo}`,
          type: "deeplink",
        };
      }
    }
    return { type: "url" };
  }

  if (SHORTHAND_PATTERN.test(trimmed)) {
    const result = parseGitHubUrl(trimmed);
    if (result.success) {
      return {
        normalized: `github.com/${result.data.owner}/${result.data.repo}`,
        type: "shorthand",
      };
    }
  }

  return { type: "invalid" };
};
