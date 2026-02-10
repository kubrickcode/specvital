package framework_test

import (
	"context"
	"fmt"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/framework"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/framework/matchers"
)

// ExampleDefinition demonstrates how to create a complete framework definition.
func Example_definition() {
	// Create a framework definition for a hypothetical test framework
	def := &framework.Definition{
		Name: "mytest",
		Languages: []domain.Language{
			domain.LanguageTypeScript,
			domain.LanguageJavaScript,
		},
		Matchers: []framework.Matcher{
			// Match import statements
			matchers.NewImportMatcher("mytest", "mytest/"),

			// Match config files
			matchers.NewConfigMatcher("mytest.config.js", "mytest.config.ts"),

			// Match content patterns (optional)
			matchers.NewContentMatcherFromStrings(
				`\bmytest\s*\(`,
				`\bmytest\.describe\s*\(`,
			),
		},
		Priority: framework.PriorityGeneric,
	}

	// Register the framework
	registry := framework.NewRegistry()
	registry.Register(def)

	// Find the framework by name
	found := registry.Find("mytest")
	fmt.Printf("Framework: %s\n", found.Name)
	fmt.Printf("Priority: %d\n", found.Priority)
	fmt.Printf("Languages: %v\n", found.Languages)
	fmt.Printf("Matchers: %d\n", len(found.Matchers))

	// Output:
	// Framework: mytest
	// Priority: 100
	// Languages: [typescript javascript]
	// Matchers: 3
}

// ExampleMatcher demonstrates how to use matchers for framework detection.
func Example_matcher() {
	// Create matchers
	importMatcher := matchers.NewImportMatcher("vitest", "vitest/")
	configMatcher := matchers.NewConfigMatcher("vitest.config.ts")

	ctx := context.Background()

	// Test import signal
	importSignal := framework.Signal{
		Type:  framework.SignalImport,
		Value: "vitest/config",
	}
	result := importMatcher.Match(ctx, importSignal)
	fmt.Printf("Import match confidence: %d\n", result.Confidence)

	// Test config signal
	configSignal := framework.Signal{
		Type:  framework.SignalConfigFile,
		Value: "/project/vitest.config.ts",
	}
	result = configMatcher.Match(ctx, configSignal)
	fmt.Printf("Config match confidence: %d\n", result.Confidence)

	// Output:
	// Import match confidence: 100
	// Config match confidence: 100
}

// ExampleRegistry demonstrates registry operations.
func Example_registry() {
	// Create a new registry
	reg := framework.NewRegistry()

	// Register multiple frameworks
	reg.Register(&framework.Definition{
		Name:      "jest",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Priority:  framework.PriorityGeneric,
	})

	reg.Register(&framework.Definition{
		Name:      "vitest",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Priority:  framework.PrioritySpecialized,
	})

	reg.Register(&framework.Definition{
		Name:      "playwright",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Priority:  framework.PriorityE2E,
	})

	// Get all frameworks (sorted by priority)
	all := reg.All()
	fmt.Println("Frameworks by priority:")
	for _, def := range all {
		fmt.Printf("  %s (priority: %d)\n", def.Name, def.Priority)
	}

	// Find by language
	tsFrameworks := reg.FindByLanguage(domain.LanguageTypeScript)
	fmt.Printf("\nTypeScript frameworks: %d\n", len(tsFrameworks))

	// Find by name
	jest := reg.Find("jest")
	fmt.Printf("Found framework: %s\n", jest.Name)

	// Output:
	// Frameworks by priority:
	//   vitest (priority: 200)
	//   playwright (priority: 150)
	//   jest (priority: 100)
	//
	// TypeScript frameworks: 3
	// Found framework: jest
}

// ExampleProjectScope demonstrates project-wide configuration management.
func Example_projectScope() {
	// Create a project scope
	scope := framework.NewProjectScope()

	// Add framework configurations
	scope.AddConfig("jest.config.js", &framework.ConfigScope{
		Framework:   "jest",
		ConfigPath:  "jest.config.js",
		GlobalsMode: true,
		TestPatterns: []string{
			"**/*.test.js",
			"**/*.spec.js",
		},
		ExcludePatterns: []string{
			"**/node_modules/**",
		},
	})

	scope.AddConfig("apps/web/vitest.config.ts", &framework.ConfigScope{
		Framework:   "vitest",
		ConfigPath:  "apps/web/vitest.config.ts",
		GlobalsMode: false,
		TestPatterns: []string{
			"**/*.test.ts",
		},
	})

	// Find configuration
	jestConfig := scope.FindConfig("jest.config.js")
	fmt.Printf("Jest config globals mode: %t\n", jestConfig.GlobalsMode)

	// Check if config exists
	hasVitest := scope.HasConfigFile("apps/web/vitest.config.ts")
	fmt.Printf("Has Vitest config: %t\n", hasVitest)

	// Output:
	// Jest config globals mode: true
	// Has Vitest config: true
}
