package entity

type AnalysisStatus string

const (
	AnalysisStatusCompleted AnalysisStatus = "completed"
	AnalysisStatusFailed    AnalysisStatus = "failed"
	AnalysisStatusPending   AnalysisStatus = "pending"
	AnalysisStatusRunning   AnalysisStatus = "running"
)

func (s AnalysisStatus) IsTerminal() bool {
	return s == AnalysisStatusCompleted || s == AnalysisStatusFailed
}

func (s AnalysisStatus) IsPending() bool {
	return s == AnalysisStatusPending || s == AnalysisStatusRunning
}

func (s AnalysisStatus) String() string {
	return string(s)
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

type UpdateStatus string

const (
	UpdateStatusNewCommits UpdateStatus = "new-commits"
	UpdateStatusUnknown    UpdateStatus = "unknown"
	UpdateStatusUpToDate   UpdateStatus = "up-to-date"
)

func (s UpdateStatus) String() string {
	return string(s)
}
