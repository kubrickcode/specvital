package entity

import "time"

type RepositoryCursor struct {
	AnalyzedAt time.Time
	ID         string
	Name       string
	SortBy     SortBy
	TestCount  int
}
