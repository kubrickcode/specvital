import { expect, test } from "@playwright/test";

test.describe("Analysis Page", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to a cached analysis page (facebook/react has cached data)
    await page.goto("/en/analyze/facebook/react");
    // Wait for analysis to complete (or load from cache)
    await page.waitForSelector("text=Test Statistics", { timeout: 30000 });
  });

  test("should display test statistics", async ({ page }) => {
    // Verify Test Statistics heading
    await expect(
      page.getByRole("heading", { name: "Test Statistics" })
    ).toBeVisible();

    // Verify statistics numbers exist
    await expect(page.getByText("Total")).toBeVisible();
    await expect(page.getByText("Active")).toBeVisible();
    await expect(page.getByText("Skipped")).toBeVisible();
    await expect(page.getByText("Todo")).toBeVisible();

    // Verify By Framework section
    await expect(
      page.getByRole("heading", { name: "By Framework" })
    ).toBeVisible();

    // Verify progressbar elements exist
    const progressbars = page.getByRole("progressbar");
    await expect(progressbars.first()).toBeVisible();
  });

  test("should display repository info header", async ({ page }) => {
    // Verify repository heading
    await expect(
      page.getByRole("heading", { name: "facebook/react" })
    ).toBeVisible();

    // Verify View on GitHub link
    await expect(
      page.getByRole("link", { name: "View on GitHub" })
    ).toBeVisible();
    await expect(
      page.getByRole("link", { name: "View on GitHub" })
    ).toHaveAttribute("href", "https://github.com/facebook/react");
  });

  test("should display test suites section", async ({ page }) => {
    // Verify Test Suites heading
    await expect(
      page.getByRole("heading", { name: "Test Suites" })
    ).toBeVisible();

    // Verify search input exists
    await expect(
      page.getByRole("textbox", { name: /search/i })
    ).toBeVisible();

    // Verify filter buttons exist
    await expect(page.getByRole("button", { name: "Status" })).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Framework" })
    ).toBeVisible();

    // Verify view toggle exists
    await expect(page.getByRole("radio", { name: "List view" })).toBeVisible();
    await expect(page.getByRole("radio", { name: "Tree view" })).toBeVisible();
  });

  test("should toggle between list and tree view", async ({ page }) => {
    // List view is default
    await expect(page.getByRole("radio", { name: "List view" })).toBeChecked();
    await expect(page.getByRole("radio", { name: "Tree view" })).not.toBeChecked();

    // Switch to tree view
    await page.getByRole("radio", { name: "Tree view" }).click();

    // Verify URL updated
    await expect(page).toHaveURL(/\?view=tree/);

    // Verify tree view is selected
    await expect(page.getByRole("radio", { name: "Tree view" })).toBeChecked();

    // Verify folder structure is displayed
    await expect(
      page.getByRole("button", { name: /expand|collapse/i }).first()
    ).toBeVisible();

    // Switch back to list view
    await page.getByRole("radio", { name: "List view" }).click();

    // Verify URL updated (no view param)
    await expect(page).not.toHaveURL(/\?view=tree/);

    // Verify list view is selected
    await expect(page.getByRole("radio", { name: "List view" })).toBeChecked();
  });

  test("should expand and collapse folders in tree view", async ({ page }) => {
    // Switch to tree view
    await page.getByRole("radio", { name: "Tree view" }).click();

    // Find and expand a folder
    const expandButton = page.getByRole("button", { name: /expand/i }).first();
    await expect(expandButton).toBeVisible();
    await expandButton.click();

    // Verify folder is now in expanded state (button text changes to Collapse)
    await expect(
      page.getByRole("button", { name: /collapse/i }).first()
    ).toBeVisible();
  });

  test("should search and filter tests", async ({ page }) => {
    const searchInput = page.getByRole("textbox", { name: /search/i });

    // Search for "useState"
    await searchInput.fill("useState");

    // Verify URL updated with search query
    await expect(page).toHaveURL(/\?q=useState/);

    // Verify filtered results (status shows "X of Y tests")
    await expect(page.getByRole("status")).toContainText(/of.*tests/i);

    // Verify clear button appears
    await expect(
      page.getByRole("button", { name: "Clear search" })
    ).toBeVisible();

    // Clear search
    await page.getByRole("button", { name: "Clear search" }).click();

    // Verify URL cleared
    await expect(page).not.toHaveURL(/\?q=/);

    // Verify search input is empty
    await expect(searchInput).toHaveValue("");
  });

  test("should filter by status", async ({ page }) => {
    // Click Status filter button
    await page.getByRole("button", { name: "Status" }).click();

    // Verify filter dialog is open
    await expect(page.getByRole("checkbox", { name: "Skipped" })).toBeVisible();

    // Select Skipped status
    await page.getByRole("checkbox", { name: "Skipped" }).click();

    // Verify URL updated
    await expect(page).toHaveURL(/\?statuses=skipped/);

    // Verify button text updated with count
    await expect(
      page.getByRole("button", { name: /Status \(1\)/ })
    ).toBeVisible();

    // Verify filtered results (status shows "X of Y tests")
    await expect(page.getByRole("status")).toContainText(/of.*tests/i);

    // Close dialog
    await page.keyboard.press("Escape");
  });

  test("should share analysis link", async ({ page }) => {
    // Click Share button
    await page.getByRole("button", { name: "Share analysis link" }).click();

    // Verify button text changes to "Link copied!"
    await expect(page.getByText("Link copied!")).toBeVisible();

    // Verify toast notification appears
    await expect(page.getByRole("listitem")).toContainText("Link copied!");
  });

  test("should display export options", async ({ page }) => {
    // Click Export button
    await page.getByRole("button", { name: "Export analysis results" }).click();

    // Verify export menu is visible
    await expect(
      page.getByRole("menu", { name: "Export analysis results" })
    ).toBeVisible();

    // Verify export options
    await expect(
      page.getByRole("menuitem", { name: "Download Markdown" })
    ).toBeVisible();
    await expect(
      page.getByRole("menuitem", { name: "Copy Markdown" })
    ).toBeVisible();

    // Close menu by pressing Escape
    await page.keyboard.press("Escape");

    // Verify menu is closed
    await expect(
      page.getByRole("menu", { name: "Export analysis results" })
    ).not.toBeVisible();
  });

  test("should filter by framework", async ({ page }) => {
    // Click Framework filter button
    await page.getByRole("button", { name: "Framework" }).click();

    // Verify filter dialog is open with framework checkboxes
    await expect(page.getByRole("checkbox", { name: "Jest" })).toBeVisible();
    await expect(
      page.getByRole("checkbox", { name: "Playwright" })
    ).toBeVisible();

    // Select Jest framework
    await page.getByRole("checkbox", { name: "Jest" }).click();

    // Verify URL updated with frameworks parameter
    await expect(page).toHaveURL(/\?frameworks=jest/);

    // Verify button text updated with count
    await expect(
      page.getByRole("button", { name: /Framework \(1\)/ })
    ).toBeVisible();

    // Verify filtered results (status shows "X of Y tests")
    await expect(page.getByRole("status")).toContainText(/of.*tests/i);

    // Close dialog
    await page.keyboard.press("Escape");
  });

  test("should show login prompt for AI Spec generation when unauthenticated", async ({
    page,
  }) => {
    // Click Generate AI Spec button
    await page
      .getByRole("button", { name: "Generate test specification document with AI" })
      .click();

    // Verify login dialog is shown
    await expect(
      page.getByRole("dialog", { name: "Sign in to Generate Specs" })
    ).toBeVisible();

    // Verify dialog content
    await expect(
      page.getByRole("heading", { name: "Sign in to Generate Specs" })
    ).toBeVisible();
    await expect(
      page.getByText(/Transform your test cases into organized specification/i)
    ).toBeVisible();

    // Verify Continue with GitHub button
    await expect(
      page.getByRole("button", { name: "Continue with GitHub" })
    ).toBeVisible();

    // Close dialog
    await page.getByRole("button", { name: "Close" }).click();
    await expect(
      page.getByRole("dialog", { name: "Sign in to Generate Specs" })
    ).not.toBeVisible();
  });

  test("should display commit info in header", async ({ page }) => {
    // Verify branch name is displayed
    await expect(page.getByText("main")).toBeVisible();

    // Verify commit SHA button exists (short hash format)
    const commitButton = page.getByRole("button", { name: /^[a-f0-9]{7}$/i });
    await expect(commitButton).toBeVisible();

    // Verify View on GitHub link
    const githubLink = page.getByRole("link", { name: "View on GitHub" });
    await expect(githubLink).toBeVisible();
    await expect(githubLink).toHaveAttribute(
      "href",
      "https://github.com/facebook/react"
    );
  });
});
