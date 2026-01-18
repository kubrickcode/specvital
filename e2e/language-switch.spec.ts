import { expect, test } from "@playwright/test";

test.describe("Language Switching", () => {
  test("should switch from English to Korean", async ({ page }) => {
    await page.goto("/en");

    // Click language selector
    await page.getByRole("button", { name: "Select language" }).click();

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

    // Click language selector
    await page.getByRole("button", { name: "언어 선택" }).click();

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
