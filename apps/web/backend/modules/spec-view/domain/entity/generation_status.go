package entity

import "time"

type GenerationStatus string

const (
	StatusPending   GenerationStatus = "pending"
	StatusRunning   GenerationStatus = "running"
	StatusCompleted GenerationStatus = "completed"
	StatusFailed    GenerationStatus = "failed"
	StatusNotFound  GenerationStatus = "not_found"
)

type SpecGenerationStatus struct {
	AnalysisID   string
	CompletedAt  *time.Time
	ErrorMessage *string
	StartedAt    *time.Time
	Status       GenerationStatus
}
