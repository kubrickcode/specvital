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

  test("should redirect to homepage for unauthenticated direct access", async ({
    page,
  }) => {
    // Navigate directly to account page via URL
    await page.goto("/en/account");

    // Unauthenticated users should be redirected to homepage
    await expect(page).toHaveURL(/\/en$/);

    // Verify we're on the homepage
    await expect(
      page.getByRole("heading", { name: /test suite/i })
    ).toBeVisible();
  });
});
