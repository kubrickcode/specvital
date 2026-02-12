package cypress

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kubrickcode/specvital/lib/parser/domain"
	"github.com/kubrickcode/specvital/lib/parser/framework"
)

func TestNewDefinition(t *testing.T) {
	def := NewDefinition()

	assert.Equal(t, "cypress", def.Name)
	assert.Equal(t, framework.PriorityE2E, def.Priority)
	assert.ElementsMatch(t,
		[]domain.Language{domain.LanguageTypeScript, domain.LanguageJavaScript},
		def.Languages,
	)
	assert.NotNil(t, def.ConfigParser)
	assert.NotNil(t, def.Parser)
	assert.Len(t, def.Matchers, 4) // ImportMatcher + ConfigMatcher + FilenameMatcher + ContentMatcher
}

func TestCypressFilenameMatcher_Match(t *testing.T) {
	matcher := &CypressFilenameMatcher{}
	ctx := context.Background()

	tests := []struct {
		name        string
		filename    string
		shouldMatch bool
	}{
		{
			name:        "cypress e2e test .cy.ts",
			filename:    "login.cy.ts",
			shouldMatch: true,
		},
		{
			name:        "cypress e2e test .cy.js",
			filename:    "dashboard.cy.js",
			shouldMatch: true,
		},
		{
			name:        "cypress component test .cy.tsx",
			filename:    "Button.cy.tsx",
			shouldMatch: true,
		},
		{
			name:        "cypress component test .cy.jsx",
			filename:    "Card.cy.jsx",
			shouldMatch: true,
		},
		{
			name:        "full path with .cy.ts",
			filename:    "cypress/e2e/login.cy.ts",
			shouldMatch: true,
		},
		{
			name:        "regular test file (not Cypress)",
			filename:    "login.test.ts",
			shouldMatch: false,
		},
		{
			name:        "spec file (not Cypress pattern)",
			filename:    "login.spec.ts",
			shouldMatch: false,
		},
		{
			name:        "just .cy without extension",
			filename:    "login.cy",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:  framework.SignalFileName,
				Value: tt.filename,
			}

			result := matcher.Match(ctx, signal)

			if tt.shouldMatch {
				assert.Equal(t, 100, result.Confidence, "Expected definite match")
				assert.NotEmpty(t, result.Evidence)
			} else {
				assert.Equal(t, 0, result.Confidence, "Should not match")
			}
		})
	}
}

func TestCypressFilenameMatcher_NonFilenameSignal(t *testing.T) {
	matcher := &CypressFilenameMatcher{}
	ctx := context.Background()

	signal := framework.Signal{
		Type:  framework.SignalImport,
		Value: "cypress",
	}

	result := matcher.Match(ctx, signal)
	assert.Equal(t, 0, result.Confidence)
}

func TestCypressContentMatcher_Match(t *testing.T) {
	matcher := &CypressContentMatcher{}
	ctx := context.Background()

	tests := []struct {
		name               string
		content            string
		expectedConfidence int
		shouldMatch        bool
	}{
		{
			name: "cy.visit() pattern",
			content: `
describe('Login Page', () => {
  it('should display login form', () => {
    cy.visit('/login');
    cy.get('input[name="email"]').should('be.visible');
  });
});
`,
			expectedConfidence: 40,
			shouldMatch:        true,
		},
		{
			name: "cy.get() pattern",
			content: `
describe('Dashboard', () => {
  it('should show user name', () => {
    cy.get('.user-name').should('contain', 'John');
  });
});
`,
			expectedConfidence: 40,
			shouldMatch:        true,
		},
		{
			name: "cy.intercept() pattern",
			content: `
describe('API Tests', () => {
  it('should intercept API call', () => {
    cy.intercept('GET', '/api/users', { fixture: 'users.json' });
    cy.visit('/dashboard');
  });
});
`,
			expectedConfidence: 40,
			shouldMatch:        true,
		},
		{
			name: "cy.request() pattern",
			content: `
describe('API', () => {
  it('should make request', () => {
    cy.request('GET', '/api/health').its('status').should('eq', 200);
  });
});
`,
			expectedConfidence: 40,
			shouldMatch:        true,
		},
		{
			name: "cy.wait() pattern",
			content: `
describe('Async', () => {
  it('should wait for request', () => {
    cy.intercept('GET', '/api/data').as('getData');
    cy.visit('/page');
    cy.wait('@getData');
  });
});
`,
			expectedConfidence: 40,
			shouldMatch:        true,
		},
		{
			name: "cy.fixture() pattern",
			content: `
describe('Fixtures', () => {
  it('should use fixture', () => {
    cy.fixture('user.json').then((user) => {
      expect(user.name).to.equal('John');
    });
  });
});
`,
			expectedConfidence: 40,
			shouldMatch:        true,
		},
		{
			name: "Cypress.Commands.add() pattern",
			content: `
Cypress.Commands.add('login', (email, password) => {
  cy.visit('/login');
  cy.get('#email').type(email);
  cy.get('#password').type(password);
  cy.get('button[type="submit"]').click();
});
`,
			expectedConfidence: 40,
			shouldMatch:        true,
		},
		{
			name: "Cypress.env() pattern",
			content: `
describe('Environment', () => {
  it('should read env variable', () => {
    const apiUrl = Cypress.env('API_URL');
    cy.visit(apiUrl);
  });
});
`,
			expectedConfidence: 40,
			shouldMatch:        true,
		},
		{
			name: "no Cypress patterns (plain Jest)",
			content: `
import { describe, test, expect } from '@jest/globals';

describe('Calculator', () => {
  test('adds numbers', () => {
    expect(1 + 1).toBe(2);
  });
});
`,
			expectedConfidence: 0,
			shouldMatch:        false,
		},
		{
			name: "cy as variable name (should not match)",
			content: `
const cy = { visit: () => {} };
describe('tests', () => {
  test('test case', () => {
    expect(true).toBe(true);
  });
});
`,
			expectedConfidence: 0,
			shouldMatch:        false,
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

			if tt.shouldMatch {
				assert.Equal(t, tt.expectedConfidence, result.Confidence, "Expected confidence mismatch")
				assert.NotEmpty(t, result.Evidence, "Expected evidence for match")
			} else {
				assert.Equal(t, 0, result.Confidence, "Should not match")
			}
		})
	}
}

func TestCypressContentMatcher_NonContentSignal(t *testing.T) {
	matcher := &CypressContentMatcher{}
	ctx := context.Background()

	signal := framework.Signal{
		Type:  framework.SignalImport,
		Value: "cypress",
	}

	result := matcher.Match(ctx, signal)
	assert.Equal(t, 0, result.Confidence)
}

func TestCypressConfigParser_Parse(t *testing.T) {
	tests := []struct {
		name                    string
		configContent           string
		configPath              string
		expectedGlobalsMode     bool
		expectedTestPatterns    []string
		expectedExcludePatterns []string
	}{
		{
			name: "default config without specPattern",
			configContent: `
const { defineConfig } = require('cypress');

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
  },
});
`,
			configPath:              "/project/cypress.config.js",
			expectedGlobalsMode:     true,
			expectedTestPatterns:    nil,
			expectedExcludePatterns: nil,
		},
		{
			name: "e2e specPattern single string",
			configContent: `
import { defineConfig } from 'cypress';

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
    specPattern: 'cypress/e2e/**/*.cy.ts',
  },
});
`,
			configPath:              "/project/cypress.config.ts",
			expectedGlobalsMode:     true,
			expectedTestPatterns:    []string{"cypress/e2e/**/*.cy.ts"},
			expectedExcludePatterns: nil,
		},
		{
			name: "e2e specPattern array",
			configContent: `
export default defineConfig({
  e2e: {
    specPattern: ['cypress/e2e/**/*.cy.ts', 'cypress/e2e/**/*.cy.js'],
  },
});
`,
			configPath:              "/project/cypress.config.ts",
			expectedGlobalsMode:     true,
			expectedTestPatterns:    []string{"cypress/e2e/**/*.cy.ts", "cypress/e2e/**/*.cy.js"},
			expectedExcludePatterns: nil,
		},
		{
			name: "component specPattern",
			configContent: `
export default defineConfig({
  component: {
    specPattern: 'src/**/*.cy.tsx',
  },
});
`,
			configPath:              "/project/cypress.config.ts",
			expectedGlobalsMode:     true,
			expectedTestPatterns:    []string{"src/**/*.cy.tsx"},
			expectedExcludePatterns: nil,
		},
		{
			name: "both e2e and component specPatterns",
			configContent: `
export default defineConfig({
  e2e: {
    specPattern: 'cypress/e2e/**/*.cy.ts',
  },
  component: {
    specPattern: 'src/**/*.cy.tsx',
  },
});
`,
			configPath:              "/project/cypress.config.ts",
			expectedGlobalsMode:     true,
			expectedTestPatterns:    []string{"cypress/e2e/**/*.cy.ts", "src/**/*.cy.tsx"},
			expectedExcludePatterns: nil,
		},
		{
			name: "excludeSpecPattern single",
			configContent: `
export default defineConfig({
  e2e: {
    specPattern: 'cypress/e2e/**/*.cy.ts',
    excludeSpecPattern: '**/ignore/**',
  },
});
`,
			configPath:              "/project/cypress.config.ts",
			expectedGlobalsMode:     true,
			expectedTestPatterns:    []string{"cypress/e2e/**/*.cy.ts"},
			expectedExcludePatterns: []string{"**/ignore/**"},
		},
		{
			name: "excludeSpecPattern array",
			configContent: `
export default defineConfig({
  e2e: {
    excludeSpecPattern: ['**/ignore/**', '**/skip/**'],
  },
});
`,
			configPath:              "/project/cypress.config.ts",
			expectedGlobalsMode:     true,
			expectedTestPatterns:    nil,
			expectedExcludePatterns: []string{"**/ignore/**", "**/skip/**"},
		},
		{
			name: "e2e with setupNodeEvents (nested braces)",
			configContent: `
export default defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      on('task', {
        log(message) {
          console.log(message);
          return null;
        },
      });
    },
    specPattern: 'cypress/e2e/**/*.cy.ts',
    baseUrl: 'http://localhost:3000',
  },
});
`,
			configPath:              "/project/cypress.config.ts",
			expectedGlobalsMode:     true,
			expectedTestPatterns:    []string{"cypress/e2e/**/*.cy.ts"},
			expectedExcludePatterns: nil,
		},
		{
			name: "component with devServer (nested braces)",
			configContent: `
export default defineConfig({
  component: {
    devServer: {
      framework: 'react',
      bundler: 'vite',
    },
    specPattern: 'src/**/*.cy.tsx',
  },
});
`,
			configPath:              "/project/cypress.config.ts",
			expectedGlobalsMode:     true,
			expectedTestPatterns:    []string{"src/**/*.cy.tsx"},
			expectedExcludePatterns: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &CypressConfigParser{}
			ctx := context.Background()

			scope, err := parser.Parse(ctx, tt.configPath, []byte(tt.configContent))

			require.NoError(t, err)
			assert.Equal(t, "cypress", scope.Framework)
			assert.Equal(t, tt.configPath, scope.ConfigPath)
			assert.Equal(t, tt.expectedGlobalsMode, scope.GlobalsMode)
			assert.Equal(t, tt.expectedTestPatterns, scope.TestPatterns)
			assert.Equal(t, tt.expectedExcludePatterns, scope.ExcludePatterns)
		})
	}
}

func TestParseSpecPattern(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		section  string
		expected []string
	}{
		{
			name:     "e2e single quote",
			content:  `e2e: { specPattern: 'cypress/e2e/**/*.cy.ts' }`,
			section:  "e2e",
			expected: []string{"cypress/e2e/**/*.cy.ts"},
		},
		{
			name:     "e2e double quote",
			content:  `e2e: { specPattern: "cypress/e2e/**/*.cy.ts" }`,
			section:  "e2e",
			expected: []string{"cypress/e2e/**/*.cy.ts"},
		},
		{
			name:     "e2e array",
			content:  `e2e: { specPattern: ['a.cy.ts', 'b.cy.ts'] }`,
			section:  "e2e",
			expected: []string{"a.cy.ts", "b.cy.ts"},
		},
		{
			name:     "component single",
			content:  `component: { specPattern: 'src/**/*.cy.tsx' }`,
			section:  "component",
			expected: []string{"src/**/*.cy.tsx"},
		},
		{
			name:     "no specPattern",
			content:  `e2e: { baseUrl: 'http://localhost' }`,
			section:  "e2e",
			expected: nil,
		},
		{
			name:     "invalid section",
			content:  `e2e: { specPattern: 'test' }`,
			section:  "invalid",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSpecPattern([]byte(tt.content), tt.section)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseExcludeSpecPattern(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "single pattern",
			content:  `excludeSpecPattern: '**/ignore/**'`,
			expected: []string{"**/ignore/**"},
		},
		{
			name:     "array pattern",
			content:  `excludeSpecPattern: ['**/a/**', '**/b/**']`,
			expected: []string{"**/a/**", "**/b/**"},
		},
		{
			name:     "no excludeSpecPattern",
			content:  `specPattern: 'test'`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseExcludeSpecPattern([]byte(tt.content))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCypressParser_Parse(t *testing.T) {
	testSource := `
describe('Login Page', () => {
  it('should display login form', () => {
    cy.visit('/login');
    cy.get('input[name="email"]').should('be.visible');
  });

  it.skip('skipped test', () => {
    cy.visit('/skip');
  });

  describe('Form Validation', () => {
    it('should show error for invalid email', () => {
      cy.get('#email').type('invalid');
      cy.get('.error').should('contain', 'Invalid email');
    });
  });
});

it('top-level test', () => {
  cy.visit('/home');
});
`

	parser := &CypressParser{}
	ctx := context.Background()

	testFile, err := parser.Parse(ctx, []byte(testSource), "login.cy.ts")

	require.NoError(t, err)
	assert.Equal(t, "login.cy.ts", testFile.Path)
	assert.Equal(t, "cypress", testFile.Framework)
	assert.Equal(t, domain.LanguageTypeScript, testFile.Language)

	// Verify suite structure
	require.Len(t, testFile.Suites, 1)
	suite := testFile.Suites[0]
	assert.Equal(t, "Login Page", suite.Name)
	assert.Len(t, suite.Tests, 2)

	// Verify tests within suite
	assert.Equal(t, "should display login form", suite.Tests[0].Name)
	assert.Equal(t, domain.TestStatusActive, suite.Tests[0].Status)

	assert.Equal(t, "skipped test", suite.Tests[1].Name)
	assert.Equal(t, domain.TestStatusSkipped, suite.Tests[1].Status)

	// Verify nested suite
	require.Len(t, suite.Suites, 1)
	nestedSuite := suite.Suites[0]
	assert.Equal(t, "Form Validation", nestedSuite.Name)
	assert.Len(t, nestedSuite.Tests, 1)

	// Verify top-level test
	require.Len(t, testFile.Tests, 1)
	assert.Equal(t, "top-level test", testFile.Tests[0].Name)
}

func TestCypressParser_ParseJavaScript(t *testing.T) {
	testSource := `
describe('Dashboard', () => {
  it('should load dashboard', () => {
    cy.visit('/dashboard');
  });
});
`

	parser := &CypressParser{}
	ctx := context.Background()

	testFile, err := parser.Parse(ctx, []byte(testSource), "dashboard.cy.js")

	require.NoError(t, err)
	assert.Equal(t, "dashboard.cy.js", testFile.Path)
	assert.Equal(t, "cypress", testFile.Framework)
	assert.Equal(t, domain.LanguageJavaScript, testFile.Language)
}
