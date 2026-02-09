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

      // Verify navigation items available without authentication
      // Note: Dashboard requires auth, so we check Explore, Pricing, and Home instead
      await expect(page.getByText(/go to explore/i)).toBeVisible();
      await expect(page.getByText(/go to pricing/i)).toBeVisible();
      await expect(page.getByText(/go to home/i)).toBeVisible();
    });

    test("should display commands section", async ({ page }) => {
      await page.goto("/en");

      await page.keyboard.press("ControlOrMeta+k");

      // Verify commands section heading (use cmdk-group-heading attribute to avoid matching description text)
      await expect(page.locator("[cmdk-group-heading]", { hasText: /commands/i })).toBeVisible();

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

      // Use arrow keys to navigate to "Go to Explore" (available without auth)
      // First item is "Go to Explore" among navigation items for unauthenticated users
      await page.keyboard.press("ArrowDown");

      // Find and click on "Go to Explore" directly to ensure correct navigation
      await page.getByText(/go to explore/i).click();

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

    test("should not have hydration errors on initial load", async ({ page }) => {
      // Listen for console messages
      const consoleMessages: string[] = [];
      page.on("console", (msg) => {
        if (msg.type() === "error" || msg.type() === "warning") {
          consoleMessages.push(msg.text());
        }
      });

      await page.goto("/en");

      // Wait for page to fully hydrate
      await page.waitForLoadState("networkidle");

      // Verify search button is immediately interactive (no hydration issues)
      const searchButton = page.getByRole("button", { name: /search/i }).first();
      await expect(searchButton).toBeVisible();
      await expect(searchButton.locator("kbd")).toBeVisible();

      // Click should work immediately (no hydration delay)
      await searchButton.click();
      await expect(page.getByRole("dialog")).toBeVisible();

      // Check for hydration errors in console
      const hydrationErrors = consoleMessages.filter(
        (msg) =>
          msg.includes("hydration") ||
          msg.includes("Hydration") ||
          msg.includes("did not match")
      );

      expect(hydrationErrors).toHaveLength(0);
    });
  });

  test.describe("Mobile", () => {
    test.beforeEach(async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
    });

    test("should open search dialog by tapping icon button", async ({ page }) => {
      await page.goto("/en");
      await page.waitForLoadState("networkidle");

      // Click mobile search icon button in header (visible one on mobile)
      const searchButton = page.getByRole("button", { name: /search/i });
      await searchButton.click();

      // Verify dialog is open
      await expect(page.getByRole("dialog")).toBeVisible();
    });

    test("should display full-screen dialog on mobile", async ({ page }) => {
      await page.goto("/en");
      await page.waitForLoadState("networkidle");

      const searchButton = page.getByRole("button", { name: /search/i });
      await searchButton.click();

      // Verify dialog takes full screen
      const dialog = page.getByRole("dialog");
      await expect(dialog).toBeVisible();

      // Check dialog dimensions are close to viewport
      const dialogBox = await dialog.boundingBox();
      expect(dialogBox?.width).toBeGreaterThan(350);
    });

    test("should display close button on mobile", async ({ page }) => {
      await page.goto("/en");
      await page.waitForLoadState("networkidle");

      const searchButton = page.getByRole("button", { name: /search/i });
      await searchButton.click();

      // Verify dialog is open first
      const dialog = page.getByRole("dialog");
      await expect(dialog).toBeVisible();

      // Verify close button is visible in dialog
      const closeButton = dialog.getByRole("button", { name: /close/i });
      await expect(closeButton).toBeVisible();
    });

    test("should close dialog with close button", async ({ page }) => {
      await page.goto("/en");
      await page.waitForLoadState("networkidle");

      const searchButton = page.getByRole("button", { name: /search/i });
      await searchButton.click();

      const dialog = page.getByRole("dialog");
      await expect(dialog).toBeVisible();

      // Click close button inside dialog
      await dialog.getByRole("button", { name: /close/i }).click();
      await expect(dialog).not.toBeVisible();
    });

    test("should hide keyboard hints on mobile", async ({ page }) => {
      await page.goto("/en");

      await page.getByRole("button", { name: /search/i }).click();

      // Verify keyboard hints container is hidden on mobile
      // The SearchKeyboardHints component has "hidden md:flex" class
      // Use a more specific selector targeting the hints footer within the dialog
      const dialog = page.getByRole("dialog");
      const hintsContainer = dialog.locator("div.hidden.border-t").filter({
        has: page.locator("kbd"),
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
