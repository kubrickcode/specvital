/**
 * Analysis Page E2E Tests with API Mocking
 * Tier 1: Critical cross-page flows that require E2E (not testable via component tests)
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockAnalysisCompleted,
  mockSpecDocumentNotFound,
  mockSpecGenerationAccepted,
  mockUsageNormal,
} from "./fixtures/api-responses";

test.describe("Analysis Page - Spec Generation Progress (Mocked API)", () => {
  test("should display progress modal during generation", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      onSpecGeneration: () => mockSpecGenerationAccepted,
      onSpecGenerationStatus: () => mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();
    const confirmButton = dialog.getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    await expect(page.getByText(/generating|processing/i)).toBeVisible({ timeout: 5000 });
  });
});

test.describe("Analysis Page - URL Backward Compatibility (Mocked API)", () => {
  test("should redirect ?view=document to ?tab=spec", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?view=document");
    await page.waitForLoadState("networkidle");

    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await expect(specTab).toHaveAttribute("data-state", "active");

    await expect(page).toHaveURL(/tab=spec/);
    await expect(page).not.toHaveURL(/view=document/);
  });
});

test.describe("Analysis Page - Spec 403 Handling (Mocked API)", () => {
  test("should display subscription guidance when 403 error occurs", async ({ page }) => {
    await page.route("**/api/spec-view/generate", async (route) => {
      return route.fulfill({
        status: 403,
        contentType: "application/json",
        body: JSON.stringify({
          error: "Subscription required",
          detail: "Upgrade your plan to use AI Spec generation",
        }),
      });
    });

    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();
    const confirmButton = dialog.getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    await expect(page.getByText("Subscription Required")).toBeVisible({ timeout: 5000 });
  });
});

test.describe("Analysis Page - Update Banner Polling (Mocked API)", () => {
  test("should dismiss banner after analysis completion", async ({ page }) => {
    let statusCallCount = 0;

    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      updateStatus: {
        status: "new-commits",
        parserOutdated: false,
      },
      onReanalyze: () => ({ status: "queued" as const }),
      onAnalysisStatus: () => {
        statusCallCount++;
        return {
          owner: "test-owner",
          repo: "test-repo",
          status: statusCallCount >= 2 ? "completed" : "analyzing",
        };
      },
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    const updateBanner = page.locator('[class*="bg-amber"]');
    await expect(updateBanner.first()).toBeVisible({ timeout: 10000 });

    const updateNowButton = page.getByRole("button", { name: /update now/i });
    await updateNowButton.click();

    await expect(updateBanner).not.toBeVisible({ timeout: 10000 });
  });
});
