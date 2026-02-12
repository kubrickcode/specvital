/**
 * User Menu E2E Tests (Authenticated)
 * Tier 1: Auth-dependent navigation and logout flow
 */

import { expect, test } from "@playwright/test";

test.describe("User Menu (Authenticated)", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/dashboard");
    await page.waitForLoadState("networkidle");
  });

  test("should navigate to account page from user menu", async ({ page }) => {
    const userMenuButton = page.getByRole("button", { name: /open user menu/i });
    await userMenuButton.click();

    await page.getByRole("menuitem", { name: /account/i }).click();

    await expect(page).toHaveURL(/\/account/);
  });

  test("should logout and redirect to home page", async ({ page }) => {
    const userMenuButton = page.getByRole("button", { name: /open user menu/i });
    await expect(userMenuButton).toBeVisible({ timeout: 10000 });

    await userMenuButton.click();

    await page.getByRole("menuitem", { name: /sign out/i }).click();

    await page.waitForLoadState("networkidle");

    await expect(page).toHaveURL(/\/(en)?$/);

    const signInButton = page.getByRole("button", { name: "Sign in" });
    await expect(signInButton).toBeVisible({ timeout: 10000 });
  });
});
