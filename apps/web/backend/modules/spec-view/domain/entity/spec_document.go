package entity

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

// DefaultLanguage is the default language for spec generation
const DefaultLanguage = "English"

// MaxRepositoryNameLength is the maximum length for GitHub owner/repo names
const MaxRepositoryNameLength = 100

// repositoryNamePattern validates GitHub owner/repo names (alphanumeric, hyphen, underscore, dot)
var repositoryNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)

// IsValidRepositoryName checks if the given name is a valid GitHub owner or repo name
func IsValidRepositoryName(name string) bool {
	if len(name) == 0 || len(name) > MaxRepositoryNameLength {
		return false
	}
	return repositoryNamePattern.MatchString(name)
}

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

// IsValidDocumentID checks if the given document ID is a valid UUID format
func IsValidDocumentID(documentID string) bool {
	_, err := uuid.Parse(documentID)
	return err == nil
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
	CreatedAt       time.Time
	HasPreviousSpec bool
	Language        string
	LatestVersion   int
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

// RepoSpecDocument extends SpecDocument with repository context (commit SHA)
// Used for repository-based queries that provide cross-analysis version access
type RepoSpecDocument struct {
	AnalysisID         string
	AvailableLanguages []AvailableLanguageInfo
	CommitSHA          string
	CreatedAt          time.Time
	Domains            []SpecDomain
	ExecutiveSummary   *string
	ID                 string
	Language           string
	ModelID            string
	Version            int
}

// RepoVersionInfo extends VersionInfo with commit SHA for repository-based queries
type RepoVersionInfo struct {
	AnalysisID string
	CommitSHA  string
	CreatedAt  time.Time
	ID         string
	Language   string
	ModelID    string
	Version    int
}

// CachePredictionData contains cache prediction statistics
type CachePredictionData struct {
	// TotalBehaviors is the total number of behaviors from the previous spec
	TotalBehaviors int
	// CacheableBehaviors is the number of behaviors that can be cached
	// (matching test case name + file path exists in current analysis)
	CacheableBehaviors int
}
