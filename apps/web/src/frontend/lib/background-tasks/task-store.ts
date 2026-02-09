"use client";

import { create } from "zustand";

export type BackgroundTaskStatus = "completed" | "failed" | "processing" | "queued";

export type BackgroundTaskType = "analysis" | "spec-generation";

export type BackgroundTaskMetadata = {
  analysisId?: string;
  language?: string;
  owner?: string;
  repo?: string;
};

export type BackgroundTask = {
  createdAt: string;
  id: string;
  metadata: BackgroundTaskMetadata;
  startedAt: string | null;
  status: BackgroundTaskStatus;
  type: BackgroundTaskType;
};

const STORAGE_KEY = "specvital:background-tasks";

const loadFromStorage = (): Map<string, BackgroundTask> => {
  if (typeof window === "undefined") return new Map();

  try {
    const stored = sessionStorage.getItem(STORAGE_KEY);
    if (stored) {
      const parsed = JSON.parse(stored) as [string, BackgroundTask][];
      return new Map(parsed);
    }
  } catch {
    // Corrupted storage
  }
  return new Map();
};

const saveToStorage = (tasks: Map<string, BackgroundTask>): void => {
  if (typeof window === "undefined") return;

  try {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(Array.from(tasks.entries())));
  } catch {
    // Storage quota exceeded or unavailable
  }
};

type TaskStoreState = {
  tasks: Map<string, BackgroundTask>;
};

type TaskStoreActions = {
  addTask: (task: Omit<BackgroundTask, "createdAt"> & { createdAt?: string }) => void;
  clearCompletedTasks: () => void;
  removeTask: (id: string) => void;
  updateTask: (id: string, updates: Partial<Omit<BackgroundTask, "id">>) => void;
};

const useTaskStore = create<TaskStoreState & TaskStoreActions>((set, get) => ({
  addTask: (task) => {
    const taskWithCreatedAt: BackgroundTask = {
      ...task,
      createdAt: task.createdAt ?? new Date().toISOString(),
    };
    const newTasks = new Map(get().tasks).set(task.id, taskWithCreatedAt);
    saveToStorage(newTasks);
    set({ tasks: newTasks });
  },

  clearCompletedTasks: () => {
    const newTasks = new Map<string, BackgroundTask>();
    get().tasks.forEach((task, id) => {
      if (task.status !== "completed" && task.status !== "failed") {
        newTasks.set(id, task);
      }
    });
    if (newTasks.size !== get().tasks.size) {
      saveToStorage(newTasks);
      set({ tasks: newTasks });
    }
  },

  removeTask: (id) => {
    if (!get().tasks.has(id)) return;
    const newTasks = new Map(get().tasks);
    newTasks.delete(id);
    saveToStorage(newTasks);
    set({ tasks: newTasks });
  },

  tasks: loadFromStorage(),

  updateTask: (id, updates) => {
    const existing = get().tasks.get(id);
    if (!existing) return;
    const updated: BackgroundTask = { ...existing, ...updates };
    const newTasks = new Map(get().tasks).set(id, updated);
    saveToStorage(newTasks);
    set({ tasks: newTasks });
  },
}));

// Standalone functions for non-React usage
export const addTask = (task: Omit<BackgroundTask, "createdAt"> & { createdAt?: string }): void =>
  useTaskStore.getState().addTask(task);

export const updateTask = (id: string, updates: Partial<Omit<BackgroundTask, "id">>): void =>
  useTaskStore.getState().updateTask(id, updates);

export const removeTask = (id: string): void => useTaskStore.getState().removeTask(id);

export const getTask = (id: string): BackgroundTask | null =>
  useTaskStore.getState().tasks.get(id) ?? null;

export const getAllTasks = (): BackgroundTask[] =>
  Array.from(useTaskStore.getState().tasks.values());

export const clearCompletedTasks = (): void => useTaskStore.getState().clearCompletedTasks();

// Export store for hook usage
export { useTaskStore };
