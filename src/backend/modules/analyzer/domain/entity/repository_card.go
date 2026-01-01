package entity

import "time"

type RepositoryCard struct {
	FullName       string
	ID             string
	IsAnalyzedByMe bool
	IsBookmarked   bool
	LatestAnalysis *AnalysisSummary
	Name           string
	Owner          string
	UpdateStatus   UpdateStatus
}

type AnalysisSummary struct {
	AnalyzedAt  time.Time
	Change      int
	CommitSHA   string
	TestCount   int
	TestSummary *TestStatusSummary
}

type TestStatusSummary struct {
	Active  int
	Focused int
	Skipped int
	Todo    int
	Xfail   int
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
