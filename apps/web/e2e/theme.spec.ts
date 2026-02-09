import { expect, test } from "@playwright/test";

test.describe("Theme Toggle", () => {
  test("should toggle between light and dark theme", async ({ page }) => {
    await page.goto("/en");

    // Get initial theme
    const initialTheme = await page.evaluate(() =>
      document.documentElement.classList.contains("dark") ? "dark" : "light"
    );

    // Click theme toggle
    await page.getByRole("button", { name: "Toggle theme" }).click();

    // Wait for theme transition (300ms animation + setTheme delay)
    await page.waitForTimeout(400);

    // Verify theme changed
    const newTheme = await page.evaluate(() =>
      document.documentElement.classList.contains("dark") ? "dark" : "light"
    );
    expect(newTheme).not.toBe(initialTheme);

    // Click again to toggle back
    await page.getByRole("button", { name: "Toggle theme" }).click();

    // Wait for theme transition
    await page.waitForTimeout(400);

    // Verify theme reverted
    const finalTheme = await page.evaluate(() =>
      document.documentElement.classList.contains("dark") ? "dark" : "light"
    );
    expect(finalTheme).toBe(initialTheme);
  });
});
