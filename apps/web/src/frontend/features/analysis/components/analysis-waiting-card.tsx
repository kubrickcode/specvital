"use client";

import { m } from "motion/react";
import { useTranslations } from "next-intl";

import { PulseRing, RotatingMessages, ShimmerBar } from "@/components/feedback";
import type { ShimmerBarColor } from "@/components/feedback";
import { useElapsedTime } from "@/lib/hooks";
import { cn } from "@/lib/utils";

import type { WaitingStatus } from "../types";

type AnalysisWaitingCardProps = {
  owner: string;
  repo: string;
  startedAt: string | null;
  status: WaitingStatus;
};

const STATUS_CONFIG: Record<
  WaitingStatus,
  {
    color: string;
    label: string;
    shimmerColor: ShimmerBarColor;
    shimmerDuration: number;
  }
> = {
  analyzing: {
    color: "text-chart-1",
    label: "analyzingTitle",
    shimmerColor: "var(--chart-1)",
    shimmerDuration: 2,
  },
  queued: {
    color: "text-chart-2",
    label: "queuedTitle",
    shimmerColor: "var(--chart-2)",
    shimmerDuration: 3,
  },
};

const MESSAGES = {
  analyzing: ["0", "1", "2", "3", "4"],
  queued: ["0", "1", "2", "3", "4"],
};

const MESSAGE_INTERVALS = [0, 5000, 15000, 30000, 60000];

const calculateMessageIndex = (elapsedSeconds: number): number => {
  const elapsedMs = elapsedSeconds * 1000;
  for (let i = MESSAGE_INTERVALS.length - 1; i >= 0; i--) {
    if (elapsedMs >= MESSAGE_INTERVALS[i]) {
      return i;
    }
  }
  return 0;
};

export const AnalysisWaitingCard = ({
  owner,
  repo,
  startedAt,
  status,
}: AnalysisWaitingCardProps) => {
  const t = useTranslations("analyze");
  const elapsed = useElapsedTime(startedAt);
  const config = STATUS_CONFIG[status];

  const messageIndex = elapsed ? calculateMessageIndex(elapsed.seconds) : 0;
  const showLongWaitGuidance = elapsed && elapsed.seconds >= 60;

  const messages = MESSAGES[status].map((key) => t(`waiting.${status}.${key}`));
  const statusLabel = t(`status.${config.label}`);
  const elapsedLabel = t("waiting.elapsed");
  const longWaitText = t("waiting.longWait");

  return (
    <main className="container mx-auto px-4 py-8">
      <m.div
        animate={{ opacity: 1 }}
        className="space-y-6"
        initial={{ opacity: 0 }}
        transition={{ duration: 0.3 }}
      >
        {/* Repository header */}
        <header>
          <h1 className="text-2xl font-bold text-foreground">
            {owner}/{repo}
          </h1>
        </header>

        {/* Waiting status card */}
        <div aria-live="polite" className="rounded-lg border bg-card p-6 space-y-4" role="status">
          {/* Status indicator row */}
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-3 min-w-0">
              <PulseRing className={cn("shrink-0", config.color)} size="sm" />
              <span className="font-medium text-foreground truncate">{statusLabel}</span>
            </div>
            {elapsed && (
              <time
                aria-label={elapsed.ariaLabel}
                className="text-sm text-muted-foreground tabular-nums shrink-0"
              >
                {elapsedLabel}: {elapsed.formatted}
              </time>
            )}
          </div>

          {/* ShimmerBar */}
          <ShimmerBar color={config.shimmerColor} duration={config.shimmerDuration} height="xs" />

          {/* Rotating messages */}
          <div aria-atomic="true" aria-live="polite" className="min-h-[24px]" role="status">
            <RotatingMessages currentIndex={messageIndex} messages={messages} />
          </div>

          {/* Long wait guidance */}
          {showLongWaitGuidance && (
            <m.div
              animate={{ opacity: 1, y: 0 }}
              className="text-sm text-muted-foreground border-t pt-4"
              initial={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.3 }}
            >
              {longWaitText}
            </m.div>
          )}
        </div>
      </m.div>
    </main>
  );
};
