package extraction

import (
	"context"
	"reflect"
	"testing"
)

func TestExtractSwiftImports(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "single import",
			content: `import XCTest
`,
			expected: []string{"XCTest"},
		},
		{
			name: "multiple imports",
			content: `import XCTest
import Foundation
import UIKit
`,
			expected: []string{"XCTest", "Foundation", "UIKit"},
		},
		{
			name: "testable import",
			content: `@testable import MyApp
import XCTest
`,
			expected: []string{"MyApp", "XCTest"},
		},
		{
			name: "multiple testable imports",
			content: `@testable import MyApp
@testable import MyAppCore
import XCTest
`,
			expected: []string{"MyApp", "MyAppCore", "XCTest"},
		},
		{
			name: "no imports",
			content: `class Calculator {
    func add(_ a: Int, _ b: Int) -> Int {
        return a + b
    }
}
`,
			expected: nil,
		},
		{
			name: "dedup duplicate imports",
			content: `import XCTest
import Foundation
import XCTest
`,
			expected: []string{"XCTest", "Foundation"},
		},
		{
			name: "import with attribute",
			content: `@_exported import Foundation
import XCTest
`,
			expected: []string{"Foundation", "XCTest"},
		},
		{
			name: "import in middle of file",
			content: `// Header comment
// Copyright 2024

import XCTest
import Foundation

class MyTests: XCTestCase {
}
`,
			expected: []string{"XCTest", "Foundation"},
		},
		{
			name: "underscore in module name",
			content: `import My_Module
import XCTest
`,
			expected: []string{"My_Module", "XCTest"},
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSwiftImports(ctx, []byte(tt.content))
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
