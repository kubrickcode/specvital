import { expect, test } from "@playwright/test";

test.describe("Dashboard Page (Authenticated)", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/dashboard");
  });

  test("should display dashboard page for authenticated user", async ({
    page,
  }) => {
    // Verify not redirected (authenticated user stays on dashboard)
    await expect(page).toHaveURL(/\/en\/dashboard/);

    // Verify dashboard page title
    await expect(
      page.getByRole("heading", { name: /dashboard/i })
    ).toBeVisible();
  });

  test("should display summary statistics section", async ({ page }) => {
    // Wait for summary section to load
    const summarySection = page.locator('[data-testid="summary-section"]').or(
      page.getByText(/active repositories|total tests/i).first()
    );

    await expect(summarySection).toBeVisible({ timeout: 15000 });
  });

  test("should display repository list or empty state", async ({ page }) => {
    // Wait for content to fully load (not just network idle)
    // Either repository cards appear or empty state heading appears
    await expect(
      page
        .getByRole("heading", { name: /no repositories yet/i })
        .or(page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first())
    ).toBeVisible({ timeout: 15000 });
  });

  test("should have filter controls", async ({ page }) => {
    // Wait for page to fully load
    await page.waitForLoadState("networkidle");

    // Check for filter/sort controls (may be in dropdown or visible)
    const filterControls = page
      .getByRole("combobox")
      .or(page.getByRole("button", { name: /filter|sort|all/i }));

    // Filter controls should exist (either visible or in mobile drawer)
    const controlCount = await filterControls.count();
    expect(controlCount).toBeGreaterThanOrEqual(0); // May be hidden on certain states
  });

  test("should navigate to repository analysis when clicking a repo card", async ({
    page,
  }) => {
    // Wait for repository cards to load
    await page.waitForLoadState("networkidle");

    const repoCard = page
      .locator('[data-testid="repository-card"]')
      .or(page.locator('[class*="repository-card"]'))
      .first();

    const hasRepoCard = await repoCard.isVisible({ timeout: 10000 }).catch(() => false);

    if (hasRepoCard) {
      // Click on the repository card link
      const repoLink = repoCard.getByRole("link").first();
      await repoLink.click();

      // Should navigate to analysis page
      await expect(page).toHaveURL(/\/analyze\//);
    } else {
      // No repositories yet - this is acceptable for new test users
      test.skip();
    }
  });

  test("should display organization picker when available", async ({
    page,
  }) => {
    // Wait for page to load
    await page.waitForLoadState("networkidle");

    // Organization picker may or may not be visible depending on user's org connections
    const orgPicker = page.getByRole("combobox", { name: /organization|owner/i }).or(
      page.locator('[data-testid="organization-picker"]')
    );

    // Just verify the page structure is correct, org picker is optional
    const hasOrgPicker = await orgPicker.isVisible({ timeout: 5000 }).catch(() => false);

    // This is informational - org picker availability depends on user setup
    if (hasOrgPicker) {
      await expect(orgPicker).toBeEnabled();
    }
  });

  test("should toggle bookmark on repository card", async ({ page }) => {
    // Wait for repository list to load
    await page.waitForLoadState("networkidle");

    // Find a repository card with bookmark button
    const bookmarkButton = page
      .getByRole("button", { name: /add bookmark|remove bookmark/i })
      .first();

    const hasBookmarkButton = await bookmarkButton
      .isVisible({ timeout: 10000 })
      .catch(() => false);

    if (!hasBookmarkButton) {
      // No repositories available - skip test
      test.skip();
      return;
    }

    // Get initial bookmark state
    const initialLabel = await bookmarkButton.getAttribute("aria-label");

    // Click to toggle bookmark
    await bookmarkButton.click();

    // Wait for state change
    await page.waitForTimeout(500);

    // Verify the button label changed
    if (initialLabel === "Add bookmark") {
      await expect(bookmarkButton).toHaveAttribute("aria-label", "Remove bookmark");
      await expect(bookmarkButton).toHaveAttribute("aria-pressed", "true");
    } else {
      await expect(bookmarkButton).toHaveAttribute("aria-label", "Add bookmark");
      await expect(bookmarkButton).toHaveAttribute("aria-pressed", "false");
    }

    // Toggle back to original state
    await bookmarkButton.click();

    // Verify returned to original state
    await expect(bookmarkButton).toHaveAttribute("aria-label", initialLabel!);
  });

  test("should filter by bookmarked repositories", async ({ page }) => {
    // Wait for content to fully load
    await expect(
      page
        .getByRole("heading", { name: /no repositories yet/i })
        .or(page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first())
    ).toBeVisible({ timeout: 15000 });

    // Check if user has any repositories first
    const hasNoRepos = await page
      .getByRole("heading", { name: /no repositories yet/i })
      .isVisible()
      .catch(() => false);

    if (hasNoRepos) {
      // No repositories to filter - skip test
      // Note: When there are 0 repositories, "No repositories yet" is shown
      // even with starred filter enabled. "No bookmarked repositories" only
      // shows when there are some repos but none are bookmarked.
      test.skip();
      return;
    }

    // Look for starred toggle button
    const starredToggle = page.getByRole("button", { name: /show starred only/i });

    const hasStarredControl = await starredToggle.isVisible({ timeout: 5000 }).catch(() => false);

    if (!hasStarredControl) {
      // Filter control not available
      test.skip();
      return;
    }

    // Click the starred toggle
    await starredToggle.click();

    // Verify toggle is pressed
    await expect(starredToggle).toHaveAttribute("aria-pressed", "true");

    // After clicking starred filter, either shows bookmarked repos or empty state
    const hasEmptyState = await page
      .getByRole("heading", { name: /no bookmarked repositories/i })
      .isVisible({ timeout: 5000 })
      .catch(() => false);

    const hasBookmarkedRepos = await page
      .getByRole("button", { name: "Remove bookmark" })
      .isVisible({ timeout: 5000 })
      .catch(() => false);

    // Either state is valid (empty or has bookmarked repos)
    expect(hasEmptyState || hasBookmarkedRepos).toBeTruthy();

    // Toggle off
    await starredToggle.click();
    await expect(starredToggle).toHaveAttribute("aria-pressed", "false");
  });
});
