package entity

import "time"

// SupportedLanguages defines all valid languages for spec generation (synced with OpenAPI SpecLanguage enum)
var SupportedLanguages = map[string]bool{
	"Arabic":     true,
	"Chinese":    true,
	"Czech":      true,
	"Danish":     true,
	"Dutch":      true,
	"English":    true,
	"Finnish":    true,
	"French":     true,
	"German":     true,
	"Greek":      true,
	"Hindi":      true,
	"Indonesian": true,
	"Italian":    true,
	"Japanese":   true,
	"Korean":     true,
	"Polish":     true,
	"Portuguese": true,
	"Russian":    true,
	"Spanish":    true,
	"Swedish":    true,
	"Thai":       true,
	"Turkish":    true,
	"Ukrainian":  true,
	"Vietnamese": true,
}

// IsValidLanguage checks if the given language is supported
func IsValidLanguage(language string) bool {
	return SupportedLanguages[language]
}

type SpecDocument struct {
	AnalysisID         string
	AvailableLanguages []AvailableLanguageInfo
	CreatedAt          time.Time
	Domains            []SpecDomain
	ExecutiveSummary   *string
	ID                 string
	Language           string
	ModelID            string
	Version            int
}

type AvailableLanguageInfo struct {
	CreatedAt     time.Time
	Language      string
	LatestVersion int
}

type VersionInfo struct {
	CreatedAt time.Time
	ModelID   string
	Version   int
}

type SpecDomain struct {
	ClassificationConfidence *float64
	Description              *string
	Features                 []SpecFeature
	ID                       string
	Name                     string
	SortOrder                int
}

type SpecFeature struct {
	Behaviors   []SpecBehavior
	Description *string
	ID          string
	Name        string
	SortOrder   int
}

type SpecBehavior struct {
	ConvertedDescription string
	ID                   string
	OriginalName         string
	SortOrder            int
	SourceInfo           *BehaviorSourceInfo
	SourceTestCaseID     *string
}

type BehaviorSourceInfo struct {
	FilePath   string
	Framework  string
	LineNumber int
	Status     string
}
