import { expect, test } from "@playwright/test";

test.describe("Homepage", () => {
  test("should display URL input form", async ({ page }) => {
    await page.goto("/");

    // Verify main heading
    await expect(
      page.getByRole("heading", { name: /test suite/i }).first()
    ).toBeVisible();

    // Verify URL input exists
    await expect(
      page.getByRole("textbox", { name: /github/i })
    ).toBeVisible();

    // Verify submit button exists
    await expect(
      page.getByRole("button", { name: /analyze|분석 시작/i })
    ).toBeVisible();
  });

  test("should navigate to analysis page on valid URL input", async ({
    page,
  }) => {
    await page.goto("/");

    // Enter GitHub URL
    const input = page.getByRole("textbox", { name: /github/i });
    await input.fill("facebook/react");

    // Wait for validation to complete (debounced 500ms)
    await page.waitForTimeout(600);

    // Click start button and wait for navigation
    const button = page.getByRole("button", { name: /analyze|분석 시작/i });
    await expect(button).toBeEnabled();

    // Use Promise.all to wait for navigation after click
    await Promise.all([
      page.waitForURL(/\/analyze\/facebook\/react/, { timeout: 30000 }),
      button.click(),
    ]);
  });

  test("should show error for invalid URL format", async ({ page }) => {
    await page.goto("/");

    // Enter invalid URL
    await page.getByRole("textbox", { name: /github/i }).fill("invalid-url");

    // Click Analyze button
    await page.getByRole("button", { name: /analyze|분석 시작/i }).click();

    // Verify error message is shown (in alert, not sr-only)
    await expect(
      page.locator("[data-slot='alert-description']").filter({ hasText: /invalid github/i })
    ).toBeVisible();

    // Verify URL hasn't changed (still on homepage)
    await expect(page).not.toHaveURL(/\/analyze/);
  });

  test("should display URL format help button", async ({ page }) => {
    await page.goto("/");

    // Verify help button exists with proper aria-label
    const helpButton = page.getByRole("button", {
      name: /view supported formats|지원되는 형식/i,
    });
    await expect(helpButton).toBeVisible();

    // Verify button is accessible (has icon)
    await expect(helpButton.locator("svg")).toBeVisible();
  });

  test("should display 20+ frameworks button", async ({ page }) => {
    await page.goto("/");

    // Verify frameworks button exists
    const frameworksButton = page.getByRole("button", {
      name: /20\+ frameworks/i,
    });
    await expect(frameworksButton).toBeVisible();

    // Click the button to show framework list
    await frameworksButton.click();

    // Verify popover/dialog with framework list appears
    await expect(page.getByRole("dialog")).toBeVisible();

    // Verify some framework names are listed
    await expect(page.getByText(/jest/i)).toBeVisible();
    await expect(page.getByText(/vitest/i)).toBeVisible();
    await expect(page.getByText(/playwright/i)).toBeVisible();

    // Close dialog
    await page.keyboard.press("Escape");
    await expect(page.getByRole("dialog")).not.toBeVisible();
  });
});
