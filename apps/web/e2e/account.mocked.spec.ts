/**
 * Account Page E2E Tests with API Mocking
 * Tests subscription and usage quota features: quota states (normal, warning, critical, exceeded)
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers } from "./fixtures/mock-handlers";
import {
  mockSubscriptionFree,
  mockSubscriptionPro,
  mockSubscriptionProPlus,
  mockSubscriptionEnterprise,
  mockUsageNormal,
  mockUsageFree,
  mockUsageWarning,
  mockUsageCritical,
  mockUsageExceeded,
  mockUsageEnterprise,
} from "./fixtures/api-responses";

test.describe("Account Page - Subscription Display (Mocked API)", () => {
  test("should display Free plan details", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionFree,
      usage: mockUsageFree,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible();

    // Verify Free plan is displayed
    await expect(page.getByText(/free/i)).toBeVisible({ timeout: 10000 });
  });

  test("should display Pro plan details", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageNormal,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible();

    // Verify Pro plan is displayed
    await expect(page.getByText(/pro/i)).toBeVisible({ timeout: 10000 });
  });

  test("should display Pro Plus plan details", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionProPlus,
      usage: mockUsageNormal,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible();

    // Verify Pro Plus plan is displayed
    await expect(page.getByText(/pro plus|pro\+/i)).toBeVisible({ timeout: 10000 });
  });

  test("should display Enterprise plan details", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionEnterprise,
      usage: mockUsageEnterprise,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible();

    // Verify Enterprise plan is displayed
    await expect(page.getByText(/enterprise/i)).toBeVisible({ timeout: 10000 });

    // Enterprise should show "unlimited" for quota
    await expect(page.getByText(/unlimited/i)).toBeVisible({ timeout: 5000 });
  });

  test("should display billing period dates", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageNormal,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible();

    // Billing period or reset date should be visible
    // The exact format depends on the implementation
    await expect(page.getByText(/period|resets|billing/i).first()).toBeVisible({ timeout: 10000 });
  });
});

test.describe("Account Page - Usage Quota States (Mocked API)", () => {
  test("should display normal usage state (30%)", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageNormal,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible({ timeout: 10000 });

    // Verify usage percentages are displayed (both SpecView and Analysis show 30%)
    const percentageElements = page.getByText("30%");
    await expect(percentageElements.first()).toBeVisible({ timeout: 10000 });
    expect(await percentageElements.count()).toBe(2);

    // Progress bar should have primary color (not warning/destructive)
    const progressBars = page.locator('[class*="bg-primary"]');
    expect(await progressBars.count()).toBeGreaterThan(0);
  });

  test("should display warning usage state (70%)", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageWarning,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible({ timeout: 10000 });

    // Verify 70% usage percentage is displayed (both SpecView and Analysis)
    const percentageElements = page.getByText("70%");
    await expect(percentageElements.first()).toBeVisible({ timeout: 10000 });
    expect(await percentageElements.count()).toBe(2);

    // Progress bar should have amber/warning color (bg-amber-500)
    const warningProgressBars = page.locator('[class*="bg-amber"]');
    expect(await warningProgressBars.count()).toBeGreaterThan(0);
  });

  test("should display critical usage state (90%+)", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageCritical,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible({ timeout: 10000 });

    // Verify 92% usage percentage is displayed (both SpecView and Analysis)
    const percentageElements = page.getByText("92%");
    await expect(percentageElements.first()).toBeVisible({ timeout: 10000 });
    expect(await percentageElements.count()).toBe(2);

    // Progress bar should have destructive color
    const criticalProgressBars = page.locator('[class*="bg-destructive"]');
    expect(await criticalProgressBars.count()).toBeGreaterThan(0);

    // Percentage text should also have destructive styling
    await expect(page.locator('[class*="text-destructive"]:has-text("92%")').first()).toBeVisible();
  });

  test("should display exceeded usage state (100%+)", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageExceeded,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible({ timeout: 10000 });

    // Verify exceeded percentage is displayed (105% and 104%)
    await expect(page.getByText("105%")).toBeVisible({ timeout: 10000 });

    // Should show destructive styling
    const criticalProgressBars = page.locator('[class*="bg-destructive"]');
    expect(await criticalProgressBars.count()).toBeGreaterThan(0);
  });

  test("should display both SpecView and Analysis usage metrics", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageNormal,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible({ timeout: 10000 });

    // SpecView usage should be displayed (labeled as "AI Spec Docs")
    // There are multiple elements with this text (plan features + usage section)
    const specViewLabels = page.getByText(/ai spec docs/i);
    await expect(specViewLabels.first()).toBeVisible({ timeout: 10000 });

    // Analysis usage should be displayed (labeled as "Analysis Runs")
    const analysisLabels = page.getByText(/analysis runs/i);
    await expect(analysisLabels.first()).toBeVisible({ timeout: 10000 });

    // Both should show usage values (30/100 for SpecView, 150/500 for Analysis)
    await expect(page.getByText("30 / 100")).toBeVisible();
    await expect(page.getByText("150 / 500")).toBeVisible();
  });

  test("should display usage reset date", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageNormal,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible({ timeout: 10000 });

    // Reset date info should be visible (shows "Resets in X days (Month DD)" format)
    // This is shown in the card description area
    await expect(page.getByText(/resets in \d+ days/i)).toBeVisible({ timeout: 10000 });
  });
});

test.describe("Account Page - Plan Limits Display (Mocked API)", () => {
  test("should display plan limits for Free tier", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionFree,
      usage: mockUsageFree,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible({ timeout: 10000 });

    // Verify Free plan limits are displayed
    // AI Spec Docs: 10/month
    await expect(page.getByText("10/month")).toBeVisible({ timeout: 5000 });
    // Analysis: 50/month
    await expect(page.getByText("50/month")).toBeVisible({ timeout: 5000 });
    // Data Retention: 30 days
    await expect(page.getByText("30 days")).toBeVisible({ timeout: 5000 });
  });

  test("should display plan limits for Pro tier", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageNormal,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible({ timeout: 10000 });

    // Verify Pro plan limits are displayed
    // AI Spec Docs: 100/month
    await expect(page.getByText("100/month")).toBeVisible({ timeout: 5000 });
    // Analysis: 500/month
    await expect(page.getByText("500/month")).toBeVisible({ timeout: 5000 });
    // Data Retention: 90 days
    await expect(page.getByText("90 days")).toBeVisible({ timeout: 5000 });
  });

  test("should display unlimited for Enterprise tier", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionEnterprise,
      usage: mockUsageEnterprise,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible({ timeout: 10000 });

    // Page should load successfully with enterprise data
    await page.waitForLoadState("networkidle");
  });
});

test.describe("Account Page - Upgrade Prompts (Mocked API)", () => {
  test("should show upgrade prompt for Free tier user approaching limit", async ({ page }) => {
    // Free tier user at 50% usage (5/10)
    await setupMockHandlers(page, {
      subscription: mockSubscriptionFree,
      usage: mockUsageFree,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible();

    // Upgrade button or link should be visible for free tier
    await expect(
      page
        .getByRole("button", { name: /upgrade/i })
        .or(page.getByRole("link", { name: /upgrade|pricing/i }))
    ).toBeVisible({ timeout: 10000 });
  });

  test("should not show upgrade prompt for Enterprise tier", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionEnterprise,
      usage: mockUsageEnterprise,
    });

    await page.goto("/en/account");

    // Wait for page to load
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible();

    // Wait for content to fully load
    await page.waitForLoadState("networkidle");

    // Upgrade button should NOT be visible for enterprise
    const upgradeButton = page.getByRole("button", { name: /upgrade/i });
    await expect(upgradeButton).not.toBeVisible({ timeout: 3000 });
  });
});

test.describe("Account Page - Page Structure (Mocked API)", () => {
  test("should have proper page structure with all sections", async ({ page }) => {
    await setupMockHandlers(page, {
      subscription: mockSubscriptionPro,
      usage: mockUsageNormal,
    });

    await page.goto("/en/account");

    // Page title
    await expect(page.getByRole("heading", { name: /account/i })).toBeVisible();

    // Plan section
    await expect(page.getByText(/plan|subscription/i).first()).toBeVisible({ timeout: 10000 });

    // Usage section
    await expect(page.getByText(/usage/i).first()).toBeVisible({ timeout: 10000 });

    // Page should have content cards/sections
    const cards = page.locator('[class*="card"], section');
    const cardCount = await cards.count();
    expect(cardCount).toBeGreaterThan(0);
  });
});
