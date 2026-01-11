package analysis

type Inventory struct {
	Files []TestFile
}

type TestFile struct {
	Path        string
	Framework   string
	DomainHints *DomainHints
	Suites      []TestSuite
	Tests       []Test
}

type DomainHints struct {
	Calls   []string
	Imports []string
}

type TestSuite struct {
	Name     string
	Location Location
	Suites   []TestSuite
	Tests    []Test
}

type Test struct {
	Name     string
	Location Location
	Status   TestStatus
}

type Location struct {
	StartLine int
	EndLine   int
}

type TestStatus string

const (
	TestStatusActive  TestStatus = "active"
	TestStatusFocused TestStatus = "focused"
	TestStatusSkipped TestStatus = "skipped"
	TestStatusTodo    TestStatus = "todo"
	TestStatusXfail   TestStatus = "xfail"
)
