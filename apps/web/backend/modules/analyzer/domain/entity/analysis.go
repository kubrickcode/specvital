package entity

import "time"

type Analysis struct {
	BranchName    *string
	CommitSHA     string
	CommittedAt   *time.Time
	CompletedAt   time.Time
	ID            string
	Owner         string
	ParserVersion *string
	Repo          string
	TestSuites    []TestSuite
	TotalSuites   int
	TotalTests    int
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

type AnalysisProgress struct {
	CommitSHA    string
	CompletedAt  *time.Time
	CreatedAt    time.Time
	ErrorMessage *string
	StartedAt    *time.Time
	Status       AnalysisStatus
}
