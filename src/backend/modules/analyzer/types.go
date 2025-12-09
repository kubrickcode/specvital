package analyzer

type TestStatus string

const (
	TestStatusActive  TestStatus = "active"
	TestStatusSkipped TestStatus = "skipped"
	TestStatusTodo    TestStatus = "todo"
)

type Framework = string

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

// AnalysisStatusType represents the status of an analysis job.
type AnalysisStatusType string

const (
	StatusCompleted AnalysisStatusType = "completed"
	StatusAnalyzing AnalysisStatusType = "analyzing"
	StatusQueued    AnalysisStatusType = "queued"
	StatusFailed    AnalysisStatusType = "failed"
)

// AnalysisResponse is the API response for analysis endpoints.
type AnalysisResponse struct {
	Data   *AnalysisResult    `json:"data,omitempty"`
	Error  string             `json:"error,omitempty"`
	Status AnalysisStatusType `json:"status"`
}
