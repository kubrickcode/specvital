package domain

// TestFile represents a parsed test file.
type TestFile struct {
	// Framework is the detected test framework (e.g., "jest", "vitest").
	Framework string `json:"framework"`
	// Language is the programming language of this file.
	Language Language `json:"language"`
	// Path is the file path.
	Path string `json:"path"`
	// Suites contains the test suites in this file.
	Suites []TestSuite `json:"suites,omitempty"`
	// Tests contains the top-level tests in this file (outside any suite).
	Tests []Test `json:"tests,omitempty"`
}

// CountTests returns the total number of tests in this file.
func (f *TestFile) CountTests() int {
	count := len(f.Tests)
	for _, s := range f.Suites {
		count += s.CountTests()
	}
	return count
}

// Inventory represents a collection of test files in a project.
type Inventory struct {
	// Files contains all parsed test files.
	Files []TestFile `json:"files"`
	// RootPath is the root directory path of the scanned project.
	RootPath string `json:"rootPath"`
}

// CountTests returns the total number of tests across all files.
func (inv Inventory) CountTests() int {
	count := 0
	for _, f := range inv.Files {
		count += f.CountTests()
	}
	return count
}
