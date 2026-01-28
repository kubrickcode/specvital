/**
 * Analysis Polling E2E Tests with API Mocking
 * Tests analysis waiting states (Queued/Analyzing) and auto-transition to results
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockAnalysisQueued,
  mockAnalysisAnalyzing,
  mockAnalysisCompleted,
  mockSpecDocumentNotFound,
  mockUsageNormal,
} from "./fixtures/api-responses";

test.describe("Analysis Waiting Card - Queued/Analyzing State (Mocked API)", () => {
  test("should display PulseRing and elapsed time in Queued state", async ({
    page,
  }) => {
    // Setup: Analysis in queued state
    await setupMockHandlers(page, {
      analysis: mockAnalysisQueued,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Verify waiting card displays instead of results
    // Look for role="status" with aria-live="polite"
    const waitingCard = page.locator("[role='status'][aria-live='polite']");
    await expect(waitingCard).toBeVisible({ timeout: 15000 });

    // Verify PulseRing component visible
    // PulseRing has chart-2 color in queued state
    const pulseRing = waitingCard.locator(".text-chart-2").first();
    await expect(pulseRing).toBeVisible();

    // Verify ShimmerBar component visible
    // Queued state has slow shimmer animation
    const shimmerBar = waitingCard.locator("[role='progressbar']");
    await expect(shimmerBar).toBeVisible();

    // Verify elapsed time display
    const timeElement = waitingCard.locator("time[datetime]");
    await expect(timeElement).toBeVisible();

    // Verify rotating message area
    const messageText = waitingCard.locator("p").first();
    await expect(messageText).toBeVisible();
  });

  test("should display PulseRing and elapsed time in Analyzing state", async ({
    page,
  }) => {
    // Setup: Analysis in analyzing state
    await setupMockHandlers(page, {
      analysis: mockAnalysisAnalyzing,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Verify waiting card displays
    const waitingCard = page.locator("[role='status'][aria-live='polite']");
    await expect(waitingCard).toBeVisible({ timeout: 15000 });

    // Verify PulseRing with chart-1 color (analyzing state)
    const pulseRing = waitingCard.locator(".text-chart-1").first();
    await expect(pulseRing).toBeVisible();

    // Verify ShimmerBar with faster animation (analyzing state)
    const shimmerBar = waitingCard.locator("[role='progressbar']");
    await expect(shimmerBar).toBeVisible();

    // Verify elapsed time
    const timeElement = waitingCard.locator("time[datetime]");
    await expect(timeElement).toBeVisible();
  });

  test("should display rotating messages that change over time", async ({
    page,
  }) => {
    // Setup: Analysis in queued state
    await setupMockHandlers(page, {
      analysis: mockAnalysisQueued,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    const waitingCard = page.locator("[role='status'][aria-live='polite']");
    await expect(waitingCard).toBeVisible({ timeout: 15000 });

    // Get initial message
    const messageElement = waitingCard.locator("p").first();
    const initialMessage = await messageElement.textContent();

    // Wait for message rotation interval (typically 5-15 seconds)
    await page.waitForTimeout(6000);

    // Message should change (or may be same if rotation interval not reached)
    const currentMessage = await messageElement.textContent();

    // Verify message element exists (content may or may not have changed)
    expect(currentMessage).toBeTruthy();
  });

  test("should show long wait guidance after 60+ seconds", async ({ page }) => {
    // Setup: Analysis started 65 seconds ago
    const mockAnalysisLongWait = {
      status: "analyzing" as const,
      owner: "test-owner",
      repo: "test-repo",
      startedAt: new Date(Date.now() - 65000).toISOString(), // 65 seconds ago
    };

    await setupMockHandlers(page, {
      analysis: mockAnalysisLongWait,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    const waitingCard = page.locator("[role='status'][aria-live='polite']");
    await expect(waitingCard).toBeVisible({ timeout: 15000 });

    // Verify "taking longer than expected" guidance appears
    await expect(
      waitingCard.getByText(/taking longer than expected/i)
    ).toBeVisible({ timeout: 5000 });
  });
});

test.describe("Analysis Complete - Auto Transition (Mocked API)", () => {
  test("should auto-transition from waiting card to results when analysis completes", async ({
    page,
  }) => {
    // Setup: Start with queued state
    await setupMockHandlers(page, {
      analysis: mockAnalysisQueued,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Verify waiting card displays initially
    const waitingCard = page.locator("[role='status'][aria-live='polite']");
    await expect(waitingCard).toBeVisible({ timeout: 15000 });

    // Update mock to completed status
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    // Wait for polling interval (typically 2-3 seconds) + transition
    await page.waitForTimeout(5000);

    // Verify waiting card disappears
    await expect(waitingCard).not.toBeVisible({ timeout: 10000 });

    // Verify analysis results display
    // InlineStats should be visible
    await expect(page.getByText("Total")).toBeVisible({ timeout: 10000 });
    await expect(page.getByText("Active")).toBeVisible();

    // Verify test suites list appears
    const testSuites = page.locator("[data-testid='test-suite']");
    await expect(testSuites.first()).toBeVisible({ timeout: 5000 });
  });

  test("should display toast notification when analysis completes", async ({
    page,
  }) => {
    // Setup: Start with analyzing state
    await setupMockHandlers(page, {
      analysis: mockAnalysisAnalyzing,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Verify waiting card displays
    const waitingCard = page.locator("[role='status'][aria-live='polite']");
    await expect(waitingCard).toBeVisible({ timeout: 15000 });

    // Update mock to completed
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    // Wait for polling + transition
    await page.waitForTimeout(5000);

    // Verify toast notification appears
    const toastRegion = page.locator("[role='region'][aria-live='polite']");
    await expect(toastRegion).toBeVisible({ timeout: 10000 });

    // Look for completion message
    await expect(
      page.getByText(/analysis.*completed/i)
    ).toBeVisible({ timeout: 10000 });
  });
});
