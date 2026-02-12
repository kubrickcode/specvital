/**
 * Spec Generation E2E Tests with API Mocking
 * Tier 1: 3-step pipeline progress modal (browser-dependent UI state machine)
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockAnalysisCompleted,
  mockSpecDocumentNotFound,
  mockSpecGenerationAccepted,
  mockUsageNormal,
} from "./fixtures/api-responses";

test.describe("Spec Generation - Pipeline Progress (Mocked API)", () => {
  test("should display 3-step pipeline with correct status indicators", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    const generateButton = page.getByRole("button", {
      name: /generate document/i,
    });
    await expect(generateButton).toBeVisible({ timeout: 15000 });
    await generateButton.click();

    const confirmButton = page
      .getByRole("dialog")
      .getByRole("button", { name: /generate document/i });
    await expect(confirmButton).toBeVisible();
    await confirmButton.click();

    const progressModal = page.getByRole("dialog", {
      name: /generation progress/i,
    });
    await expect(progressModal).toBeVisible();

    const pipelineList = progressModal.getByRole("list", {
      name: /spec generation progress steps/i,
    });
    await expect(pipelineList).toBeVisible();

    const activeStep = pipelineList.locator("[aria-current='step']");
    await expect(activeStep).toBeVisible();

    const steps = pipelineList.getByRole("listitem");
    await expect(steps).toHaveCount(3);
  });
});
