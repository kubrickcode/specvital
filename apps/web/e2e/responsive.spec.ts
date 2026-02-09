import { expect, test } from "@playwright/test";

test.describe("Responsive Layout", () => {
  test("should display mobile navigation on small viewport", async ({
    page,
  }) => {
    await page.goto("/en");

    // Resize to mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });

    // Verify mobile navigation is visible
    await expect(
      page.getByRole("navigation", { name: "Mobile navigation" })
    ).toBeVisible();

    // Verify main navigation is hidden (or not visible)
    await expect(
      page.getByRole("navigation", { name: "Main navigation" })
    ).not.toBeVisible();

    // Verify mobile nav elements
    await expect(page.getByRole("link", { name: "Explore" })).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Toggle theme" })
    ).toBeVisible();
    await expect(page.getByRole("button", { name: "Sign in" })).toBeVisible();
  });

  test("should display desktop navigation on large viewport", async ({
    page,
  }) => {
    await page.goto("/en");

    // Ensure desktop viewport
    await page.setViewportSize({ width: 1280, height: 720 });

    // Verify main navigation is visible
    await expect(
      page.getByRole("navigation", { name: "Main navigation" })
    ).toBeVisible();

    // Verify desktop header links
    await expect(page.getByRole("link", { name: "Explore" })).toBeVisible();
    await expect(page.getByRole("link", { name: "Pricing" })).toBeVisible();
  });

  test("should show More menu with Docs and Pricing in mobile bottom bar", async ({
    page,
  }) => {
    await page.goto("/en");

    // Resize to mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });

    // Verify More button exists in mobile navigation
    const moreButton = page.getByRole("button", { name: "More" });
    await expect(moreButton).toBeVisible();

    // Click More button
    await moreButton.click();

    // Verify menu appears with Docs and Pricing
    const menu = page.getByRole("menu", { name: "More" });
    await expect(menu).toBeVisible();

    // Verify Docs menuitem
    const docsItem = menu.getByRole("menuitem", { name: "Docs" });
    await expect(docsItem).toBeVisible();

    // Verify Pricing menuitem
    const pricingItem = menu.getByRole("menuitem", { name: "Pricing" });
    await expect(pricingItem).toBeVisible();

    // Test navigation by clicking Docs
    await docsItem.click();
    await expect(page).toHaveURL("/en/docs/test-detection");
  });
});
