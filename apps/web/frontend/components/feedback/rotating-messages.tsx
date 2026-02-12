"use client";

import { AnimatePresence, m } from "motion/react";

import { useReducedMotion } from "@/lib/motion";

type RotatingMessagesProps = {
  className?: string;
  currentIndex: number;
  messages: string[];
};

export const RotatingMessages = ({ className, currentIndex, messages }: RotatingMessagesProps) => {
  const shouldReduceMotion = useReducedMotion();
  const safeIndex = Math.min(currentIndex, messages.length - 1);
  const message = messages[safeIndex];

  if (!message) {
    return null;
  }

  // aria-live="off": supplementary status text, not authoritative state.
  // Parent's role="status" handles critical announcements.
  return (
    <div aria-live="off" className={className}>
      {shouldReduceMotion ? (
        <p className="text-sm text-muted-foreground">{message}</p>
      ) : (
        <AnimatePresence mode="wait">
          <m.p
            animate={{ opacity: 1, y: 0 }}
            className="text-sm text-muted-foreground"
            exit={{ opacity: 0, y: -8 }}
            initial={{ opacity: 0, y: 8 }}
            key={safeIndex}
            transition={{ duration: 0.3 }}
          >
            {message}
          </m.p>
        </AnimatePresence>
      )}
    </div>
  );
};
