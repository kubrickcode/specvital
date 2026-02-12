package rspec

import (
	"context"
	"testing"

	"github.com/kubrickcode/specvital/lib/parser/domain"
	"github.com/kubrickcode/specvital/lib/parser/framework"
)

func TestRSpecParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		wantSuites int
		wantTests  int
		checkFunc  func(t *testing.T, file *domain.TestFile)
	}{
		{
			name: "basic describe and it blocks",
			source: `
RSpec.describe User do
  it "creates a user" do
    expect(User.new).to be_valid
  end
end
`,
			wantSuites: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				if len(file.Suites) != 1 {
					t.Errorf("expected 1 suite, got %d", len(file.Suites))
					return
				}
				suite := file.Suites[0]
				if suite.Name != "User" {
					t.Errorf("expected suite name 'User', got %q", suite.Name)
				}
				if len(suite.Tests) != 1 {
					t.Errorf("expected 1 test, got %d", len(suite.Tests))
					return
				}
				if suite.Tests[0].Name != "creates a user" {
					t.Errorf("expected test name 'creates a user', got %q", suite.Tests[0].Name)
				}
			},
		},
		{
			name: "nested describe and context",
			source: `
RSpec.describe Calculator do
  describe "#add" do
    context "with positive numbers" do
      it "returns the sum" do
        expect(Calculator.add(1, 2)).to eq(3)
      end
    end
  end
end
`,
			wantSuites: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				if len(file.Suites) != 1 {
					t.Errorf("expected 1 top-level suite, got %d", len(file.Suites))
					return
				}
				suite := file.Suites[0]
				if suite.Name != "Calculator" {
					t.Errorf("expected suite name 'Calculator', got %q", suite.Name)
				}
				if len(suite.Suites) != 1 {
					t.Errorf("expected 1 nested suite, got %d", len(suite.Suites))
					return
				}
				addSuite := suite.Suites[0]
				if addSuite.Name != "#add" {
					t.Errorf("expected nested suite name '#add', got %q", addSuite.Name)
				}
				if len(addSuite.Suites) != 1 {
					t.Errorf("expected 1 context suite, got %d", len(addSuite.Suites))
					return
				}
				contextSuite := addSuite.Suites[0]
				if contextSuite.Name != "with positive numbers" {
					t.Errorf("expected context name 'with positive numbers', got %q", contextSuite.Name)
				}
				if len(contextSuite.Tests) != 1 {
					t.Errorf("expected 1 test in context, got %d", len(contextSuite.Tests))
				}
			},
		},
		{
			name: "skipped tests with xit",
			source: `
RSpec.describe User do
  xit "is skipped" do
    expect(true).to be true
  end

  it "runs normally" do
    expect(true).to be true
  end
end
`,
			wantSuites: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				suite := file.Suites[0]
				if len(suite.Tests) != 2 {
					t.Errorf("expected 2 tests, got %d", len(suite.Tests))
					return
				}

				skipped := suite.Tests[0]
				if skipped.Name != "is skipped" {
					t.Errorf("expected test name 'is skipped', got %q", skipped.Name)
				}
				if skipped.Status != domain.TestStatusSkipped {
					t.Errorf("expected status skipped, got %q", skipped.Status)
				}

				normal := suite.Tests[1]
				if normal.Status != domain.TestStatusActive {
					t.Errorf("expected active status, got %q", normal.Status)
				}
			},
		},
		{
			name: "skipped suite with xdescribe",
			source: `
xdescribe "skipped suite" do
  it "is in skipped suite" do
    expect(true).to be true
  end
end
`,
			wantSuites: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				suite := file.Suites[0]
				if suite.Name != "skipped suite" {
					t.Errorf("expected suite name 'skipped suite', got %q", suite.Name)
				}
				if suite.Status != domain.TestStatusSkipped {
					t.Errorf("expected suite status skipped, got %q", suite.Status)
				}
			},
		},
		{
			name: "specify and example keywords",
			source: `
RSpec.describe User do
  specify "user is valid" do
    expect(User.new).to be_valid
  end

  example "another test" do
    expect(true).to be true
  end
end
`,
			wantSuites: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				suite := file.Suites[0]
				if len(suite.Tests) != 2 {
					t.Errorf("expected 2 tests, got %d", len(suite.Tests))
					return
				}
				if suite.Tests[0].Name != "user is valid" {
					t.Errorf("expected test name 'user is valid', got %q", suite.Tests[0].Name)
				}
				if suite.Tests[1].Name != "another test" {
					t.Errorf("expected test name 'another test', got %q", suite.Tests[1].Name)
				}
			},
		},
		{
			name: "string with double quotes",
			source: `
RSpec.describe "String Utils" do
  it "handles strings" do
    expect("hello").to eq("hello")
  end
end
`,
			wantSuites: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				suite := file.Suites[0]
				if suite.Name != "String Utils" {
					t.Errorf("expected suite name 'String Utils', got %q", suite.Name)
				}
			},
		},
		{
			name: "multiple tests in one suite",
			source: `
RSpec.describe Array do
  it "starts empty" do
    expect([]).to be_empty
  end

  it "can have elements" do
    expect([1, 2]).not_to be_empty
  end

  it "has a length" do
    expect([1, 2, 3].length).to eq(3)
  end
end
`,
			wantSuites: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				suite := file.Suites[0]
				if len(suite.Tests) != 3 {
					t.Errorf("expected 3 tests, got %d", len(suite.Tests))
				}
			},
		},
	}

	parser := &RSpecParser{}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.Parse(ctx, []byte(tt.source), "test_spec.rb")
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if file.Framework != frameworkName {
				t.Errorf("expected framework %q, got %q", frameworkName, file.Framework)
			}

			if file.Language != domain.LanguageRuby {
				t.Errorf("expected language Ruby, got %q", file.Language)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, file)
			}
		})
	}
}

func TestRSpecFileMatcher_Match(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		wantMatch bool
	}{
		{"spec file", "user_spec.rb", true},
		{"spec in directory", "spec/models/user_spec.rb", true},
		{"regular ruby file", "user.rb", false},
		{"test file not spec", "user_test.rb", false},
		{"file in spec directory", "spec/support/helpers.rb", false},
	}

	matcher := &RSpecFileMatcher{}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:  framework.SignalFileName,
				Value: tt.filename,
			}
			result := matcher.Match(ctx, signal)
			if (result.Confidence > 0) != tt.wantMatch {
				t.Errorf("Match() = %v, want match = %v", result.Confidence > 0, tt.wantMatch)
			}
		})
	}
}

func TestRSpecContentMatcher_Match(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantMatch bool
	}{
		{"RSpec.describe", "RSpec.describe User do", true},
		{"describe block", `describe "something" do`, true},
		{"context block", `context "when valid" do`, true},
		{"it block", `it "works" do`, true},
		{"specify block", `specify "something" do`, true},
		{"expect assertion", "expect(user).to be_valid", true},
		{"let definition", "let(:user) { User.new }", true},
		{"before hook", "before(:each) do", true},
		{"plain ruby code", "class User; end", false},
	}

	matcher := &RSpecContentMatcher{}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:    framework.SignalFileContent,
				Value:   tt.content,
				Context: []byte(tt.content),
			}
			result := matcher.Match(ctx, signal)
			if (result.Confidence > 0) != tt.wantMatch {
				t.Errorf("Match() = %v, want match = %v", result.Confidence > 0, tt.wantMatch)
			}
		})
	}
}

func TestRSpecParser_LoopBlockDetection(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		wantTests int
		checkFunc func(t *testing.T, file *domain.TestFile)
	}{
		{
			name: "times loop with it blocks",
			source: `
RSpec.describe "Group" do
  3.times do |i|
    it "test #{i}" do
    end
  end
end
`,
			wantTests: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				if len(file.Suites) != 1 {
					t.Fatalf("expected 1 suite, got %d", len(file.Suites))
				}
				suite := file.Suites[0]
				if len(suite.Tests) != 1 {
					t.Errorf("expected 1 test inside times loop, got %d", len(suite.Tests))
				}
			},
		},
		{
			name: "each loop with it blocks",
			source: `
RSpec.describe "Array tests" do
  [1, 2, 3].each do |n|
    it "handles #{n}" do
      expect(n).to be > 0
    end
  end
end
`,
			wantTests: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				if len(file.Suites) != 1 {
					t.Fatalf("expected 1 suite, got %d", len(file.Suites))
				}
				suite := file.Suites[0]
				if len(suite.Tests) != 1 {
					t.Errorf("expected 1 test inside each loop, got %d", len(suite.Tests))
				}
			},
		},
		{
			name: "loop in nested describe",
			source: `
RSpec.describe "Outer" do
  describe "Inner" do
    2.times do |i|
      it "nested test #{i}" do
      end
    end
  end
end
`,
			wantTests: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				if len(file.Suites) != 1 {
					t.Fatalf("expected 1 top-level suite, got %d", len(file.Suites))
				}
				outer := file.Suites[0]
				if len(outer.Suites) != 1 {
					t.Fatalf("expected 1 nested suite, got %d", len(outer.Suites))
				}
				inner := outer.Suites[0]
				if len(inner.Tests) != 1 {
					t.Errorf("expected 1 test in nested suite loop, got %d", len(inner.Tests))
				}
			},
		},
		{
			name: "multiple tests in loop",
			source: `
RSpec.describe "Multiple" do
  2.times do |i|
    it "passes #{i}" do
    end
    it "fails #{i}" do
    end
  end
end
`,
			wantTests: 2,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				if len(file.Suites) != 1 {
					t.Fatalf("expected 1 suite, got %d", len(file.Suites))
				}
				suite := file.Suites[0]
				if len(suite.Tests) != 2 {
					t.Errorf("expected 2 tests inside loop, got %d", len(suite.Tests))
				}
			},
		},
		{
			name: "upto loop with it blocks",
			source: `
RSpec.describe "Upto" do
  1.upto(3) do |n|
    it "test #{n}" do
    end
  end
end
`,
			wantTests: 1,
			checkFunc: func(t *testing.T, file *domain.TestFile) {
				if len(file.Suites) != 1 {
					t.Fatalf("expected 1 suite, got %d", len(file.Suites))
				}
				suite := file.Suites[0]
				if len(suite.Tests) != 1 {
					t.Errorf("expected 1 test inside upto loop, got %d", len(suite.Tests))
				}
			},
		},
	}

	parser := &RSpecParser{}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.Parse(ctx, []byte(tt.source), "loop_spec.rb")
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, file)
			}
		})
	}
}
