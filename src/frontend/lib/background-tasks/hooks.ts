"use client";

import { useShallow } from "zustand/react/shallow";

import { type BackgroundTask, useTaskStore } from "./task-store";

export const useBackgroundTasks = (): BackgroundTask[] =>
  useTaskStore(useShallow((state) => Array.from(state.tasks.values())));

export const useActiveTaskCount = (): number =>
  useTaskStore((state) => {
    let count = 0;
    state.tasks.forEach((task) => {
      if (task.status === "queued" || task.status === "processing") {
        count++;
      }
    });
    return count;
  });
