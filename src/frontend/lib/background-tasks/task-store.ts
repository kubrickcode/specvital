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

type TaskStore = {
  listeners: Set<() => void>;
  tasks: Map<string, BackgroundTask>;
};

const STORAGE_KEY = "specvital:background-tasks";

const store: TaskStore = {
  listeners: new Set(),
  tasks: new Map(),
};

const loadFromStorage = (): void => {
  if (typeof window === "undefined") {
    return;
  }

  try {
    const stored = sessionStorage.getItem(STORAGE_KEY);
    if (stored) {
      const parsed = JSON.parse(stored) as [string, BackgroundTask][];
      store.tasks = new Map(parsed);
    }
  } catch {
    store.tasks = new Map();
  }
};

const saveToStorage = (): void => {
  if (typeof window === "undefined") {
    return;
  }

  try {
    const entries = Array.from(store.tasks.entries());
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(entries));
  } catch {
    // Storage quota exceeded or unavailable
  }
};

const notifyListeners = (): void => {
  store.listeners.forEach((listener) => listener());
};

loadFromStorage();

export const subscribe = (listener: () => void): (() => void) => {
  store.listeners.add(listener);
  return () => store.listeners.delete(listener);
};

export const getTasksSnapshot = (): Map<string, BackgroundTask> => store.tasks;

export const getServerSnapshot = (): Map<string, BackgroundTask> => new Map();

export const addTask = (task: Omit<BackgroundTask, "createdAt"> & { createdAt?: string }): void => {
  const taskWithCreatedAt: BackgroundTask = {
    ...task,
    createdAt: task.createdAt ?? new Date().toISOString(),
  };
  store.tasks = new Map(store.tasks).set(task.id, taskWithCreatedAt);
  saveToStorage();
  notifyListeners();
};

export const updateTask = (id: string, updates: Partial<Omit<BackgroundTask, "id">>): void => {
  const existing = store.tasks.get(id);
  if (!existing) {
    return;
  }

  const updated: BackgroundTask = { ...existing, ...updates };
  store.tasks = new Map(store.tasks).set(id, updated);
  saveToStorage();
  notifyListeners();
};

export const removeTask = (id: string): void => {
  if (!store.tasks.has(id)) {
    return;
  }

  const newTasks = new Map(store.tasks);
  newTasks.delete(id);
  store.tasks = newTasks;
  saveToStorage();
  notifyListeners();
};

export const getTask = (id: string): BackgroundTask | null => {
  return store.tasks.get(id) ?? null;
};

export const getAllTasks = (): BackgroundTask[] => {
  return Array.from(store.tasks.values());
};

export const clearCompletedTasks = (): void => {
  const newTasks = new Map<string, BackgroundTask>();
  store.tasks.forEach((task, id) => {
    if (task.status !== "completed" && task.status !== "failed") {
      newTasks.set(id, task);
    }
  });

  if (newTasks.size !== store.tasks.size) {
    store.tasks = newTasks;
    saveToStorage();
    notifyListeners();
  }
};
