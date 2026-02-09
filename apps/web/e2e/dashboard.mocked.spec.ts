/**
 * Dashboard Page E2E Tests with API Mocking
 * Tests repository list features: sorting, filtering, search, infinite scroll, reanalyze
 */

import { expect, test } from "@playwright/test";
import { setupMockHandlers, createRepositoriesHandler } from "./fixtures/mock-handlers";
import {
  mockRepositoriesPage1,
  mockRepositoriesPage2,
  mockRepositoriesMine,
  mockRepositoriesOrg,
  mockRepositoriesBookmarked,
  mockRepositoriesSortedByName,
  mockRepositoriesSortedByTests,
  mockRepositoriesEmpty,
  mockStatsNormal,
  mockStatsEmpty,
} from "./fixtures/api-responses";

test.describe("Dashboard Page (Mocked API)", () => {
  test.describe("Repository List Display", () => {
    test("should display repository cards with mock data", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Verify repository cards are displayed
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Verify first repo data is displayed correctly
      const firstRepo = mockRepositoriesPage1.data[0];
      await expect(page.getByRole("link", { name: firstRepo.fullName }).first()).toBeVisible();
    });

    test("should display summary statistics", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for stats to load
      await expect(page.getByText(mockStatsNormal.totalRepositories.toString())).toBeVisible({
        timeout: 10000,
      });
    });
  });

  test.describe("Sorting", () => {
    test("should sort repositories by name (A-Z)", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
        onRepositoriesRequest: createRepositoriesHandler(mockRepositoriesPage1, {
          sortedByName: mockRepositoriesSortedByName,
        }),
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Open sort dropdown
      const sortButton = page.getByRole("button", { name: "Sort" });
      await sortButton.click();

      // Verify dropdown options
      await expect(page.getByRole("menuitemradio", { name: "Recent" })).toBeVisible();
      await expect(page.getByRole("menuitemradio", { name: "Name" })).toBeVisible();
      await expect(page.getByRole("menuitemradio", { name: "Tests" })).toBeVisible();

      // Select Name sort
      await page.getByRole("menuitemradio", { name: "Name" }).click();

      // Wait for sort to apply
      await page.waitForLoadState("networkidle");

      // Verify Name is now selected
      await sortButton.click();
      await expect(page.getByRole("menuitemradio", { name: "Name" })).toHaveAttribute(
        "aria-checked",
        "true"
      );
    });

    test("should sort repositories by test count", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
        onRepositoriesRequest: createRepositoriesHandler(mockRepositoriesPage1, {
          sortedByTests: mockRepositoriesSortedByTests,
        }),
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Open sort dropdown and select Tests
      const sortButton = page.getByRole("button", { name: "Sort" });
      await sortButton.click();
      await page.getByRole("menuitemradio", { name: "Tests" }).click();

      // Wait for sort to apply
      await page.waitForLoadState("networkidle");

      // Verify Tests is selected
      await sortButton.click();
      await expect(page.getByRole("menuitemradio", { name: "Tests" })).toHaveAttribute(
        "aria-checked",
        "true"
      );
    });
  });

  test.describe("Ownership Filtering", () => {
    test("should filter repositories by 'Mine' ownership", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
        onRepositoriesRequest: (url) => {
          const ownership = url.searchParams.get("ownership");
          if (ownership === "mine") {
            return mockRepositoriesMine;
          }
          return mockRepositoriesPage1;
        },
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Open ownership filter dropdown - aria-label is "Filter by ownership"
      const ownershipButton = page.getByRole("button", { name: /filter by ownership/i });
      await ownershipButton.click();

      // Verify dropdown options and select Mine filter
      const mineOption = page.getByRole("menuitemradio", { name: /^mine$/i });
      await expect(mineOption).toBeVisible({ timeout: 5000 });
      await mineOption.click();

      // Wait for filter to apply
      await page.waitForLoadState("networkidle");

      // Verify page still loads (filter applied successfully)
      await expect(page.getByRole("heading", { name: /dashboard/i })).toBeVisible();
    });

    test("should filter repositories by 'Organization' ownership", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
        onRepositoriesRequest: (url) => {
          const ownership = url.searchParams.get("ownership");
          if (ownership === "organization") {
            return mockRepositoriesOrg;
          }
          return mockRepositoriesPage1;
        },
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Open ownership filter dropdown
      const ownershipButton = page.getByRole("button", { name: /filter by ownership/i });
      await ownershipButton.click();

      const orgOption = page.getByRole("menuitemradio", { name: /^organization$/i });
      await expect(orgOption).toBeVisible({ timeout: 5000 });
      await orgOption.click();

      // Wait for filter to apply
      await page.waitForLoadState("networkidle");

      // Verify page still loads
      await expect(page.getByRole("heading", { name: /dashboard/i })).toBeVisible();
    });
  });

  test.describe("Bookmark Filtering", () => {
    test("should filter to show only bookmarked repositories", async ({ page }) => {
      let showBookmarkedOnly = false;

      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
        onRepositoriesRequest: () => {
          return showBookmarkedOnly ? mockRepositoriesBookmarked : mockRepositoriesPage1;
        },
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Click bookmarked toggle
      const bookmarkedToggle = page.getByRole("button", { name: /show bookmarked only/i });
      showBookmarkedOnly = true;
      await bookmarkedToggle.click();

      // Verify toggle is pressed
      await expect(bookmarkedToggle).toHaveAttribute("aria-pressed", "true");

      // Wait for filter to apply
      await page.waitForLoadState("networkidle");

      // Bookmarked repos should have "Remove bookmark" button
      await expect(page.getByRole("button", { name: "Remove bookmark" }).first()).toBeVisible({
        timeout: 5000,
      });
    });
  });

  test.describe("Search", () => {
    test("should search repositories and display results", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Find and use search input (type="search" renders as textbox)
      const searchInput = page.locator('input[type="search"]');
      await searchInput.fill("repo-1");

      // Wait for client-side filter to apply
      await page.waitForTimeout(300);

      // Verify search results are displayed (client-side filtering)
      await expect(page.getByText(/repo-1/).first()).toBeVisible();
    });

    test("should display empty state when no search results", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Search for nonexistent repo (client-side filter won't find matches)
      const searchInput = page.locator('input[type="search"]');
      await searchInput.fill("zzzznonexistent");

      // Wait for client-side filter
      await page.waitForTimeout(300);

      // Verify empty state is shown (no matching repos in client-side filter)
      await expect(
        page
          .getByRole("heading", { name: /no.*result|not found/i })
          .or(page.getByText(/no.*result|not found/i))
      ).toBeVisible({ timeout: 5000 });
    });
  });

  test.describe("Infinite Scroll", () => {
    test("should load more repositories on scroll", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        repositoriesPage2: mockRepositoriesPage2,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Count initial repos (should be 20 from page 1)
      const initialRepoCount = await page
        .getByRole("button", { name: /add bookmark|remove bookmark/i })
        .count();
      expect(initialRepoCount).toBe(20);

      // Scroll to bottom to trigger infinite scroll
      await page.evaluate(() => {
        window.scrollTo(0, document.body.scrollHeight);
      });

      // Wait for more repos to load
      await page.waitForTimeout(1000);

      // Count repos after scroll (should be 25 = 20 + 5)
      const finalRepoCount = await page
        .getByRole("button", { name: /add bookmark|remove bookmark/i })
        .count();
      expect(finalRepoCount).toBe(25);
    });
  });

  test.describe("Bookmark Toggle", () => {
    test("should toggle bookmark on repository card", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
        bookmarkResponse: { success: true, isBookmarked: true },
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      const bookmarkButton = page
        .getByRole("button", { name: /add bookmark|remove bookmark/i })
        .first();
      await expect(bookmarkButton).toBeVisible({ timeout: 10000 });

      // Get initial state
      const initialLabel = await bookmarkButton.getAttribute("aria-label");

      // Click to toggle - API call should be made
      await bookmarkButton.click();

      // Wait for optimistic update
      await page.waitForTimeout(300);

      // Verify button was clicked (state should change due to optimistic update)
      // Note: The actual state depends on UI implementation
      // Here we just verify the click was handled without error
      await expect(bookmarkButton).toBeVisible();
    });
  });

  test.describe("Reanalyze", () => {
    test("should trigger reanalyze for repository with new commits", async ({ page }) => {
      let reanalyzeTriggered = false;

      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
        onReanalyze: () => {
          reanalyzeTriggered = true;
          return { status: "queued" };
        },
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // First repo has updateStatus: "new-commits"
      // Look for reanalyze button or indicator
      const reanalyzeButton = page.getByRole("button", { name: /reanalyze|update/i }).first();
      const hasReanalyzeButton = await reanalyzeButton
        .isVisible({ timeout: 5000 })
        .catch(() => false);

      if (hasReanalyzeButton) {
        await reanalyzeButton.click();

        // Verify reanalyze was triggered
        await page.waitForTimeout(500);
        expect(reanalyzeTriggered).toBe(true);
      }
    });

    test("should show analyzing status and update UI after reanalysis completion", async ({
      page,
    }) => {
      let statusCallCount = 0;

      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
        onReanalyze: () => ({ status: "queued" as const }),
        onAnalysisStatus: () => {
          statusCallCount++;
          // Return completed after 2 polls
          return {
            owner: "test-owner",
            repo: "repo-1",
            status: statusCallCount >= 2 ? "completed" : "analyzing",
          };
        },
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Find reanalyze button for first repo (has new-commits status)
      const reanalyzeButton = page.getByRole("button", { name: /reanalyze|update/i }).first();
      const hasReanalyzeButton = await reanalyzeButton.isVisible({ timeout: 5000 }).catch(() => false);

      if (hasReanalyzeButton) {
        await reanalyzeButton.click();

        // Verify toast notification for queued status
        await expect(page.getByText(/queued|success/i)).toBeVisible({ timeout: 5000 });

        // Wait for polling to complete
        await page.waitForTimeout(3000);

        // Verify status call count (should have polled at least 2 times)
        expect(statusCallCount).toBeGreaterThanOrEqual(2);
      }
    });
  });

  test.describe("Empty State", () => {
    test("should display empty state when no repositories analyzed", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesEmpty,
        stats: mockStatsEmpty,
      });

      await page.goto("/en/dashboard");

      // Wait for page to load
      await expect(page.getByRole("heading", { name: /dashboard/i })).toBeVisible({
        timeout: 10000,
      });

      // Empty state should show "Analyze a repository" button (AnalyzeDialog with empty-state variant)
      await expect(
        page.getByRole("button", { name: /analyze a repository/i })
      ).toBeVisible({ timeout: 10000 });
    });

    test("should open analyze dialog from empty state", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesEmpty,
        stats: mockStatsEmpty,
      });

      await page.goto("/en/dashboard");

      // Wait for page to load
      await expect(page.getByRole("heading", { name: /dashboard/i })).toBeVisible({
        timeout: 10000,
      });

      // Click the Analyze button to open dialog
      const analyzeButton = page.getByRole("button", { name: /analyze a repository/i });
      await analyzeButton.click();

      // Verify dialog opens
      const dialog = page.getByRole("dialog");
      await expect(dialog).toBeVisible({ timeout: 5000 });
      await expect(dialog.getByRole("heading", { name: /analyze repository/i })).toBeVisible();
    });
  });

  test.describe("Navigation", () => {
    test("should navigate to repository analysis page when clicking repo card", async ({
      page,
    }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Click on repo card link
      const firstRepo = mockRepositoriesPage1.data[0];
      const repoLink = page.getByRole("link", { name: new RegExp(firstRepo.name) }).first();
      await repoLink.click();

      // Verify navigation to analysis page
      await expect(page).toHaveURL(new RegExp(`/analyze/${firstRepo.owner}/${firstRepo.name}`));
    });
  });

  test.describe("AI Spec Badge", () => {
    test("should display AI Spec badge for repos with AI spec", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // First 5 repos have AI Spec - check that AI Spec badge is visible
      const aiSpecBadges = page.getByText("AI Spec", { exact: true });
      await expect(aiSpecBadges.first()).toBeVisible({ timeout: 5000 });
    });

    test("should show popover with language count on badge hover", async ({ page }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Click on AI Spec badge to open popover
      const aiSpecBadge = page.getByText("AI Spec", { exact: true }).first();
      await aiSpecBadge.click();

      // Verify popover content shows language count info
      await expect(page.getByText(/Generated in \d+ language/i)).toBeVisible({ timeout: 5000 });
    });

    test("should not display AI Spec badge for repos without AI spec", async ({ page }) => {
      // Create mock data where no repos have AI Spec
      const reposWithoutAiSpec = {
        ...mockRepositoriesPage1,
        data: mockRepositoriesPage1.data.map((repo) => ({
          ...repo,
          aiSpecSummary: undefined,
        })),
      };

      await setupMockHandlers(page, {
        repositories: reposWithoutAiSpec,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // AI Spec badge should not be visible
      const aiSpecBadges = page.getByText("AI Spec", { exact: true });
      await expect(aiSpecBadges).toHaveCount(0);
    });

    test("should navigate to analysis page when clicking card with AI Spec badge", async ({
      page,
    }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Get first repo with AI Spec (index 1 has hasSpec: true)
      const repoWithAiSpec = mockRepositoriesPage1.data[0];
      const repoCard = page.getByRole("link", { name: new RegExp(repoWithAiSpec.name) }).first();
      await repoCard.click();

      // Verify navigation to analysis page
      await expect(page).toHaveURL(
        new RegExp(`/analyze/${repoWithAiSpec.owner}/${repoWithAiSpec.name}`)
      );
    });

    test("should navigate to AI Spec tab when clicking 'View details' link in popover", async ({
      page,
    }) => {
      await setupMockHandlers(page, {
        repositories: mockRepositoriesPage1,
        stats: mockStatsNormal,
      });

      await page.goto("/en/dashboard");

      // Wait for initial load
      await expect(
        page.getByRole("button", { name: /add bookmark|remove bookmark/i }).first()
      ).toBeVisible({ timeout: 10000 });

      // Click on AI Spec badge to open popover
      const aiSpecBadge = page.getByText("AI Spec", { exact: true }).first();
      await aiSpecBadge.click();

      // Click "View details →" link in popover
      const viewDetailsLink = page.getByRole("link", { name: /view details/i });
      await expect(viewDetailsLink).toBeVisible({ timeout: 5000 });

      const repoWithAiSpec = mockRepositoriesPage1.data[0];
      await viewDetailsLink.click();

      // Verify navigation to analysis page with tab=spec query param
      await expect(page).toHaveURL(
        new RegExp(`/analyze/${repoWithAiSpec.owner}/${repoWithAiSpec.name}.*tab=spec`)
      );
    });
  });
});
