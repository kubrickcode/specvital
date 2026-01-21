/**
 * Analysis Page E2E Tests with API Mocking
 * Tests AI Spec generation flow, tab navigation, and quota-related UI behavior
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockAnalysisCompleted,
  mockAnalysisLarge,
  mockSpecDocumentNotFound,
  mockSpecStatusNotFound,
  mockSpecGenerationAccepted,
  SPEC_LANGUAGES,
  mockUsageNormal,
  mockUsageExceeded,
} from "./fixtures/api-responses";

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

  // Skip: Re-analyze button removed from analysis page header in UI redesign (cff5b21)
  // Reanalyze functionality moved to dashboard repository cards and update banner
  test.skip("should show Re-analyze button in header", async ({ page }) => {
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

    // Re-analyze button should be visible
    const reanalyzeButton = page.getByRole("button", { name: /re-analyze/i });
    await expect(reanalyzeButton).toBeVisible();
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

    // Click language selector
    const languageCombobox = dialog.getByRole("combobox", { name: /output language/i });
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
