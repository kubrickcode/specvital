import { describe, expect, it } from "vitest";

import { isValidGitHubUrl, parseGitHubUrl } from "./github-url";

describe("parseGitHubUrl", () => {
  describe("valid URLs", () => {
    it("parses HTTPS URL", () => {
      const result = parseGitHubUrl("https://github.com/facebook/react");
      expect(result).toEqual({
        data: { owner: "facebook", repo: "react" },
        success: true,
      });
    });

    it("parses HTTP URL", () => {
      const result = parseGitHubUrl("http://github.com/owner/repo");
      expect(result.success).toBe(true);
    });

    it("handles trailing slash", () => {
      const result = parseGitHubUrl("https://github.com/owner/repo/");
      expect(result).toEqual({
        data: { owner: "owner", repo: "repo" },
        success: true,
      });
    });

    it("handles repo with dots", () => {
      const result = parseGitHubUrl("https://github.com/owner/repo.js");
      expect(result.success).toBe(true);
    });

    it("handles repo with underscores", () => {
      const result = parseGitHubUrl("https://github.com/owner/my_repo");
      expect(result.success).toBe(true);
    });

    it("handles owner with hyphens", () => {
      const result = parseGitHubUrl("https://github.com/my-org/repo");
      expect(result.success).toBe(true);
    });

    it("handles single char owner", () => {
      const result = parseGitHubUrl("https://github.com/a/repo");
      expect(result.success).toBe(true);
    });
  });

  describe("invalid URLs", () => {
    it("rejects empty URL", () => {
      const result = parseGitHubUrl("");
      expect(result).toEqual({
        error: "URL is required",
        success: false,
      });
    });

    it("rejects whitespace only", () => {
      const result = parseGitHubUrl("   ");
      expect(result).toEqual({
        error: "URL is required",
        success: false,
      });
    });

    it("rejects non-GitHub URL", () => {
      const result = parseGitHubUrl("https://gitlab.com/owner/repo");
      expect(result.success).toBe(false);
    });

    it("rejects URL with extra path segments", () => {
      const result = parseGitHubUrl("https://github.com/owner/repo/tree/main");
      expect(result.success).toBe(false);
    });

    it("rejects URL without repo", () => {
      const result = parseGitHubUrl("https://github.com/owner");
      expect(result.success).toBe(false);
    });

    it("rejects owner starting with hyphen", () => {
      const result = parseGitHubUrl("https://github.com/-owner/repo");
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error).toBe("Invalid GitHub username format");
      }
    });

    it("rejects owner ending with hyphen", () => {
      const result = parseGitHubUrl("https://github.com/owner-/repo");
      expect(result.success).toBe(false);
    });

    it("rejects repo ending with .git", () => {
      const result = parseGitHubUrl("https://github.com/owner/repo.git");
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error).toBe("Invalid repository name format");
      }
    });

    it("rejects double dots in path (path traversal)", () => {
      const result = parseGitHubUrl("https://github.com/../etc/passwd");
      expect(result.success).toBe(false);
    });

    it('rejects repo named ".."', () => {
      const result = parseGitHubUrl("https://github.com/owner/..");
      expect(result.success).toBe(false);
    });
  });
});

describe("isValidGitHubUrl", () => {
  it("returns true for valid URL", () => {
    expect(isValidGitHubUrl("https://github.com/facebook/react")).toBe(true);
  });

  it("returns false for invalid URL", () => {
    expect(isValidGitHubUrl("invalid")).toBe(false);
  });
});
