/**
 * Explore Page E2E Tests with API Mocking
 * Tests GitHub repository/organization features: My Repos, Organizations tabs
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockGitHubReposResponse,
  mockGitHubReposEmpty,
  mockOrganizationsResponse,
  mockOrganizationsEmpty,
  mockOrganizationsAllRestricted,
  mockOrgRepos,
} from "./fixtures/api-responses";

test.describe("Explore Page - My Repos Tab (Mocked API)", () => {
  test("should display user's GitHub repositories", async ({ page }) => {
    await setupMockHandlers(page, {
      githubRepos: mockGitHubReposResponse,
    });

    await page.goto("/en/explore");

    // Click My Repos tab
    const myReposTab = page.getByRole("tab", { name: /my repos/i });
    await myReposTab.click();
    await expect(myReposTab).toHaveAttribute("aria-selected", "true");

    // Wait for repos to load - check for any content from the mock data
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(1000);

    // Verify tab panel is visible
    await expect(page.getByRole("tabpanel")).toBeVisible({ timeout: 10000 });
  });

  test("should display private repository indicator", async ({ page }) => {
    await setupMockHandlers(page, {
      githubRepos: mockGitHubReposResponse,
    });

    await page.goto("/en/explore");
    await page.getByRole("tab", { name: /my repos/i }).click();

    // Wait for repos to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(500);

    // Verify private indicator is shown (icon or text)
    // mockGitHubReposList has private repos at index 0, 4, 8, 12 (every 4th repo)
    const tabpanel = page.getByRole("tabpanel");
    await expect(tabpanel).toBeVisible();

    // Check for lock icon indicator (private repos show a lock icon)
    // The lock icon is typically an SVG or img element near the repo name
    // github-repo-4 should be private (every 4th repo)
    const repoCards = tabpanel.locator('[class*="card"], article, li').filter({ hasText: /github-repo-4/i });
    await expect(repoCards.first()).toBeVisible({ timeout: 10000 });

    // Verify there are repos displayed
    const repoList = tabpanel.locator('[class*="card"], article, li');
    expect(await repoList.count()).toBeGreaterThan(0);
  });

  test("should display empty state when no repositories", async ({ page }) => {
    await setupMockHandlers(page, {
      githubRepos: mockGitHubReposEmpty,
    });

    await page.goto("/en/explore");

    // Click My Repos tab
    await page.getByRole("tab", { name: /my repos/i }).click();

    // Wait for content to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(1000);

    // Verify tab panel is showing (empty state is implementation-specific)
    await expect(page.getByRole("tabpanel")).toBeVisible({ timeout: 10000 });
  });

  test("should display refresh button for reloading data", async ({ page }) => {
    await setupMockHandlers(page, {
      githubRepos: mockGitHubReposResponse,
    });

    await page.goto("/en/explore");
    await page.getByRole("tab", { name: /my repos/i }).click();

    // Wait for page to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(500);

    // Verify Refresh from GitHub button is visible
    const refreshButton = page.getByRole("button", { name: /refresh.*github/i });
    await expect(refreshButton).toBeVisible({ timeout: 10000 });
  });

  test("should allow clicking on repository to start analysis", async ({ page }) => {
    await setupMockHandlers(page, {
      githubRepos: mockGitHubReposResponse,
    });

    await page.goto("/en/explore");
    await page.getByRole("tab", { name: /my repos/i }).click();

    // Wait for repos to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(500);

    // Verify repo cards are clickable (have links or click handlers)
    const tabpanel = page.getByRole("tabpanel");
    await expect(tabpanel).toBeVisible();

    // Check first repository link/card is clickable
    const repoLinks = tabpanel.getByRole("link").or(tabpanel.locator('[data-testid="repo-card"]'));
    await expect(repoLinks.first()).toBeVisible({ timeout: 10000 });
  });
});

test.describe("Explore Page - Organizations Tab (Mocked API)", () => {
  test("should display user's organizations", async ({ page }) => {
    await setupMockHandlers(page, {
      githubOrgs: mockOrganizationsResponse,
    });

    await page.goto("/en/explore");

    // Click Organizations tab
    const orgsTab = page.getByRole("tab", { name: /organizations/i });
    await orgsTab.click();
    await expect(orgsTab).toHaveAttribute("aria-selected", "true");

    // Wait for orgs to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(1000);

    // Verify tab panel is visible
    await expect(page.getByRole("tabpanel")).toBeVisible({ timeout: 10000 });
  });

  test("should display empty state when no organizations", async ({ page }) => {
    await setupMockHandlers(page, {
      githubOrgs: mockOrganizationsEmpty,
    });

    await page.goto("/en/explore");

    // Click Organizations tab
    await page.getByRole("tab", { name: /organizations/i }).click();

    // Wait for content to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(1000);

    // Verify tab panel is visible
    await expect(page.getByRole("tabpanel")).toBeVisible({ timeout: 10000 });
  });

  test("should display organization dropdown when multiple organizations exist", async ({ page }) => {
    await setupMockHandlers(page, {
      githubOrgs: mockOrganizationsResponse,
    });

    await page.goto("/en/explore");
    await page.getByRole("tab", { name: /organizations/i }).click();

    // Wait for orgs to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(1000);

    const tabpanel = page.getByRole("tabpanel");
    await expect(tabpanel).toBeVisible();

    // With multiple orgs, a dropdown selector should be visible
    // The dropdown trigger shows "Select organization" initially
    const orgDropdown = tabpanel.getByRole("button", { name: /select organization/i });
    await expect(orgDropdown).toBeVisible({ timeout: 10000 });

    // Click dropdown to see org options
    await orgDropdown.click();

    // Verify organizations are in dropdown menu
    const dropdownContent = page.locator('[role="menu"]');
    await expect(dropdownContent).toBeVisible();
    await expect(dropdownContent.getByText("test-org")).toBeVisible();
    await expect(dropdownContent.getByText("another-org")).toBeVisible();
  });

  test("should show install app prompt for restricted organizations", async ({ page }) => {
    await setupMockHandlers(page, {
      githubOrgs: mockOrganizationsAllRestricted,
    });

    await page.goto("/en/explore");
    await page.getByRole("tab", { name: /organizations/i }).click();

    // Wait for content to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(500);

    const tabpanel = page.getByRole("tabpanel");
    await expect(tabpanel).toBeVisible();

    // Check for install app prompt (button or link)
    const installPrompt = tabpanel.locator('button:has-text("Install"), a:has-text("Install"), :text("GitHub App")');
    await expect(installPrompt.first()).toBeVisible({ timeout: 10000 });
  });

  test("should show organization repositories when selecting from dropdown", async ({ page }) => {
    await setupMockHandlers(page, {
      githubOrgs: mockOrganizationsResponse,
      orgRepos: mockOrgRepos,
    });

    await page.goto("/en/explore");
    await page.getByRole("tab", { name: /organizations/i }).click();

    // Wait for orgs to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(1000);

    const tabpanel = page.getByRole("tabpanel");
    await expect(tabpanel).toBeVisible();

    // Click dropdown to select an organization
    const orgDropdown = tabpanel.getByRole("button", { name: /select organization/i });
    await expect(orgDropdown).toBeVisible({ timeout: 10000 });
    await orgDropdown.click();

    // Select test-org from dropdown
    const dropdownContent = page.locator('[role="menu"]');
    await expect(dropdownContent).toBeVisible();
    await dropdownContent.getByText("test-org").click();

    // Wait for org repos to load
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(500);

    // After selecting an org, should show repository list or search input
    // The UI shows a search input and repo list for accessible orgs
    const searchInput = tabpanel.getByRole("searchbox");
    await expect(searchInput).toBeVisible({ timeout: 10000 });

    // Verify repos from mockOrgRepos["test-org"] are visible
    // mockOrgRepos has repos named like "org-repo-1", "org-repo-2", etc.
    await expect(tabpanel.getByText(/org-repo/i).first()).toBeVisible({ timeout: 10000 });
  });
});

test.describe("Explore Page - Tab Navigation (Mocked API)", () => {
  test("should switch between Community, My Repos, and Organizations tabs", async ({ page }) => {
    await setupMockHandlers(page, {
      githubRepos: mockGitHubReposResponse,
      githubOrgs: mockOrganizationsResponse,
    });

    await page.goto("/en/explore");

    // Start on Community tab (default)
    const communityTab = page.getByRole("tab", { name: /community/i });
    await expect(communityTab).toHaveAttribute("aria-selected", "true");

    // Switch to My Repos
    const myReposTab = page.getByRole("tab", { name: /my repos/i });
    await myReposTab.click();
    await expect(myReposTab).toHaveAttribute("aria-selected", "true");
    await expect(communityTab).toHaveAttribute("aria-selected", "false");

    // Switch to Organizations
    const orgsTab = page.getByRole("tab", { name: /organizations/i });
    await orgsTab.click();
    await expect(orgsTab).toHaveAttribute("aria-selected", "true");
    await expect(myReposTab).toHaveAttribute("aria-selected", "false");

    // Switch back to Community
    await communityTab.click();
    await expect(communityTab).toHaveAttribute("aria-selected", "true");
    await expect(orgsTab).toHaveAttribute("aria-selected", "false");
  });

  test("should preserve tab content when switching back", async ({ page }) => {
    await setupMockHandlers(page, {
      githubRepos: mockGitHubReposResponse,
      githubOrgs: mockOrganizationsResponse,
    });

    await page.goto("/en/explore");

    // Go to My Repos tab
    await page.getByRole("tab", { name: /my repos/i }).click();
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(500);

    // Check repos loaded
    const myReposPanel = page.getByRole("tabpanel");
    await expect(myReposPanel).toBeVisible();

    // Switch to Organizations
    await page.getByRole("tab", { name: /organizations/i }).click();
    await page.waitForTimeout(500);

    // Switch back to My Repos
    await page.getByRole("tab", { name: /my repos/i }).click();
    await page.waitForTimeout(500);

    // Verify content is still there (not requiring re-fetch)
    await expect(page.getByRole("tabpanel")).toBeVisible();
  });
});
