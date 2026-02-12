package entity

import "testing"

func TestIsValidDocumentID(t *testing.T) {
	tests := []struct {
		name       string
		documentID string
		want       bool
	}{
		{"valid UUID v4", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid UUID v1", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},
		{"valid UUID uppercase", "550E8400-E29B-41D4-A716-446655440000", true},
		{"valid UUID no hyphens", "550e8400e29b41d4a716446655440000", true}, // uuid.Parse accepts this
		{"empty string", "", false},
		{"invalid format - too short", "550e8400-e29b-41d4-a716", false},
		{"invalid format - random string", "invalid-uuid-format", false},
		{"invalid format - wrong characters", "550e8400-e29b-41d4-a716-44665544000g", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDocumentID(tt.documentID); got != tt.want {
				t.Errorf("IsValidDocumentID(%q) = %v, want %v", tt.documentID, got, tt.want)
			}
		})
	}
}

func TestIsValidLanguage(t *testing.T) {
	tests := []struct {
		name     string
		language string
		want     bool
	}{
		{"English", "English", true},
		{"Korean", "Korean", true},
		{"Japanese", "Japanese", true},
		{"invalid language", "InvalidLang", false},
		{"empty string", "", false},
		{"lowercase english", "english", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidLanguage(tt.language); got != tt.want {
				t.Errorf("IsValidLanguage(%q) = %v, want %v", tt.language, got, tt.want)
			}
		})
	}
}

func TestIsValidRepositoryName(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		want     bool
	}{
		{"valid simple name", "react", true},
		{"valid with hyphen", "my-repo", true},
		{"valid with underscore", "my_repo", true},
		{"valid with dot", "my.repo", true},
		{"valid with numbers", "react123", true},
		{"empty string", "", false},
		{"starts with hyphen", "-react", false},
		{"starts with dot", ".react", false},
		{"contains slash", "my/repo", false},
		{"contains path traversal", "../etc", false},
		{"too long", string(make([]byte, 101)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidRepositoryName(tt.repoName); got != tt.want {
				t.Errorf("IsValidRepositoryName(%q) = %v, want %v", tt.repoName, got, tt.want)
			}
		})
	}
}
