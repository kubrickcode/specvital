import { expect, test } from "@playwright/test";

test.describe("Documentation - Test Writing Guide", () => {
  test("should display Test Writing Guide page with all sections", async ({
    page,
  }) => {
    await page.goto("/en/docs/writing-guide");

    // Verify page title
    await expect(
      page.getByRole("heading", { name: "Test Writing Guide", level: 1 })
    ).toBeVisible();

    // Verify page description
    await expect(
      page.getByText(
        "Structure your tests for accurate analysis and high-quality spec documents."
      )
    ).toBeVisible();
  });

  test("should display How SpecVital Works section with 4-step pipeline", async ({
    page,
  }) => {
    await page.goto("/en/docs/writing-guide");

    // Verify section heading
    const heading = page.getByRole("heading", {
      name: "How SpecVital Works",
      level: 2,
    });
    await expect(heading).toBeVisible();

    // Verify all 4 pipeline steps are displayed together
    await expect(page.getByText("Test Detection").first()).toBeVisible();
    await expect(page.getByText("AST Parsing").first()).toBeVisible();
    await expect(page.getByText("Analysis").first()).toBeVisible();
    await expect(page.getByText("Spec Document").first()).toBeVisible();
  });

  test("should display Naming Tests section with comparison table", async ({
    page,
  }) => {
    await page.goto("/en/docs/writing-guide");

    // Verify section heading
    await expect(
      page.getByRole("heading", { name: "Naming Tests", level: 2 })
    ).toBeVisible();

    // Verify table exists and contains expected data
    const table = page.getByRole("table");
    await expect(table).toBeVisible();

    // Verify key sample data that should appear in the table
    // Use first() to avoid strict mode violation when text appears in both table and code blocks
    await expect(page.getByText("tc_01").first()).toBeVisible();
    await expect(page.getByText("AuthService > Login").first()).toBeVisible();
    await expect(page.getByText("OrderService > Sorting").first()).toBeVisible();
    await expect(
      page.getByText("should reject with expired token").first()
    ).toBeVisible();
  });

  test("should display Structuring Tests section with Good/Bad examples", async ({
    page,
  }) => {
    await page.goto("/en/docs/writing-guide");

    // Verify section heading and scroll to it
    const heading = page.getByRole("heading", {
      name: "Structuring Tests",
      level: 2,
    });
    await expect(heading).toBeVisible();
    await heading.scrollIntoViewIfNeeded();

    // Verify Flat Structure example (Bad)
    await expect(page.getByText("Flat Structure")).toBeVisible();
    await expect(
      page.getByText("Without suite hierarchy, analysis lacks grouping context.")
    ).toBeVisible();

    // Verify Nested Structure example (Good)
    await expect(page.getByText("Nested Structure")).toBeVisible();
    await expect(
      page.getByText(
        "Nested test groups provide clear structural context for analysis."
      )
    ).toBeVisible();
  });

  test("should display Organizing Test Files section with examples", async ({
    page,
  }) => {
    await page.goto("/en/docs/writing-guide");

    // Verify section heading
    await expect(
      page.getByRole("heading", { name: "Organizing Test Files", level: 2 })
    ).toBeVisible();

    // Verify both example types are shown
    // Use first() to avoid strict mode violation when text appears multiple times
    await expect(page.getByText("Avoid").first()).toBeVisible();
    await expect(page.getByText("Recommended").first()).toBeVisible();

    // Verify best practice list items exist
    await expect(page.getByText("Follow Framework Conventions")).toBeVisible();
    await expect(page.getByText("Co-locate Tests with Source")).toBeVisible();
    await expect(page.getByText("Separate by Domain")).toBeVisible();
  });

  test("should display Project Architecture section with signal strength", async ({
    page,
  }) => {
    await page.goto("/en/docs/writing-guide");

    // Verify section heading
    await expect(
      page.getByRole("heading", { name: "Project Architecture", level: 2 })
    ).toBeVisible();

    // Verify signal strength indicators
    await expect(page.getByText("Analysis Signal Strength")).toBeVisible();
    await expect(
      page.getByText("Imports and function calls (strongest)")
    ).toBeVisible();
    await expect(page.getByText("File path (strong)")).toBeVisible();
    await expect(page.getByText("Test names (moderate)")).toBeVisible();
    await expect(page.getByText("Test content patterns (weak)")).toBeVisible();

    // Verify Weak Signals example
    await expect(page.getByText("Weak Signals")).toBeVisible();
    await expect(
      page.getByText("Generic utility imports provide no domain context.")
    ).toBeVisible();

    // Verify Strong Signals example
    await expect(page.getByText("Strong Signals")).toBeVisible();
    await expect(
      page.getByText("Domain-specific imports clearly indicate business area.")
    ).toBeVisible();
  });

  test("should display AI Tool Configuration section", async ({ page }) => {
    await page.goto("/en/docs/writing-guide");

    // Verify section heading
    await expect(
      page.getByRole("heading", { name: "AI Tool Configuration", level: 2 })
    ).toBeVisible();

    // Verify configuration cards
    await expect(page.getByText("Require Suite Hierarchy")).toBeVisible();
    await expect(page.getByText("Enforce Behavior Names")).toBeVisible();
    await expect(page.getByText("Mirror Domain in Paths")).toBeVisible();
    await expect(page.getByText("Import Real Modules")).toBeVisible();
  });

  test("should show Recommended Prompt Template dialog when button clicked", async ({
    page,
  }) => {
    await page.goto("/en/docs/writing-guide");

    // Click the accordion button
    const promptButton = page.getByRole("button", {
      name: "View Recommended Prompt Template",
    });
    await expect(promptButton).toBeVisible();
    await promptButton.click();

    // Verify dialog appears
    const dialog = page.getByRole("dialog", {
      name: "Recommended Prompt Template",
    });
    await expect(dialog).toBeVisible();

    // Verify dialog content
    await expect(
      dialog.getByRole("heading", {
        name: "Recommended Prompt Template",
        level: 2,
      })
    ).toBeVisible();
    await expect(
      dialog.getByText(
        "Copy this into your AI tool's rule or configuration file"
      )
    ).toBeVisible();

    // Verify prompt template code block
    await expect(
      dialog.getByText("# Test Writing Rules for SpecVital Compatibility")
    ).toBeVisible();

    // Verify Copy to Clipboard button
    await expect(
      dialog.getByRole("button", { name: "Copy to Clipboard" })
    ).toBeVisible();

    // Close dialog
    await dialog.getByRole("button", { name: "Close" }).click();
    await expect(dialog).not.toBeVisible();
  });

  test("should display Coming Soon section with MCP tools", async ({
    page,
  }) => {
    await page.goto("/en/docs/writing-guide");

    // Verify section heading
    await expect(
      page.getByRole("heading", {
        name: "Coming Soon: AI-Guided Tools",
        level: 2,
      })
    ).toBeVisible();

    // Verify MCP integration description
    await expect(
      page.getByText(
        "We're building MCP (Model Context Protocol) integration"
      )
    ).toBeVisible();

    // Verify all 4 coming soon features
    await expect(page.getByText("Test Name Improvement")).toBeVisible();
    await expect(page.getByText("Structure Refactoring")).toBeVisible();
    await expect(page.getByText("Analysis Preview")).toBeVisible();
    await expect(page.getByText("Architecture Diagnostics")).toBeVisible();

    // Verify "Coming Soon" badges
    const comingSoonBadges = page.getByText("Coming Soon");
    await expect(comingSoonBadges.first()).toBeVisible();
  });
});
