import { expect, test } from "@playwright/test";

test.describe("User Menu (Authenticated)", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/dashboard");
  });

  test("should display user menu with avatar", async ({ page }) => {
    // Wait for page to load
    await page.waitForLoadState("networkidle");

    // Find user menu button (avatar button with "Open user menu" label)
    const userMenuButton = page.getByRole("button", { name: /open user menu/i });
    await expect(userMenuButton).toBeVisible({ timeout: 10000 });
  });

  test("should open user menu dropdown on click", async ({ page }) => {
    await page.waitForLoadState("networkidle");

    // Click user menu button
    const userMenuButton = page.getByRole("button", { name: /open user menu/i });
    await userMenuButton.click();

    // Verify dropdown menu is visible
    const dropdownContent = page.locator('[role="menu"]');
    await expect(dropdownContent).toBeVisible();

    // Verify user info is displayed (username with @ prefix)
    await expect(page.getByText(/@\w+/)).toBeVisible();

    // Verify Account menu item
    await expect(page.getByRole("menuitem", { name: /account/i })).toBeVisible();

    // Verify Sign out menu item
    await expect(page.getByRole("menuitem", { name: /sign out/i })).toBeVisible();
  });

  test("should navigate to account page from user menu", async ({ page }) => {
    await page.waitForLoadState("networkidle");

    // Open user menu
    const userMenuButton = page.getByRole("button", { name: /open user menu/i });
    await userMenuButton.click();

    // Click Account menu item
    await page.getByRole("menuitem", { name: /account/i }).click();

    // Verify navigation to account page
    await expect(page).toHaveURL(/\/account/);
  });

  test("should logout and redirect to home page", async ({ page }) => {
    await page.waitForLoadState("networkidle");

    // Verify we are authenticated (user menu visible)
    const userMenuButton = page.getByRole("button", { name: /open user menu/i });
    await expect(userMenuButton).toBeVisible({ timeout: 10000 });

    // Open user menu
    await userMenuButton.click();

    // Click Sign out
    await page.getByRole("menuitem", { name: /sign out/i }).click();

    // Wait for logout to complete and redirect
    await page.waitForLoadState("networkidle");

    // Verify redirected to home page (or locale root)
    await expect(page).toHaveURL(/\/(en)?$/);

    // Verify user is logged out - Sign In button should be visible instead of user menu
    const signInButton = page.getByRole("button", { name: "Sign in" });
    await expect(signInButton).toBeVisible({ timeout: 10000 });
  });

  test("should close user menu when clicking outside", async ({ page }) => {
    await page.waitForLoadState("networkidle");

    // Open user menu
    const userMenuButton = page.getByRole("button", { name: /open user menu/i });
    await userMenuButton.click();

    // Verify dropdown is open
    const dropdownContent = page.locator('[role="menu"]');
    await expect(dropdownContent).toBeVisible();

    // Click outside (press Escape to close dropdown)
    await page.keyboard.press("Escape");

    // Verify dropdown is closed
    await expect(dropdownContent).not.toBeVisible();
  });
});
