"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useEffect, useRef } from "react";
import { toast } from "sonner";

import { paginatedRepositoriesKeys } from "@/features/dashboard";
import type { AnalysisResponse } from "@/lib/api/types";
import { getTask, removeTask, useBackgroundTasks } from "@/lib/background-tasks";
import type { BackgroundTask } from "@/lib/background-tasks";

import { fetchAnalysis } from "../api";
import { analysisKeys } from "../hooks/use-analysis";
import { updateStatusKeys } from "../hooks/use-update-status";

const POLL_INTERVAL_MS = 1000;

const isTerminalStatus = (status: AnalysisResponse["status"]): boolean =>
  status === "completed" || status === "failed";

type AnalysisTaskPollerProps = {
  owner: string;
  repo: string;
  taskId: string;
};

/**
 * Polls analysis status for a single background task.
 *
 * Uses `analysisKeys.detail` + `fetchAnalysis` (same as useAnalysis) so React Query
 * deduplicates requests when both are mounted. When both are active, the shorter
 * refetchInterval wins — useAnalysis backoff (300ms→5000ms) is capped at 2000ms.
 * This is acceptable: the monitor exists primarily for when useAnalysis is unmounted.
 *
 * fetchAnalysis can trigger a new analysis if none exists, but the monitor only
 * runs for tasks already in the store (queued/processing), so no side effect.
 *
 * Deduplication: getTask() check before removeTask() ensures only one
 * handler (this or local component) processes the completion.
 */
const AnalysisTaskPoller = ({ owner, repo, taskId }: AnalysisTaskPollerProps) => {
  const t = useTranslations("backgroundTasks.toast");
  const queryClient = useQueryClient();
  const completedRef = useRef(false);

  const pollingQuery = useQuery({
    queryFn: () => fetchAnalysis(owner, repo),
    queryKey: analysisKeys.detail(owner, repo),
    refetchInterval: (query) => {
      const data = query.state.data;
      if (!data) return POLL_INTERVAL_MS;
      return isTerminalStatus(data.status) ? false : POLL_INTERVAL_MS;
    },
    staleTime: 0,
  });

  useEffect(() => {
    if (!pollingQuery.data || completedRef.current) return;

    const { status } = pollingQuery.data;
    if (!isTerminalStatus(status)) return;

    completedRef.current = true;

    queryClient.invalidateQueries({ queryKey: paginatedRepositoriesKeys.all });
    queryClient.invalidateQueries({ queryKey: updateStatusKeys.detail(owner, repo) });

    const task = getTask(taskId);
    if (task) {
      if (status === "completed") {
        toast.success(t("analysisComplete", { repo: `${owner}/${repo}` }));
      } else {
        toast.error(t("analysisFailed", { repo: `${owner}/${repo}` }));
      }

      setTimeout(() => {
        removeTask(taskId);
      }, 100);
    }
  }, [pollingQuery.data, taskId, owner, repo, queryClient, t]);

  return null;
};

/**
 * Global monitor for analysis background tasks.
 * Renders at layout level to continue polling even when analysis page is unmounted.
 */
export const AnalysisMonitor = () => {
  const tasks = useBackgroundTasks();

  const activeAnalysisTasks = tasks.filter(
    (task): task is BackgroundTask & { metadata: { owner: string; repo: string } } =>
      task.type === "analysis" &&
      (task.status === "queued" || task.status === "processing") &&
      Boolean(task.metadata.owner) &&
      Boolean(task.metadata.repo)
  );

  return (
    <>
      {activeAnalysisTasks.map((task) => (
        <AnalysisTaskPoller
          key={task.id}
          owner={task.metadata.owner}
          repo={task.metadata.repo}
          taskId={task.id}
        />
      ))}
    </>
  );
};
