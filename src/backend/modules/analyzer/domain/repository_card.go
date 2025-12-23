package domain

import "time"

type UpdateStatus string

const (
	UpdateStatusNewCommits UpdateStatus = "new-commits"
	UpdateStatusUnknown    UpdateStatus = "unknown"
	UpdateStatusUpToDate   UpdateStatus = "up-to-date"
)

func (s UpdateStatus) String() string {
	return string(s)
}

type RepositoryCard struct {
	FullName       string
	ID             string
	IsBookmarked   bool
	LatestAnalysis *AnalysisSummary
	Name           string
	Owner          string
	UpdateStatus   UpdateStatus
}

type AnalysisSummary struct {
	AnalyzedAt time.Time
	Change     int
	CommitSHA  string
	TestCount  int
}

type RepositoryStats struct {
	TotalRepositories int
	TotalTests        int
}

type UpdateStatusResult struct {
	AnalyzedCommitSHA string
	LatestCommitSHA   string
	Status            UpdateStatus
}
