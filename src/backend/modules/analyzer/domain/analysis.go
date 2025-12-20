package domain

import "time"

type Analysis struct {
	CommitSHA   string
	CompletedAt time.Time
	ID          string
	Owner       string
	Repo        string
	TestSuites  []TestSuite
	TotalSuites int
	TotalTests  int
}

type TestSuite struct {
	FilePath  string
	Framework string
	ID        string
	Name      string
	TestCases []TestCase
}

type TestCase struct {
	Line   int
	Name   string
	Status TestStatus
}

type TestStatus string

const (
	TestStatusActive  TestStatus = "active"
	TestStatusFocused TestStatus = "focused"
	TestStatusSkipped TestStatus = "skipped"
	TestStatusTodo    TestStatus = "todo"
	TestStatusXfail   TestStatus = "xfail"
)

func (s TestStatus) String() string {
	return string(s)
}

type AnalysisProgress struct {
	CommitSHA    string
	CompletedAt  *time.Time
	CreatedAt    time.Time
	ErrorMessage *string
	Status       Status
}
