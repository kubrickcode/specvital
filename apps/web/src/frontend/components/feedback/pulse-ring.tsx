"use client";

import { m } from "motion/react";

import { useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

type PulseRingSize = "xs" | "sm" | "md";

type PulseRingProps = {
  className?: string;
  size?: PulseRingSize;
};

const SIZE_MAP: Record<PulseRingSize, { container: string; dot: string }> = {
  md: { container: "size-5", dot: "size-3" },
  sm: { container: "size-4", dot: "size-2.5" },
  xs: { container: "size-2", dot: "size-2" },
};

export const PulseRing = ({ className, size = "sm" }: PulseRingProps) => {
  const shouldReduceMotion = useReducedMotion();
  const sizes = SIZE_MAP[size];

  return (
    <span
      className={cn("relative inline-flex items-center justify-center", sizes.container, className)}
    >
      {!shouldReduceMotion && (
        <m.span
          animate={{ opacity: [0, 0.4, 0], scale: [0.8, 1, 2] }}
          className="absolute inset-0 rounded-full bg-current"
          transition={{ duration: 2, ease: "easeOut", repeat: Infinity }}
        />
      )}
      <span className={cn("rounded-full bg-current", sizes.dot)} />
    </span>
  );
};
