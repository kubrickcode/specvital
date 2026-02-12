"use client";

import { useEffect, useState } from "react";

type UseRotatingIndexOptions = {
  count: number;
  intervalMs?: number;
};

export const useRotatingIndex = ({ count, intervalMs = 6000 }: UseRotatingIndexOptions): number => {
  const [index, setIndex] = useState(0);

  useEffect(() => {
    if (count <= 1) return;

    const timer = setInterval(() => {
      setIndex((prev) => (prev + 1) % count);
    }, intervalMs);

    return () => clearInterval(timer);
  }, [count, intervalMs]);

  return index;
};
