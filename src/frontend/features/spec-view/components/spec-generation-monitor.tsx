"use client";

import { useTranslations } from "next-intl";
import { toast } from "sonner";

import { getTask, removeTask, useBackgroundTasks } from "@/lib/background-tasks";
import type { BackgroundTask } from "@/lib/background-tasks";

import { useSpecGenerationStatus } from "../hooks/use-spec-generation-status";
import type { SpecLanguage } from "../types";

type SpecGenerationTaskPollerProps = {
  analysisId: string;
  language: SpecLanguage | undefined;
  owner: string | undefined;
  repo: string | undefined;
  taskId: string;
};

/**
 * Polls spec generation status for a single background task.
 *
 * Uses the same `specGenerationStatusKeys.status` queryKey as SpecPanel's
 * `useSpecGenerationStatus`, so React Query deduplicates HTTP requests
 * when both are mounted.
 *
 * Deduplication: getTask() check before removeTask() ensures only one
 * handler (this or SpecPanel) processes the completion — execution order
 * between siblings is irrelevant.
 */
const SpecGenerationTaskPoller = ({
  analysisId,
  language,
  owner,
  repo,
  taskId,
}: SpecGenerationTaskPollerProps) => {
  const t = useTranslations("specView.toast");

  useSpecGenerationStatus(analysisId, {
    enabled: Boolean(analysisId),
    language,
    onCompleted: () => {
      const task = getTask(taskId);
      if (task) {
        removeTask(taskId);
        toast.success(t("generationComplete.title"), {
          description: t("generationComplete.description"),
        });
      }
    },
    onFailed: () => {
      const task = getTask(taskId);
      if (task) {
        removeTask(taskId);
        toast.error(t("generateFailed.title"));
      }
    },
    owner,
    repo,
  });

  return null;
};

/**
 * Global monitor for spec generation background tasks.
 * Renders at layout level to continue polling even when SpecPanel is unmounted.
 */
export const SpecGenerationMonitor = () => {
  const tasks = useBackgroundTasks();

  const activeSpecGenerationTasks = tasks.filter(
    (
      task
    ): task is BackgroundTask & {
      metadata: { analysisId: string };
    } =>
      task.type === "spec-generation" &&
      (task.status === "queued" || task.status === "processing") &&
      Boolean(task.metadata.analysisId)
  );

  return (
    <>
      {activeSpecGenerationTasks.map((task) => (
        <SpecGenerationTaskPoller
          analysisId={task.metadata.analysisId}
          key={task.id}
          language={task.metadata.language as SpecLanguage | undefined}
          owner={task.metadata.owner}
          repo={task.metadata.repo}
          taskId={task.id}
        />
      ))}
    </>
  );
};
