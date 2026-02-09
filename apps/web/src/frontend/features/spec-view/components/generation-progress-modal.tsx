"use client";

import { AlertCircle, CheckCircle2, FileText, Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";

import { RotatingMessages } from "@/components/feedback/rotating-messages";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useElapsedTime } from "@/lib/hooks/use-elapsed-time";
import { cn } from "@/lib/utils";

import { GenerationPipeline } from "./generation-pipeline";
import { getStatusDisplayInfo } from "./generation-progress-utils";
import { useGenerationProgress } from "../hooks/use-generation-progress";
import type { SpecGenerationStatusEnum } from "../types";

type GenerationProgressContentProps = {
  onContinueBrowsing: () => void;
  onRetry?: () => void;
  startedAt: string | null;
  status: SpecGenerationStatusEnum | null;
};

// Goal: Provide progressive reassurance without message fatigue
const MESSAGE_INTERVALS = {
  PENDING: {
    INITIAL: 5, // Show "waiting for worker" for first 5 seconds
    QUEUED: 15, // Show "in queue" until 15 seconds
    // After 15s, show "starting soon"
  },
  RUNNING: {
    ANALYZING: 30, // First phase: structure analysis (0-30s)
    CLASSIFYING: 60, // Second phase: domain classification (30-60s)
    WRITING: 120, // Third phase: description generation (60-120s)
    // After 120s, show "almost done"
  },
} as const;

const getMessageIndex = (seconds: number, status: SpecGenerationStatusEnum | null): number => {
  // not_found: Backend job doesn't exist yet (polling started before job creation)
  // Treat as pending since user experience is identical (waiting in queue)
  if (!status || status === "pending" || status === "not_found") {
    if (seconds < MESSAGE_INTERVALS.PENDING.INITIAL) return 0;
    if (seconds < MESSAGE_INTERVALS.PENDING.QUEUED) return 1;
    return 2;
  }

  if (status === "running") {
    if (seconds < MESSAGE_INTERVALS.RUNNING.ANALYZING) return 0;
    if (seconds < MESSAGE_INTERVALS.RUNNING.CLASSIFYING) return 1;
    if (seconds < MESSAGE_INTERVALS.RUNNING.WRITING) return 2;
    return 3;
  }

  return 0;
};

const GenerationProgressContent = ({
  onContinueBrowsing,
  onRetry,
  startedAt,
  status,
}: GenerationProgressContentProps) => {
  const t = useTranslations("specView.generationProgress");
  const { ariaLabel, formatted, seconds } = useElapsedTime(startedAt);

  const isFailed = status === "failed";
  const isCompleted = status === "completed";
  // not_found included: Job may not exist yet if polling starts before backend creates it
  const isInProgress = status === "pending" || status === "running" || status === "not_found";
  const displayInfo = getStatusDisplayInfo(status);

  const messages: string[] = (() => {
    if (status === "pending" || status === "not_found" || !status) {
      return [t("pipeline.pending.0"), t("pipeline.pending.1"), t("pipeline.pending.2")];
    }
    if (status === "running") {
      return [
        t("pipeline.running.0"),
        t("pipeline.running.1"),
        t("pipeline.running.2"),
        t("pipeline.running.3"),
      ];
    }
    return [];
  })();

  const messageIndex = getMessageIndex(seconds, status);

  return (
    <div className="flex flex-col items-center gap-5 py-4">
      {/* Icon */}
      <div
        className={cn(
          "flex h-16 w-16 items-center justify-center rounded-full",
          isFailed ? "bg-destructive/10" : "bg-primary/10"
        )}
      >
        {isFailed ? (
          <AlertCircle className="h-8 w-8 text-destructive" />
        ) : isCompleted ? (
          <CheckCircle2 className="h-8 w-8 text-primary" />
        ) : (
          <div className="relative">
            <FileText className="h-8 w-8 text-primary" />
            <Sparkles className="absolute -right-1 -top-1 h-4 w-4 text-primary" />
          </div>
        )}
      </div>

      {/* Title, Description, and Elapsed Time */}
      <div className="text-center">
        <h3 className="text-lg font-semibold">{t(displayInfo.titleKey)}</h3>
        {!isInProgress && (
          <p className="mt-1 text-sm text-muted-foreground">{t(displayInfo.descriptionKey)}</p>
        )}
        {isInProgress && startedAt && (
          <div className="mt-2 flex items-center justify-center gap-1.5 text-xs text-muted-foreground">
            <span>{t("elapsed")}</span>
            <span aria-hidden>·</span>
            <time aria-label={ariaLabel} className="font-mono tabular-nums" dateTime={startedAt}>
              {formatted}
            </time>
          </div>
        )}
      </div>

      {/* Pipeline visualization (only when in progress) */}
      {isInProgress && (
        <div className="flex w-full flex-col gap-4" role="status">
          <GenerationPipeline status={status} />
          <RotatingMessages
            className="text-center"
            currentIndex={messageIndex}
            messages={messages}
          />
        </div>
      )}

      {/* Actions */}
      <div className="flex w-full flex-col gap-2">
        {isFailed && onRetry && (
          <Button className="w-full" onClick={onRetry} variant="default">
            {t("retryButton")}
          </Button>
        )}
        {!isCompleted && (
          <Button className="w-full" onClick={onContinueBrowsing} variant="outline">
            {t("continueBrowsing")}
          </Button>
        )}
        {isCompleted && (
          <Button className="w-full" onClick={onContinueBrowsing} variant="default">
            {t("viewDocument")}
          </Button>
        )}
      </div>

      {/* Tip */}
      {isInProgress && <p className="text-center text-xs text-muted-foreground">{t("tip")}</p>}
    </div>
  );
};

type GenerationProgressModalProps = {
  startedAt?: string | null;
  status: SpecGenerationStatusEnum | null;
};

export const GenerationProgressModal = ({
  startedAt = null,
  status,
}: GenerationProgressModalProps) => {
  const { close, isOpen, switchToBackground } = useGenerationProgress();

  const handleContinueBrowsing = () => {
    if (status === "completed") {
      close();
    } else {
      switchToBackground();
    }
  };

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      if (status === "pending" || status === "running") {
        switchToBackground();
      } else {
        close();
      }
    }
  };

  return (
    <Dialog onOpenChange={handleOpenChange} open={isOpen}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader className="sr-only">
          <DialogTitle>Generation Progress</DialogTitle>
          <DialogDescription>Track the progress of your spec document generation</DialogDescription>
        </DialogHeader>
        <GenerationProgressContent
          onContinueBrowsing={handleContinueBrowsing}
          startedAt={startedAt}
          status={status}
        />
      </DialogContent>
    </Dialog>
  );
};
