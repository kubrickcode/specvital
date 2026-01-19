import { expect, test } from "@playwright/test";

test.describe("Analysis Page - AI Spec Generation (Authenticated)", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to analysis page for a known repository
    await page.goto("/en/analyze/facebook/react");
  });

  test("should display AI Spec generate button for authenticated user", async ({
    page,
  }) => {
    // Wait for analysis data to load
    await expect(
      page.getByRole("heading", { name: "facebook/react" })
    ).toBeVisible();

    // Verify Generate AI Spec button is visible
    const generateButton = page.getByRole("button", {
      name: "Generate test specification document with AI",
    });
    await expect(generateButton).toBeVisible({ timeout: 30000 });
  });

  test("should show quota confirm dialog when clicking generate button", async ({
    page,
  }) => {
    // Wait for analysis data to load
    await expect(
      page.getByRole("heading", { name: "Test Statistics" })
    ).toBeVisible({ timeout: 30000 });

    // Click Generate AI Spec button
    const generateButton = page.getByRole("button", {
      name: "Generate test specification document with AI",
    });
    await generateButton.click();

    // Verify quota confirm dialog is shown
    await expect(
      page.getByRole("dialog", { name: "Generate AI Specification" })
    ).toBeVisible();

    // Verify dialog content
    await expect(
      page.getByText("This will use your AI Spec Doc generation quota")
    ).toBeVisible();

    // Verify language selection is present
    await expect(page.getByText("Output Language")).toBeVisible();
    await expect(page.getByRole("combobox")).toBeVisible();

    // Verify action buttons
    await expect(page.getByRole("button", { name: "Cancel" })).toBeVisible();
    await expect(page.getByRole("button", { name: "Generate" })).toBeVisible();
  });

  test("should allow language selection in quota confirm dialog", async ({
    page,
  }) => {
    // Wait for analysis data to load
    await expect(
      page.getByRole("heading", { name: "Test Statistics" })
    ).toBeVisible({ timeout: 30000 });

    // Click Generate AI Spec button
    await page
      .getByRole("button", {
        name: "Generate test specification document with AI",
      })
      .click();

    // Wait for dialog
    await expect(
      page.getByRole("dialog", { name: "Generate AI Specification" })
    ).toBeVisible();

    // Click language selector
    const languageSelector = page.getByRole("combobox");
    await languageSelector.click();

    // Verify language options are available
    // Language options depend on backend configuration
    const listbox = page.getByRole("listbox");
    await expect(listbox).toBeVisible();
  });

  test("should close quota confirm dialog on cancel", async ({ page }) => {
    // Wait for analysis data to load
    await expect(
      page.getByRole("heading", { name: "Test Statistics" })
    ).toBeVisible({ timeout: 30000 });

    // Click Generate AI Spec button
    await page
      .getByRole("button", {
        name: "Generate test specification document with AI",
      })
      .click();

    // Wait for dialog
    const dialog = page.getByRole("dialog", {
      name: "Generate AI Specification",
    });
    await expect(dialog).toBeVisible();

    // Click Cancel button
    await page.getByRole("button", { name: "Cancel" }).click();

    // Verify dialog is closed
    await expect(dialog).not.toBeVisible();
  });

  test("should show generate document button in quota confirm dialog", async ({
    page,
  }) => {
    // Wait for analysis data to load
    await expect(
      page.getByRole("heading", { name: "Test Statistics" })
    ).toBeVisible({ timeout: 30000 });

    // Click Generate AI Spec button
    await page
      .getByRole("button", {
        name: "Generate test specification document with AI",
      })
      .click();

    // Wait for dialog
    await expect(
      page.getByRole("dialog", { name: "Generate AI Specification" })
    ).toBeVisible();

    // Verify Generate Document button is present and enabled
    const generateDocButton = page.getByRole("button", {
      name: /Generate Document|Generate/i,
    });
    await expect(generateDocButton).toBeVisible();
    await expect(generateDocButton).toBeEnabled();
  });
});
