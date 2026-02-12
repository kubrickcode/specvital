/**
 * Dashboard Page E2E Tests with API Mocking
 * Tier 1: Critical cross-page navigation and infinite scroll
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockRepositoriesPage1,
  mockRepositoriesPage2,
  mockStatsNormal,
} from "./fixtures/api-responses";

test.describe("Dashboard Page (Mocked API)", () => {
  test.describe("Navigation", () => {
    test("should navigate to repository analysis page when clicking repo card", async ({
      page,
    }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      const firstRepo = mockRepositoriesPage1.data[0];
      const repoLink = page.getByRole("link", { name: new RegExp(firstRepo.name) }).first();
      await repoLink.click();

      await expect(page).toHaveURL(new RegExp(`/analyze/${firstRepo.owner}/${firstRepo.name}`));
    });
  });

  test.describe("Infinite Scroll", () => {
    test("should load more repositories on scroll", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        repositoriesPage2: mockRepositoriesPage2,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      const initialRepoCount = await page
        .getByRole("button", { name: /add bookmark|remove bookmark/i })
        .count();
      expect(initialRepoCount).toBe(20);

      await page.evaluate(() => {
        window.scrollTo(0, document.body.scrollHeight);
      });

      await page.waitForTimeout(1000);

      const finalRepoCount = await page
        .getByRole("button", { name: /add bookmark|remove bookmark/i })
        .count();
      expect(finalRepoCount).toBe(25);
    });
  });
});
