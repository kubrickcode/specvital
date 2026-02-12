/**
 * Analysis Polling E2E Tests with API Mocking
 * Tier 1: Auto-transition from waiting to results (polling state machine)
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockAnalysisQueued,
  mockAnalysisCompleted,
  mockSpecDocumentNotFound,
  mockUsageNormal,
} from "./fixtures/api-responses";

test.describe("Analysis Complete - Auto Transition (Mocked API)", () => {
  test("should auto-transition from waiting card to results when analysis completes", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisQueued,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    const waitingCard = page
      .locator("[role='status'][aria-live='polite']")
      .filter({ has: page.locator("time") });
    await expect(waitingCard).toBeVisible({ timeout: 15000 });

    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.waitForTimeout(5000);

    await expect(waitingCard).not.toBeVisible({ timeout: 10000 });

    await expect(page.getByText("Total")).toBeVisible({ timeout: 10000 });
    await expect(page.getByText("Active")).toBeVisible();
  });
});
