"use client";

import { throttle } from "es-toolkit";
import { motion } from "motion/react";
import { useEffect, useState } from "react";

import { useReducedMotion } from "@/lib/motion";

const THROTTLE_MS = 16;

export const ReadingProgressBar = () => {
  const [progress, setProgress] = useState(0);
  const shouldReduceMotion = useReducedMotion();

  useEffect(() => {
    const calculateProgress = () => {
      const scrollTop = window.scrollY;
      const docHeight = document.documentElement.scrollHeight - window.innerHeight;

      if (docHeight <= 0) {
        setProgress(0);
        return;
      }

      const scrollProgress = Math.min(Math.max((scrollTop / docHeight) * 100, 0), 100);
      setProgress(scrollProgress);
    };

    const throttledCalculate = throttle(calculateProgress, THROTTLE_MS);

    calculateProgress();
    window.addEventListener("scroll", throttledCalculate, { passive: true });

    return () => {
      window.removeEventListener("scroll", throttledCalculate);
      throttledCalculate.cancel();
    };
  }, []);

  return (
    <motion.div
      aria-label="Reading progress"
      aria-valuemax={100}
      aria-valuemin={0}
      aria-valuenow={Math.round(progress)}
      className="fixed top-0 left-0 right-0 z-50 h-1 bg-muted"
      role="progressbar"
    >
      <motion.div
        animate={{ width: `${progress}%` }}
        className="h-full bg-primary"
        initial={{ width: 0 }}
        transition={shouldReduceMotion ? { duration: 0 } : { duration: 0.1, ease: "linear" }}
      />
    </motion.div>
  );
};
