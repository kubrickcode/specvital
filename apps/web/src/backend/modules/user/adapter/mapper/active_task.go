package mapper

import (
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/user/domain/entity"
)

func ToActiveTasksResponse(tasks []entity.ActiveTask) api.UserActiveTasksResponse {
	apiTasks := make([]api.ActiveTask, len(tasks))
	for i, task := range tasks {
		apiTasks[i] = toActiveTask(task)
	}
	return api.UserActiveTasksResponse{Tasks: apiTasks}
}

func toActiveTask(task entity.ActiveTask) api.ActiveTask {
	return api.ActiveTask{
		CreatedAt: task.CreatedAt,
		ID:        task.ID,
		Owner:     task.Owner,
		Repo:      task.Repo,
		StartedAt: task.StartedAt,
		Status:    mapTaskStatus(task.Status),
		Type:      mapTaskType(task.Type),
	}
}

func mapTaskStatus(status entity.TaskStatus) api.ActiveTaskStatus {
	switch status {
	case entity.TaskStatusAnalyzing:
		return api.Analyzing
	default:
		return api.Queued
	}
}

func mapTaskType(taskType entity.TaskType) api.ActiveTaskType {
	switch taskType {
	default:
		return api.ActiveTaskTypeAnalysis
	}
}
