package domain

// TestStatus represents the execution status of a test.
type TestStatus string

// Test status values.
const (
	// TestStatusFixme indicates a test marked with .fixme (Playwright).
	TestStatusFixme TestStatus = "fixme"
	// TestStatusOnly indicates a test marked to run exclusively (.only).
	TestStatusOnly TestStatus = "only"
	// TestStatusPending indicates a test without implementation.
	TestStatusPending TestStatus = "pending"
	// TestStatusSkipped indicates a test marked to skip (.skip, xit, etc.).
	TestStatusSkipped TestStatus = "skipped"
)
