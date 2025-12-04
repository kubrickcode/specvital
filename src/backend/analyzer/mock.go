package analyzer

import "time"

const mockCommitSHA = "abc1234def5678"

func calculateSummary(suites []TestSuite) Summary {
	frameworkStats := make(map[Framework]*FrameworkSummary)
	var totalActive, totalSkipped, totalTodo int

	for _, suite := range suites {
		stats, exists := frameworkStats[suite.Framework]
		if !exists {
			stats = &FrameworkSummary{Framework: suite.Framework}
			frameworkStats[suite.Framework] = stats
		}

		for _, test := range suite.Tests {
			stats.Total++
			switch test.Status {
			case TestStatusActive:
				stats.Active++
				totalActive++
			case TestStatusSkipped:
				stats.Skipped++
				totalSkipped++
			case TestStatusTodo:
				stats.Todo++
				totalTodo++
			}
		}
	}

	frameworks := make([]FrameworkSummary, 0, len(frameworkStats))
	for _, stats := range frameworkStats {
		frameworks = append(frameworks, *stats)
	}

	return Summary{
		Active:     totalActive,
		Frameworks: frameworks,
		Skipped:    totalSkipped,
		Todo:       totalTodo,
		Total:      totalActive + totalSkipped + totalTodo,
	}
}

// GenerateMockResult creates mock analysis data for development and testing.
func GenerateMockResult(owner, repo string) AnalysisResult {
	suites := []TestSuite{
		{
			FilePath:  "src/components/Button.test.tsx",
			Framework: FrameworkJest,
			Tests: []TestCase{
				{FilePath: "src/components/Button.test.tsx", Framework: FrameworkJest, Line: 5, Name: "renders correctly", Status: TestStatusActive},
				{FilePath: "src/components/Button.test.tsx", Framework: FrameworkJest, Line: 15, Name: "handles click events", Status: TestStatusActive},
				{FilePath: "src/components/Button.test.tsx", Framework: FrameworkJest, Line: 25, Name: "applies custom className", Status: TestStatusActive},
				{FilePath: "src/components/Button.test.tsx", Framework: FrameworkJest, Line: 35, Name: "supports disabled state", Status: TestStatusSkipped},
			},
		},
		{
			FilePath:  "src/hooks/useAuth.test.ts",
			Framework: FrameworkJest,
			Tests: []TestCase{
				{FilePath: "src/hooks/useAuth.test.ts", Framework: FrameworkJest, Line: 10, Name: "returns user when authenticated", Status: TestStatusActive},
				{FilePath: "src/hooks/useAuth.test.ts", Framework: FrameworkJest, Line: 20, Name: "returns null when not authenticated", Status: TestStatusActive},
				{FilePath: "src/hooks/useAuth.test.ts", Framework: FrameworkJest, Line: 30, Name: "handles login flow", Status: TestStatusTodo},
			},
		},
		{
			FilePath:  "src/utils/format.test.ts",
			Framework: FrameworkVitest,
			Tests: []TestCase{
				{FilePath: "src/utils/format.test.ts", Framework: FrameworkVitest, Line: 5, Name: "formats date correctly", Status: TestStatusActive},
				{FilePath: "src/utils/format.test.ts", Framework: FrameworkVitest, Line: 15, Name: "formats currency correctly", Status: TestStatusActive},
				{FilePath: "src/utils/format.test.ts", Framework: FrameworkVitest, Line: 25, Name: "handles edge cases", Status: TestStatusSkipped},
			},
		},
		{
			FilePath:  "e2e/login.spec.ts",
			Framework: FrameworkPlaywright,
			Tests: []TestCase{
				{FilePath: "e2e/login.spec.ts", Framework: FrameworkPlaywright, Line: 8, Name: "user can log in with valid credentials", Status: TestStatusActive},
				{FilePath: "e2e/login.spec.ts", Framework: FrameworkPlaywright, Line: 20, Name: "shows error for invalid credentials", Status: TestStatusActive},
				{FilePath: "e2e/login.spec.ts", Framework: FrameworkPlaywright, Line: 35, Name: "redirects to dashboard after login", Status: TestStatusActive},
			},
		},
		{
			FilePath:  "e2e/checkout.spec.ts",
			Framework: FrameworkPlaywright,
			Tests: []TestCase{
				{FilePath: "e2e/checkout.spec.ts", Framework: FrameworkPlaywright, Line: 10, Name: "completes purchase flow", Status: TestStatusActive},
				{FilePath: "e2e/checkout.spec.ts", Framework: FrameworkPlaywright, Line: 25, Name: "validates payment information", Status: TestStatusSkipped},
				{FilePath: "e2e/checkout.spec.ts", Framework: FrameworkPlaywright, Line: 40, Name: "handles out of stock items", Status: TestStatusTodo},
			},
		},
	}

	summary := calculateSummary(suites)

	return AnalysisResult{
		AnalyzedAt: time.Now().UTC().Format(time.RFC3339),
		CommitSHA:  mockCommitSHA,
		Owner:      owner,
		Repo:       repo,
		Suites:     suites,
		Summary:    summary,
	}
}
