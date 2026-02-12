package entity

import "time"

type Repository struct {
	Archived      bool
	DefaultBranch string
	Description   string
	Disabled      bool
	Fork          bool
	FullName      string
	HTMLURL       string
	ID            int64
	Language      string
	Name          string
	Owner         string
	Private       bool
	PushedAt      *time.Time
	StarCount     int
	Visibility    string
}
