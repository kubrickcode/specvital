"use client";

import { Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";

import type { BackgroundTaskStatus } from "../task-store";
import { formatElapsedTime } from "../utils";

type TaskStatusBadgeProps = {
  size?: "sm" | "md";
  startedAt: string | null;
  status: BackgroundTaskStatus;
};

export const TaskStatusBadge = ({ size = "sm", startedAt, status }: TaskStatusBadgeProps) => {
  const t = useTranslations("backgroundTasks.status");
  const [elapsed, setElapsed] = useState(() => formatElapsedTime(startedAt));

  useEffect(() => {
    setElapsed(formatElapsedTime(startedAt));
  }, [startedAt]);

  useEffect(() => {
    if (status !== "processing" || !startedAt) return;

    const interval = setInterval(() => {
      setElapsed(formatElapsedTime(startedAt));
    }, 1000);

    return () => clearInterval(interval);
  }, [status, startedAt]);

  const iconSize = size === "sm" ? "size-3" : "size-4";
  const dotSize = size === "sm" ? "size-1.5" : "size-2";
  const textSize = size === "sm" ? "text-xs" : "text-sm";

  if (status === "processing") {
    return (
      <span className={`flex items-center gap-1 ${textSize} text-muted-foreground`}>
        <Loader2 className={`${iconSize} animate-spin`} />
        {t("processing", { time: elapsed })}
      </span>
    );
  }

  if (status === "queued") {
    return (
      <span className={`flex items-center gap-1 ${textSize} text-muted-foreground`}>
        <span className={`${dotSize} rounded-full bg-amber-500`} />
        {t("queued")}
      </span>
    );
  }

  return null;
};
