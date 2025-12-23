package domain

import "time"

type BookmarkedRepository struct {
	AnalyzedAt   *time.Time
	BookmarkedAt time.Time
	Change       int
	CodebaseID   string
	CommitSHA    string
	Name         string
	Owner        string
	TotalTests   int
}
