import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import {
  addTask,
  type BackgroundTask,
  clearCompletedTasks,
  getAllTasks,
  getServerSnapshot,
  getTask,
  getTasksSnapshot,
  removeTask,
  subscribe,
  updateTask,
} from "./task-store";

const createTask = (overrides: Partial<BackgroundTask> = {}): BackgroundTask => ({
  createdAt: "2026-01-23T12:00:00.000Z",
  id: "task-1",
  metadata: { analysisId: "analysis-1", owner: "owner", repo: "repo" },
  startedAt: null,
  status: "queued",
  type: "spec-generation",
  ...overrides,
});

describe("task-store", () => {
  beforeEach(() => {
    vi.stubGlobal("sessionStorage", {
      getItem: vi.fn(() => null),
      setItem: vi.fn(),
    });

    getAllTasks().forEach((task) => removeTask(task.id));
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  describe("addTask", () => {
    it("should add task to store", () => {
      const task = createTask();
      addTask(task);

      expect(getTask("task-1")).toEqual(task);
    });

    it("should auto-generate createdAt if not provided", () => {
      const now = "2026-01-23T15:00:00.000Z";
      vi.setSystemTime(new Date(now));

      addTask({
        id: "task-2",
        metadata: {},
        startedAt: null,
        status: "queued",
        type: "analysis",
      });

      expect(getTask("task-2")?.createdAt).toBe(now);
      vi.useRealTimers();
    });

    it("should notify listeners on add", () => {
      const listener = vi.fn();
      const unsubscribe = subscribe(listener);

      addTask(createTask());

      expect(listener).toHaveBeenCalledTimes(1);
      unsubscribe();
    });

    it("should persist to sessionStorage", () => {
      const setItem = vi.fn();
      vi.stubGlobal("sessionStorage", { getItem: vi.fn(() => null), setItem });

      addTask(createTask());

      expect(setItem).toHaveBeenCalledWith("specvital:background-tasks", expect.any(String));
    });
  });

  describe("updateTask", () => {
    it("should update existing task", () => {
      addTask(createTask());

      updateTask("task-1", { startedAt: "2026-01-23T12:05:00.000Z", status: "processing" });

      const updated = getTask("task-1");
      expect(updated?.status).toBe("processing");
      expect(updated?.startedAt).toBe("2026-01-23T12:05:00.000Z");
    });

    it("should not modify store for non-existent task", () => {
      const listener = vi.fn();
      const unsubscribe = subscribe(listener);

      updateTask("non-existent", { status: "completed" });

      expect(listener).not.toHaveBeenCalled();
      unsubscribe();
    });

    it("should notify listeners on update", () => {
      addTask(createTask());

      const listener = vi.fn();
      const unsubscribe = subscribe(listener);

      updateTask("task-1", { status: "completed" });

      expect(listener).toHaveBeenCalledTimes(1);
      unsubscribe();
    });
  });

  describe("removeTask", () => {
    it("should remove task from store", () => {
      addTask(createTask());

      removeTask("task-1");

      expect(getTask("task-1")).toBeNull();
    });

    it("should not modify store for non-existent task", () => {
      const listener = vi.fn();
      const unsubscribe = subscribe(listener);

      removeTask("non-existent");

      expect(listener).not.toHaveBeenCalled();
      unsubscribe();
    });
  });

  describe("getAllTasks", () => {
    it("should return all tasks as array", () => {
      addTask(createTask({ id: "task-1" }));
      addTask(createTask({ id: "task-2", status: "processing" }));

      const tasks = getAllTasks();

      expect(tasks).toHaveLength(2);
    });
  });

  describe("clearCompletedTasks", () => {
    it("should remove completed and failed tasks", () => {
      addTask(createTask({ id: "task-1", status: "completed" }));
      addTask(createTask({ id: "task-2", status: "failed" }));
      addTask(createTask({ id: "task-3", status: "processing" }));
      addTask(createTask({ id: "task-4", status: "queued" }));

      clearCompletedTasks();

      expect(getAllTasks()).toHaveLength(2);
      expect(getTask("task-1")).toBeNull();
      expect(getTask("task-2")).toBeNull();
      expect(getTask("task-3")).not.toBeNull();
      expect(getTask("task-4")).not.toBeNull();
    });

    it("should not notify listeners if nothing to clear", () => {
      addTask(createTask({ id: "task-1", status: "processing" }));

      const listener = vi.fn();
      const unsubscribe = subscribe(listener);

      clearCompletedTasks();

      expect(listener).not.toHaveBeenCalled();
      unsubscribe();
    });
  });

  describe("subscribe", () => {
    it("should return unsubscribe function", () => {
      const listener = vi.fn();
      const unsubscribe = subscribe(listener);

      addTask(createTask());
      expect(listener).toHaveBeenCalledTimes(1);

      unsubscribe();
      addTask(createTask({ id: "task-2" }));
      expect(listener).toHaveBeenCalledTimes(1);
    });
  });

  describe("snapshots", () => {
    it("getTasksSnapshot should return current tasks Map", () => {
      addTask(createTask());

      const snapshot = getTasksSnapshot();

      expect(snapshot).toBeInstanceOf(Map);
      expect(snapshot.get("task-1")).toBeDefined();
    });

    it("getServerSnapshot should return empty Map for SSR", () => {
      const snapshot = getServerSnapshot();

      expect(snapshot).toBeInstanceOf(Map);
      expect(snapshot.size).toBe(0);
    });
  });
});
