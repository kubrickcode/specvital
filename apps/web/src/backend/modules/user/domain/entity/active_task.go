package entity

import "time"

type ActiveTask struct {
	CreatedAt time.Time
	ID        string
	Owner     string
	Repo      string
	StartedAt *time.Time
	Status    TaskStatus
	Type      TaskType
}

type TaskType string

const (
	TaskTypeAnalysis TaskType = "analysis"
)

type TaskStatus string

const (
	TaskStatusQueued    TaskStatus = "queued"
	TaskStatusAnalyzing TaskStatus = "analyzing"
)

func MapRiverStateToTaskStatus(state string) TaskStatus {
	switch state {
	case "running":
		return TaskStatusAnalyzing
	default:
		return TaskStatusQueued
	}
}
