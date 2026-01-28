import { expect, test } from "@playwright/test";

test.describe("Documentation Pages", () => {
  test.describe("Header Navigation", () => {
    test("should navigate to docs page from header Docs link", async ({
      page,
    }) => {
      await page.goto("/en");

      // Verify Docs link exists in header navigation
      const docsLink = page
        .getByRole("navigation", { name: /main/i })
        .getByRole("link", { name: "Docs" });
      await expect(docsLink).toBeVisible();

      // Click Docs link - redirects to first doc page
      await docsLink.click();
      await expect(page).toHaveURL("/en/docs/test-detection");

      // Verify Docs link shows active state
      await expect(docsLink).toHaveAttribute("aria-current", "page");
    });

    test("should show Docs link alongside Explore and Pricing", async ({
      page,
    }) => {
      await page.goto("/en");

      const nav = page.getByRole("navigation", { name: /main/i });

      // Verify all three navigation links exist
      await expect(nav.getByRole("link", { name: "Explore" })).toBeVisible();
      await expect(nav.getByRole("link", { name: "Pricing" })).toBeVisible();
      await expect(nav.getByRole("link", { name: "Docs" })).toBeVisible();
    });

    test("should display navigation links in correct order: Explore, Docs, Pricing", async ({
      page,
    }) => {
      await page.goto("/en");

      const nav = page.getByRole("navigation", { name: /main/i });
      const links = nav.getByRole("link");

      // Get all link texts
      const linkTexts = await links.allTextContents();

      // Filter to only navigation links (Explore, Docs, Pricing)
      const navLinks = linkTexts.filter((text) =>
        ["Explore", "Docs", "Pricing"].includes(text)
      );

      // Verify order: Docs should come before Pricing
      const docsIndex = navLinks.indexOf("Docs");
      const pricingIndex = navLinks.indexOf("Pricing");

      expect(docsIndex).toBeGreaterThan(-1);
      expect(pricingIndex).toBeGreaterThan(-1);
      expect(docsIndex).toBeLessThan(pricingIndex);

      // Verify complete order
      expect(navLinks).toEqual(["Explore", "Docs", "Pricing"]);
    });
  });

  test.describe("Sidebar Navigation", () => {
    // Desktop: sidebar is always visible (lg:block)
    // Mobile: sidebar is hidden, use sheet/dialog

    test("should display sidebar with all navigation links on desktop", async ({
      page,
    }) => {
      await page.goto("/en/docs/test-detection");

      // Desktop sidebar should be visible
      const docNav = page.getByRole("navigation", {
        name: /documentation navigation/i,
      });
      await expect(docNav).toBeVisible();

      // Verify all 4 document links are displayed
      await expect(docNav.getByRole("link", { name: "Test Detection" })).toBeVisible();
      await expect(docNav.getByRole("link", { name: "Usage & Billing" })).toBeVisible();
      await expect(docNav.getByRole("link", { name: "AI Spec Generation" })).toBeVisible();
      await expect(docNav.getByRole("link", { name: "Test Writing Guide" })).toBeVisible();
    });

    test("should navigate to another doc page from sidebar", async ({
      page,
    }) => {
      await page.goto("/en/docs/test-detection");

      // Click Usage & Billing link in desktop sidebar
      const docNav = page.getByRole("navigation", {
        name: /documentation navigation/i,
      });
      await docNav.getByRole("link", { name: "Usage & Billing" }).click();

      // Verify navigation
      await expect(page).toHaveURL("/en/docs/usage-billing");
      await expect(
        page.getByRole("heading", { name: "Usage & Billing", level: 1 })
      ).toBeVisible();
    });

    test("should highlight current page in sidebar", async ({ page }) => {
      await page.goto("/en/docs/test-detection");

      // Verify Test Detection link is highlighted (uses gradient for active state)
      const docNav = page.getByRole("navigation", {
        name: /documentation navigation/i,
      });
      const testDetectionLink = docNav.getByRole("link", { name: "Test Detection" });
      await expect(testDetectionLink).toHaveClass(/bg-gradient-to-r/);
    });
  });

  test.describe("Individual Page Content", () => {
    test("should display Test Detection page with supported frameworks table", async ({
      page,
    }) => {
      await page.goto("/en/docs/test-detection");

      // Verify page heading
      await expect(
        page.getByRole("heading", { name: "Test Detection", level: 1 })
      ).toBeVisible();

      // Verify supported frameworks section
      await expect(
        page.getByRole("heading", { name: "Supported Frameworks", level: 2 })
      ).toBeVisible();

      // Verify frameworks table exists with content
      const table = page.getByRole("table").first();
      await expect(table).toBeVisible();
      await expect(table.getByText("JavaScript / TypeScript")).toBeVisible();
      await expect(table.getByText(/Jest.*Vitest/)).toBeVisible();
    });

    test("should display Usage & Billing page with plan limits table", async ({
      page,
    }) => {
      await page.goto("/en/docs/usage-billing");

      // Verify page heading
      await expect(
        page.getByRole("heading", { name: "Usage & Billing", level: 1 })
      ).toBeVisible();

      // Verify plan limits section
      await expect(
        page.getByRole("heading", { name: "Plan Limits", level: 2 })
      ).toBeVisible();

      // Verify plan limits table
      const table = page.locator("table").filter({ hasText: "Plan" });
      await expect(table.getByText("Free")).toBeVisible();
      // Use exact match to avoid matching "Pro+"
      await expect(table.getByText("Pro", { exact: true })).toBeVisible();
    });


    test("should display AI Spec Generation page with test classification table", async ({
      page,
    }) => {
      await page.goto("/en/docs/specview-generation");

      // Verify page heading
      await expect(
        page.getByRole("heading", { name: "AI Spec Generation", level: 1 })
      ).toBeVisible();

      // Verify test classification section
      await expect(
        page.getByRole("heading", { name: "Test Classification", level: 2 })
      ).toBeVisible();

      // Verify classification types are displayed in table
      const table = page.locator("table").filter({ hasText: "Functional" });
      await expect(table).toBeVisible();
      await expect(table.getByText("Edge Case")).toBeVisible();
      await expect(table.getByText("Integration")).toBeVisible();
    });
  });
});
