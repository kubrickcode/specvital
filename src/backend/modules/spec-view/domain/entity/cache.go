package entity

import (
	"crypto/sha256"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CacheEntry struct {
	CacheKeyHash   []byte
	CodebaseID     uuid.UUID
	ConvertedName  string
	CreatedAt      time.Time
	FilePath       string
	Framework      string
	ID             uuid.UUID
	Language       Language
	ModelID        string
	OriginalName   string
	SuiteHierarchy string
}

func GenerateCacheKey(codebaseID uuid.UUID, filePath, suiteHierarchy, testName string, language Language) []byte {
	key := strings.Join([]string{
		codebaseID.String(),
		strings.ToLower(filePath),
		strings.ToLower(normalizeSuiteHierarchy(suiteHierarchy)),
		strings.ToLower(strings.TrimSpace(testName)),
		language.String(),
	}, ":")

	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func normalizeSuiteHierarchy(hierarchy string) string {
	parts := strings.Split(hierarchy, ">")
	normalized := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}

	return strings.Join(normalized, " > ")
}
