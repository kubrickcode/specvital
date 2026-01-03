package entity

import "time"

type OwnershipFilter string

const (
	OwnershipAll          OwnershipFilter = "all"
	OwnershipMine         OwnershipFilter = "mine"
	OwnershipOrganization OwnershipFilter = "organization"
	OwnershipOthers       OwnershipFilter = "others"
)

func ParseOwnershipFilter(s string) OwnershipFilter {
	switch s {
	case "mine":
		return OwnershipMine
	case "organization":
		return OwnershipOrganization
	case "others":
		return OwnershipOthers
	default:
		return OwnershipAll
	}
}

type AnalyzedRepository struct {
	CodebaseID  string
	CommitSHA   string
	CompletedAt time.Time
	HistoryID   string
	Name        string
	Owner       string
	TotalTests  int
	UpdatedAt   time.Time
}

type AnalyzedReposResult struct {
	Data       []*AnalyzedRepository
	HasNext    bool
	NextCursor *string
}
