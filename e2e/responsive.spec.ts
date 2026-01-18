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
});
