package kotest

import (
	"context"
	"testing"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/framework"
)

func TestNewDefinition(t *testing.T) {
	def := NewDefinition()

	if def.Name != framework.FrameworkKotest {
		t.Errorf("expected name %s, got %s", framework.FrameworkKotest, def.Name)
	}

	if len(def.Languages) != 1 || def.Languages[0] != domain.LanguageKotlin {
		t.Errorf("expected [kotlin], got %v", def.Languages)
	}

	if def.Parser == nil {
		t.Error("expected parser to be non-nil")
	}

	if len(def.Matchers) == 0 {
		t.Error("expected at least one matcher")
	}
}

func TestKotestFileMatcher(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		wantMatch bool
		wantConf  int
	}{
		{"spec file", "UserSpec.kt", true, 20},
		{"test file", "UserTest.kt", true, 20},
		{"tests file", "UserTests.kt", true, 20},
		{"kts script test", "UserSpec.kts", true, 20},
		{"non-test file", "User.kt", false, 0},
		{"java file", "UserTest.java", false, 0},
	}

	matcher := &KotestFileMatcher{}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:  framework.SignalFileName,
				Value: tt.filename,
			}

			result := matcher.Match(ctx, signal)

			if tt.wantMatch && result.Confidence != tt.wantConf {
				t.Errorf("expected confidence %d, got %d", tt.wantConf, result.Confidence)
			}
			if !tt.wantMatch && result.Confidence != 0 {
				t.Errorf("expected no match, got confidence %d", result.Confidence)
			}
		})
	}
}

func TestKotestContentMatcher(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantMatch bool
	}{
		{
			name:      "FunSpec class",
			content:   "class MyTest : FunSpec({",
			wantMatch: true,
		},
		{
			name:      "StringSpec class",
			content:   "class MyTest : StringSpec({",
			wantMatch: true,
		},
		{
			name:      "BehaviorSpec class",
			content:   "class MyTest : BehaviorSpec({",
			wantMatch: true,
		},
		{
			name:      "shouldBe matcher",
			content:   `result shouldBe 5`,
			wantMatch: true,
		},
		{
			name:      "kotest import",
			content:   "import io.kotest.core.spec.style.FunSpec",
			wantMatch: true,
		},
		{
			name:      "no kotest patterns",
			content:   "class MyClass { fun doSomething() {} }",
			wantMatch: false,
		},
	}

	matcher := &KotestContentMatcher{}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:    framework.SignalFileContent,
				Value:   tt.content,
				Context: []byte(tt.content),
			}

			result := matcher.Match(ctx, signal)

			if tt.wantMatch && result.Confidence == 0 {
				t.Error("expected match but got no match")
			}
			if !tt.wantMatch && result.Confidence != 0 {
				t.Errorf("expected no match but got confidence %d", result.Confidence)
			}
		})
	}
}

func TestKotestParser_FunSpec(t *testing.T) {
	source := `
package com.example

import io.kotest.core.spec.style.FunSpec
import io.kotest.matchers.shouldBe

class CalculatorTest : FunSpec({
    test("addition works") {
        1 + 1 shouldBe 2
    }

    test("subtraction works") {
        5 - 3 shouldBe 2
    }

    context("multiplication") {
        test("basic multiplication") {
            2 * 3 shouldBe 6
        }
    }

    xtest("skipped test") {
    }
})
`

	parser := &KotestParser{}
	ctx := context.Background()

	result, err := parser.Parse(ctx, []byte(source), "CalculatorTest.kt")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if result.Framework != frameworkName {
		t.Errorf("expected framework %s, got %s", frameworkName, result.Framework)
	}

	if result.Language != domain.LanguageKotlin {
		t.Errorf("expected language kotlin, got %s", result.Language)
	}

	if len(result.Suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(result.Suites))
	}

	suite := result.Suites[0]
	if suite.Name != "CalculatorTest" {
		t.Errorf("expected suite name CalculatorTest, got %s", suite.Name)
	}

	// Validate top-level tests (addition works, subtraction works, skipped test)
	if len(suite.Tests) < 2 {
		t.Errorf("expected at least 2 top-level tests, got %d", len(suite.Tests))
	}

	expectedTests := []struct {
		name   string
		status domain.TestStatus
	}{
		{"addition works", domain.TestStatusActive},
		{"subtraction works", domain.TestStatusActive},
	}

	for i, expected := range expectedTests {
		if i >= len(suite.Tests) {
			break
		}
		if suite.Tests[i].Name != expected.name {
			t.Errorf("test[%d]: expected name %q, got %q", i, expected.name, suite.Tests[i].Name)
		}
		if suite.Tests[i].Status != expected.status {
			t.Errorf("test[%d]: expected status %v, got %v", i, expected.status, suite.Tests[i].Status)
		}
	}

	// Validate nested suites (context "multiplication")
	if len(suite.Suites) != 1 {
		t.Errorf("expected 1 nested suite, got %d", len(suite.Suites))
	} else {
		nestedSuite := suite.Suites[0]
		if nestedSuite.Name != "multiplication" {
			t.Errorf("expected nested suite name 'multiplication', got %s", nestedSuite.Name)
		}
		if len(nestedSuite.Tests) != 1 {
			t.Errorf("expected 1 test in nested suite, got %d", len(nestedSuite.Tests))
		} else if nestedSuite.Tests[0].Name != "basic multiplication" {
			t.Errorf("expected nested test name 'basic multiplication', got %s", nestedSuite.Tests[0].Name)
		}
	}
}

func TestKotestParser_StringSpec(t *testing.T) {
	source := `
package com.example

import io.kotest.core.spec.style.StringSpec
import io.kotest.matchers.shouldBe

class StringSpecTest : StringSpec({
    "length of hello should be 5" {
        "hello".length shouldBe 5
    }

    "startsWith should test for prefix" {
        "world".startsWith("wor") shouldBe true
    }
})
`

	parser := &KotestParser{}
	ctx := context.Background()

	result, err := parser.Parse(ctx, []byte(source), "StringSpecTest.kt")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(result.Suites))
	}

	suite := result.Suites[0]
	if suite.Name != "StringSpecTest" {
		t.Errorf("expected suite name StringSpecTest, got %s", suite.Name)
	}
}

func TestKotestParser_BehaviorSpec(t *testing.T) {
	source := `
package com.example

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe

class BehaviorSpecTest : BehaviorSpec({
    Given("a calculator") {
        When("adding numbers") {
            Then("should return correct sum") {
                1 + 1 shouldBe 2
            }
        }
    }
})
`

	parser := &KotestParser{}
	ctx := context.Background()

	result, err := parser.Parse(ctx, []byte(source), "BehaviorSpecTest.kt")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(result.Suites))
	}

	suite := result.Suites[0]
	if suite.Name != "BehaviorSpecTest" {
		t.Errorf("expected suite name BehaviorSpecTest, got %s", suite.Name)
	}
}

func TestKotestParser_DescribeSpec(t *testing.T) {
	source := `
package com.example

import io.kotest.core.spec.style.DescribeSpec
import io.kotest.matchers.shouldBe

class DescribeSpecTest : DescribeSpec({
    describe("a calculator") {
        it("should add numbers") {
            1 + 1 shouldBe 2
        }

        context("when subtracting") {
            it("should return difference") {
                5 - 3 shouldBe 2
            }
        }
    }
})
`

	parser := &KotestParser{}
	ctx := context.Background()

	result, err := parser.Parse(ctx, []byte(source), "DescribeSpecTest.kt")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(result.Suites))
	}

	suite := result.Suites[0]
	if suite.Name != "DescribeSpecTest" {
		t.Errorf("expected suite name DescribeSpecTest, got %s", suite.Name)
	}
}

func TestKotestParser_AnnotationSpec(t *testing.T) {
	source := `
package com.example

import io.kotest.core.spec.style.AnnotationSpec

class AnnotationSpecTest : AnnotationSpec() {
    @Test
    fun testAddition() {
        assert(1 + 1 == 2)
    }

    @Test
    @Disabled
    fun testDisabled() {
    }
}
`

	parser := &KotestParser{}
	ctx := context.Background()

	result, err := parser.Parse(ctx, []byte(source), "AnnotationSpecTest.kt")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(result.Suites))
	}

	suite := result.Suites[0]
	if suite.Name != "AnnotationSpecTest" {
		t.Errorf("expected suite name AnnotationSpecTest, got %s", suite.Name)
	}
}

func TestKotestParser_NonKotestClass(t *testing.T) {
	source := `
package com.example

class RegularClass {
    fun doSomething(): Int {
        return 42
    }
}
`

	parser := &KotestParser{}
	ctx := context.Background()

	result, err := parser.Parse(ctx, []byte(source), "RegularClass.kt")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Suites) != 0 {
		t.Errorf("expected 0 suites for non-kotest class, got %d", len(result.Suites))
	}
}
