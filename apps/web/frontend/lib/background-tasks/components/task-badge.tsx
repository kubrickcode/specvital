"use client";

import { Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";

import { cn } from "@/lib/utils";

import { useActiveTasks } from "../hooks";
import { useUserActiveTasks } from "../use-user-active-tasks";

type TaskBadgeProps = {
  className?: string;
};

/**
 * Badge showing active task count for authenticated users.
 * Combines server API tasks (analysis) and sessionStorage tasks (spec-generation).
 */
export const TaskBadge = ({ className }: TaskBadgeProps) => {
  const t = useTranslations("backgroundTasks.badge");
  // Always enabled - this component is only rendered for authenticated users
  const { activeTaskCount: serverTaskCount } = useUserActiveTasks({ enabled: true });
  // Get spec-generation tasks from sessionStorage (not available via server API)
  const localTasks = useActiveTasks();
  const specGenTaskCount = localTasks.filter((task) => task.type === "spec-generation").length;

  const activeTaskCount = serverTaskCount + specGenTaskCount;

  if (activeTaskCount === 0) {
    return null;
  }

  return (
    <span
      aria-label={t("ariaLabel", { count: activeTaskCount })}
      className={cn(
        "absolute -right-0.5 -top-0.5 flex size-3.5 items-center justify-center rounded-full bg-primary text-[8px] font-medium text-primary-foreground",
        className
      )}
    >
      <Loader2 className="size-2 animate-spin" />
    </span>
  );
};
