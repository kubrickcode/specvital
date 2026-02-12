"use client";

import { ChevronDown, ExternalLink, ListTodo } from "lucide-react";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Link } from "@/i18n/navigation";
import type { ActiveTask } from "@/lib/api/types";
import {
  type BackgroundTask,
  TaskStatusBadge,
  useActiveTasks,
  useUserActiveTasks,
} from "@/lib/background-tasks";
import { cn } from "@/lib/utils";

// Unified task type for display
type DisplayTask = {
  id: string;
  owner: string;
  repo: string;
  startedAt: string | null;
  status: "processing" | "queued";
  type: "analysis" | "spec-generation";
};

// Map server status to UI status
const mapServerStatusToUiStatus = (status: ActiveTask["status"]): "processing" | "queued" => {
  if (status === "analyzing") return "processing";
  return "queued";
};

// Convert server task to display task
const serverTaskToDisplayTask = (task: ActiveTask): DisplayTask => ({
  id: task.id,
  owner: task.owner,
  repo: task.repo,
  startedAt: task.startedAt ?? null,
  status: mapServerStatusToUiStatus(task.status),
  type: task.type,
});

// Convert local task to display task (spec-generation only)
const localTaskToDisplayTask = (task: BackgroundTask): DisplayTask | null => {
  if (task.type !== "spec-generation") return null;
  if (!task.metadata.owner || !task.metadata.repo) return null;
  return {
    id: task.id,
    owner: task.metadata.owner,
    repo: task.metadata.repo,
    startedAt: task.startedAt,
    status: task.status === "processing" ? "processing" : "queued",
    type: task.type,
  };
};

type TaskRowProps = {
  task: DisplayTask;
};

const TaskRow = ({ task }: TaskRowProps) => {
  const tType = useTranslations("backgroundTasks.taskType");
  const tActiveTasks = useTranslations("backgroundTasks.activeTasks");

  const pageUrl = `/analyze/${task.owner}/${task.repo}`;
  const displayName = `${task.owner}/${task.repo}`;

  return (
    <div className="flex items-center justify-between gap-4 rounded-lg border bg-muted/30 px-4 py-3">
      <div className="flex min-w-0 flex-col gap-0.5">
        <span className="truncate text-sm font-medium">{displayName}</span>
        <span className="text-xs text-muted-foreground">{tType(task.type)}</span>
      </div>
      <div className="flex shrink-0 items-center gap-3">
        <TaskStatusBadge size="sm" startedAt={task.startedAt} status={task.status} />
        <Button asChild size="sm" variant="ghost">
          <Link href={pageUrl}>
            {tActiveTasks("viewTask")}
            <ExternalLink className="ml-1 size-3" />
          </Link>
        </Button>
      </div>
    </div>
  );
};

/**
 * Dashboard section showing active tasks for authenticated users.
 * Combines server API tasks (analysis) and sessionStorage tasks (spec-generation).
 */
export const ActiveTasksSection = () => {
  const t = useTranslations("backgroundTasks.activeTasks");
  // Always enabled - this component is only rendered for authenticated users
  const { data } = useUserActiveTasks({ enabled: true });
  const localTasks = useActiveTasks();
  const searchParams = useSearchParams();

  // Combine server tasks and local spec-generation tasks
  const serverTasks: DisplayTask[] = (data?.tasks ?? []).map(serverTaskToDisplayTask);
  const specGenTasks: DisplayTask[] = localTasks
    .map(localTaskToDisplayTask)
    .filter((task): task is DisplayTask => task !== null);

  const activeTasks = [...serverTasks, ...specGenTasks];

  const shouldAutoExpand = searchParams.get("section") === "tasks";
  const [isOpen, setIsOpen] = useState(shouldAutoExpand);

  useEffect(() => {
    if (shouldAutoExpand) {
      setIsOpen(true);
    }
  }, [shouldAutoExpand]);

  if (activeTasks.length === 0) {
    return null;
  }

  return (
    <Collapsible onOpenChange={setIsOpen} open={isOpen}>
      <Card className="overflow-hidden py-0" depth="raised">
        <CollapsibleTrigger asChild>
          <button
            className="flex w-full cursor-pointer items-center justify-between px-6 py-4 hover:bg-muted/50 transition-colors"
            type="button"
          >
            <div className="flex items-center gap-3">
              <div className="flex size-10 items-center justify-center rounded-lg bg-gradient-to-br from-violet-500/20 to-violet-500/5 text-violet-500 ring-1 ring-violet-500/10">
                <ListTodo aria-hidden="true" className="size-5" />
              </div>
              <span className="text-base font-semibold">{t("title")}</span>
              <Badge variant="secondary">{activeTasks.length}</Badge>
            </div>
            <ChevronDown
              className={cn(
                "size-5 text-muted-foreground transition-transform duration-200",
                isOpen && "rotate-180"
              )}
            />
          </button>
        </CollapsibleTrigger>
        <CollapsibleContent>
          <CardContent className="space-y-2 pb-4">
            {activeTasks.map((task) => (
              <TaskRow key={task.id} task={task} />
            ))}
          </CardContent>
        </CollapsibleContent>
      </Card>
    </Collapsible>
  );
};
