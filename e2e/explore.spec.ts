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

  // Note: Data-dependent tests moved to explore.mocked.spec.ts
});
