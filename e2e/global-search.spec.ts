import { expect, test } from "@playwright/test";

test.describe("Global Search", () => {
  test.describe("Desktop", () => {
    test.beforeEach(async ({ page }) => {
      await page.setViewportSize({ width: 1280, height: 720 });
    });

    test("should open search dialog with Cmd/Ctrl+K shortcut", async ({ page }) => {
      await page.goto("/en");

      // Open with keyboard shortcut
      await page.keyboard.press("ControlOrMeta+k");

      // Verify dialog is open
      await expect(page.getByRole("dialog")).toBeVisible();
      await expect(page.getByPlaceholder(/search/i)).toBeVisible();
      await expect(page.getByPlaceholder(/search/i)).toBeFocused();
    });

    test("should close search dialog with Escape", async ({ page }) => {
      await page.goto("/en");

      // Open dialog
      await page.keyboard.press("ControlOrMeta+k");
      await expect(page.getByRole("dialog")).toBeVisible();

      // Close with Escape
      await page.keyboard.press("Escape");
      await expect(page.getByRole("dialog")).not.toBeVisible();
    });

    test("should open search dialog by clicking trigger button", async ({ page }) => {
      await page.goto("/en");

      // Click search trigger button
      await page
        .getByRole("button", { name: /search/i })
        .first()
        .click();

      // Verify dialog is open
      await expect(page.getByRole("dialog")).toBeVisible();
    });

    test("should display navigation items when dialog opens", async ({ page }) => {
      await page.goto("/en");

      await page.keyboard.press("ControlOrMeta+k");

      // Verify navigation section exists
      await expect(page.getByText(/navigation/i)).toBeVisible();

      // Verify some navigation items
      await expect(page.getByText(/go to dashboard/i)).toBeVisible();
      await expect(page.getByText(/go to explore/i)).toBeVisible();
    });

    test("should display commands section", async ({ page }) => {
      await page.goto("/en");

      await page.keyboard.press("ControlOrMeta+k");

      // Verify commands section
      await expect(page.getByText(/commands/i)).toBeVisible();

      // Verify theme toggle command
      await expect(page.getByText(/switch to (light|dark) mode/i)).toBeVisible();
    });

    test("should navigate with keyboard arrows", async ({ page }) => {
      await page.goto("/en");

      await page.keyboard.press("ControlOrMeta+k");

      // Navigate down
      await page.keyboard.press("ArrowDown");
      await page.keyboard.press("ArrowDown");

      // Navigate up
      await page.keyboard.press("ArrowUp");

      // Verify an item is selected (has data-selected attribute)
      await expect(page.locator("[data-selected=true]")).toBeVisible();
    });

    test("should navigate to page on Enter", async ({ page }) => {
      await page.goto("/en");

      await page.keyboard.press("ControlOrMeta+k");

      // Type to filter to Explore
      await page.getByPlaceholder(/search/i).fill("explore");

      // Wait for filtering
      await page.waitForTimeout(200);

      // Select and navigate
      await page.keyboard.press("Enter");

      // Verify navigation
      await expect(page).toHaveURL(/\/explore/);
    });

    test("should display keyboard hints on desktop", async ({ page }) => {
      await page.goto("/en");

      await page.keyboard.press("ControlOrMeta+k");

      // Verify keyboard hints are visible
      await expect(page.getByText(/navigate/i).last()).toBeVisible();
      await expect(page.getByText(/select/i).last()).toBeVisible();
      await expect(page.getByText(/close/i).last()).toBeVisible();
    });

    test("should display shortcut badge in search button", async ({ page }) => {
      await page.goto("/en");

      // Verify shortcut hint in button
      const searchButton = page.getByRole("button", { name: /search/i }).first();
      await expect(searchButton.locator("kbd")).toBeVisible();
    });
  });

  test.describe("Mobile", () => {
    test.beforeEach(async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
    });

    test("should open search dialog by tapping icon button", async ({ page }) => {
      await page.goto("/en");

      // Click mobile search icon button in header
      await page.getByRole("button", { name: /search/i }).click();

      // Verify dialog is open
      await expect(page.getByRole("dialog")).toBeVisible();
    });

    test("should display full-screen dialog on mobile", async ({ page }) => {
      await page.goto("/en");

      await page.getByRole("button", { name: /search/i }).click();

      // Verify dialog takes full screen
      const dialog = page.getByRole("dialog");
      await expect(dialog).toBeVisible();

      // Check dialog dimensions are close to viewport
      const dialogBox = await dialog.boundingBox();
      expect(dialogBox?.width).toBeGreaterThan(350);
    });

    test("should display close button on mobile", async ({ page }) => {
      await page.goto("/en");

      await page.getByRole("button", { name: /search/i }).click();

      // Verify close button is visible
      const closeButton = page.getByRole("button", { name: /close/i });
      await expect(closeButton).toBeVisible();
    });

    test("should close dialog with close button", async ({ page }) => {
      await page.goto("/en");

      await page.getByRole("button", { name: /search/i }).click();
      await expect(page.getByRole("dialog")).toBeVisible();

      // Click close button
      await page.getByRole("button", { name: /close/i }).click();
      await expect(page.getByRole("dialog")).not.toBeVisible();
    });

    test("should hide keyboard hints on mobile", async ({ page }) => {
      await page.goto("/en");

      await page.getByRole("button", { name: /search/i }).click();

      // Verify keyboard hints container is hidden
      // The hints container has "hidden md:flex" class
      const hintsContainer = page.locator("[class*='hidden'][class*='md:flex']").filter({
        hasText: /navigate|select|close/i,
      });
      await expect(hintsContainer).toBeHidden();
    });

    test("should hide shortcut badges on mobile", async ({ page }) => {
      await page.goto("/en");

      await page.getByRole("button", { name: /search/i }).click();

      // Verify shortcut badges are hidden (CommandShortcut has hidden md:inline)
      const shortcuts = page.locator("[data-slot='command-shortcut']");
      const count = await shortcuts.count();
      if (count > 0) {
        await expect(shortcuts.first()).toBeHidden();
      }
    });
  });

  test.describe("Search Functionality", () => {
    test("should show empty state when no results", async ({ page }) => {
      await page.goto("/en");

      await page.keyboard.press("ControlOrMeta+k");

      // Type a query that won't match anything
      await page.getByPlaceholder(/search/i).fill("xyznonexistent12345");

      // Wait for search debounce
      await page.waitForTimeout(200);

      // Verify empty state message
      await expect(page.getByText(/no results/i)).toBeVisible();
    });

    test("should clear search and show static actions", async ({ page }) => {
      await page.goto("/en");

      await page.keyboard.press("ControlOrMeta+k");

      // Type something
      await page.getByPlaceholder(/search/i).fill("test");
      await page.waitForTimeout(200);

      // Clear input
      await page.getByPlaceholder(/search/i).fill("");
      await page.waitForTimeout(200);

      // Verify static actions are shown again
      await expect(page.getByText(/navigation/i)).toBeVisible();
    });
  });
});
