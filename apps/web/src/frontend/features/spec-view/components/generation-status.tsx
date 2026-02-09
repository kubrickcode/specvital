"use client";

import { useTranslations } from "next-intl";

import { PulseRing, RotatingMessages, ShimmerBar } from "@/components/feedback";
import type { ShimmerBarColor } from "@/components/feedback";
import { useElapsedTime } from "@/lib/hooks";
import { cn } from "@/lib/utils";

import type { SpecGenerationStatusEnum } from "../types";

type GenerationStatusProps = {
  startedAt: string | null;
  status: SpecGenerationStatusEnum;
};

type ActiveStatus = "pending" | "running";

const STATUS_CONFIG: Record<
  ActiveStatus,
  {
    color: string;
    labelKey: string;
    shimmerColor: ShimmerBarColor;
    shimmerDuration: number;
  }
> = {
  pending: {
    color: "text-chart-2",
    labelKey: "status.pending",
    shimmerColor: "var(--chart-2)",
    shimmerDuration: 3,
  },
  running: {
    color: "text-[hsl(var(--ai-primary))]",
    labelKey: "status.running",
    shimmerColor: "var(--ai-primary)",
    shimmerDuration: 2,
  },
};

const MESSAGE_INTERVALS_SEC = [0, 5, 15, 30, 60];

const calculateMessageIndex = (elapsedSeconds: number, maxIndex: number): number => {
  for (let i = MESSAGE_INTERVALS_SEC.length - 1; i >= 0; i--) {
    if (elapsedSeconds >= MESSAGE_INTERVALS_SEC[i]) {
      return Math.min(i, maxIndex);
    }
  }
  return 0;
};

export const GenerationStatus = ({ startedAt, status }: GenerationStatusProps) => {
  const t = useTranslations("specView.generationProgress");
  const elapsed = useElapsedTime(startedAt);

  const activeStatus: ActiveStatus = status === "running" ? "running" : "pending";
  const config = STATUS_CONFIG[activeStatus];

  const messages =
    activeStatus === "pending"
      ? [t("pipeline.pending.0"), t("pipeline.pending.1"), t("pipeline.pending.2")]
      : [
          t("pipeline.running.0"),
          t("pipeline.running.1"),
          t("pipeline.running.2"),
          t("pipeline.running.3"),
        ];

  const messageIndex = elapsed ? calculateMessageIndex(elapsed.seconds, messages.length - 1) : 0;

  return (
    <div aria-live="polite" className="rounded-lg border bg-card p-6 space-y-4" role="status">
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-3 min-w-0">
          <PulseRing className={cn("shrink-0", config.color)} size="sm" />
          <span className="font-medium text-foreground truncate">{t(config.labelKey)}</span>
        </div>
        {elapsed && (
          <time
            aria-label={elapsed.ariaLabel}
            className="text-sm text-muted-foreground tabular-nums shrink-0"
          >
            {t("elapsed")}: {elapsed.formatted}
          </time>
        )}
      </div>

      <ShimmerBar color={config.shimmerColor} duration={config.shimmerDuration} height="xs" />

      <div aria-atomic="true" aria-live="polite" className="min-h-[24px]" role="status">
        <RotatingMessages currentIndex={messageIndex} messages={messages} />
      </div>
    </div>
  );
};
