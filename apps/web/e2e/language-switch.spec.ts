import { expect, test } from "@playwright/test";

test.describe("Language Switching", () => {
  test("should switch from English to Korean", async ({ page }) => {
    await page.goto("/en");

    // Wait for hydration to complete (client components mount)
    const languageButton = page.getByRole("button", { name: "Select language" });
    await expect(languageButton).toBeEnabled();
    await page.waitForTimeout(300); // Additional wait for event handlers to attach

    // Click language selector
    await languageButton.click();

    // Wait for dropdown menu to appear
    await expect(page.getByRole("menu")).toBeVisible();

    // Verify English is disabled (current language)
    await expect(page.getByRole("menuitem", { name: "English" })).toBeDisabled();

    // Verify Korean is enabled
    await expect(page.getByRole("menuitem", { name: "한국어" })).toBeEnabled();

    // Switch to Korean
    await page.getByRole("menuitem", { name: "한국어" }).click();

    // Verify URL changed to Korean
    await expect(page).toHaveURL(/\/ko/);

    // Verify page title is in Korean
    await expect(page).toHaveTitle(/테스트 스펙 분석기/);

    // Verify UI text is in Korean
    await expect(page.getByRole("button", { name: "로그인" })).toBeVisible();
    await expect(page.getByRole("button", { name: "분석 시작" })).toBeVisible();
  });

  test("should switch from Korean to English", async ({ page }) => {
    await page.goto("/ko");

    // Wait for hydration to complete (client components mount)
    const languageButton = page.getByRole("button", { name: "언어 선택" });
    await expect(languageButton).toBeEnabled();
    await page.waitForTimeout(300); // Additional wait for event handlers to attach

    // Click language selector
    await languageButton.click();

    // Wait for dropdown menu to appear
    await expect(page.getByRole("menu")).toBeVisible();

    // Verify Korean is disabled (current language)
    await expect(page.getByRole("menuitem", { name: "한국어" })).toBeDisabled();

    // Verify English is enabled
    await expect(page.getByRole("menuitem", { name: "English" })).toBeEnabled();

    // Switch to English
    await page.getByRole("menuitem", { name: "English" }).click();

    // Verify URL changed to English
    await expect(page).toHaveURL(/\/en/);

    // Verify page title is in English
    await expect(page).toHaveTitle(/Test Spec Analyzer/);

    // Verify UI text is in English
    await expect(page.getByRole("button", { name: "Sign in" })).toBeVisible();
    await expect(page.getByRole("button", { name: "Analyze" })).toBeVisible();
  });
});
