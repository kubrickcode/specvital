import { Check, Circle, CircleDashed, Crosshair, XCircle } from "lucide-react";

import type { TestStatus } from "@/lib/api/types";

export type StatusConfigItem = {
  bgColor: string;
  color: string;
  descriptionKey: string;
  icon: typeof Check;
  labelKey: string;
  ringColor: string;
};

export const STATUS_CONFIG: Record<TestStatus, StatusConfigItem> = {
  active: {
    bgColor: "bg-emerald-50 dark:bg-emerald-950/40",
    color: "text-emerald-600 dark:text-emerald-400",
    descriptionKey: "activeDescription",
    icon: Check,
    labelKey: "active",
    ringColor: "ring-emerald-200 dark:ring-emerald-800",
  },
  focused: {
    bgColor: "bg-violet-50 dark:bg-violet-950/40",
    color: "text-violet-600 dark:text-violet-400",
    descriptionKey: "focusedDescription",
    icon: Crosshair,
    labelKey: "focused",
    ringColor: "ring-violet-200 dark:ring-violet-800",
  },
  skipped: {
    bgColor: "bg-amber-50 dark:bg-amber-950/40",
    color: "text-amber-600 dark:text-amber-400",
    descriptionKey: "skippedDescription",
    icon: CircleDashed,
    labelKey: "skipped",
    ringColor: "ring-amber-200 dark:ring-amber-800",
  },
  todo: {
    bgColor: "bg-sky-50 dark:bg-sky-950/40",
    color: "text-sky-600 dark:text-sky-400",
    descriptionKey: "todoDescription",
    icon: Circle,
    labelKey: "todo",
    ringColor: "ring-sky-200 dark:ring-sky-800",
  },
  xfail: {
    bgColor: "bg-red-50 dark:bg-red-950/40",
    color: "text-red-500 dark:text-red-400",
    descriptionKey: "xfailDescription",
    icon: XCircle,
    labelKey: "xfail",
    ringColor: "ring-red-200 dark:ring-red-800",
  },
} as const;

export const STATUS_ORDER: TestStatus[] = ["active", "skipped", "todo", "focused", "xfail"];
