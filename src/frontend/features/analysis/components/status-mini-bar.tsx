"use client";

import { cn } from "@/lib/utils";

import type { StatusCounts } from "../utils/calculate-status-counts";

type StatusMiniBarProps = {
  className?: string;
  counts: StatusCounts;
};

export const StatusMiniBar = ({ className, counts }: StatusMiniBarProps) => {
  const { active, skipped, todo } = counts;
  const total = active + skipped + todo;

  if (total === 0) {
    return null;
  }

  const activePercent = (active / total) * 100;
  const skippedPercent = (skipped / total) * 100;
  const todoPercent = (todo / total) * 100;

  const ariaLabel = `${active} active, ${skipped} skipped, ${todo} todo out of ${total} tests`;

  return (
    <div
      aria-label={ariaLabel}
      className={cn("flex h-1.5 w-16 overflow-hidden rounded-full bg-muted", className)}
      role="img"
    >
      {activePercent > 0 && (
        <div className="bg-status-active" style={{ width: `${activePercent}%` }} />
      )}
      {skippedPercent > 0 && (
        <div className="bg-status-skipped" style={{ width: `${skippedPercent}%` }} />
      )}
      {todoPercent > 0 && <div className="bg-status-todo" style={{ width: `${todoPercent}%` }} />}
    </div>
  );
};
