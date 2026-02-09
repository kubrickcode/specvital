import { expect, test } from "@playwright/test";

test.describe("Explore Page (Authenticated)", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/explore");
  });

  test("should display My Repos tab for authenticated user", async ({ page }) => {
    // Verify My Repos tab is visible
    const myReposTab = page.getByRole("tab", { name: "My Repos" });
    await expect(myReposTab).toBeVisible();

    // Click My Repos tab
    await myReposTab.click();

    // Verify My Repos tab is selected
    await expect(myReposTab).toHaveAttribute("aria-selected", "true");

    // Since dev-login user has no GitHub token, expect error state or loading
    // Either shows error message with retry button OR loading state
    await expect(
      page.getByRole("button", { name: /refresh/i }).or(page.getByText(/error|unauthorized/i))
    ).toBeVisible({ timeout: 15000 });
  });

  test("should display Organizations tab for authenticated user", async ({ page }) => {
    // Verify Organizations tab is visible
    const orgsTab = page.getByRole("tab", { name: "Organizations" });
    await expect(orgsTab).toBeVisible();

    // Click Organizations tab
    await orgsTab.click();

    // Verify Organizations tab is selected
    await expect(orgsTab).toHaveAttribute("aria-selected", "true");

    // Wait for Organizations tab panel to be visible
    await expect(page.getByRole("tabpanel", { name: "Organizations" })).toBeVisible();

    // Verify authenticated state - should NOT show login required state
    // Instead should show: loading skeleton, empty state, error state, or org list
    await expect(page.getByText(/sign in to view/i)).not.toBeVisible({ timeout: 5000 });
  });

  test("should not show login required state on My Repos for authenticated user", async ({
    page,
  }) => {
    // Click My Repos tab
    await page.getByRole("tab", { name: "My Repos" }).click();

    // Wait for content to load
    await page.waitForLoadState("networkidle");

    // Should NOT show "Sign in to view your repositories" message
    await expect(page.getByText(/sign in to view your repositories/i)).not.toBeVisible();
  });

  test("should not show login required state on Organizations for authenticated user", async ({
    page,
  }) => {
    // Click Organizations tab
    await page.getByRole("tab", { name: "Organizations" }).click();

    // Wait for content to load
    await page.waitForLoadState("networkidle");

    // Should NOT show "Sign in to view organizations" message
    await expect(page.getByText(/sign in to view organizations/i)).not.toBeVisible();
  });

  test("should allow tab switching between Community and My Repos", async ({ page }) => {
    // Start on Community tab (default)
    const communityTab = page.getByRole("tab", { name: "Community" });
    await expect(communityTab).toHaveAttribute("aria-selected", "true");

    // Switch to My Repos
    await page.getByRole("tab", { name: "My Repos" }).click();
    await expect(page.getByRole("tab", { name: "My Repos" })).toHaveAttribute(
      "aria-selected",
      "true"
    );

    // Switch back to Community
    await communityTab.click();
    await expect(communityTab).toHaveAttribute("aria-selected", "true");

    // Verify Community content is shown
    await expect(page.getByRole("tabpanel", { name: "Community" })).toBeVisible();
  });
});
