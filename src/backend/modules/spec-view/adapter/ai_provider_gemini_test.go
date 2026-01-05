package adapter

import (
	"testing"

	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		name     string
		input    port.ConvertInput
		contains []string
	}{
		{
			name: "basic prompt structure",
			input: port.ConvertInput{
				FilePath: "src/utils/auth.spec.ts",
				Language: entity.LanguageEn,
				Suites: []port.SuiteInput{
					{
						Hierarchy: "Auth > Login",
						Tests:     []string{"should create session", "should handle error"},
					},
				},
			},
			contains: []string{
				"English",
				"Auth > Login",
				"1|should create session",
				"2|should handle error",
				"<rules>",
				"<output_format>",
			},
		},
		{
			name: "korean language",
			input: port.ConvertInput{
				FilePath: "test.spec.ts",
				Language: entity.LanguageKo,
				Suites: []port.SuiteInput{
					{
						Hierarchy: "TestSuite",
						Tests:     []string{"test case"},
					},
				},
			},
			contains: []string{
				"Korean",
			},
		},
		{
			name: "multiple suites",
			input: port.ConvertInput{
				FilePath: "multi.spec.ts",
				Language: entity.LanguageJa,
				Suites: []port.SuiteInput{
					{
						Hierarchy: "Suite A",
						Tests:     []string{"test 1", "test 2"},
					},
					{
						Hierarchy: "Suite B",
						Tests:     []string{"test 3"},
					},
				},
			},
			contains: []string{
				"Japanese",
				"Suite A",
				"Suite B",
				"1|test 1",
				"2|test 2",
				"3|test 3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := buildPrompt(tt.input)

			for _, s := range tt.contains {
				if !contains(prompt, s) {
					t.Errorf("expected prompt to contain %q, but it didn't\nPrompt: %s", s, prompt)
				}
			}
		})
	}
}

func TestLanguageDisplayName(t *testing.T) {
	tests := []struct {
		lang     entity.Language
		expected string
	}{
		{entity.LanguageEn, "English"},
		{entity.LanguageKo, "Korean"},
		{entity.LanguageJa, "Japanese"},
		{entity.LanguageZh, "Chinese"},
		{entity.LanguageFr, "French"},
		{entity.LanguageDe, "German"},
		{entity.LanguageEs, "Spanish"},
		{entity.LanguagePt, "Portuguese"},
		{entity.LanguageRu, "Russian"},
		{entity.LanguageAr, "Arabic"},
		{entity.Language("unknown"), "English"},
	}

	for _, tt := range tests {
		t.Run(string(tt.lang), func(t *testing.T) {
			result := languageDisplayName(tt.lang)
			if result != tt.expected {
				t.Errorf("expected %s for language %s, got %s", tt.expected, tt.lang, result)
			}
		})
	}
}

func TestNewGeminiProvider_MissingAPIKey(t *testing.T) {
	_, err := NewGeminiProvider(t.Context(), GeminiConfig{})

	if err == nil {
		t.Error("expected error when API key is missing")
	}

	if !contains(err.Error(), "API key is required") {
		t.Errorf("expected error message about API key, got: %v", err)
	}
}

func TestGeminiProvider_ModelID(t *testing.T) {
	tests := []struct {
		name     string
		modelID  string
		expected string
	}{
		{
			name:     "default model",
			modelID:  "",
			expected: "gemini-2.5-flash-lite",
		},
		{
			name:     "custom model",
			modelID:  "gemini-pro",
			expected: "gemini-pro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &GeminiProvider{
				modelID: tt.modelID,
			}
			if tt.modelID == "" {
				provider.modelID = defaultModelID
			}

			if provider.ModelID() != tt.expected {
				t.Errorf("expected model ID %s, got %s", tt.expected, provider.ModelID())
			}
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected bool
	}{
		{
			name:     "nil error",
			errMsg:   "",
			expected: false,
		},
		{
			name:     "rate limit error",
			errMsg:   "rate limit exceeded",
			expected: true,
		},
		{
			name:     "timeout error",
			errMsg:   "request timeout",
			expected: true,
		},
		{
			name:     "temporarily unavailable",
			errMsg:   "service temporarily unavailable",
			expected: true,
		},
		{
			name:     "non-retryable error",
			errMsg:   "invalid request format",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.errMsg != "" {
				err = &testError{msg: tt.errMsg}
			}

			result := isRetryableError(err)
			if result != tt.expected {
				t.Errorf("expected %v for error %q, got %v", tt.expected, tt.errMsg, result)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchSubstring(s, substr)))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
