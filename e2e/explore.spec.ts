import { expect, test } from "@playwright/test";

test.describe("Explore Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/explore");
  });

  test("should display tab navigation", async ({ page }) => {
    // Verify Community tab is selected by default
    const communityTab = page.getByRole("tab", { name: "Community" });
    await expect(communityTab).toHaveAttribute("aria-selected", "true");

    // Verify My Repos tab exists
    await expect(page.getByRole("tab", { name: "My Repos" })).toBeVisible();

    // Verify Organizations tab exists
    await expect(
      page.getByRole("tab", { name: "Organizations" })
    ).toBeVisible();
  });

  test("should switch to My Repos tab and show login prompt", async ({
    page,
  }) => {
    // Click My Repos tab
    await page.getByRole("tab", { name: "My Repos" }).click();

    // Verify My Repos tab is selected
    await expect(page.getByRole("tab", { name: "My Repos" })).toHaveAttribute(
      "aria-selected",
      "true"
    );

    // Verify login prompt is shown (when not authenticated)
    await expect(
      page.getByRole("heading", { name: /sign in to view/i })
    ).toBeVisible();
    await expect(
      page.getByRole("button", { name: /sign in to continue/i })
    ).toBeVisible();
  });

  test("should switch back to Community tab", async ({ page }) => {
    // First switch to My Repos
    await page.getByRole("tab", { name: "My Repos" }).click();

    // Then switch back to Community
    await page.getByRole("tab", { name: "Community" }).click();

    // Verify Community tab is selected
    await expect(page.getByRole("tab", { name: "Community" })).toHaveAttribute(
      "aria-selected",
      "true"
    );

    // Verify Community content is shown (tabpanel)
    await expect(
      page.getByRole("tabpanel", { name: "Community" })
    ).toBeVisible();
  });

  test("should filter repositories by search", async ({ page }) => {
    // Wait for repository list to load
    await expect(page.getByText(/Showing all \d+ repositories/)).toBeVisible({
      timeout: 15000,
    });

    const searchbox = page.getByRole("searchbox", {
      name: "Search repositories...",
    });

    // Search for specific repository
    await searchbox.fill("react");

    // Verify filtered results show matching repository
    await expect(page.getByRole("heading", { name: "facebook/react" })).toBeVisible();

    // Clear search
    await searchbox.fill("");

    // Verify all repositories are shown again
    await expect(page.getByText(/Showing all \d+ repositories/)).toBeVisible();
  });

  test("should show no results message for non-matching search", async ({
    page,
  }) => {
    // Wait for repository list to load
    await expect(page.getByText(/Showing all \d+ repositories/)).toBeVisible({
      timeout: 15000,
    });

    // Search for non-existent repository
    await page
      .getByRole("searchbox", { name: "Search repositories..." })
      .fill("nonexistent");

    // Verify no results message
    await expect(
      page.getByRole("heading", { name: "No matching repositories" })
    ).toBeVisible();
    await expect(page.getByText(/No results found for "nonexistent"/)).toBeVisible();
  });

  test("should navigate to analysis page when clicking repository card", async ({
    page,
  }) => {
    // Wait for repository list to load
    await expect(page.getByText(/Showing all \d+ repositories/)).toBeVisible({
      timeout: 15000,
    });

    // Click on first repository card (facebook/react)
    await page
      .getByRole("link", { name: /facebook\/react.*tests/i })
      .click();

    // Verify navigation to analysis page
    await expect(page).toHaveURL(/\/analyze\/facebook\/react/);

    // Verify analysis page content
    await expect(
      page.getByRole("heading", { name: "facebook/react" })
    ).toBeVisible();
    await expect(
      page.getByRole("heading", { name: "Test Statistics" })
    ).toBeVisible({ timeout: 30000 });
  });

  test("should sort repositories by different criteria", async ({ page }) => {
    // Wait for repository list to load
    await expect(page.getByText(/Showing all \d+ repositories/)).toBeVisible({
      timeout: 15000,
    });

    // Verify default sort is "Recent"
    const sortButton = page.getByRole("button", { name: /Sort: Recent/i });
    await expect(sortButton).toBeVisible();

    // Open sort dropdown
    await sortButton.click();

    // Verify dropdown options are visible
    await expect(page.getByRole("menuitemradio", { name: "Recent" })).toBeVisible();
    await expect(page.getByRole("menuitemradio", { name: "Name" })).toBeVisible();
    await expect(page.getByRole("menuitemradio", { name: "Tests" })).toBeVisible();

    // Select "Name" sort option
    await page.getByRole("menuitemradio", { name: "Name" }).click();

    // Verify sort button text updated
    await expect(page.getByRole("button", { name: /Sort: Name/i })).toBeVisible();

    // Sort by "Tests"
    await page.getByRole("button", { name: /Sort: Name/i }).click();
    await page.getByRole("menuitemradio", { name: "Tests" }).click();

    // Verify sort button text updated to "Tests"
    await expect(page.getByRole("button", { name: /Sort: Tests/i })).toBeVisible();
  });

  test("should maintain sort selection while searching", async ({ page }) => {
    // Wait for repository list to load
    await expect(page.getByText(/Showing all \d+ repositories/)).toBeVisible({
      timeout: 15000,
    });

    // Change sort to "Name"
    await page.getByRole("button", { name: /Sort: Recent/i }).click();
    await page.getByRole("menuitemradio", { name: "Name" }).click();
    await expect(page.getByRole("button", { name: /Sort: Name/i })).toBeVisible();

    // Search for repository
    const searchbox = page.getByRole("searchbox", {
      name: "Search repositories...",
    });
    await searchbox.fill("react");

    // Verify sort is still "Name"
    await expect(page.getByRole("button", { name: /Sort: Name/i })).toBeVisible();

    // Clear search
    await searchbox.fill("");

    // Verify sort is still "Name"
    await expect(page.getByRole("button", { name: /Sort: Name/i })).toBeVisible();
  });
});
