package xunit

import (
	"context"
	"testing"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/framework"
)

func TestNewDefinition(t *testing.T) {
	def := NewDefinition()

	if def.Name != "xunit" {
		t.Errorf("expected Name='xunit', got '%s'", def.Name)
	}
	if def.Priority != framework.PriorityGeneric {
		t.Errorf("expected Priority=%d, got %d", framework.PriorityGeneric, def.Priority)
	}
	if len(def.Languages) != 1 || def.Languages[0] != domain.LanguageCSharp {
		t.Errorf("expected Languages=[csharp], got %v", def.Languages)
	}
	if def.Parser == nil {
		t.Error("expected Parser to be non-nil")
	}
	if len(def.Matchers) != 3 {
		t.Errorf("expected 3 Matchers, got %d", len(def.Matchers))
	}
}

func TestXUnitFileMatcher_Match(t *testing.T) {
	matcher := &XUnitFileMatcher{}
	ctx := context.Background()

	tests := []struct {
		name               string
		filename           string
		expectedConfidence int
	}{
		{"Test suffix", "CalculatorTest.cs", 20},
		{"Tests suffix", "CalculatorTests.cs", 20},
		{"Test prefix", "TestCalculator.cs", 20},
		{"Test suffix with path", "src/Tests/UserServiceTests.cs", 20},
		{"regular cs file", "Calculator.cs", 0},
		{"non-cs file", "CalculatorTest.java", 0},
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

func TestXUnitContentMatcher_Match(t *testing.T) {
	matcher := &XUnitContentMatcher{}
	ctx := context.Background()

	tests := []struct {
		name               string
		content            string
		expectedConfidence int
	}{
		{
			name: "[Fact] attribute",
			content: `
using Xunit;

public class CalculatorTests
{
    [Fact]
    public void Add_ReturnsSum()
    {
        Assert.Equal(3, 1 + 2);
    }
}
`,
			expectedConfidence: 40,
		},
		{
			name: "[Theory] attribute",
			content: `
using Xunit;

public class CalculatorTests
{
    [Theory]
    [InlineData(1, 2, 3)]
    public void Add_ReturnsSum(int a, int b, int expected)
    {
        Assert.Equal(expected, a + b);
    }
}
`,
			expectedConfidence: 40,
		},
		{
			name: "[InlineData] attribute",
			content: `
using Xunit;

public class CalculatorTests
{
    [InlineData(1, 2)]
    public void TestMethod() {}
}
`,
			expectedConfidence: 40,
		},
		{
			name: "[MemberData] attribute",
			content: `
using Xunit;

public class CalculatorTests
{
    [MemberData(nameof(TestData))]
    public void TestMethod() {}
}
`,
			expectedConfidence: 40,
		},
		{
			name: "using Xunit",
			content: `
using Xunit;

public class SomeClass {}
`,
			expectedConfidence: 40,
		},
		{
			name: "no xUnit patterns",
			content: `
public class Calculator
{
    public int Add(int a, int b)
    {
        return a + b;
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

func TestXUnitParser_Parse(t *testing.T) {
	p := &XUnitParser{}
	ctx := context.Background()

	t.Run("basic [Fact] test methods", func(t *testing.T) {
		source := `
using Xunit;

public class CalculatorTests
{
    [Fact]
    public void Add_ReturnsSum()
    {
        Assert.Equal(3, 1 + 2);
    }

    [Fact]
    public void Subtract_ReturnsDifference()
    {
        Assert.Equal(1, 3 - 2);
    }

    public void HelperMethod()
    {
        // not a test
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "CalculatorTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if testFile.Path != "CalculatorTests.cs" {
			t.Errorf("expected Path='CalculatorTests.cs', got '%s'", testFile.Path)
		}
		if testFile.Framework != "xunit" {
			t.Errorf("expected Framework='xunit', got '%s'", testFile.Framework)
		}
		if testFile.Language != domain.LanguageCSharp {
			t.Errorf("expected Language=csharp, got '%s'", testFile.Language)
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
		if suite.Tests[0].Name != "Add_ReturnsSum" {
			t.Errorf("expected Tests[0].Name='Add_ReturnsSum', got '%s'", suite.Tests[0].Name)
		}
		if suite.Tests[1].Name != "Subtract_ReturnsDifference" {
			t.Errorf("expected Tests[1].Name='Subtract_ReturnsDifference', got '%s'", suite.Tests[1].Name)
		}
	})

	t.Run("[Theory] with [InlineData] counts each attribute", func(t *testing.T) {
		source := `
using Xunit;

public class MathTests
{
    [Theory]
    [InlineData(1, 2, 3)]
    [InlineData(2, 3, 5)]
    public void Add_WithValues_ReturnsSum(int a, int b, int expected)
    {
        Assert.Equal(expected, a + b);
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "MathTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 2 {
			t.Fatalf("expected 2 Tests (one per [InlineData]), got %d", len(suite.Tests))
		}

		if suite.Tests[0].Name != "Add_WithValues_ReturnsSum" {
			t.Errorf("expected Name='Add_WithValues_ReturnsSum', got '%s'", suite.Tests[0].Name)
		}
		if suite.Tests[1].Name != "Add_WithValues_ReturnsSum" {
			t.Errorf("expected Name='Add_WithValues_ReturnsSum', got '%s'", suite.Tests[1].Name)
		}
	})

	t.Run("[Fact(Skip = ...)] marks test as skipped", func(t *testing.T) {
		source := `
using Xunit;

public class SkippedTests
{
    [Fact(Skip = "Not implemented yet")]
    public void SkippedTest()
    {
    }

    [Fact]
    public void NormalTest()
    {
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "SkippedTests.cs")
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

		if suite.Tests[0].Name != "SkippedTest" {
			t.Errorf("expected Tests[0].Name='SkippedTest', got '%s'", suite.Tests[0].Name)
		}
		if suite.Tests[0].Status != domain.TestStatusSkipped {
			t.Errorf("expected Tests[0].Status='skipped', got '%s'", suite.Tests[0].Status)
		}

		if suite.Tests[1].Name != "NormalTest" {
			t.Errorf("expected Tests[1].Name='NormalTest', got '%s'", suite.Tests[1].Name)
		}
		if suite.Tests[1].Status != domain.TestStatusActive {
			t.Errorf("expected Tests[1].Status='active', got '%s'", suite.Tests[1].Status)
		}
	})

	t.Run("[Fact(DisplayName = ...)] uses display name", func(t *testing.T) {
		source := `
using Xunit;

public class DisplayNameTests
{
    [Fact(DisplayName = "Addition should work correctly")]
    public void TestAdd()
    {
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "DisplayNameTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 1 {
			t.Fatalf("expected 1 Test, got %d", len(suite.Tests))
		}

		if suite.Tests[0].Name != "Addition should work correctly" {
			t.Errorf("expected Name='Addition should work correctly', got '%s'", suite.Tests[0].Name)
		}
	})

	t.Run("[Theory(Skip = ...)] with [InlineData] marks all as skipped", func(t *testing.T) {
		source := `
using Xunit;

public class SkippedTheoryTests
{
    [Theory(Skip = "Database not available")]
    [InlineData(1)]
    [InlineData(2)]
    public void SkippedTheory(int value)
    {
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "SkippedTheoryTests.cs")
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

		for i, test := range suite.Tests {
			if test.Status != domain.TestStatusSkipped {
				t.Errorf("expected Tests[%d].Status='skipped', got '%s'", i, test.Status)
			}
		}
	})

	t.Run("nested classes", func(t *testing.T) {
		source := `
using Xunit;

public class OuterTests
{
    [Fact]
    public void OuterTest()
    {
    }

    public class InnerTests
    {
        [Fact]
        public void InnerTest()
        {
        }
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "OuterTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if suite.Name != "OuterTests" {
			t.Errorf("expected Suite.Name='OuterTests', got '%s'", suite.Name)
		}
		if len(suite.Tests) != 1 {
			t.Fatalf("expected 1 Test in outer suite, got %d", len(suite.Tests))
		}
		if suite.Tests[0].Name != "OuterTest" {
			t.Errorf("expected outer test name='OuterTest', got '%s'", suite.Tests[0].Name)
		}

		if len(suite.Suites) != 1 {
			t.Fatalf("expected 1 nested Suite, got %d", len(suite.Suites))
		}

		nested := suite.Suites[0]
		if nested.Name != "InnerTests" {
			t.Errorf("expected nested Suite.Name='InnerTests', got '%s'", nested.Name)
		}
		if len(nested.Tests) != 1 {
			t.Fatalf("expected 1 Test in nested suite, got %d", len(nested.Tests))
		}
		if nested.Tests[0].Name != "InnerTest" {
			t.Errorf("expected nested test name='InnerTest', got '%s'", nested.Tests[0].Name)
		}
	})

	t.Run("[Fact] with both Skip and DisplayName", func(t *testing.T) {
		source := `
using Xunit;

public class CombinedTests
{
    [Fact(Skip = "Not ready", DisplayName = "Should skip with custom name")]
    public void TestMethod()
    {
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "CombinedTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 1 {
			t.Fatalf("expected 1 Test, got %d", len(suite.Tests))
		}

		if suite.Tests[0].Status != domain.TestStatusSkipped {
			t.Errorf("expected Status='skipped', got '%s'", suite.Tests[0].Status)
		}
		if suite.Tests[0].Name != "Should skip with custom name" {
			t.Errorf("expected Name='Should skip with custom name', got '%s'", suite.Tests[0].Name)
		}
	})

	t.Run("multiple classes in file", func(t *testing.T) {
		source := `
using Xunit;

public class FirstTests
{
    [Fact]
    public void FirstTest()
    {
    }
}

public class SecondTests
{
    [Fact]
    public void SecondTest()
    {
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "MultipleClasses.cs")
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

	t.Run("preprocessor directive #if wrapping nested class", func(t *testing.T) {
		source := `
using Xunit;

public class TaskCompletionSourceAssertionSpecs
{
#if NET6_0_OR_GREATER
    public class NonGeneric
    {
        [Fact]
        public void Test1() { }

        [Fact]
        public void Test2() { }
    }
#endif

    public class Generic
    {
        [Fact]
        public void Test3() { }
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "TaskCompletionSourceAssertionSpecs.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if suite.Name != "TaskCompletionSourceAssertionSpecs" {
			t.Errorf("expected Suite.Name='TaskCompletionSourceAssertionSpecs', got '%s'", suite.Name)
		}

		if len(suite.Suites) != 2 {
			t.Fatalf("expected 2 nested Suites, got %d", len(suite.Suites))
		}

		nonGeneric := suite.Suites[0]
		if nonGeneric.Name != "NonGeneric" {
			t.Errorf("expected nested[0].Name='NonGeneric', got '%s'", nonGeneric.Name)
		}
		if len(nonGeneric.Tests) != 2 {
			t.Errorf("expected 2 tests in NonGeneric, got %d", len(nonGeneric.Tests))
		}

		generic := suite.Suites[1]
		if generic.Name != "Generic" {
			t.Errorf("expected nested[1].Name='Generic', got '%s'", generic.Name)
		}
		if len(generic.Tests) != 1 {
			t.Errorf("expected 1 test in Generic, got %d", len(generic.Tests))
		}
	})

	t.Run("preprocessor directive #if wrapping test methods", func(t *testing.T) {
		source := `
using Xunit;

public class DateTimePropertiesSpecs
{
    [Fact]
    public void CommonTest1() { }

#if NET6_0_OR_GREATER
    [Fact]
    public void Net6Test1() { }

    [Fact]
    public void Net6Test2() { }
#endif

    [Fact]
    public void CommonTest2() { }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "DateTimePropertiesSpecs.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 4 {
			t.Fatalf("expected 4 tests, got %d", len(suite.Tests))
		}

		expectedNames := []string{"CommonTest1", "Net6Test1", "Net6Test2", "CommonTest2"}
		for i, name := range expectedNames {
			if suite.Tests[i].Name != name {
				t.Errorf("expected Tests[%d].Name='%s', got '%s'", i, name, suite.Tests[i].Name)
			}
		}
	})

	t.Run("preprocessor directive #if with #else", func(t *testing.T) {
		source := `
using Xunit;

public class ConditionalTests
{
#if NETFRAMEWORK
    [Fact]
    public void FrameworkOnlyTest() { }
#else
    [Fact]
    public void CoreOnlyTest() { }
#endif
}
`
		testFile, err := p.Parse(ctx, []byte(source), "ConditionalTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 2 {
			t.Fatalf("expected 2 tests, got %d", len(suite.Tests))
		}
	})

	t.Run("preprocessor directive #if with #elif", func(t *testing.T) {
		source := `
using Xunit;

public class MultiConditionTests
{
#if NET8_0
    [Fact]
    public void Net8Test() { }
#elif NET6_0
    [Fact]
    public void Net6Test() { }
#else
    [Fact]
    public void LegacyTest() { }
#endif
}
`
		testFile, err := p.Parse(ctx, []byte(source), "MultiConditionTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 3 {
			t.Fatalf("expected 3 tests, got %d", len(suite.Tests))
		}
	})

	t.Run("[Theory] with [MemberData] counts as single test", func(t *testing.T) {
		source := `
using Xunit;

public class MemberDataTests
{
    public static IEnumerable<object[]> TestData => new[] { new object[] { 1 }, new object[] { 2 } };

    [Theory]
    [MemberData(nameof(TestData))]
    public void TestWithMemberData(int value)
    {
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "MemberDataTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 1 {
			t.Fatalf("expected 1 Test (MemberData is runtime-expanded), got %d", len(suite.Tests))
		}

		if suite.Tests[0].Name != "TestWithMemberData" {
			t.Errorf("expected Name='TestWithMemberData', got '%s'", suite.Tests[0].Name)
		}
	})

	t.Run("[Theory] with [ClassData] counts as single test", func(t *testing.T) {
		source := `
using Xunit;

public class ClassDataTests
{
    [Theory]
    [ClassData(typeof(TestDataClass))]
    public void TestWithClassData(int value)
    {
    }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "ClassDataTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 1 {
			t.Fatalf("expected 1 Test (ClassData is runtime-expanded), got %d", len(suite.Tests))
		}

		if suite.Tests[0].Name != "TestWithClassData" {
			t.Errorf("expected Name='TestWithClassData', got '%s'", suite.Tests[0].Name)
		}
	})

	t.Run("custom xUnit attributes [UIFact] [StaFact] [UITheory]", func(t *testing.T) {
		source := `
using Xunit;

public class CustomAttributeTests
{
    [UIFact]
    public void UIFactTest() { }

    [StaFact]
    public void StaFactTest() { }

    [WpfFact]
    public void WpfFactTest() { }

    [UITheory]
    [InlineData(1)]
    [InlineData(2)]
    public void UITheoryTest(int value) { }

    [Fact]
    public void NormalFactTest() { }
}
`
		testFile, err := p.Parse(ctx, []byte(source), "CustomAttributeTests.cs")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(testFile.Suites) != 1 {
			t.Fatalf("expected 1 Suite, got %d", len(testFile.Suites))
		}

		suite := testFile.Suites[0]
		if len(suite.Tests) != 6 {
			t.Fatalf("expected 6 Tests (UIFact + StaFact + WpfFact + 2 InlineData + Fact), got %d", len(suite.Tests))
		}

		expectedNames := []string{"UIFactTest", "StaFactTest", "WpfFactTest", "UITheoryTest", "UITheoryTest", "NormalFactTest"}
		for i, name := range expectedNames {
			if suite.Tests[i].Name != name {
				t.Errorf("expected Tests[%d].Name='%s', got '%s'", i, name, suite.Tests[i].Name)
			}
		}
	})
}

func TestIsFactAttribute(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Fact", true},
		{"FactAttribute", true},
		{"UIFact", true},
		{"UIFactAttribute", true},
		{"StaFact", true},
		{"WpfFact", true},
		{"Theory", false},
		{"Test", false},
		{"Artifact", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFactAttribute(tt.name); got != tt.expected {
				t.Errorf("isFactAttribute(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestIsTheoryAttribute(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Theory", true},
		{"TheoryAttribute", true},
		{"UITheory", true},
		{"UITheoryAttribute", true},
		{"Fact", false},
		{"Test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTheoryAttribute(tt.name); got != tt.expected {
				t.Errorf("isTheoryAttribute(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}
