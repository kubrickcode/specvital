package entity

import (
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateCacheKey(t *testing.T) {
	codebaseID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name           string
		codebaseID     uuid.UUID
		filePath       string
		suiteHierarchy string
		testName       string
		language       Language
		wantSameAs     *struct {
			filePath       string
			suiteHierarchy string
			testName       string
			language       Language
		}
	}{
		{
			name:           "basic input",
			codebaseID:     codebaseID,
			filePath:       "src/auth.spec.ts",
			suiteHierarchy: "AuthService > Login",
			testName:       "should create session",
			language:       LanguageEn,
		},
		{
			name:           "case insensitive - uppercase file path",
			codebaseID:     codebaseID,
			filePath:       "SRC/AUTH.SPEC.TS",
			suiteHierarchy: "AuthService > Login",
			testName:       "should create session",
			language:       LanguageEn,
			wantSameAs: &struct {
				filePath       string
				suiteHierarchy string
				testName       string
				language       Language
			}{
				filePath:       "src/auth.spec.ts",
				suiteHierarchy: "AuthService > Login",
				testName:       "should create session",
				language:       LanguageEn,
			},
		},
		{
			name:           "case insensitive - uppercase hierarchy",
			codebaseID:     codebaseID,
			filePath:       "src/auth.spec.ts",
			suiteHierarchy: "AUTHSERVICE > LOGIN",
			testName:       "should create session",
			language:       LanguageEn,
			wantSameAs: &struct {
				filePath       string
				suiteHierarchy string
				testName       string
				language       Language
			}{
				filePath:       "src/auth.spec.ts",
				suiteHierarchy: "AuthService > Login",
				testName:       "should create session",
				language:       LanguageEn,
			},
		},
		{
			name:           "hierarchy normalization - extra spaces",
			codebaseID:     codebaseID,
			filePath:       "src/auth.spec.ts",
			suiteHierarchy: "AuthService  >   Login",
			testName:       "should create session",
			language:       LanguageEn,
			wantSameAs: &struct {
				filePath       string
				suiteHierarchy string
				testName       string
				language       Language
			}{
				filePath:       "src/auth.spec.ts",
				suiteHierarchy: "AuthService > Login",
				testName:       "should create session",
				language:       LanguageEn,
			},
		},
		{
			name:           "test name trimming",
			codebaseID:     codebaseID,
			filePath:       "src/auth.spec.ts",
			suiteHierarchy: "AuthService > Login",
			testName:       "  should create session  ",
			language:       LanguageEn,
			wantSameAs: &struct {
				filePath       string
				suiteHierarchy string
				testName       string
				language       Language
			}{
				filePath:       "src/auth.spec.ts",
				suiteHierarchy: "AuthService > Login",
				testName:       "should create session",
				language:       LanguageEn,
			},
		},
		{
			name:           "different language produces different key",
			codebaseID:     codebaseID,
			filePath:       "src/auth.spec.ts",
			suiteHierarchy: "AuthService > Login",
			testName:       "should create session",
			language:       LanguageKo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateCacheKey(tt.codebaseID, tt.filePath, tt.suiteHierarchy, tt.testName, tt.language)

			if len(got) != 32 {
				t.Errorf("GenerateCacheKey() returned %d bytes, want 32 (SHA-256)", len(got))
			}

			if tt.wantSameAs != nil {
				want := GenerateCacheKey(
					tt.codebaseID,
					tt.wantSameAs.filePath,
					tt.wantSameAs.suiteHierarchy,
					tt.wantSameAs.testName,
					tt.wantSameAs.language,
				)
				gotHex := hex.EncodeToString(got)
				wantHex := hex.EncodeToString(want)
				if gotHex != wantHex {
					t.Errorf("GenerateCacheKey() = %s, want same as %s", gotHex, wantHex)
				}
			}
		})
	}
}

func TestGenerateCacheKey_DifferentInputsProduceDifferentKeys(t *testing.T) {
	codebaseID1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	codebaseID2 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	baseKey := GenerateCacheKey(codebaseID1, "src/auth.spec.ts", "AuthService > Login", "test", LanguageEn)

	tests := []struct {
		name           string
		codebaseID     uuid.UUID
		filePath       string
		suiteHierarchy string
		testName       string
		language       Language
	}{
		{
			name:           "different codebase ID",
			codebaseID:     codebaseID2,
			filePath:       "src/auth.spec.ts",
			suiteHierarchy: "AuthService > Login",
			testName:       "test",
			language:       LanguageEn,
		},
		{
			name:           "different file path",
			codebaseID:     codebaseID1,
			filePath:       "src/user.spec.ts",
			suiteHierarchy: "AuthService > Login",
			testName:       "test",
			language:       LanguageEn,
		},
		{
			name:           "different suite hierarchy",
			codebaseID:     codebaseID1,
			filePath:       "src/auth.spec.ts",
			suiteHierarchy: "UserService > Profile",
			testName:       "test",
			language:       LanguageEn,
		},
		{
			name:           "different test name",
			codebaseID:     codebaseID1,
			filePath:       "src/auth.spec.ts",
			suiteHierarchy: "AuthService > Login",
			testName:       "different test",
			language:       LanguageEn,
		},
		{
			name:           "different language",
			codebaseID:     codebaseID1,
			filePath:       "src/auth.spec.ts",
			suiteHierarchy: "AuthService > Login",
			testName:       "test",
			language:       LanguageKo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateCacheKey(tt.codebaseID, tt.filePath, tt.suiteHierarchy, tt.testName, tt.language)
			gotHex := hex.EncodeToString(got)
			baseHex := hex.EncodeToString(baseKey)

			if gotHex == baseHex {
				t.Errorf("GenerateCacheKey() produced same key as base case, expected different")
			}
		})
	}
}

func TestLanguage_IsValid(t *testing.T) {
	tests := []struct {
		language Language
		want     bool
	}{
		{LanguageAr, true},
		{LanguageCs, true},
		{LanguageDa, true},
		{LanguageDe, true},
		{LanguageEl, true},
		{LanguageEn, true},
		{LanguageEs, true},
		{LanguageFi, true},
		{LanguageFr, true},
		{LanguageHi, true},
		{LanguageId, true},
		{LanguageIt, true},
		{LanguageJa, true},
		{LanguageKo, true},
		{LanguageNl, true},
		{LanguagePl, true},
		{LanguagePt, true},
		{LanguageRu, true},
		{LanguageSv, true},
		{LanguageTh, true},
		{LanguageTr, true},
		{LanguageUk, true},
		{LanguageVi, true},
		{LanguageZh, true},
		{Language("invalid"), false},
		{Language(""), false},
		{Language("EN"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.language), func(t *testing.T) {
			if got := tt.language.IsValid(); got != tt.want {
				t.Errorf("Language(%q).IsValid() = %v, want %v", tt.language, got, tt.want)
			}
		})
	}
}

func TestLanguage_String(t *testing.T) {
	if got := LanguageEn.String(); got != "en" {
		t.Errorf("LanguageEn.String() = %q, want %q", got, "en")
	}
	if got := LanguageKo.String(); got != "ko" {
		t.Errorf("LanguageKo.String() = %q, want %q", got, "ko")
	}
}
