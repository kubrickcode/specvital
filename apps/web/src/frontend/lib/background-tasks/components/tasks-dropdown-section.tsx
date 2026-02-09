"use client";

import { FileText, Search } from "lucide-react";
import { useTranslations } from "next-intl";

import { ShimmerBar } from "@/components/feedback";
import { DropdownMenuItem, DropdownMenuSeparator } from "@/components/ui/dropdown-menu";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { Link } from "@/i18n/navigation";
import type { ActiveTask } from "@/lib/api/types";

import { useActiveTasks } from "../hooks";
import type { BackgroundTask } from "../task-store";
import { useUserActiveTasks } from "../use-user-active-tasks";
import { TaskStatusBadge } from "./task-status-badge";

// Unified task type for display
type DisplayTask = {
  id: string;
  language?: string;
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
    language: task.metadata.language,
    owner: task.metadata.owner,
    repo: task.metadata.repo,
    startedAt: task.startedAt,
    status: task.status === "processing" ? "processing" : "queued",
    type: task.type,
  };
};

type TaskItemProps = {
  task: DisplayTask;
};

const TaskItem = ({ task }: TaskItemProps) => {
  const tType = useTranslations("backgroundTasks.taskType");
  const tDropdown = useTranslations("backgroundTasks.dropdown");

  const pageUrl = `/analyze/${task.owner}/${task.repo}`;
  const displayName = `${task.owner}/${task.repo}`;
  const TaskIcon = task.type === "spec-generation" ? FileText : Search;

  const content = (
    <div className="flex w-full flex-col gap-0.5">
      <div className="flex items-center justify-between gap-2">
        <Tooltip>
          <TooltipTrigger asChild>
            <div className="flex min-w-0 items-center gap-2">
              <TaskIcon className="size-3.5 shrink-0 text-muted-foreground" />
              <span className="truncate text-sm font-medium">{displayName}</span>
            </div>
          </TooltipTrigger>
          <TooltipContent side="top">{displayName}</TooltipContent>
        </Tooltip>
        <span className="shrink-0 text-xs text-primary">{tDropdown("viewPage")}</span>
      </div>
      <div className="flex items-center justify-between gap-2">
        <span className="text-xs text-muted-foreground">{tType(task.type)}</span>
        <TaskStatusBadge startedAt={task.startedAt} status={task.status} />
      </div>
      <div className="mt-1">
        <ShimmerBar aria-hidden="true" color="var(--primary)" duration={2} height="xs" />
      </div>
    </div>
  );

  return (
    <DropdownMenuItem asChild>
      <Link href={pageUrl}>{content}</Link>
    </DropdownMenuItem>
  );
};

type TasksDropdownSectionProps = {
  className?: string;
};

/**
 * Dropdown section showing active tasks for authenticated users.
 * Combines server API tasks (analysis) and sessionStorage tasks (spec-generation).
 */
export const TasksDropdownSection = ({ className }: TasksDropdownSectionProps) => {
  const t = useTranslations("backgroundTasks.dropdown");
  // Always enabled - this component is only rendered for authenticated users
  const { data } = useUserActiveTasks({ enabled: true });
  const localTasks = useActiveTasks();

  // Combine server tasks and local spec-generation tasks
  const serverTasks: DisplayTask[] = (data?.tasks ?? []).map(serverTaskToDisplayTask);
  const specGenTasks: DisplayTask[] = localTasks
    .map(localTaskToDisplayTask)
    .filter((task): task is DisplayTask => task !== null);

  const activeTasks = [...serverTasks, ...specGenTasks];

  if (activeTasks.length === 0) {
    return null;
  }

  const MAX_DISPLAY_TASKS = 3;
  const displayedTasks = activeTasks.slice(0, MAX_DISPLAY_TASKS);
  const remainingCount = activeTasks.length - MAX_DISPLAY_TASKS;

  return (
    <div className={className}>
      <div className="px-2 py-1.5">
        <p className="text-xs font-medium text-muted-foreground">{t("title")}</p>
      </div>
      {displayedTasks.map((task) => (
        <TaskItem key={task.id} task={task} />
      ))}
      {remainingCount > 0 && (
        <DropdownMenuItem asChild>
          <Link href="/dashboard?section=tasks">
            <span className="text-xs text-muted-foreground">
              {t("viewAll", { count: remainingCount })}
            </span>
          </Link>
        </DropdownMenuItem>
      )}
      <DropdownMenuSeparator />
    </div>
  );
};
