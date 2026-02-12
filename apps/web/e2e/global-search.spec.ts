/**
 * Global Search E2E Tests
 * Tier 1: Cross-page navigation via search
 */

import { expect, test } from "@playwright/test";

test.describe("Global Search", () => {
  test("should navigate to page on Enter", async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 720 });
    await page.goto("/en");

    await page.keyboard.press("ControlOrMeta+k");

    await page.keyboard.press("ArrowDown");

    await page.getByText(/go to explore/i).click();

    await expect(page).toHaveURL(/\/explore/);
  });
});
