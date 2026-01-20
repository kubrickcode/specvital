"use client";

import { AlertCircle, CheckCircle2, FileText, Loader2, Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";

import { getStatusDisplayInfo } from "./generation-progress-utils";
import { useGenerationProgress } from "../hooks/use-generation-progress";
import type { SpecGenerationStatusEnum } from "../types";

type GenerationProgressContentProps = {
  onContinueBrowsing: () => void;
  onRetry?: () => void;
  status: SpecGenerationStatusEnum | null;
};

const GenerationProgressContent = ({
  onContinueBrowsing,
  onRetry,
  status,
}: GenerationProgressContentProps) => {
  const t = useTranslations("specView.generationProgress");

  const isFailed = status === "failed";
  const isCompleted = status === "completed";
  const isInProgress = status === "pending" || status === "running";
  const displayInfo = getStatusDisplayInfo(status);

  return (
    <div className="flex flex-col items-center gap-6 py-4">
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

      {/* Title and Description */}
      <div className="text-center">
        <h3 className="text-lg font-semibold">{t(displayInfo.titleKey)}</h3>
        <p className="mt-1 text-sm text-muted-foreground">{t(displayInfo.descriptionKey)}</p>
      </div>

      {/* Spinner + Status (only when in progress) */}
      {isInProgress && (
        <div className="flex flex-col items-center gap-3">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="text-center text-sm text-muted-foreground">{t("estimatedTime")}</p>
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

export const GenerationProgressModal = () => {
  const { close, isOpen, status, switchToBackground } = useGenerationProgress();

  const handleContinueBrowsing = () => {
    if (status === "completed") {
      close();
    } else {
      switchToBackground();
    }
  };

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      // When closing modal, switch to background instead of fully closing
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
        <GenerationProgressContent onContinueBrowsing={handleContinueBrowsing} status={status} />
      </DialogContent>
    </Dialog>
  );
};
