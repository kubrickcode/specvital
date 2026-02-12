"use client";

import { cn } from "@/lib/utils";

import type { StatusCounts } from "../utils/calculate-status-counts";

type StatusMiniBarProps = {
  className?: string;
  counts: StatusCounts;
};

export const StatusMiniBar = ({ className, counts }: StatusMiniBarProps) => {
  const { active, focused, skipped, todo, xfail } = counts;
  const total = active + focused + skipped + todo + xfail;

  if (total === 0) {
    return null;
  }

  const segments = [
    { colorClass: "bg-status-active", count: active, label: "active" },
    { colorClass: "bg-status-focused", count: focused, label: "focused" },
    { colorClass: "bg-status-skipped", count: skipped, label: "skipped" },
    { colorClass: "bg-status-xfail", count: xfail, label: "xfail" },
    { colorClass: "bg-status-todo", count: todo, label: "todo" },
  ];

  const ariaLabel = segments
    .map((s) => `${s.count} ${s.label}`)
    .join(", ")
    .concat(` out of ${total} tests`);

  return (
    <div
      aria-label={ariaLabel}
      className={cn("flex h-1.5 w-16 overflow-hidden rounded-full bg-muted", className)}
      role="img"
    >
      {segments.map(
        (segment) =>
          segment.count > 0 && (
            <div
              className={segment.colorClass}
              key={segment.label}
              style={{ width: `${(segment.count / total) * 100}%` }}
            />
          )
      )}
    </div>
  );
};
