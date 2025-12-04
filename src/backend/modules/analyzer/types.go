package analyzer

type TestStatus string

const (
	TestStatusActive  TestStatus = "active"
	TestStatusSkipped TestStatus = "skipped"
	TestStatusTodo    TestStatus = "todo"
)

type Framework string

const (
	FrameworkJest       Framework = "jest"
	FrameworkVitest     Framework = "vitest"
	FrameworkPlaywright Framework = "playwright"
	FrameworkGo         Framework = "go"
)

type TestCase struct {
	FilePath  string     `json:"filePath"`
	Framework Framework  `json:"framework"`
	Line      int        `json:"line"`
	Name      string     `json:"name"`
	Status    TestStatus `json:"status"`
}

type TestSuite struct {
	FilePath  string     `json:"filePath"`
	Framework Framework  `json:"framework"`
	Tests     []TestCase `json:"tests"`
}

type FrameworkSummary struct {
	Active    int       `json:"active"`
	Framework Framework `json:"framework"`
	Skipped   int       `json:"skipped"`
	Todo      int       `json:"todo"`
	Total     int       `json:"total"`
}

type Summary struct {
	Active     int                `json:"active"`
	Frameworks []FrameworkSummary `json:"frameworks"`
	Skipped    int                `json:"skipped"`
	Todo       int                `json:"todo"`
	Total      int                `json:"total"`
}

type AnalysisResult struct {
	AnalyzedAt string      `json:"analyzedAt"`
	CommitSHA  string      `json:"commitSha"`
	Owner      string      `json:"owner"`
	Repo       string      `json:"repo"`
	Suites     []TestSuite `json:"suites"`
	Summary    Summary     `json:"summary"`
}
