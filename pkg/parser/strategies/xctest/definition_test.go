package xctest

import (
	"context"
	"testing"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/framework"
)

func TestNewDefinition(t *testing.T) {
	def := NewDefinition()

	if def.Name != "xctest" {
		t.Errorf("expected Name='xctest', got '%s'", def.Name)
	}
	if def.Priority != framework.PriorityGeneric {
		t.Errorf("expected Priority=%d, got %d", framework.PriorityGeneric, def.Priority)
	}
	if len(def.Languages) != 1 || def.Languages[0] != domain.LanguageSwift {
		t.Errorf("expected Languages=[swift], got %v", def.Languages)
	}
	if def.Parser == nil {
		t.Error("expected Parser to be non-nil")
	}
	if len(def.Matchers) != 3 {
		t.Errorf("expected 3 Matchers, got %d", len(def.Matchers))
	}
}

func TestXCTestFileMatcher_Match(t *testing.T) {
	matcher := &XCTestFileMatcher{}
	ctx := context.Background()

	tests := []struct {
		name               string
		filename           string
		expectedConfidence int
	}{
		{"Test suffix", "CalculatorTest.swift", 20},
		{"Tests suffix", "CalculatorTests.swift", 20},
		{"Test suffix with path", "Tests/AppTests/CalculatorTests.swift", 20},
		{"regular swift file", "Calculator.swift", 0},
		{"non-swift file", "CalculatorTest.java", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:  framework.SignalFileName,
				Value: tt.filename,
			}

			result := matcher.Match(ctx, signal)

			if result.Confidence != tt.expectedConfidence {
				t.Errorf("expected Confidence=%d, got %d", tt.expectedConfidence, result.Confidence)
			}
		})
	}
}

func TestXCTestContentMatcher_Match(t *testing.T) {
	matcher := &XCTestContentMatcher{}
	ctx := context.Background()

	tests := []struct {
		name               string
		content            string
		expectedConfidence int
	}{
		{
			name: "XCTestCase class",
			content: `
import XCTest

class CalculatorTests: XCTestCase {
    func testAdd() {
        XCTAssertEqual(3, 1 + 2)
    }
}
`,
			expectedConfidence: 40,
		},
		{
			name: "test method pattern",
			content: `
class MyTests {
    func testSomething() {
    }
}
`,
			expectedConfidence: 40,
		},
		{
			name: "XCTAssert",
			content: `
XCTAssertTrue(result)
`,
			expectedConfidence: 40,
		},
		{
			name: "XCTFail",
			content: `
XCTFail("should not reach here")
`,
			expectedConfidence: 40,
		},
		{
			name: "XCTest import",
			content: `
import XCTest
import Foundation
`,
			expectedConfidence: 40,
		},
		{
			name: "no XCTest patterns",
			content: `
class Calculator {
    func add(_ a: Int, _ b: Int) -> Int {
        return a + b
    }
}
`,
			expectedConfidence: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:    framework.SignalFileContent,
				Value:   tt.content,
				Context: []byte(tt.content),
			}

			result := matcher.Match(ctx, signal)

			if result.Confidence != tt.expectedConfidence {
				t.Errorf("expected Confidence=%d, got %d", tt.expectedConfidence, result.Confidence)
			}
		})
	}
}

func TestXCTestParser_Parse(t *testing.T) {
	p := &XCTestParser{}
	ctx := context.Background()

	t.Run("basic test methods", func(t *testing.T) {
		source := `
import XCTest

class CalculatorTests: XCTestCase {
    func testAdd() {
        XCTAssertEqual(3, 1 + 2)
    }

    func testSubtract() {
        XCTAssertEqual(1, 3 - 2)
    }

    func helperMethod() {
        // not a test
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "CalculatorTests.swift")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if testFile.Path != "CalculatorTests.swift" {
			t.Errorf("expected Path='CalculatorTests.swift', got '%s'", testFile.Path)
		}
		if testFile.Framework != "xctest" {
			t.Errorf("expected Framework='xctest', got '%s'", testFile.Framework)
		}
		if testFile.Language != domain.LanguageSwift {
			t.Errorf("expected Language=swift, got '%s'", testFile.Language)
		}
		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if suite.Name != "CalculatorTests" {
			t.Errorf("expected Suite.Name='CalculatorTests', got '%s'", suite.Name)
		}
		if len(suite.Tests) != 2 {
			t.Fatalf("expected 2 Tests in suite, got %d", len(suite.Tests))
		}
		if suite.Tests[0].Name != "testAdd" {
			t.Errorf("expected Tests[0].Name='testAdd', got '%s'", suite.Tests[0].Name)
		}
		if suite.Tests[1].Name != "testSubtract" {
			t.Errorf("expected Tests[1].Name='testSubtract', got '%s'", suite.Tests[1].Name)
		}
	})

	t.Run("XCTSkip handling", func(t *testing.T) {
		source := `
import XCTest

class SkipTests: XCTestCase {
    func testSkipped() throws {
        throw XCTSkip("not implemented")
    }

    func testNormal() {
        XCTAssertTrue(true)
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "SkipTests.swift")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 2 {
			t.Fatalf("expected 2 Tests, got %d", len(suite.Tests))
		}

		if suite.Tests[0].Name != "testSkipped" {
			t.Errorf("expected Tests[0].Name='testSkipped', got '%s'", suite.Tests[0].Name)
		}
		if suite.Tests[0].Status != domain.TestStatusSkipped {
			t.Errorf("expected Tests[0].Status='skipped', got '%s'", suite.Tests[0].Status)
		}

		if suite.Tests[1].Name != "testNormal" {
			t.Errorf("expected Tests[1].Name='testNormal', got '%s'", suite.Tests[1].Name)
		}
		if suite.Tests[1].Status != domain.TestStatusActive {
			t.Errorf("expected Tests[1].Status='active', got '%s'", suite.Tests[1].Status)
		}
	})

	t.Run("async test methods", func(t *testing.T) {
		source := `
import XCTest

class AsyncTests: XCTestCase {
    func testAsync() async {
        let result = await fetchData()
        XCTAssertNotNil(result)
    }

    func testAsyncThrows() async throws {
        let result = try await fetchDataThrows()
        XCTAssertNotNil(result)
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "AsyncTests.swift")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 2 {
			t.Fatalf("expected 2 Tests, got %d", len(suite.Tests))
		}

		if suite.Tests[0].Name != "testAsync" {
			t.Errorf("expected Tests[0].Name='testAsync', got '%s'", suite.Tests[0].Name)
		}
		if suite.Tests[0].Modifier != "async" {
			t.Errorf("expected Tests[0].Modifier='async', got '%s'", suite.Tests[0].Modifier)
		}
	})

	t.Run("multiple test classes", func(t *testing.T) {
		source := `
import XCTest

class FirstTests: XCTestCase {
    func testFirst() {
        XCTAssertTrue(true)
    }
}

class SecondTests: XCTestCase {
    func testSecond() {
        XCTAssertTrue(true)
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "MultipleTests.swift")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 2 {
			t.Fatalf("expected 2 Suites, got %d", len(testFile.Suites))
		}

		if testFile.Suites[0].Name != "FirstTests" {
			t.Errorf("expected Suites[0].Name='FirstTests', got '%s'", testFile.Suites[0].Name)
		}
		if testFile.Suites[1].Name != "SecondTests" {
			t.Errorf("expected Suites[1].Name='SecondTests', got '%s'", testFile.Suites[1].Name)
		}
	})

	t.Run("non-XCTestCase class ignored", func(t *testing.T) {
		source := `
import XCTest

class NotATestClass {
    func testSomething() {
        // this should not be detected
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "NotATest.swift")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 0 {
			t.Errorf("expected 0 Suites, got %d", len(testFile.Suites))
		}
	})

	t.Run("setUp and tearDown not detected as tests", func(t *testing.T) {
		source := `
import XCTest

class SetupTests: XCTestCase {
    override func setUp() {
        super.setUp()
    }

    override func tearDown() {
        super.tearDown()
    }

    func testActual() {
        XCTAssertTrue(true)
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "SetupTests.swift")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 1 {
			t.Fatalf("expected 1 Test (setUp/tearDown excluded), got %d", len(suite.Tests))
		}
		if suite.Tests[0].Name != "testActual" {
			t.Errorf("expected Tests[0].Name='testActual', got '%s'", suite.Tests[0].Name)
		}
	})
}
