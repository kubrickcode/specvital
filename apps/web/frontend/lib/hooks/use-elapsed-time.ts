"use client";

import { useEffect, useState } from "react";

type ElapsedTimeResult = {
  ariaLabel: string;
  formatted: string;
  seconds: number;
};

const EMPTY_RESULT: ElapsedTimeResult = {
  ariaLabel: "",
  formatted: "0:00",
  seconds: 0,
};

const formatAriaLabel = (minutes: number, secs: number): string => {
  if (minutes > 0) {
    const minLabel = minutes === 1 ? "minute" : "minutes";
    const secLabel = secs === 1 ? "second" : "seconds";
    return `${minutes} ${minLabel} ${secs} ${secLabel}`;
  }
  const secLabel = secs === 1 ? "second" : "seconds";
  return `${secs} ${secLabel}`;
};

export const useElapsedTime = (startedAt: string | null): ElapsedTimeResult => {
  const [seconds, setSeconds] = useState(0);

  useEffect(() => {
    if (!startedAt) return;

    const startMs = new Date(startedAt).getTime();

    if (Number.isNaN(startMs)) return;

    const computeElapsed = () => Math.max(0, Math.floor((Date.now() - startMs) / 1000));

    setSeconds(computeElapsed());

    const interval = setInterval(() => {
      setSeconds(computeElapsed());
    }, 1000);

    return () => clearInterval(interval);
  }, [startedAt]);

  if (!startedAt) {
    return EMPTY_RESULT;
  }

  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;
  const formatted = `${minutes}:${secs.toString().padStart(2, "0")}`;
  const ariaLabel = formatAriaLabel(minutes, secs);

  return { ariaLabel, formatted, seconds };
};
