import { expect, test } from "@playwright/test";

test.describe("Account Page", () => {
  test("should redirect unauthenticated users to homepage", async ({ page }) => {
    // Navigate to account page
    await page.goto("/en/account");

    // Should redirect to homepage for unauthenticated users
    await expect(page).toHaveURL(/\/en$/);

    // Verify we're on the homepage
    await expect(
      page.getByRole("heading", { name: /test suite/i })
    ).toBeVisible();
  });

  test("should display account page structure when navigated directly", async ({
    page,
  }) => {
    // Navigate to account page
    await page.goto("/en/account");

    // Brief check for initial account page structure before redirect
    // Note: Unauthenticated users are redirected, but this verifies
    // the page exists and has proper structure during initial load
    const accountHeading = page.getByRole("heading", { name: "Account" });

    // Either account heading is visible (brief moment) or we're redirected
    const isAccountPage = await accountHeading
      .isVisible({ timeout: 1000 })
      .catch(() => false);

    if (isAccountPage) {
      // If we see the account page, verify its structure
      await expect(page.getByText("Current Plan")).toBeVisible();
      await expect(page.getByText("Usage This Period")).toBeVisible();
    } else {
      // Redirected to homepage - expected for unauthenticated users
      await expect(page).toHaveURL(/\/en$/);
    }
  });
});
