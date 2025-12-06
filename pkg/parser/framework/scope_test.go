package framework

import (
	"testing"
)

func TestConfigScope_Contains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		scope    *ConfigScope
		filePath string
		want     bool
	}{
		{
			name: "should match file in base directory",
			scope: &ConfigScope{
				ConfigPath: "src/extension/vitest.config.ts",
				BaseDir:    "src",
			},
			filePath: "src/internal/adapters.spec.ts",
			want:     true,
		},
		{
			name: "should match file in same directory as config",
			scope: &ConfigScope{
				ConfigPath: "src/extension/vitest.config.ts",
				BaseDir:    "src/extension",
			},
			filePath: "src/extension/test.spec.ts",
			want:     true,
		},
		{
			name: "should not match file outside base directory",
			scope: &ConfigScope{
				ConfigPath: "src/extension/vitest.config.ts",
				BaseDir:    "src/extension",
			},
			filePath: "other/test.spec.ts",
			want:     false,
		},
		{
			name: "should not match file in parent directory",
			scope: &ConfigScope{
				ConfigPath: "src/extension/vitest.config.ts",
				BaseDir:    "src/extension",
			},
			filePath: "src/test.spec.ts",
			want:     false,
		},
		{
			name: "should match with include patterns - spec files",
			scope: &ConfigScope{
				BaseDir: "src",
				Include: []string{"**/*.spec.ts", "**/*.test.ts"},
			},
			filePath: "src/foo/bar.spec.ts",
			want:     true,
		},
		{
			name: "should match with include patterns - test files",
			scope: &ConfigScope{
				BaseDir: "src",
				Include: []string{"**/*.spec.ts", "**/*.test.ts"},
			},
			filePath: "src/foo/bar.test.ts",
			want:     true,
		},
		{
			name: "should not match with include patterns - non-matching",
			scope: &ConfigScope{
				BaseDir: "src",
				Include: []string{"**/*.spec.ts", "**/*.test.ts"},
			},
			filePath: "src/foo/bar.ts",
			want:     false,
		},
		{
			name: "should exclude node_modules",
			scope: &ConfigScope{
				BaseDir: "src",
				Exclude: []string{"**/node_modules/**"},
			},
			filePath: "src/node_modules/foo/test.spec.ts",
			want:     false,
		},
		{
			name: "should exclude __mocks__",
			scope: &ConfigScope{
				BaseDir: "src",
				Exclude: []string{"**/__mocks__/**"},
			},
			filePath: "src/__mocks__/adapter.ts",
			want:     false,
		},
		{
			name: "should exclude multiple patterns",
			scope: &ConfigScope{
				BaseDir: "src",
				Exclude: []string{"**/node_modules/**", "**/__mocks__/**", "**/dist/**"},
			},
			filePath: "src/dist/bundle.js",
			want:     false,
		},
		{
			name: "should match with include and exclude combined",
			scope: &ConfigScope{
				BaseDir: "src",
				Include: []string{"**/*.spec.ts"},
				Exclude: []string{"**/__mocks__/**"},
			},
			filePath: "src/foo/bar.spec.ts",
			want:     true,
		},
		{
			name: "should exclude even if matches include pattern",
			scope: &ConfigScope{
				BaseDir: "src",
				Include: []string{"**/*.spec.ts"},
				Exclude: []string{"**/__mocks__/**"},
			},
			filePath: "src/__mocks__/foo.spec.ts",
			want:     false,
		},
		{
			name: "should handle root directory",
			scope: &ConfigScope{
				ConfigPath: "vitest.config.ts",
				BaseDir:    ".",
			},
			filePath: "src/test.spec.ts",
			want:     true,
		},
		{
			name: "should handle nested subdirectories",
			scope: &ConfigScope{
				BaseDir: "src/features",
			},
			filePath: "src/features/auth/login/test.spec.ts",
			want:     true,
		},
		{
			name:     "should return false for nil scope",
			scope:    nil,
			filePath: "src/test.spec.ts",
			want:     false,
		},
		{
			name: "should handle windows-style backslashes in ConfigPath",
			scope: &ConfigScope{
				ConfigPath: "src\\extension\\vitest.config.ts",
				BaseDir:    "src",
			},
			filePath: "src/internal/adapters.spec.ts",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.scope.Contains(tt.filePath)
			if got != tt.want {
				t.Errorf("ConfigScope.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigScope_Depth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		scope *ConfigScope
		want  int
	}{
		{
			name: "should return 0 for root directory",
			scope: &ConfigScope{
				BaseDir: ".",
			},
			want: 0,
		},
		{
			name: "should return 1 for one-level directory",
			scope: &ConfigScope{
				BaseDir: "src",
			},
			want: 0, // "src" has no slashes
		},
		{
			name: "should return 2 for two-level directory",
			scope: &ConfigScope{
				BaseDir: "src/extension",
			},
			want: 1, // one slash
		},
		{
			name: "should return 3 for three-level directory",
			scope: &ConfigScope{
				BaseDir: "src/features/auth",
			},
			want: 2, // two slashes
		},
		{
			name: "should return 0 for nil scope",
			scope: &ConfigScope{
				BaseDir: "",
			},
			want: 0,
		},
		{
			name:  "should return 0 for completely nil scope",
			scope: nil,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.scope.Depth()
			if got != tt.want {
				t.Errorf("ConfigScope.Depth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigScope_FindMatchingProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		scope    *ConfigScope
		filePath string
		want     *ProjectScope
	}{
		{
			name: "should find matching project",
			scope: &ConfigScope{
				Projects: []ProjectScope{
					{
						Name:    "backend",
						BaseDir: "packages/backend",
					},
					{
						Name:    "frontend",
						BaseDir: "packages/frontend",
					},
				},
			},
			filePath: "packages/backend/src/test.spec.ts",
			want: &ProjectScope{
				Name:    "backend",
				BaseDir: "packages/backend",
			},
		},
		{
			name: "should find most specific project",
			scope: &ConfigScope{
				Projects: []ProjectScope{
					{
						Name:    "root",
						BaseDir: "src",
					},
					{
						Name:    "nested",
						BaseDir: "src/features",
					},
				},
			},
			filePath: "src/features/auth/test.spec.ts",
			want: &ProjectScope{
				Name:    "nested",
				BaseDir: "src/features",
			},
		},
		{
			name: "should return nil for no matching project",
			scope: &ConfigScope{
				Projects: []ProjectScope{
					{
						Name:    "backend",
						BaseDir: "packages/backend",
					},
				},
			},
			filePath: "packages/frontend/src/test.spec.ts",
			want:     nil,
		},
		{
			name: "should respect include patterns",
			scope: &ConfigScope{
				Projects: []ProjectScope{
					{
						Name:    "tests",
						BaseDir: "src",
						Include: []string{"**/*.test.ts"},
					},
				},
			},
			filePath: "src/foo/bar.test.ts",
			want: &ProjectScope{
				Name:    "tests",
				BaseDir: "src",
				Include: []string{"**/*.test.ts"},
			},
		},
		{
			name: "should exclude by patterns",
			scope: &ConfigScope{
				Projects: []ProjectScope{
					{
						Name:    "src",
						BaseDir: "src",
						Exclude: []string{"**/__mocks__/**"},
					},
				},
			},
			filePath: "src/__mocks__/test.ts",
			want:     nil,
		},
		{
			name: "should return nil for nil scope",
			scope: &ConfigScope{
				Projects: nil,
			},
			filePath: "src/test.ts",
			want:     nil,
		},
		{
			name:     "should return nil for empty projects",
			scope:    &ConfigScope{},
			filePath: "src/test.ts",
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.scope.FindMatchingProject(tt.filePath)

			if tt.want == nil {
				if got != nil {
					t.Errorf("ConfigScope.FindMatchingProject() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Errorf("ConfigScope.FindMatchingProject() = nil, want %v", tt.want)
				return
			}

			if got.Name != tt.want.Name || got.BaseDir != tt.want.BaseDir {
				t.Errorf("ConfigScope.FindMatchingProject() = {Name: %v, BaseDir: %v}, want {Name: %v, BaseDir: %v}",
					got.Name, got.BaseDir, tt.want.Name, tt.want.BaseDir)
			}
		})
	}
}

func TestNewConfigScope(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		configPath string
		root       string
		wantBase   string
	}{
		{
			name:       "should use config directory when no root specified",
			configPath: "src/extension/vitest.config.ts",
			root:       "",
			wantBase:   "src/extension",
		},
		{
			name:       "should resolve root relative to config directory",
			configPath: "src/extension/vitest.config.ts",
			root:       "..",
			wantBase:   "src",
		},
		{
			name:       "should handle root with multiple levels",
			configPath: "src/features/auth/config.ts",
			root:       "../..",
			wantBase:   "src",
		},
		{
			name:       "should handle absolute-like root",
			configPath: "src/extension/vitest.config.ts",
			root:       "../../",
			wantBase:   ".",
		},
		{
			name:       "should handle subdirectory root",
			configPath: "vitest.config.ts",
			root:       "src",
			wantBase:   "src",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scope := NewConfigScope(tt.configPath, tt.root)

			if scope.ConfigPath != tt.configPath {
				t.Errorf("NewConfigScope().ConfigPath = %v, want %v", scope.ConfigPath, tt.configPath)
			}

			if scope.BaseDir != tt.wantBase {
				t.Errorf("NewConfigScope().BaseDir = %v, want %v", scope.BaseDir, tt.wantBase)
			}

			if scope.Settings == nil {
				t.Error("NewConfigScope().Settings should not be nil")
			}
		})
	}
}

// Integration test: Real-world scenario from the bug report
func TestConfigScope_RealWorldScenario(t *testing.T) {
	t.Parallel()

	// Config at src/extension/vitest.config.ts with root: ".."
	scope := NewConfigScope("src/extension/vitest.config.ts", "..")

	testCases := []struct {
		filePath string
		want     bool
	}{
		// Should match: files in src/ (parent of src/extension/)
		{"src/internal/adapters.spec.ts", true},
		{"src/extension/test.spec.ts", true},
		{"src/other/module.spec.ts", true},

		// Should not match: files outside src/
		{"other/test.spec.ts", false},
		{"test.spec.ts", false},
		{"packages/foo/test.spec.ts", false},
	}

	for _, tc := range testCases {
		t.Run(tc.filePath, func(t *testing.T) {
			got := scope.Contains(tc.filePath)
			if got != tc.want {
				t.Errorf("scope.Contains(%q) = %v, want %v (BaseDir=%q)",
					tc.filePath, got, tc.want, scope.BaseDir)
			}
		})
	}
}
