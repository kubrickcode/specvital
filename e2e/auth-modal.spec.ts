import { expect, test } from "@playwright/test";

test.describe("Login Modal", () => {
  test("should open and close login modal", async ({ page }) => {
    await page.goto("/");

    // Click Sign in button
    await page.getByRole("button", { name: "Sign in" }).click();

    // Verify modal is displayed
    const modal = page.getByRole("dialog", { name: "Sign in to SpecVital" });
    await expect(modal).toBeVisible();

    // Verify modal heading
    await expect(
      modal.getByRole("heading", { name: "Sign in to SpecVital" })
    ).toBeVisible();

    // Verify Continue with GitHub button
    await expect(
      modal.getByRole("button", { name: "Continue with GitHub" })
    ).toBeVisible();

    // Verify terms text
    await expect(
      modal.getByText("By signing in, you agree to our Terms of Service")
    ).toBeVisible();

    // Close modal
    await modal.getByRole("button", { name: "Close" }).click();

    // Verify modal is closed
    await expect(modal).not.toBeVisible();
  });

  test("should close modal on escape key", async ({ page }) => {
    await page.goto("/");

    // Open modal
    await page.getByRole("button", { name: "Sign in" }).click();
    const modal = page.getByRole("dialog", { name: "Sign in to SpecVital" });
    await expect(modal).toBeVisible();

    // Press Escape
    await page.keyboard.press("Escape");

    // Verify modal is closed
    await expect(modal).not.toBeVisible();
  });
});
