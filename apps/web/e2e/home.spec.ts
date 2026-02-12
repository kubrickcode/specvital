/**
 * Homepage E2E Tests
 * Tier 1: Critical user flow - URL input to analysis navigation
 */

import { expect, test } from "@playwright/test";

test.describe("Homepage", () => {
  test("should navigate to analysis page on valid URL input", async ({ page }) => {
    await page.goto("/");

    const input = page.getByRole("textbox", { name: /github/i });
    await input.fill("facebook/react");

    await page.waitForTimeout(600);

    const button = page.getByRole("button", { name: /analyze|분석 시작/i });
    await expect(button).toBeEnabled();

    await Promise.all([
      page.waitForURL(/\/analyze\/facebook\/react/, { timeout: 30000 }),
      button.click(),
    ]);
  });
});
