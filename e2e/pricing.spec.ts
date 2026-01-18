import { expect, test } from "@playwright/test";

test.describe("Pricing Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/pricing");
  });

  test("should display all pricing plans", async ({ page }) => {
    // Verify page heading
    await expect(
      page.getByRole("heading", { name: "Simple, transparent pricing" })
    ).toBeVisible();

    // Verify Free plan
    await expect(page.getByText("Free")).toBeVisible();
    await expect(page.getByText("$0")).toBeVisible();

    // Verify Pro plan
    await expect(page.getByText("Pro").first()).toBeVisible();
    await expect(page.getByText("$15")).toBeVisible();
    await expect(page.getByText("Most Popular")).toBeVisible();

    // Verify Pro+ plan
    await expect(page.getByText("Pro+")).toBeVisible();
    await expect(page.getByText("$59")).toBeVisible();

    // Verify Enterprise plan
    await expect(page.getByText("Enterprise")).toBeVisible();
    await expect(page.getByText("Custom")).toBeVisible();
  });

  test("should expand FAQ accordion", async ({ page }) => {
    const faqButton = page.getByRole("button", {
      name: "What is an AI Spec Document?",
    });

    // Verify FAQ button exists and is not expanded
    await expect(faqButton).toBeVisible();
    await expect(faqButton).not.toHaveAttribute("aria-expanded", "true");

    // Click to expand
    await faqButton.click();

    // Verify expanded state
    await expect(faqButton).toHaveAttribute("aria-expanded", "true");

    // Verify answer is visible
    await expect(
      page.getByText(/automatically organizes your test cases/i)
    ).toBeVisible();
  });

  test("should have all FAQ items", async ({ page }) => {
    await expect(
      page.getByRole("button", { name: "What is an AI Spec Document?" })
    ).toBeVisible();
    await expect(
      page.getByRole("button", { name: "What is Analysis?" })
    ).toBeVisible();
    await expect(
      page.getByRole("button", { name: "When will payments go live?" })
    ).toBeVisible();
    await expect(
      page.getByRole("button", { name: "What is data retention?" })
    ).toBeVisible();
  });

  test("should have Get Started buttons for plans", async ({ page }) => {
    const getStartedButtons = page.getByRole("button", { name: "Get Started" });
    await expect(getStartedButtons).toHaveCount(3); // Free, Pro, Pro+
  });

  test("should have Contact Us link for Enterprise", async ({ page }) => {
    const contactLink = page.getByRole("link", { name: "Contact Us" });
    await expect(contactLink).toBeVisible();
    await expect(contactLink).toHaveAttribute(
      "href",
      "mailto:support@specvital.com"
    );
  });
});
