"use client";

import { type BackgroundTask, useTaskStore } from "@/lib/background-tasks";

import { getSpecGenerationTaskId } from "../utils/task-ids";

export const useSpecGenerationTask = (analysisId: string): BackgroundTask | null =>
  useTaskStore((state) => {
    const task = state.tasks.get(getSpecGenerationTaskId(analysisId));
    return task && (task.status === "queued" || task.status === "processing") ? task : null;
  });
