"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";

import { useAuth } from "@/features/auth";
import { paginatedRepositoriesKeys, repositoryStatsKeys } from "@/features/dashboard";
import { useUserActiveTasks, userActiveTasksKeys } from "@/lib/background-tasks";

import { fetchAnalysisStatus } from "../api";
import { analysisKeys } from "../hooks/use-analysis";
import { updateStatusKeys } from "../hooks/use-update-status";

type TrackedTask = {
  id: string;
  owner: string;
  repo: string;
};

/**
 * Global monitor for analysis background tasks.
 * Renders at layout level to continue polling even when analysis page is unmounted.
 *
 * For authenticated users: Uses server API to get active tasks.
 * For unauthenticated users: No global monitoring (only current page polling).
 *
 * Key insight: When a task completes, the server removes it from active tasks immediately.
 * We detect completion by comparing current vs previous task lists, then poll the API
 * to confirm the final status and show appropriate toast.
 */
export const AnalysisMonitor = () => {
  const t = useTranslations("backgroundTasks.toast");
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuth();
  const { data } = useUserActiveTasks({ enabled: isAuthenticated });

  // Track tasks we're monitoring (to detect when they disappear from server)
  const [trackedTasks, setTrackedTasks] = useState<Map<string, TrackedTask>>(new Map());
  const processingRef = useRef<Set<string>>(new Set());

  // Update tracked tasks when server returns new active tasks
  useEffect(() => {
    const tasks = data?.tasks;
    if (!tasks) return;

    setTrackedTasks((prev) => {
      const next = new Map(prev);

      // Add new tasks from server
      for (const task of tasks) {
        if (!next.has(task.id)) {
          next.set(task.id, { id: task.id, owner: task.owner, repo: task.repo });
        }
      }

      return next;
    });
  }, [data?.tasks]);

  // Detect tasks that disappeared from server (completed or failed)
  useEffect(() => {
    const tasks = data?.tasks;
    if (!tasks) return;

    const currentTaskIds = new Set(tasks.map((t) => t.id));

    // Find tasks that were tracked but are no longer in server response
    const disappearedTasks: TrackedTask[] = [];
    trackedTasks.forEach((task) => {
      if (!currentTaskIds.has(task.id) && !processingRef.current.has(task.id)) {
        disappearedTasks.push(task);
      }
    });

    if (disappearedTasks.length === 0) return;

    // Mark as processing to prevent duplicate handling
    disappearedTasks.forEach((task) => processingRef.current.add(task.id));

    // Check final status and show toast for each disappeared task
    const handleDisappearedTasks = async () => {
      for (const task of disappearedTasks) {
        try {
          const result = await fetchAnalysisStatus(task.owner, task.repo);
          const repoName = `${task.owner}/${task.repo}`;

          // Invalidate queries to refresh data
          queryClient.invalidateQueries({ queryKey: analysisKeys.detail(task.owner, task.repo) });
          queryClient.invalidateQueries({ queryKey: paginatedRepositoriesKeys.all });
          queryClient.invalidateQueries({ queryKey: repositoryStatsKeys.all });
          queryClient.invalidateQueries({
            queryKey: updateStatusKeys.detail(task.owner, task.repo),
          });

          if (result.status === "completed") {
            toast.success(t("analysisComplete", { repo: repoName }));
          } else if (result.status === "failed") {
            toast.error(t("analysisFailed", { repo: repoName }));
          }
        } catch {
          // Task might have been deleted, just remove from tracking
        }

        // Remove from tracked tasks
        setTrackedTasks((prev) => {
          const next = new Map(prev);
          next.delete(task.id);
          return next;
        });
        processingRef.current.delete(task.id);
      }

      // Refresh active tasks list
      queryClient.invalidateQueries({ queryKey: userActiveTasksKeys.all });
    };

    handleDisappearedTasks();
  }, [data?.tasks, trackedTasks, queryClient, t]);

  return null;
};
