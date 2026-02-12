"use client";

import { m } from "motion/react";

import { useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

type ShimmerBarHeight = "xs" | "sm" | "md";

export type ShimmerBarColor =
  | "var(--ai-primary)"
  | "var(--ai-secondary)"
  | "var(--chart-1)"
  | "var(--chart-2)"
  | "var(--destructive)"
  | "var(--primary)";

type ShimmerBarProps = {
  className?: string;
  color?: ShimmerBarColor;
  duration?: number;
  height?: ShimmerBarHeight;
};

const HEIGHT_MAP: Record<ShimmerBarHeight, string> = {
  md: "h-1.5",
  sm: "h-1",
  xs: "h-0.5",
};

export const ShimmerBar = ({
  className,
  color = "var(--primary)",
  duration = 2,
  height = "sm",
}: ShimmerBarProps) => {
  const shouldReduceMotion = useReducedMotion();

  return (
    <div
      aria-valuemax={100}
      aria-valuemin={0}
      aria-valuetext="Loading"
      className={cn("w-full overflow-hidden rounded-full bg-muted", HEIGHT_MAP[height], className)}
      role="progressbar"
    >
      {shouldReduceMotion ? (
        <div className="h-full w-full rounded-full opacity-30" style={{ background: color }} />
      ) : (
        <m.div
          animate={{ x: ["-100%", "400%"] }}
          className="h-full w-1/3 rounded-full"
          style={{
            background: `linear-gradient(90deg, transparent, ${color}, transparent)`,
          }}
          transition={{ duration, ease: "easeInOut", repeat: Infinity }}
        />
      )}
    </div>
  );
};
