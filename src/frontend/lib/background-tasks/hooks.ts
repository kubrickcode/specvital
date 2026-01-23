"use client";

import { useSyncExternalStore } from "react";

import { type BackgroundTask, getServerSnapshot, getTasksSnapshot, subscribe } from "./task-store";

export const useBackgroundTasks = (): BackgroundTask[] => {
  const tasks = useSyncExternalStore(subscribe, getTasksSnapshot, getServerSnapshot);
  return Array.from(tasks.values());
};

export const useBackgroundTask = (id: string): BackgroundTask | null => {
  const tasks = useSyncExternalStore(subscribe, getTasksSnapshot, getServerSnapshot);
  return tasks.get(id) ?? null;
};

export const useActiveTaskCount = (): number => {
  const tasks = useSyncExternalStore(subscribe, getTasksSnapshot, getServerSnapshot);
  let count = 0;
  tasks.forEach((task) => {
    if (task.status === "queued" || task.status === "processing") {
      count++;
    }
  });
  return count;
};
