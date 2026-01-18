import { expect, test } from "@playwright/test";

test.describe("Homepage", () => {
  test("should display URL input form", async ({ page }) => {
    await page.goto("/");

    // Verify main heading
    await expect(
      page.getByRole("heading", { name: /spec/i }).first()
    ).toBeVisible();

    // Verify URL input exists
    await expect(
      page.getByRole("textbox", { name: /github/i })
    ).toBeVisible();

    // Verify submit button exists
    await expect(
      page.getByRole("button", { name: /start analysis|분석 시작/i })
    ).toBeVisible();
  });

  test("should navigate to analysis page on valid URL input", async ({
    page,
  }) => {
    await page.goto("/");

    // Enter GitHub URL
    await page.getByRole("textbox", { name: /github/i }).fill("facebook/react");

    // Click start button
    await page.getByRole("button", { name: /start analysis|분석 시작/i }).click();

    // Verify navigation to analysis page
    await expect(page).toHaveURL(/\/analyze\/facebook\/react/);
  });

  test("should show error for invalid URL format", async ({ page }) => {
    await page.goto("/");

    // Enter invalid URL
    await page.getByRole("textbox", { name: /github/i }).fill("invalid-url");

    // Click Analyze button
    await page
      .getByRole("button", { name: /start analysis|분석 시작|analyze/i })
      .click();

    // Verify error message is shown
    await expect(page.getByText(/invalid github/i)).toBeVisible();

    // Verify URL hasn't changed (still on homepage)
    await expect(page).not.toHaveURL(/\/analyze/);
  });
});
