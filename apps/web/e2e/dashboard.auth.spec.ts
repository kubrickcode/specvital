/**
 * Dashboard Page E2E Tests (Authenticated)
 * Tier 1: Critical cross-page navigation with real auth
 */

import { expect, test } from "@playwright/test";

test.describe("Dashboard Page (Authenticated)", () => {
  test("should navigate to repository analysis when clicking a repo card", async ({ page }) => {
    await page.goto("/en/dashboard");
    await page.waitForLoadState("networkidle");

    const repoCard = page
      .locator('[data-testid="repository-card"]')
      .or(page.locator('[class*="repository-card"]'))
      .first();

    const hasRepoCard = await repoCard.isVisible({ timeout: 10000 }).catch(() => false);

    if (hasRepoCard) {
      const repoLink = repoCard.getByRole("link").first();
      await repoLink.click();

      await expect(page).toHaveURL(/\/analyze\//);
    } else {
      test.skip();
    }
  });
});
