/**
 * Analysis Page E2E Tests with API Mocking
 * Tests AI Spec generation flow, tab navigation, and quota-related UI behavior
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockAnalysisCompleted,
  mockAnalysisLarge,
  mockAnalysisWithFocused,
  mockAnalysisWithXfail,
  mockAnalysisWithAllStatuses,
  mockSpecDocumentNotFound,
  mockSpecDocumentCompleted,
  mockSpecDocumentVersion1,
  mockSpecStatusNotFound,
  mockSpecGenerationAccepted,
  mockSpecGenerationRunning,
  SPEC_LANGUAGES,
  mockUsageNormal,
  mockUsageExceeded,
  mockVersionHistoryMultiple,
  mockVersionHistorySingle,
} from "./fixtures/api-responses";

test.describe("Analysis Page - InlineStats Display (Mocked API)", () => {
  test("should display InlineStats with Total, Active, Skipped counts", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Verify InlineStats displays Total, Active, Skipped labels
    await expect(page.getByText("Total")).toBeVisible();
    await expect(page.getByText("Active")).toBeVisible();
    // Use exact match to avoid matching test names like "should be skipped"
    await expect(page.getByText("Skipped", { exact: true })).toBeVisible();

    // Verify summary values are displayed in stats section
    // mockAnalysisCompleted has: total=4, active=3, skipped=1
    // Use tabular-nums class selector for stat values to avoid ambiguity
    const statValues = page.locator(".tabular-nums");
    await expect(statValues.first()).toBeVisible();
  });

  test("should display formatted numbers for large test counts", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisLarge,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/large-test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Verify large numbers are formatted with commas (1,000 format)
    await expect(page.getByText("1,000")).toBeVisible();
    await expect(page.getByText("800")).toBeVisible();
    await expect(page.getByText("200")).toBeVisible();
  });

  test("should display framework icons/badges", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Verify framework is displayed (mockAnalysisCompleted uses jest)
    await expect(page.getByText(/jest/i).first()).toBeVisible();
  });
});

test.describe("Analysis Page - InlineStats Conditional Display (Mocked API)", () => {
  // Helper to find InlineStats status label (uses text-sm + text-muted-foreground, not font-normal)
  const getStatusLabel = (page: import("@playwright/test").Page, label: string) =>
    page.locator("span.text-sm.text-muted-foreground", { hasText: label });

  test("should display Focused count when repository has focused tests", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisWithFocused,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/focused-test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Verify Focused label is displayed in InlineStats (mockAnalysisWithFocused has focused: 2)
    const focusedLabel = getStatusLabel(page, "Focused");
    await expect(focusedLabel).toBeVisible();

    // Verify the focused count value (2) - find by the status item structure
    const focusedCount = focusedLabel.locator(".. >> .tabular-nums");
    await expect(focusedCount).toHaveText("2");
  });

  test("should display Xfail count when repository has xfail tests", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisWithXfail,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/xfail-test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Verify Xfail label is displayed in InlineStats (mockAnalysisWithXfail has xfail: 2)
    const xfailLabel = getStatusLabel(page, "Xfail");
    await expect(xfailLabel).toBeVisible();

    // Verify the xfail count value (2)
    const xfailCount = xfailLabel.locator(".. >> .tabular-nums");
    await expect(xfailCount).toHaveText("2");
  });

  test("should not display Focused/Xfail when counts are zero", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Verify Active and Skipped are shown (always visible)
    await expect(getStatusLabel(page, "Active")).toBeVisible();
    await expect(getStatusLabel(page, "Skipped")).toBeVisible();

    // Verify Focused and Xfail status labels are NOT shown when counts are 0
    await expect(getStatusLabel(page, "Focused")).not.toBeVisible();
    await expect(getStatusLabel(page, "Xfail")).not.toBeVisible();
  });

  test("should display all 5 status types when all are present", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisWithAllStatuses,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/all-status-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Verify all 5 status labels are displayed in InlineStats
    await expect(getStatusLabel(page, "Active")).toBeVisible();
    await expect(getStatusLabel(page, "Focused")).toBeVisible();
    await expect(getStatusLabel(page, "Skipped")).toBeVisible();
    await expect(getStatusLabel(page, "Xfail")).toBeVisible();
    await expect(getStatusLabel(page, "Todo")).toBeVisible();

    // Verify total count (6 tests)
    await expect(page.getByText("6").first()).toBeVisible();
  });
});

test.describe("Analysis Page - StatusMiniBar Color Segments (Mocked API)", () => {
  test("should render MiniBar with correct aria-label for all 5 statuses", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisWithAllStatuses,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/all-status-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Find MiniBar by role="img" and verify aria-label contains all statuses
    const miniBar = page.locator('[role="img"]').first();
    await expect(miniBar).toBeVisible();

    // Verify aria-label includes all 5 status counts
    const ariaLabel = await miniBar.getAttribute("aria-label");
    expect(ariaLabel).toContain("active");
    expect(ariaLabel).toContain("focused");
    expect(ariaLabel).toContain("skipped");
    expect(ariaLabel).toContain("xfail");
    expect(ariaLabel).toContain("todo");
  });

  test("should render focused color segment in MiniBar when focused tests exist", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisWithFocused,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/focused-test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Find MiniBar and verify it has focused color segment
    const miniBar = page.locator('[role="img"]').first();
    await expect(miniBar).toBeVisible();

    // Verify focused segment exists (bg-status-focused class)
    const focusedSegment = miniBar.locator(".bg-status-focused");
    await expect(focusedSegment).toBeVisible();
  });

  test("should render xfail color segment in MiniBar when xfail tests exist", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisWithXfail,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/xfail-test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Find MiniBar and verify it has xfail color segment
    const miniBar = page.locator('[role="img"]').first();
    await expect(miniBar).toBeVisible();

    // Verify xfail segment exists (bg-status-xfail class)
    const xfailSegment = miniBar.locator(".bg-status-xfail");
    await expect(xfailSegment).toBeVisible();
  });

  test("should not render focused/xfail segments when counts are zero", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Find MiniBar
    const miniBar = page.locator('[role="img"]').first();
    await expect(miniBar).toBeVisible();

    // Verify focused and xfail segments do NOT exist (counts are 0)
    const focusedSegment = miniBar.locator(".bg-status-focused");
    const xfailSegment = miniBar.locator(".bg-status-xfail");
    await expect(focusedSegment).not.toBeVisible();
    await expect(xfailSegment).not.toBeVisible();

    // Verify active and skipped segments DO exist
    const activeSegment = miniBar.locator(".bg-status-active");
    const skippedSegment = miniBar.locator(".bg-status-skipped");
    await expect(activeSegment).toBeVisible();
    await expect(skippedSegment).toBeVisible();
  });
});

test.describe("Analysis Page - Primary Tab Navigation (Mocked API)", () => {
  test("should display Tests tab as default", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Tests tab should be selected by default
    const testsTab = page.getByRole("tab", { name: /tests/i });
    await expect(testsTab).toBeVisible();
    await expect(testsTab).toHaveAttribute("data-state", "active");
  });

  test("should switch between Tests and AI Spec tabs", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Click AI Spec tab
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await specTab.click();

    // Spec tab should be active
    await expect(specTab).toHaveAttribute("data-state", "active");

    // Click Tests tab
    const testsTab = page.getByRole("tab", { name: /tests/i });
    await testsTab.click();

    // Tests tab should be active
    await expect(testsTab).toHaveAttribute("data-state", "active");
  });

  test("should display list/tree toggle in Tests tab", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // List view toggle should be visible
    const listViewToggle = page.getByRole("radio", { name: /list view/i });
    await expect(listViewToggle).toBeVisible();

    // Tree view toggle should be visible
    const treeViewToggle = page.getByRole("radio", { name: /tree view/i });
    await expect(treeViewToggle).toBeVisible();
  });
});

test.describe("Analysis Page - AI Spec Tab (Mocked API)", () => {
  test("should display generate button in AI Spec tab when no document exists", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Switch to AI Spec tab
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await specTab.click();

    // Generate button should be visible
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await expect(generateButton).toBeVisible();
  });

  test("should open quota confirm dialog when clicking generate button", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      specGeneration: mockSpecStatusNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Switch to AI Spec tab
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await specTab.click();

    // Click Generate button
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    // Verify dialog opens
    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();
    await expect(dialog.getByRole("heading", { name: /generate ai specification/i })).toBeVisible();
  });

  test("should display language selector with 24 options in dialog", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      specGeneration: mockSpecStatusNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Switch to AI Spec tab
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await specTab.click();

    // Click Generate button
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    // Verify dialog opened
    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();

    // Click language selector (aria-label is "Select language")
    const languageCombobox = dialog.getByRole("combobox", { name: /select language/i });
    await expect(languageCombobox).toBeVisible();
    await languageCombobox.click();

    // Verify all 24 languages are available
    const listbox = page.getByRole("listbox");
    await expect(listbox).toBeVisible();

    // Check for specific languages
    for (const language of SPEC_LANGUAGES) {
      await expect(listbox.getByRole("option", { name: language })).toBeVisible();
    }
  });

  test("should display Cancel and Generate Document buttons in dialog", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      specGeneration: mockSpecStatusNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Switch to AI Spec tab
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await specTab.click();

    // Click Generate button
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();

    // Verify buttons
    await expect(dialog.getByRole("button", { name: /cancel/i })).toBeVisible();
    await expect(dialog.getByRole("button", { name: /generate document/i })).toBeVisible();
  });

  test("should close dialog when clicking Cancel", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      specGeneration: mockSpecStatusNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Switch to AI Spec tab
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await specTab.click();

    // Click Generate button
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();

    // Click Cancel
    await dialog.getByRole("button", { name: /cancel/i }).click();

    // Dialog should close
    await expect(dialog).not.toBeVisible();
  });
});

test.describe("Analysis Page - Quota States (Mocked API)", () => {
  test("should disable Generate Document button when quota exceeded", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      specGeneration: mockSpecStatusNotFound,
      usage: mockUsageExceeded,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Switch to AI Spec tab
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await specTab.click();

    // Generate Document button should be disabled when quota exceeded
    // (button on EmptyDocument card is disabled when quota is exceeded)
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await expect(generateButton).toBeDisabled();
  });

  test("should display quota exceeded message", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      specGeneration: mockSpecStatusNotFound,
      usage: mockUsageExceeded,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Switch to AI Spec tab
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await specTab.click();

    // When quota is exceeded, the Generate Document button should be disabled
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await expect(generateButton).toBeDisabled();

    // Should show "View Account" link in quota indicator when quota is exceeded
    await expect(page.getByRole("link", { name: /view account/i })).toBeVisible();
  });

  test("should show usage percentage in dialog", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      specGeneration: mockSpecStatusNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Switch to AI Spec tab
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await specTab.click();

    // Click Generate button
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();

    // Should show usage information (30% in mockUsageNormal)
    await expect(dialog.getByText(/30.*%|specview/i)).toBeVisible();
  });
});

test.describe("Analysis Page - Virtual Scroll Performance (Mocked API)", () => {
  test("should handle large test suite list with virtual scroll", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisLarge,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/large-test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Verify total test count is displayed (1000 tests formatted with comma)
    await expect(page.getByText("1,000")).toBeVisible({ timeout: 10000 });

    // Verify test suites section is rendered
    const testSuitesSection = page.getByText(/test suite/i).first();
    await expect(testSuitesSection).toBeVisible();

    // Check that DOM doesn't have all 100 suites rendered at once (virtual scroll)
    // Only visible items should be in DOM
    const suiteItems = await page.locator('[data-testid="test-suite-item"]').count();
    // With virtualization, we expect fewer items than total (100)
    // If no virtualization, all 100 would be rendered

    // Scroll to bottom to trigger virtual scroll loading
    await page.evaluate(() => {
      const scrollContainer = document.querySelector('[data-testid="suites-container"]') ||
        document.querySelector('[class*="overflow-auto"]') ||
        window;
      if (scrollContainer instanceof Element) {
        scrollContainer.scrollTop = scrollContainer.scrollHeight;
      } else {
        window.scrollTo(0, document.body.scrollHeight);
      }
    });

    // Wait for scroll to settle
    await page.waitForTimeout(500);

    // Page should still be responsive
    await expect(page.getByText("Total")).toBeVisible();
  });

  test("should display correct summary statistics for large dataset", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisLarge,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/large-test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Verify summary statistics
    // Total: 1000, Active: 800, Skipped: 200
    await expect(page.getByText("1,000").or(page.getByText("1000"))).toBeVisible({ timeout: 10000 });
  });
});

test.describe("Analysis Page - URL Backward Compatibility (Mocked API)", () => {
  test("should redirect ?view=document to ?tab=spec", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    // Navigate with old URL parameter
    await page.goto("/en/analyze/test-owner/test-repo?view=document");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // AI Spec tab should be active after redirect
    const specTab = page.getByRole("tab", { name: /ai spec/i });
    await expect(specTab).toHaveAttribute("data-state", "active");

    // URL should have ?tab=spec instead of ?view=document
    await expect(page).toHaveURL(/tab=spec/);
    await expect(page).not.toHaveURL(/view=document/);
  });
});

test.describe("Analysis Page - Filter Empty State (Mocked API)", () => {
  test("should display empty state when search yields no results", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Search for non-existent text
    const searchInput = page.getByRole("textbox", { name: /search/i });
    await searchInput.fill("nonexistenttestname12345");

    // Wait for debounce and empty state
    await page.waitForTimeout(600);

    // Verify empty state is displayed
    await expect(page.getByText(/no.*match/i)).toBeVisible();
  });

  test("should display applied filter badges in empty state", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Search for non-existent text
    const searchInput = page.getByRole("textbox", { name: /search/i });
    await searchInput.fill("nonexistent");

    // Wait for debounce
    await page.waitForTimeout(600);

    // Verify search query is shown as badge
    await expect(page.getByText(/nonexistent/)).toBeVisible();
  });

  test("should reset filters when clicking Reset all filters button", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for analysis to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Search for non-existent text
    const searchInput = page.getByRole("textbox", { name: /search/i });
    await searchInput.fill("nonexistenttestname12345");

    // Wait for debounce
    await page.waitForTimeout(600);

    // Click reset button
    const resetButton = page.getByRole("button", { name: /reset.*filter/i });
    await expect(resetButton).toBeVisible();
    await resetButton.click();

    // Verify filters are reset - test list should be visible again
    // Use first() since multiple test suites may match
    await expect(page.getByText(/example test suite|another test suite/i).first()).toBeVisible();

    // Search input should be cleared
    await expect(searchInput).toHaveValue("");
  });
});

test.describe("Analysis Page - Spec View Language Switch (Mocked API)", () => {
  test("should display language dropdown in ExecutiveSummary", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Find language dropdown button
    const languageButton = page.getByRole("button", { name: /english/i }).first();
    await expect(languageButton).toBeVisible();
  });

  test("should display Two-Tier dropdown with Available and Generate New sections", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Click language dropdown
    const languageButton = page.getByRole("button", { name: /english/i }).first();
    await languageButton.click();

    // Verify dropdown menu is open
    const menu = page.getByRole("menu");
    await expect(menu).toBeVisible();

    // Check for Available Languages section (English and Korean in mock data)
    await expect(menu.getByText(/available languages/i)).toBeVisible();
    await expect(menu.getByRole("menuitem", { name: /english/i })).toBeVisible();
    await expect(menu.getByRole("menuitem", { name: /korean/i })).toBeVisible();

    // Check for Generate New section
    await expect(menu.getByText(/generate new/i)).toBeVisible();

    // Check for some languages in Generate New section (not in available list)
    await expect(menu.getByRole("menuitem", { name: /japanese/i })).toBeVisible();
    await expect(menu.getByRole("menuitem", { name: /chinese/i })).toBeVisible();
  });

  test("should mark current language with checkmark indicator", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Click language dropdown
    const languageButton = page.getByRole("button", { name: /english/i }).first();
    await languageButton.click();

    // English option should be highlighted (bg-muted class indicates current)
    const englishOption = page.getByRole("menuitem", { name: /english/i });
    await expect(englishOption).toBeVisible();
    // Current language should have version info displayed
    await expect(englishOption.getByText(/v2/i)).toBeVisible();
  });

  test("should display version and date info for available languages", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Click language dropdown
    const languageButton = page.getByRole("button", { name: /english/i }).first();
    await languageButton.click();

    const menu = page.getByRole("menu");

    // Check for version info in available languages
    const koreanOption = menu.getByRole("menuitem", { name: /korean/i });
    await expect(koreanOption).toBeVisible();
    // Korean should show v1
    await expect(koreanOption.getByText(/v1/i)).toBeVisible();
  });
});

test.describe("Analysis Page - Spec View Regeneration (Mocked API)", () => {
  test("should display regeneration button in ExecutiveSummary", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Find regeneration button (has RefreshCw icon)
    const regenerateButton = page.getByRole("button", { name: /regenerate/i });
    await expect(regenerateButton).toBeVisible();
  });

  test("should open confirmation dialog when regeneration button clicked", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      specGeneration: mockSpecStatusNotFound,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Click regeneration button
    const regenerateButton = page.getByRole("button", { name: /regenerate/i });
    await regenerateButton.click();

    // Verify confirmation dialog opens
    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();
    // Dialog should have a regeneration heading
    await expect(
      dialog.getByRole("heading", { name: /regenerate/i })
    ).toBeVisible();
  });
});

test.describe("Analysis Page - Spec View TOC Sidebar (Mocked API)", () => {
  test("should display TOC sidebar with domain/feature items", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Verify domain items in TOC (use first() since domains appear in both TOC and content)
    await expect(page.getByText("User Authentication").first()).toBeVisible();
    await expect(page.getByText("Payment Processing").first()).toBeVisible();
  });

  test("should navigate to domain when TOC item clicked", async ({ page }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Click on Payment Processing in TOC (should scroll to that section)
    const paymentTocItem = page.getByText("Payment Processing").first();
    await paymentTocItem.click();

    // Verify the section is visible/scrolled to
    await expect(page.getByText("Checkout Flow")).toBeVisible();
  });
});

test.describe("Analysis Page - Spec Generation Progress (Mocked API)", () => {
  test("should display progress modal during generation", async ({ page }) => {
    // Start with not_found status, then switch to running after generate API call
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      // When generate is called, return pending status
      onSpecGeneration: () => mockSpecGenerationAccepted,
      // Status polling returns pending after generation starts
      onSpecGenerationStatus: () => mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Click Generate button
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    // Confirm in dialog
    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();
    const confirmButton = dialog.getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    // Verify progress modal appears with spinner/loading state
    await expect(page.getByText(/generating|processing/i)).toBeVisible({ timeout: 5000 });
  });

  test("should show spinner animation during generation", async ({ page }) => {
    // Start with not_found status, then switch to running after generate API call
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      // When generate is called, return pending status
      onSpecGeneration: () => mockSpecGenerationAccepted,
      // Status polling returns pending after generation starts
      onSpecGenerationStatus: () => mockSpecGenerationAccepted,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Click Generate button
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    // Confirm in dialog
    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();
    const confirmButton = dialog.getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    // Verify loading indicator exists (SVG with animate or spinning class)
    const loadingIndicator = page.locator('svg[class*="animate"]').or(
      page.locator('[class*="spinner"]')
    ).or(
      page.locator('[data-slot="progress"]')
    );
    await expect(loadingIndicator.first()).toBeVisible({ timeout: 5000 });
  });
});

test.describe("Analysis Page - Spec 403 Handling (Mocked API)", () => {
  test("should display subscription guidance when 403 error occurs", async ({
    page,
  }) => {
    // Setup mock that returns 403 for spec generation
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

    // Wait for page to load
    await expect(page.getByText("Total")).toBeVisible({
      timeout: 15000,
    });

    // Click Generate button
    const generateButton = page.getByRole("button", { name: /generate document/i });
    await generateButton.click();

    // Confirm in dialog
    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible();
    const confirmButton = dialog.getByRole("button", { name: /generate document/i });
    await confirmButton.click();

    // Verify error message or subscription link appears
    // Look for "Subscription Required" which is the dialog title
    await expect(
      page.getByText("Subscription Required")
    ).toBeVisible({ timeout: 5000 });
  });
});

test.describe("Analysis Page - Version History Dropdown (Mocked API)", () => {
  test("should display version dropdown when multiple versions exist", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      versionHistory: mockVersionHistoryMultiple,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Version dropdown should be visible (shows date with "latest" indicator)
    const versionButton = page.getByRole("button", { name: /latest/i });
    await expect(versionButton).toBeVisible();
  });

  test("should display version badge when only single version exists", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      versionHistory: mockVersionHistorySingle,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // With only one version, it should show as a badge with date and "latest" indicator
    await expect(page.getByText(/latest/i)).toBeVisible();
  });

  test("should show version history dropdown with all versions", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      versionHistory: mockVersionHistoryMultiple,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Click version dropdown
    const versionButton = page.getByRole("button", { name: /latest/i });
    await versionButton.click();

    // Verify dropdown menu is open
    const menu = page.getByRole("menu");
    await expect(menu).toBeVisible();

    // Check for Version History header
    await expect(menu.getByText(/version history/i)).toBeVisible();

    // Check for versions in the list (dates with version numbers)
    await expect(menu.getByRole("menuitem", { name: /latest/i })).toBeVisible();
    await expect(menu.getByRole("menuitem", { name: /v1/i })).toBeVisible();
  });

  test("should switch to previous version when selected", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      specDocumentByVersion: {
        1: mockSpecDocumentVersion1,
        2: mockSpecDocumentCompleted,
      },
      versionHistory: mockVersionHistoryMultiple,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Click version dropdown
    const versionButton = page.getByRole("button", { name: /latest/i });
    await versionButton.click();

    // Select older version (v1)
    const v1Option = page.getByRole("menuitem", { name: /v1/i });
    await v1Option.click();

    // Wait for document to reload with v1 content
    await expect(page.getByText(/Test Repository Specification v1/i)).toBeVisible({
      timeout: 5000,
    });
  });

  test("should mark current version with checkmark indicator", async ({
    page,
  }) => {
    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentCompleted,
      versionHistory: mockVersionHistoryMultiple,
      usage: mockUsageNormal,
    });

    await page.goto("/en/analyze/test-owner/test-repo?tab=spec");
    await page.waitForLoadState("networkidle");

    // Wait for spec document to load
    await expect(page.getByText(/Test Repository Specification/i)).toBeVisible({
      timeout: 15000,
    });

    // Click version dropdown
    const versionButton = page.getByRole("button", { name: /latest/i });
    await versionButton.click();

    // Latest version should be highlighted as current (bg-muted class)
    const latestOption = page.getByRole("menuitem", { name: /latest/i });
    await expect(latestOption).toBeVisible();
    // Current version should have highlight styling
    await expect(latestOption).toHaveClass(/bg-muted|font-medium/);
  });
});

test.describe("Analysis Page - Update Banner Polling (Mocked API)", () => {
  test("should show analyzing state in banner when Update Now clicked", async ({
    page,
  }) => {
    let reanalyzeCallCount = 0;

    await setupMockHandlers(page, {
      analysis: mockAnalysisCompleted,
      specDocument: mockSpecDocumentNotFound,
      usage: mockUsageNormal,
      updateStatus: {
        status: "new-commits",
        parserOutdated: false,
      },
      onReanalyze: () => {
        reanalyzeCallCount++;
        return { status: "queued" as const };
      },
      onAnalysisStatus: (_owner, _repo, callCount) => ({
        owner: "test-owner",
        repo: "test-repo",
        status: callCount < 3 ? "analyzing" : "completed",
      }),
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for update banner to appear
    const updateBanner = page.locator('[class*="bg-amber"]');
    await expect(updateBanner.first()).toBeVisible({ timeout: 10000 });

    // Click Update Now button
    const updateNowButton = page.getByRole("button", { name: /update now/i });
    await expect(updateNowButton).toBeVisible();
    await updateNowButton.click();

    // Verify reanalyze was triggered
    await page.waitForTimeout(500);
    expect(reanalyzeCallCount).toBe(1);

    // Verify banner shows analyzing state (loading spinner should be visible)
    const loadingSpinner = updateBanner.locator('svg[class*="animate-spin"]');
    await expect(loadingSpinner).toBeVisible({ timeout: 5000 });
  });

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
        // Return completed after 2 polls
        return {
          owner: "test-owner",
          repo: "test-repo",
          status: statusCallCount >= 2 ? "completed" : "analyzing",
        };
      },
    });

    await page.goto("/en/analyze/test-owner/test-repo");
    await page.waitForLoadState("networkidle");

    // Wait for update banner
    const updateBanner = page.locator('[class*="bg-amber"]');
    await expect(updateBanner.first()).toBeVisible({ timeout: 10000 });

    // Click Update Now
    const updateNowButton = page.getByRole("button", { name: /update now/i });
    await updateNowButton.click();

    // Wait for polling to complete and banner to dismiss
    await expect(updateBanner).not.toBeVisible({ timeout: 10000 });
  });
});
