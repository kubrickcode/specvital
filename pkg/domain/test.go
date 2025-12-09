package domain

// Test represents a single test case (it, test, func TestXxx).
type Test struct {
	// Location is the source code location of this test.
	Location Location `json:"location"`
	// Name is the test description or function name.
	Name string `json:"name"`
	// Status indicates if the test is skipped, only, etc.
	Status TestStatus `json:"status"`
	// Modifier is the original framework marker (skip, todo, fixme, @Disabled, etc.).
	Modifier string `json:"modifier,omitempty"`
}

// TestSuite represents a test suite (describe, test.describe).
type TestSuite struct {
	// Location is the source code location of this suite.
	Location Location `json:"location"`
	// Name is the suite description.
	Name string `json:"name"`
	// Status indicates if the suite is skipped, only, etc.
	Status TestStatus `json:"status"`
	// Modifier is the original framework marker (skip, todo, fixme, @Disabled, etc.).
	Modifier string `json:"modifier,omitempty"`
	// Suites contains nested test suites.
	Suites []TestSuite `json:"suites,omitempty"`
	// Tests contains the tests in this suite.
	Tests []Test `json:"tests,omitempty"`
}

// CountTests returns the total number of tests in this suite including nested suites.
func (s *TestSuite) CountTests() int {
	count := len(s.Tests)
	for _, sub := range s.Suites {
		count += sub.CountTests()
	}
	return count
}
