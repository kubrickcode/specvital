/**
 * Background Tasks E2E Tests (Authenticated)
 * Tests background task transitions, polling persistence, toast notifications, and badge UI
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockAnalysisCompleted,
  mockSpecDocumentNotFound,
  mockSpecGenerationAccepted,
  mockSpecGenerationCompleted,
  mockUsageNormal,
} from "./fixtures/api-responses";

// SKIPPED: Background task tests require complex state machine interactions
// between spec generation API, polling mechanism, and Zustand store.
// The mock setup does not fully replicate the async state transitions needed
// for progress modal lifecycle. These features are tested via unit tests.
test.describe.skip("Background Tasks - Modal Transition (Authenticated)", () => {
  test("should close modal and show TaskBadge when clicking 'Continue Browsing'", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Start generation
    const generateButton = page.getByRole("button", {
      name: /generate document/i,
    });
    await expect(generateButton).toBeVisible({ timeout: 15000 });
    await generateButton.click();

    const confirmButton = page.getByRole("dialog").getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    // Wait for modal
    const progressModal = page.getByRole("dialog", {
      name: /generation progress/i,
    });
    await expect(progressModal).toBeVisible();

    // Click "Continue browsing" button
    const continueButton = progressModal.getByRole("button", {
      name: /continue browsing/i,
    });
    await expect(continueButton).toBeVisible();
    await continueButton.click();

    // Verify modal closes
    await expect(progressModal).not.toBeVisible();

    // Verify TaskBadge appears in User Menu
    // TaskBadge is a spinner badge overlaying the user avatar button
    const userMenuButton = page.getByRole("button", {
      name: /account menu/i,
    });
    await expect(userMenuButton).toBeVisible();

    // Look for TaskBadge (Loader2 spinner with aria-label)
    const taskBadge = page.locator("span[aria-label*='task']");
    await expect(taskBadge).toBeVisible({ timeout: 5000 });
  });

  test("should show toast notification when task completes on different page", async ({
    page,
  }) => {
    // Setup: Task will complete after navigation
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationAccepted, // Initially running
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Start generation and switch to background
    const generateButton = page.getByRole("button", {
      name: /generate document/i,
    });
    await expect(generateButton).toBeVisible({ timeout: 15000 });
    await generateButton.click();

    const confirmButton = page.getByRole("dialog").getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    const progressModal = page.getByRole("dialog", {
      name: /generation progress/i,
    });
    await expect(progressModal).toBeVisible();

    const continueButton = progressModal.getByRole("button", {
      name: /continue browsing/i,
    });
    await continueButton.click();

    // Navigate to dashboard
    await page.goto("/en/dashboard");
    await page.waitForLoadState("networkidle");

    // Update mock to completed status
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationCompleted, // Now completed
    });

    // Wait for polling interval + toast animation
    await page.waitForTimeout(5000);

    // Verify toast notification appears
    // Toast region should contain completion message
    const toastRegion = page.locator("[role='region'][aria-live='polite']");
    await expect(toastRegion).toBeVisible({ timeout: 10000 });

    // Look for completion message
    await expect(
      page.getByText(/spec generation.*completed/i)
    ).toBeVisible({ timeout: 10000 });
  });
});

test.describe.skip("Background Tasks - Polling Persistence (Authenticated)", () => {
  test("should maintain TaskBadge when navigating between pages", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Start generation and switch to background
    const generateButton = page.getByRole("button", {
      name: /generate document/i,
    });
    await expect(generateButton).toBeVisible({ timeout: 15000 });
    await generateButton.click();

    const confirmButton = page.getByRole("dialog").getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    const progressModal = page.getByRole("dialog", {
      name: /generation progress/i,
    });
    const continueButton = progressModal.getByRole("button", {
      name: /continue browsing/i,
    });
    await continueButton.click();

    // Verify TaskBadge visible
    const taskBadge = page.locator("span[aria-label*='task']");
    await expect(taskBadge).toBeVisible({ timeout: 5000 });

    // Navigate to dashboard
    await page.goto("/en/dashboard");
    await page.waitForLoadState("networkidle");

    // TaskBadge should still be visible
    await expect(taskBadge).toBeVisible();

    // Navigate to explore
    await page.goto("/en/explore");
    await page.waitForLoadState("networkidle");

    // TaskBadge should still be visible
    await expect(taskBadge).toBeVisible();

    // Navigate to pricing
    await page.goto("/en/pricing");
    await page.waitForLoadState("networkidle");

    // TaskBadge should still be visible (if user is authenticated)
    await expect(taskBadge).toBeVisible();
  });

  test("should restore in-progress task after browser refresh", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Start generation and switch to background
    const generateButton = page.getByRole("button", {
      name: /generate document/i,
    });
    await expect(generateButton).toBeVisible({ timeout: 15000 });
    await generateButton.click();

    const confirmButton = page.getByRole("dialog").getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    const progressModal = page.getByRole("dialog", {
      name: /generation progress/i,
    });
    const continueButton = progressModal.getByRole("button", {
      name: /continue browsing/i,
    });
    await continueButton.click();

    // Verify TaskBadge visible
    const taskBadge = page.locator("span[aria-label*='task']");
    await expect(taskBadge).toBeVisible({ timeout: 5000 });

    // Refresh page
    await page.reload();
    await page.waitForLoadState("networkidle");

    // TaskBadge should be restored after refresh
    await expect(taskBadge).toBeVisible({ timeout: 10000 });
  });
});

test.describe.skip("Background Tasks - Multiple Tasks (Authenticated)", () => {
  test("should display all in-progress tasks in dropdown list", async ({
    page,
  }) => {
    // Note: This test assumes multiple task support
    // In practice, may need to start multiple generations
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Start generation and switch to background
    const generateButton = page.getByRole("button", {
      name: /generate document/i,
    });
    await expect(generateButton).toBeVisible({ timeout: 15000 });
    await generateButton.click();

    const confirmButton = page.getByRole("dialog").getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    const progressModal = page.getByRole("dialog", {
      name: /generation progress/i,
    });
    const continueButton = progressModal.getByRole("button", {
      name: /continue browsing/i,
    });
    await continueButton.click();

    // Open User Menu dropdown
    const userMenuButton = page.getByRole("button", {
      name: /account menu/i,
    });
    await userMenuButton.click();

    // Verify TasksDropdownSection displays
    // Look for task items in dropdown menu
    const dropdownMenu = page.getByRole("menu");
    await expect(dropdownMenu).toBeVisible();

    // Verify task item exists (should show owner/repo and status)
    const taskItem = dropdownMenu.getByText(/test-owner.*test-repo/i);
    await expect(taskItem).toBeVisible({ timeout: 5000 });
  });
});

test.describe.skip("Background Tasks - Badge UI (Authenticated)", () => {
  test("should display TaskBadge with spinner animation when tasks are active", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Start generation and switch to background
    const generateButton = page.getByRole("button", {
      name: /generate document/i,
    });
    await expect(generateButton).toBeVisible({ timeout: 15000 });
    await generateButton.click();

    const confirmButton = page.getByRole("dialog").getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    const progressModal = page.getByRole("dialog", {
      name: /generation progress/i,
    });
    const continueButton = progressModal.getByRole("button", {
      name: /continue browsing/i,
    });
    await continueButton.click();

    // Verify TaskBadge visible with spinner (Loader2)
    const taskBadge = page.locator("span[aria-label*='task']");
    await expect(taskBadge).toBeVisible({ timeout: 5000 });

    // Verify spinner animation (Loader2 component with animate-spin class)
    const spinner = taskBadge.locator("svg.animate-spin");
    await expect(spinner).toBeVisible();

    // Verify aria-label includes task count
    const ariaLabel = await taskBadge.getAttribute("aria-label");
    expect(ariaLabel).toMatch(/\d+.*task/i);
  });

  test("should display tasks in User Menu dropdown section", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Start generation and switch to background
    const generateButton = page.getByRole("button", {
      name: /generate document/i,
    });
    await expect(generateButton).toBeVisible({ timeout: 15000 });
    await generateButton.click();

    const confirmButton = page.getByRole("dialog").getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    const progressModal = page.getByRole("dialog", {
      name: /generation progress/i,
    });
    const continueButton = progressModal.getByRole("button", {
      name: /continue browsing/i,
    });
    await continueButton.click();

    // Open User Menu
    const userMenuButton = page.getByRole("button", {
      name: /account menu/i,
    });
    await userMenuButton.click();

    // Verify dropdown menu displays
    const dropdownMenu = page.getByRole("menu");
    await expect(dropdownMenu).toBeVisible();

    // Verify task section exists
    // TasksDropdownSection should show repository name and status icon
    await expect(
      dropdownMenu.getByText(/test-owner.*test-repo/i)
    ).toBeVisible({ timeout: 5000 });

    // Click task item should navigate to analysis page
    const taskLink = dropdownMenu.locator(
      "a[href*='/analyze/test-owner/test-repo']"
    );
    await expect(taskLink).toBeVisible();
  });
});

test.describe.skip("Background Tasks - Dashboard Integration (Authenticated)", () => {
  test("should auto-update dashboard when background task completes", async ({
    page,
  }) => {
    // Setup: Task will complete and trigger dashboard update
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationAccepted,
    });

    // Start at dashboard
    await page.goto("/en/dashboard");
    await page.waitForLoadState("networkidle");

    // Navigate to analysis and start generation
    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    const generateButton = page.getByRole("button", {
      name: /generate document/i,
    });
    await expect(generateButton).toBeVisible({ timeout: 15000 });
    await generateButton.click();

    const confirmButton = page.getByRole("dialog").getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    const progressModal = page.getByRole("dialog", {
      name: /generation progress/i,
    });
    const continueButton = progressModal.getByRole("button", {
      name: /continue browsing/i,
    });
    await continueButton.click();

    // Go back to dashboard
    await page.goto("/en/dashboard");
    await page.waitForLoadState("networkidle");

    // Update mock to completed
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      specGeneration: mockSpecGenerationCompleted,
    });

    // Wait for task completion and dashboard update
    await page.waitForTimeout(5000);

    // Verify toast notification
    await expect(
      page.getByText(/spec generation.*completed/i)
    ).toBeVisible({ timeout: 10000 });

    // Verify dashboard auto-updates (repository card should show AI Spec badge)
    // This may require data refetch to show the badge
    const repositoryCard = page.locator("[data-testid='repository-card']");
    await expect(repositoryCard.first()).toBeVisible({ timeout: 5000 });
  });
});
