package domain_hints

import (
	"context"
	"strings"
	"testing"

	"github.com/specvital/core/pkg/domain"
)

func TestCSharpExtractor_Extract_Usings(t *testing.T) {
	source := []byte(`
using System;
using System.Collections.Generic;
using NUnit.Framework;
using MyApp.Services;
using MyApp.Models;

namespace MyApp.Tests
{
    public class CalculatorTests
    {
    }
}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Domain usings should be included (System.* filtered as stdlib)
	expectedUsings := map[string]bool{
		"NUnit.Framework": true,
		"MyApp.Services":  true,
		"MyApp.Models":    true,
	}

	usingSet := make(map[string]bool)
	for _, u := range hints.Imports {
		usingSet[u] = true
	}

	for u := range expectedUsings {
		if !usingSet[u] {
			t.Errorf("expected using %q to be included, got %v", u, hints.Imports)
		}
	}

	// System imports should be filtered
	excludedUsings := []string{"System", "System.Collections.Generic"}
	for _, u := range excludedUsings {
		if usingSet[u] {
			t.Errorf("expected System using %q to be filtered, got %v", u, hints.Imports)
		}
	}
}

func TestCSharpExtractor_Extract_Calls(t *testing.T) {
	source := []byte(`
using NUnit.Framework;

namespace MyApp.Tests
{
    public class CalculatorTests
    {
        [Test]
        public void TestAdd()
        {
            var calculator = new Calculator();
            var result = calculator.Add(1, 2);
            userService.FindById(123);
            paymentGateway.Process(order);
        }
    }
}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedCalls := map[string]bool{
		"calculator.Add":         true,
		"userService.FindById":   true,
		"paymentGateway.Process": true,
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	for call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected call %q to be included, got %v", call, hints.Calls)
		}
	}
}

func TestCSharpExtractor_Extract_EmptyFile(t *testing.T) {
	source := []byte(`// empty file`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints != nil {
		t.Errorf("expected nil for empty file, got %+v", hints)
	}
}

func TestCSharpExtractor_Extract_TestFrameworkCalls(t *testing.T) {
	source := []byte(`
using NUnit.Framework;

namespace MyApp.Tests
{
    public class CalculatorTests
    {
        [Test]
        public void TestAdd()
        {
            var calculator = new Calculator();
            Assert.AreEqual(2, calculator.Add(1, 1));
            Assert.IsTrue(calculator.IsPositive(5));
            Assert.Throws<Exception>(() => calculator.Divide(1, 0));
            userService.Validate(user);
        }
    }
}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Test framework calls should be excluded
	excludedCalls := []string{"Assert.AreEqual", "Assert.IsTrue", "Assert.Throws"}
	for _, call := range excludedCalls {
		if callSet[call] {
			t.Errorf("expected test framework call %q to be excluded", call)
		}
	}

	// Domain calls should be included
	if !callSet["calculator.Add"] {
		t.Errorf("expected calculator.Add call, got %v", hints.Calls)
	}
	if !callSet["userService.Validate"] {
		t.Errorf("expected userService.Validate call, got %v", hints.Calls)
	}
}

func TestCSharpExtractor_Extract_UsingAlias(t *testing.T) {
	source := []byte(`
using System;
using Env = System.Environment;
using Console = System.Console;
using MyAlias = MyApp.CustomType;

namespace MyApp {}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	usingSet := make(map[string]bool)
	for _, u := range hints.Imports {
		usingSet[u] = true
	}

	// System.* should be filtered (even aliased)
	excludedUsings := []string{"System", "System.Environment", "System.Console"}
	for _, u := range excludedUsings {
		if usingSet[u] {
			t.Errorf("expected System using %q to be filtered, got %v", u, hints.Imports)
		}
	}

	// Non-System aliased usings should be included
	if !usingSet["MyApp.CustomType"] {
		t.Errorf("expected MyApp.CustomType using (aliased), got %v", hints.Imports)
	}
}

func TestCSharpExtractor_Extract_ChainedCalls(t *testing.T) {
	source := []byte(`
namespace MyApp.Tests
{
    public class Test
    {
        void TestMethod()
        {
            // Long chains should be normalized to 2 segments
            client.Api.Users.FindAll();
            response.Data.Items.First().GetValue();
        }
    }
}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Should be normalized to 2 segments
	expectedCalls := []string{"client.Api", "response.Data"}
	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected %q call (2-segment normalized), got %v", call, hints.Calls)
		}
	}
}

func TestCSharpExtractor_Extract_StaticUsing(t *testing.T) {
	source := []byte(`
using System;
using static System.Console;
using static System.Math;
using static MyApp.Helpers.StringHelper;

namespace MyApp {}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	usingSet := make(map[string]bool)
	for _, u := range hints.Imports {
		usingSet[u] = true
	}

	// System.* static usings should be filtered
	excludedUsings := []string{"System", "System.Console", "System.Math"}
	for _, u := range excludedUsings {
		if usingSet[u] {
			t.Errorf("expected System using %q to be filtered, got %v", u, hints.Imports)
		}
	}

	// Non-System static usings should be included
	if !usingSet["MyApp.Helpers.StringHelper"] {
		t.Errorf("expected MyApp.Helpers.StringHelper using (static), got %v", hints.Imports)
	}
}

func TestCSharpExtractor_Extract_GlobalUsing(t *testing.T) {
	source := []byte(`
global using System;
global using System.Linq;
global using MyApp.Common;

namespace MyApp {}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	usingSet := make(map[string]bool)
	for _, u := range hints.Imports {
		usingSet[u] = true
	}

	// System.* global usings should be filtered
	excludedUsings := []string{"System", "System.Linq"}
	for _, u := range excludedUsings {
		if usingSet[u] {
			t.Errorf("expected System using %q to be filtered, got %v", u, hints.Imports)
		}
	}

	// Non-System global usings should be included
	if !usingSet["MyApp.Common"] {
		t.Errorf("expected MyApp.Common using (global), got %v", hints.Imports)
	}
}

func TestCSharpExtractor_Extract_XUnitTest(t *testing.T) {
	source := []byte(`
using Xunit;
using FluentAssertions;
using MyApp.Services;

namespace MyApp.Tests
{
    public class UserServiceTests
    {
        [Fact]
        public void GetUser_ReturnsUser()
        {
            var service = new UserService();
            var result = service.GetUser(1);
            result.Should().NotBeNull();
            orderService.CreateOrder(result);
        }
    }
}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Check usings
	usingSet := make(map[string]bool)
	for _, u := range hints.Imports {
		usingSet[u] = true
	}

	expectedUsings := []string{"Xunit", "FluentAssertions", "MyApp.Services"}
	for _, u := range expectedUsings {
		if !usingSet[u] {
			t.Errorf("expected using %q, got %v", u, hints.Imports)
		}
	}

	// Check calls
	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Domain calls should be included
	if !callSet["service.GetUser"] {
		t.Errorf("expected service.GetUser call, got %v", hints.Calls)
	}
	if !callSet["orderService.CreateOrder"] {
		t.Errorf("expected orderService.CreateOrder call, got %v", hints.Calls)
	}

	// FluentAssertions calls should be excluded (Should is in the filter)
	if callSet["Should.NotBeNull"] {
		t.Errorf("expected FluentAssertions call to be excluded, got %v", hints.Calls)
	}
}

func TestGetExtractor_CSharp(t *testing.T) {
	ext := GetExtractor(domain.LanguageCSharp)
	if ext == nil {
		t.Error("expected extractor for CSharp, got nil")
	}

	_, ok := ext.(*CSharpExtractor)
	if !ok {
		t.Errorf("expected CSharpExtractor, got %T", ext)
	}
}

func TestCSharpExtractor_Extract_NoiseFilter(t *testing.T) {
	source := []byte(`
namespace MyApp.Tests
{
    public class Test
    {
        void TestMethod()
        {
            // Decimal literal patterns should be filtered
            var x = (0.5).ToString();
            var y = (1.0).GetType();

            // Real domain calls should be included
            userService.Create(data);
            paymentGateway.Process(amount);
        }
    }
}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	t.Run("noise patterns filtered", func(t *testing.T) {
		for call := range callSet {
			if len(call) > 0 && (call[0] == '[' || call[0] == '(') {
				t.Errorf("noise pattern %q should be filtered, got %v", call, hints.Calls)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		expectedCalls := []string{"userService.Create", "paymentGateway.Process"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}

func TestCSharpExtractor_Extract_ObjectMethods(t *testing.T) {
	source := []byte(`
namespace MyApp.Tests
{
    public class Test
    {
        void TestMethod()
        {
            // Object base methods should be filtered
            var s = obj.ToString();
            var eq = obj.Equals(other);
            var hash = obj.GetHashCode();
            var type = obj.GetType();

            // Domain calls should be included
            userService.Create(data);
            paymentGateway.Process(amount);
        }
    }
}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	t.Run("Object methods filtered", func(t *testing.T) {
		excludedCalls := []string{"obj.ToString", "obj.Equals", "obj.GetHashCode", "obj.GetType"}
		for _, call := range excludedCalls {
			if callSet[call] {
				t.Errorf("expected Object method %q to be excluded, got %v", call, hints.Calls)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		expectedCalls := []string{"userService.Create", "paymentGateway.Process"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}

func TestCSharpExtractor_Extract_MultilineGenericUsing(t *testing.T) {
	source := []byte(`
using System;
using VerifyCS = MSTest.Analyzers.Test.CSharpCodeFixVerifier<
    MSTest.Analyzers.AssemblyCleanupShouldBeValidAnalyzer,
    MSTest.Analyzers.AssemblyCleanupShouldBeValidFixer>;

namespace MyApp {}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	usingSet := make(map[string]bool)
	for _, u := range hints.Imports {
		usingSet[u] = true
	}

	t.Run("System using filtered", func(t *testing.T) {
		if usingSet["System"] {
			t.Errorf("expected System using to be filtered, got %v", hints.Imports)
		}
	})

	t.Run("no newlines in imports", func(t *testing.T) {
		for _, u := range hints.Imports {
			if strings.Contains(u, "\n") {
				t.Errorf("import should not contain newlines: %q", u)
			}
			if strings.Contains(u, "\t") {
				t.Errorf("import should not contain tabs: %q", u)
			}
		}
	})

	t.Run("multiline generic normalized", func(t *testing.T) {
		// The multiline generic type should be normalized to single line
		found := false
		for _, u := range hints.Imports {
			if strings.Contains(u, "CSharpCodeFixVerifier") {
				found = true
				// Should contain the full type but normalized
				if !strings.Contains(u, "AssemblyCleanupShouldBeValidAnalyzer") {
					t.Errorf("expected normalized generic type import, got %q", u)
				}
				break
			}
		}
		if !found {
			t.Errorf("expected CSharpCodeFixVerifier import, got %v", hints.Imports)
		}
	})
}

func TestNormalizeWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple string", "hello", "hello"},
		{"leading space", "  hello", "hello"},
		{"trailing space", "hello  ", "hello"},
		{"newlines", "hello\nworld", "hello world"},
		{"tabs", "hello\tworld", "hello world"},
		{"mixed whitespace", "hello\n\t  world", "hello world"},
		{"multiline generic", "Type<\n    A,\n    B>", "Type< A, B>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeWhitespace(tt.input)
			if got != tt.want {
				t.Errorf("normalizeWhitespace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsCSharpObjectMethod(t *testing.T) {
	tests := []struct {
		call string
		want bool
	}{
		{"obj.ToString", true},
		{"obj.Equals", true},
		{"obj.GetHashCode", true},
		{"obj.GetType", true},
		{"ReferenceEquals", true},
		{"obj.CustomMethod", false},
		{"userService.Create", false},
		{"ToString", true}, // standalone
	}

	for _, tt := range tests {
		t.Run(tt.call, func(t *testing.T) {
			got := isCSharpObjectMethod(tt.call)
			if got != tt.want {
				t.Errorf("isCSharpObjectMethod(%q) = %v, want %v", tt.call, got, tt.want)
			}
		})
	}
}

func TestIsCSharpStdlibImport(t *testing.T) {
	tests := []struct {
		importPath string
		want       bool
	}{
		// Exact matches
		{"System", true},
		{"System.Collections.Generic", true},
		{"System.Linq", true},
		{"System.IO", true},
		{"System.Threading.Tasks", true},
		// Prefix matches (System.*)
		{"System.Collections.Concurrent", true},
		{"System.Text.Json", true},
		{"System.Net.Http", true},
		{"System.Security.Cryptography", true},
		// Non-stdlib
		{"NUnit.Framework", false},
		{"MyApp.Services", false},
		{"Xunit", false},
		{"Microsoft.Extensions", false},
		{"Newtonsoft.Json", false},
		{"FluentAssertions", false},
	}

	for _, tt := range tests {
		t.Run(tt.importPath, func(t *testing.T) {
			got := isCSharpStdlibImport(tt.importPath)
			if got != tt.want {
				t.Errorf("isCSharpStdlibImport(%q) = %v, want %v", tt.importPath, got, tt.want)
			}
		})
	}
}

func TestCSharpExtractor_Extract_NameofFiltered(t *testing.T) {
	source := []byte(`
using NUnit.Framework;

namespace MyApp.Tests
{
    public class Test
    {
        void TestMethod()
        {
            // nameof is a C# compile-time operator, should be filtered
            var name = nameof(TestMethod);
            var propName = nameof(MyClass.Property);

            // Real domain calls should be included
            userService.Create(data);
            logger.LogInfo(nameof(TestMethod));
        }
    }
}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	t.Run("nameof filtered", func(t *testing.T) {
		if callSet["nameof"] {
			t.Errorf("expected nameof to be filtered, got %v", hints.Calls)
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		if !callSet["userService.Create"] {
			t.Errorf("expected userService.Create call, got %v", hints.Calls)
		}
		if !callSet["logger.LogInfo"] {
			t.Errorf("expected logger.LogInfo call, got %v", hints.Calls)
		}
	})
}

func TestCSharpExtractor_Extract_StdlibFiltering(t *testing.T) {
	source := []byte(`
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.IO;
using System.Text.RegularExpressions;
using NUnit.Framework;
using MyApp.Services;
using Newtonsoft.Json;

namespace MyApp.Tests
{
    public class Test {}
}
`)

	extractor := &CSharpExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	usingSet := make(map[string]bool)
	for _, u := range hints.Imports {
		usingSet[u] = true
	}

	t.Run("System imports filtered", func(t *testing.T) {
		excludedUsings := []string{
			"System",
			"System.Collections.Generic",
			"System.Linq",
			"System.Threading.Tasks",
			"System.IO",
			"System.Text.RegularExpressions",
		}
		for _, u := range excludedUsings {
			if usingSet[u] {
				t.Errorf("expected System using %q to be filtered, got %v", u, hints.Imports)
			}
		}
	})

	t.Run("domain imports included", func(t *testing.T) {
		expectedUsings := []string{"NUnit.Framework", "MyApp.Services", "Newtonsoft.Json"}
		for _, u := range expectedUsings {
			if !usingSet[u] {
				t.Errorf("expected domain using %q to be included, got %v", u, hints.Imports)
			}
		}
	})
}
